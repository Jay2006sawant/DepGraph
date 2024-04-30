package graph

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// StorageManager handles efficient graph storage
type StorageManager struct {
	dir      string
	nodeFile *os.File
	edgeFile *os.File
	mutex    sync.RWMutex
	nodeMap  map[string]int64 // Maps node IDs to file offsets
}

// NewStorageManager creates a new storage manager
func NewStorageManager(dir string) (*StorageManager, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	nodeFile, err := os.OpenFile(filepath.Join(dir, "nodes.bin"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open node file: %w", err)
	}

	edgeFile, err := os.OpenFile(filepath.Join(dir, "edges.bin"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		nodeFile.Close()
		return nil, fmt.Errorf("failed to open edge file: %w", err)
	}

	return &StorageManager{
		dir:      dir,
		nodeFile: nodeFile,
		edgeFile: edgeFile,
		nodeMap:  make(map[string]int64),
	}, nil
}

// Close closes the storage files
func (sm *StorageManager) Close() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if err := sm.nodeFile.Close(); err != nil {
		return err
	}
	return sm.edgeFile.Close()
}

// SaveNode saves a node to storage
func (sm *StorageManager) SaveNode(node *Node) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Get current position
	offset, err := sm.nodeFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// Write node data
	if err := sm.writeString(sm.nodeFile, node.ID); err != nil {
		return err
	}
	if err := sm.writeString(sm.nodeFile, node.Label); err != nil {
		return err
	}
	if err := sm.writeString(sm.nodeFile, node.Type); err != nil {
		return err
	}
	if err := sm.writeString(sm.nodeFile, node.Version); err != nil {
		return err
	}

	// Store offset in map
	sm.nodeMap[node.ID] = offset
	return nil
}

// LoadNode loads a node from storage
func (sm *StorageManager) LoadNode(id string) (*Node, error) {
	sm.mutex.RLock()
	offset, ok := sm.nodeMap[id]
	sm.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("node not found: %s", id)
	}

	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if _, err := sm.nodeFile.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	node := &Node{}
	var err error

	node.ID, err = sm.readString(sm.nodeFile)
	if err != nil {
		return nil, err
	}

	node.Label, err = sm.readString(sm.nodeFile)
	if err != nil {
		return nil, err
	}

	node.Type, err = sm.readString(sm.nodeFile)
	if err != nil {
		return nil, err
	}

	node.Version, err = sm.readString(sm.nodeFile)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// SaveEdge saves an edge to storage
func (sm *StorageManager) SaveEdge(edge *Edge) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if err := sm.writeString(sm.edgeFile, edge.Source); err != nil {
		return err
	}
	if err := sm.writeString(sm.edgeFile, edge.Target); err != nil {
		return err
	}
	return sm.writeString(sm.edgeFile, edge.Version)
}

// LoadEdges loads all edges from storage
func (sm *StorageManager) LoadEdges() ([]*Edge, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if _, err := sm.edgeFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var edges []*Edge
	for {
		edge := &Edge{}
		var err error

		edge.Source, err = sm.readString(sm.edgeFile)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		edge.Target, err = sm.readString(sm.edgeFile)
		if err != nil {
			return nil, err
		}

		edge.Version, err = sm.readString(sm.edgeFile)
		if err != nil {
			return nil, err
		}

		edges = append(edges, edge)
	}

	return edges, nil
}

// Helper functions for reading/writing strings
func (sm *StorageManager) writeString(f *os.File, s string) error {
	data := []byte(s)
	if err := binary.Write(f, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := f.Write(data)
	return err
}

func (sm *StorageManager) readString(f *os.File) (string, error) {
	var length uint32
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		return "", err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(f, data); err != nil {
		return "", err
	}

	return string(data), nil
} 