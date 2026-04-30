package database

import (
	"database/sql"
	"fmt"

	"github.com/dengdeng-harmonyos/server/internal/config"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) InitTables() error {
	queries := []string{
		// 设备表（简化版，去除用户关联）
		`CREATE TABLE IF NOT EXISTS devices (
			id SERIAL PRIMARY KEY,
			device_id VARCHAR(64) UNIQUE NOT NULL,
			push_token TEXT NOT NULL UNIQUE,
			device_type VARCHAR(50),
			os_version VARCHAR(50),
			app_version VARCHAR(50),
			is_active BOOLEAN DEFAULT TRUE,
			last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// 推送统计表（仅统计数据，不记录具体内容）
		`CREATE TABLE IF NOT EXISTS push_statistics (
			id SERIAL PRIMARY KEY,
			date DATE NOT NULL,
			push_type VARCHAR(20) NOT NULL,
			total_count INTEGER DEFAULT 0,
			success_count INTEGER DEFAULT 0,
			failed_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_date_type UNIQUE(date, push_type)
		)`,

		// 创建索引
		`CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_is_active ON devices(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_push_stats_date ON push_statistics(date)`,

		// App更新策略表
		`CREATE TABLE IF NOT EXISTS app_update_policies (
			platform VARCHAR(32) PRIMARY KEY DEFAULT 'harmonyos',
			latest_version_code BIGINT NOT NULL DEFAULT 0,
			latest_version_name VARCHAR(64) NOT NULL DEFAULT '',
			min_version_code BIGINT NOT NULL DEFAULT 0,
			force_update BOOLEAN NOT NULL DEFAULT TRUE,
			store_url TEXT NOT NULL DEFAULT 'store://appgallery.huawei.com/app/detail?id=top.yidingyaojizhu.dengdeng',
			release_notes TEXT NOT NULL DEFAULT '',
			enabled BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT INTO app_update_policies (platform)
			VALUES ('harmonyos')
			ON CONFLICT (platform) DO NOTHING`,
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW.updated_at = CURRENT_TIMESTAMP;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql`,
		`DROP TRIGGER IF EXISTS update_app_update_policies_updated_at ON app_update_policies`,
		`CREATE TRIGGER update_app_update_policies_updated_at
			BEFORE UPDATE ON app_update_policies
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column()`,
	}

	for _, query := range queries {
		if _, err := db.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
