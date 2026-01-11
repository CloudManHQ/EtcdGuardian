/*
Copyright 2026 EtcdGuardian Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"context"
	"os"
	"testing"

	"github.com/go-logr/logr"
)

func TestValidator_ValidateSnapshot(t *testing.T) {
	validator := NewValidator(logr.Discard())

	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test-snapshot-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some test data
	testData := []byte("test snapshot data")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	ctx := context.Background()
	result, err := validator.ValidateSnapshot(ctx, tmpFile.Name())

	if err != nil {
		t.Fatalf("ValidateSnapshot failed: %v", err)
	}

	if !result.Valid {
		t.Error("Expected valid snapshot")
	}

	if result.Hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestValidator_ValidateSnapshot_NotFound(t *testing.T) {
	validator := NewValidator(logr.Discard())

	ctx := context.Background()
	result, err := validator.ValidateSnapshot(ctx, "/nonexistent/file.db")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result for nonexistent file")
	}
}

func TestValidator_CalculateHash(t *testing.T) {
	validator := NewValidator(logr.Discard())

	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test-hash-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := []byte("consistent test data")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Calculate hash twice
	hash1, err := validator.calculateHash(tmpFile.Name())
	if err != nil {
		t.Fatalf("First hash calculation failed: %v", err)
	}

	hash2, err := validator.calculateHash(tmpFile.Name())
	if err != nil {
		t.Fatalf("Second hash calculation failed: %v", err)
	}

	// Hashes should be consistent
	if hash1 != hash2 {
		t.Errorf("Hash mismatch: %s != %s", hash1, hash2)
	}

	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}
}
/*
Copyright 2026 EtcdGuardian Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"context"
	"os"
	"testing"

	"github.com/go-logr/logr"
)

func TestValidator_ValidateSnapshot(t *testing.T) {
	validator := NewValidator(logr.Discard())

	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test-snapshot-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some test data
	testData := []byte("test snapshot data")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	ctx := context.Background()
	result, err := validator.ValidateSnapshot(ctx, tmpFile.Name())

	if err != nil {
		t.Fatalf("ValidateSnapshot failed: %v", err)
	}

	if !result.Valid {
		t.Error("Expected valid snapshot")
	}

	if result.Hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestValidator_ValidateSnapshot_NotFound(t *testing.T) {
	validator := NewValidator(logr.Discard())

	ctx := context.Background()
	result, err := validator.ValidateSnapshot(ctx, "/nonexistent/file.db")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result for nonexistent file")
	}
}

func TestValidator_CalculateHash(t *testing.T) {
	validator := NewValidator(logr.Discard())

	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test-hash-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := []byte("consistent test data")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Calculate hash twice
	hash1, err := validator.calculateHash(tmpFile.Name())
	if err != nil {
		t.Fatalf("First hash calculation failed: %v", err)
	}

	hash2, err := validator.calculateHash(tmpFile.Name())
	if err != nil {
		t.Fatalf("Second hash calculation failed: %v", err)
	}

	// Hashes should be consistent
	if hash1 != hash2 {
		t.Errorf("Hash mismatch: %s != %s", hash1, hash2)
	}

	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}
}
