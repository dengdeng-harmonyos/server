package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/dengdeng-harmonyos/server/internal/logger"
)

const expiredMessageCleanupSQL = `
DELETE FROM pending_messages
WHERE expires_at < NOW()
   OR (delivered = true AND confirmed_at < NOW() - INTERVAL '24 hours')
`

// CleanExpiredMessages removes pending messages that can no longer be delivered.
func CleanExpiredMessages(ctx context.Context, db *sql.DB) (int64, error) {
	result, err := db.ExecContext(ctx, expiredMessageCleanupSQL)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// StartExpiredMessageCleanup runs cleanup immediately, then repeats on interval.
func StartExpiredMessageCleanup(ctx context.Context, db *sql.DB, interval time.Duration) context.CancelFunc {
	cleanupCtx, cancel := context.WithCancel(ctx)

	go func() {
		runCleanup(cleanupCtx, db)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-cleanupCtx.Done():
				return
			case <-ticker.C:
				runCleanup(cleanupCtx, db)
			}
		}
	}()

	return cancel
}

func runCleanup(ctx context.Context, db *sql.DB) {
	deleted, err := CleanExpiredMessages(ctx, db)
	if err != nil {
		logger.Error("Expired pending message cleanup failed: %v", err)
		return
	}
	if deleted > 0 {
		logger.Info("Expired pending message cleanup removed %d rows", deleted)
	}
}
