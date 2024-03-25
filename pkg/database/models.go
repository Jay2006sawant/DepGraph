package database

import (
	"time"
)

// Repository represents a GitHub repository and its dependencies
type Repository struct {
	ID           int64     `db:"id"`
	Name         string    `db:"name"`
	Organization string    `db:"organization"`
	LastScanned  time.Time `db:"last_scanned"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Dependency represents a module dependency found in repositories
type Dependency struct {
	ID           int64     `db:"id"`
	RepositoryID int64     `db:"repository_id"`
	Module       string    `db:"module"`
	Version      string    `db:"version"`
	IsIndirect   bool      `db:"is_indirect"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// DependencyOverlap represents shared dependencies between repositories
type DependencyOverlap struct {
	ID              int64     `db:"id"`
	Module          string    `db:"module"`
	RepositoryCount int       `db:"repository_count"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// ScanHistory tracks repository scan history
type ScanHistory struct {
	ID           int64     `db:"id"`
	RepositoryID int64     `db:"repository_id"`
	Status       string    `db:"status"`
	Error        string    `db:"error"`
	StartedAt    time.Time `db:"started_at"`
	CompletedAt  time.Time `db:"completed_at"`
	CreatedAt    time.Time `db:"created_at"`
}

// Schema returns the SQL schema for creating the database tables
func Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		organization TEXT NOT NULL,
		last_scanned DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(name, organization)
	);

	CREATE TABLE IF NOT EXISTS dependencies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		repository_id INTEGER NOT NULL,
		module TEXT NOT NULL,
		version TEXT NOT NULL,
		is_indirect BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY(repository_id) REFERENCES repositories(id),
		UNIQUE(repository_id, module, version)
	);

	CREATE TABLE IF NOT EXISTS dependency_overlaps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		module TEXT NOT NULL,
		repository_count INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(module)
	);

	CREATE TABLE IF NOT EXISTS scan_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		repository_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		error TEXT,
		started_at DATETIME NOT NULL,
		completed_at DATETIME,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(repository_id) REFERENCES repositories(id)
	);

	CREATE INDEX IF NOT EXISTS idx_repositories_org ON repositories(organization);
	CREATE INDEX IF NOT EXISTS idx_dependencies_module ON dependencies(module);
	CREATE INDEX IF NOT EXISTS idx_scan_history_repo ON scan_history(repository_id);
	`
} 