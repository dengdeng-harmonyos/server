package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dengdeng-harmonyos/server/internal/logger"
)

const defaultAppStoreURL = "store://appgallery.huawei.com/app/detail?id=top.yidingyaojizhu.dengdeng"

type AppUpdateReleaseManifest struct {
	Platform       string `json:"platform"`
	VersionCode    int64  `json:"versionCode"`
	VersionName    string `json:"versionName"`
	MinVersionCode int64  `json:"minVersionCode"`
	ForceUpdate    *bool  `json:"forceUpdate"`
	StoreURL       string `json:"storeUrl"`
	ReleaseNotes   string `json:"releaseNotes"`
	Enabled        bool   `json:"enabled"`
}

type normalizedAppUpdateRelease struct {
	Platform       string
	VersionCode    int64
	VersionName    string
	MinVersionCode int64
	ForceUpdate    bool
	StoreURL       string
	ReleaseNotes   string
	Enabled        bool
}

func SyncAppUpdatePolicyFromManifest(ctx context.Context, db *sql.DB, policyFile string) error {
	policyFile = strings.TrimSpace(policyFile)
	if policyFile == "" {
		logger.Info("App update policy file is empty, skip release seed")
		return nil
	}

	release, found, err := loadAppUpdateReleaseManifest(policyFile)
	if err != nil {
		return err
	}
	if !found {
		logger.Info("App update policy file not found, skip release seed: %s", policyFile)
		return nil
	}
	if release.VersionCode <= 0 {
		logger.Info("App update policy manifest has no releasable version, skip release seed")
		return nil
	}

	if err := validateAppUpdateRelease(release); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin app update release seed transaction: %w", err)
	}
	defer tx.Rollback()

	if err := upsertAppUpdateRelease(ctx, tx, release); err != nil {
		return err
	}

	if release.Enabled {
		activated, err := upsertCurrentAppUpdatePolicy(ctx, tx, release)
		if err != nil {
			return err
		}
		if activated {
			logger.Info("App update policy activated from manifest: platform=%s versionCode=%d", release.Platform, release.VersionCode)
		} else {
			logger.Info("App update policy manifest is not newer than current policy, current policy kept: platform=%s versionCode=%d", release.Platform, release.VersionCode)
		}
	} else {
		logger.Info("App update release recorded but not activated: platform=%s versionCode=%d", release.Platform, release.VersionCode)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit app update release seed transaction: %w", err)
	}
	return nil
}

func loadAppUpdateReleaseManifest(policyFile string) (normalizedAppUpdateRelease, bool, error) {
	content, err := os.ReadFile(policyFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return normalizedAppUpdateRelease{}, false, nil
		}
		return normalizedAppUpdateRelease{}, false, fmt.Errorf("read app update policy file: %w", err)
	}

	var manifest AppUpdateReleaseManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return normalizedAppUpdateRelease{}, false, fmt.Errorf("parse app update policy file: %w", err)
	}

	return normalizeAppUpdateRelease(manifest), true, nil
}

func normalizeAppUpdateRelease(manifest AppUpdateReleaseManifest) normalizedAppUpdateRelease {
	platform := strings.TrimSpace(manifest.Platform)
	if platform == "" {
		platform = "harmonyos"
	}

	storeURL := strings.TrimSpace(manifest.StoreURL)
	if storeURL == "" {
		storeURL = defaultAppStoreURL
	}

	forceUpdate := true
	if manifest.ForceUpdate != nil {
		forceUpdate = *manifest.ForceUpdate
	}

	return normalizedAppUpdateRelease{
		Platform:       platform,
		VersionCode:    manifest.VersionCode,
		VersionName:    strings.TrimSpace(manifest.VersionName),
		MinVersionCode: manifest.MinVersionCode,
		ForceUpdate:    forceUpdate,
		StoreURL:       storeURL,
		ReleaseNotes:   strings.TrimSpace(manifest.ReleaseNotes),
		Enabled:        manifest.Enabled,
	}
}

func validateAppUpdateRelease(release normalizedAppUpdateRelease) error {
	if release.VersionName == "" {
		return fmt.Errorf("app update policy versionName is required when versionCode is %d", release.VersionCode)
	}
	if release.MinVersionCode < 0 {
		return fmt.Errorf("app update policy minVersionCode cannot be negative")
	}
	if release.MinVersionCode > release.VersionCode {
		return fmt.Errorf("app update policy minVersionCode cannot exceed versionCode")
	}
	return nil
}

func upsertAppUpdateRelease(ctx context.Context, tx *sql.Tx, release normalizedAppUpdateRelease) error {
	const query = `
		INSERT INTO app_update_releases (
			platform, version_code, version_name, min_version_code, force_update,
			store_url, release_notes, enabled, source
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'manifest')
		ON CONFLICT (platform, version_code) DO UPDATE SET
			version_name = EXCLUDED.version_name,
			min_version_code = EXCLUDED.min_version_code,
			force_update = EXCLUDED.force_update,
			store_url = EXCLUDED.store_url,
			release_notes = EXCLUDED.release_notes,
			enabled = EXCLUDED.enabled,
			source = EXCLUDED.source,
			updated_at = CURRENT_TIMESTAMP
	`

	if _, err := tx.ExecContext(ctx, query,
		release.Platform,
		release.VersionCode,
		release.VersionName,
		release.MinVersionCode,
		release.ForceUpdate,
		release.StoreURL,
		release.ReleaseNotes,
		release.Enabled,
	); err != nil {
		return fmt.Errorf("upsert app update release: %w", err)
	}
	return nil
}

func upsertCurrentAppUpdatePolicy(ctx context.Context, tx *sql.Tx, release normalizedAppUpdateRelease) (bool, error) {
	const query = `
		INSERT INTO app_update_policies (
			platform, latest_version_code, latest_version_name, min_version_code,
			force_update, store_url, release_notes, enabled
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)
		ON CONFLICT (platform) DO UPDATE SET
			latest_version_code = EXCLUDED.latest_version_code,
			latest_version_name = EXCLUDED.latest_version_name,
			min_version_code = EXCLUDED.min_version_code,
			force_update = EXCLUDED.force_update,
			store_url = EXCLUDED.store_url,
			release_notes = EXCLUDED.release_notes,
			enabled = TRUE,
			updated_at = CURRENT_TIMESTAMP
		WHERE app_update_policies.latest_version_code <= EXCLUDED.latest_version_code
	`

	result, err := tx.ExecContext(ctx, query,
		release.Platform,
		release.VersionCode,
		release.VersionName,
		release.MinVersionCode,
		release.ForceUpdate,
		release.StoreURL,
		release.ReleaseNotes,
	)
	if err != nil {
		return false, fmt.Errorf("upsert current app update policy: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return true, nil
	}
	return affected > 0, nil
}
