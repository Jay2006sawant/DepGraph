package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Manager handles database operations
type Manager struct {
	db *sql.DB
}

// NewManager creates a new database manager
func NewManager(dbPath string) (*Manager, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	manager := &Manager{db: db}
	if err := manager.Initialize(); err != nil {
		return nil, err
	}

	return manager, nil
}

// Initialize creates the database schema if it doesn't exist
func (m *Manager) Initialize() error {
	_, err := m.db.Exec(Schema())
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	return nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	return m.db.Close()
}

// SaveRepository saves or updates a repository
func (m *Manager) SaveRepository(repo *Repository) error {
	now := time.Now()
	if repo.ID == 0 {
		repo.CreatedAt = now
	}
	repo.UpdatedAt = now

	query := `
		INSERT INTO repositories (name, organization, last_scanned, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(name, organization) DO UPDATE SET
			last_scanned = excluded.last_scanned,
			updated_at = excluded.updated_at
	`

	result, err := m.db.Exec(query,
		repo.Name,
		repo.Organization,
		repo.LastScanned,
		repo.CreatedAt,
		repo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save repository: %v", err)
	}

	if repo.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %v", err)
		}
		repo.ID = id
	}

	return nil
}

// SaveDependency saves or updates a dependency
func (m *Manager) SaveDependency(dep *Dependency) error {
	now := time.Now()
	if dep.ID == 0 {
		dep.CreatedAt = now
	}
	dep.UpdatedAt = now

	query := `
		INSERT INTO dependencies (repository_id, module, version, is_indirect, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(repository_id, module, version) DO UPDATE SET
			is_indirect = excluded.is_indirect,
			updated_at = excluded.updated_at
	`

	result, err := m.db.Exec(query,
		dep.RepositoryID,
		dep.Module,
		dep.Version,
		dep.IsIndirect,
		dep.CreatedAt,
		dep.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save dependency: %v", err)
	}

	if dep.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %v", err)
		}
		dep.ID = id
	}

	return nil
}

// UpdateDependencyOverlaps updates the dependency overlap counts
func (m *Manager) UpdateDependencyOverlaps() error {
	query := `
		INSERT INTO dependency_overlaps (module, repository_count, created_at, updated_at)
		SELECT 
			module,
			COUNT(DISTINCT repository_id) as repo_count,
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		FROM dependencies
		GROUP BY module
		HAVING repo_count > 1
		ON CONFLICT(module) DO UPDATE SET
			repository_count = excluded.repository_count,
			updated_at = excluded.updated_at
	`

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to update dependency overlaps: %v", err)
	}

	return nil
}

// LogScan logs a repository scan event
func (m *Manager) LogScan(scan *ScanHistory) error {
	scan.CreatedAt = time.Now()

	query := `
		INSERT INTO scan_history (repository_id, status, error, started_at, completed_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := m.db.Exec(query,
		scan.RepositoryID,
		scan.Status,
		scan.Error,
		scan.StartedAt,
		scan.CompletedAt,
		scan.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to log scan: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}
	scan.ID = id

	return nil
} 