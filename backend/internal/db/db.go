package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func DefaultConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "geo"),
		Password: getEnv("DB_PASSWORD", "geopass"),
		DBName:   getEnv("DB_NAME", "geoproject"),
	}
}

func Connect(cfg Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("db open: %w", err)
	}

	if err = DB.Ping(); err != nil {
		DB = nil
		return fmt.Errorf("db ping: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return migrate()
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tracker_history (
		id SERIAL PRIMARY KEY,
		topic TEXT NOT NULL,
		page_views INT DEFAULT 0,
		wiki_page_id INT DEFAULT 0,
		description TEXT DEFAULT '',
		checked_at TIMESTAMP DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_tracker_topic ON tracker_history(topic);

	CREATE TABLE IF NOT EXISTS audit_history (
		id SERIAL PRIMARY KEY,
		brand TEXT NOT NULL,
		total_score INT DEFAULT 0,
		label TEXT DEFAULT '',
		details JSONB DEFAULT '{}',
		audited_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS review_cache (
		id SERIAL PRIMARY KEY,
		business TEXT NOT NULL,
		sentiment JSONB DEFAULT '{}',
		review_count INT DEFAULT 0,
		cached_at TIMESTAMP DEFAULT NOW()
	);
	`
	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	log.Println("Database migrated")
	return nil
}

func SaveTrackerRecord(topic string, pageViews, wikiPageID int, description string) error {
	if DB == nil {
		return nil
	}
	_, err := DB.Exec(
		`INSERT INTO tracker_history (topic, page_views, wiki_page_id, description) VALUES ($1, $2, $3, $4)`,
		topic, pageViews, wikiPageID, description,
	)
	return err
}

func GetTrackerHistory(topic string, limit int) ([]struct {
	Topic      string
	PageViews  int
	CheckedAt  string
}, error) {
	if DB == nil {
		return nil, nil
	}
	rows, err := DB.Query(
		`SELECT topic, page_views, checked_at FROM tracker_history WHERE topic = $1 ORDER BY checked_at DESC LIMIT $2`,
		topic, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []struct {
		Topic      string
		PageViews  int
		CheckedAt  string
	}
	for rows.Next() {
		var r struct {
			Topic      string
			PageViews  int
			CheckedAt  string
		}
		if err := rows.Scan(&r.Topic, &r.PageViews, &r.CheckedAt); err == nil {
			result = append(result, r)
		}
	}
	return result, nil
}

func IsConnected() bool {
	return DB != nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
