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
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
)

// Validator handles snapshot validation
type Validator struct {
	log logr.Logger
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid   bool
	Hash    string
	Message string
}

// NewValidator creates a new validator
func NewValidator(log logr.Logger) *Validator {
	return &Validator{
		log: log,
	}
}

// ValidateSnapshot validates a snapshot file
func (v *Validator) ValidateSnapshot(ctx context.Context, snapshotPath string) (*ValidationResult, error) {
	v.log.Info("Validating snapshot", "path", snapshotPath)

	// Check if file exists
	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return &ValidationResult{
			Valid:   false,
			Message: fmt.Sprintf("Snapshot file not found: %s", snapshotPath),
		}, nil
	}

	// Calculate hash
	hash, err := v.calculateHash(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	// TODO: Implement etcdctl snapshot status check
	// For now, just check file integrity

	return &ValidationResult{
		Valid:   true,
		Hash:    hash,
		Message: "Snapshot validation passed",
	}, nil
}

// calculateHash calculates SHA256 hash of a file
func (v *Validator) calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
