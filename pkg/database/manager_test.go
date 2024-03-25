package database

import (
	"os"
	"testing"
	"time"
)

func TestDatabaseOperations(t *testing.T) {
	// Create a temporary database file
	dbPath := "test.db"
	defer os.Remove(dbPath)

	// Create manager
	manager, err := NewManager(dbPath)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test saving repository
	repo := &Repository{
		Name:         "test-repo",
		Organization: "test-org",
		LastScanned:  time.Now(),
	}

	if err := manager.SaveRepository(repo); err != nil {
		t.Errorf("Failed to save repository: %v", err)
	}

	if repo.ID == 0 {
		t.Error("Repository ID should not be 0 after save")
	}

	// Test saving dependency
	dep := &Dependency{
		RepositoryID: repo.ID,
		Module:      "github.com/test/module",
		Version:     "v1.0.0",
		IsIndirect:  false,
	}

	if err := manager.SaveDependency(dep); err != nil {
		t.Errorf("Failed to save dependency: %v", err)
	}

	if dep.ID == 0 {
		t.Error("Dependency ID should not be 0 after save")
	}

	// Test updating dependency overlaps
	if err := manager.UpdateDependencyOverlaps(); err != nil {
		t.Errorf("Failed to update dependency overlaps: %v", err)
	}

	// Test logging scan
	scan := &ScanHistory{
		RepositoryID: repo.ID,
		Status:      "completed",
		StartedAt:   time.Now(),
		CompletedAt: time.Now(),
	}

	if err := manager.LogScan(scan); err != nil {
		t.Errorf("Failed to log scan: %v", err)
	}

	if scan.ID == 0 {
		t.Error("Scan ID should not be 0 after save")
	}
} 