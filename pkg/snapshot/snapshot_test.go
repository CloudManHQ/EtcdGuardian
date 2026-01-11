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

package snapshot

import (
	"context"
	"os"
	"testing"

	"github.com/go-logr/logr"
	etcdguardianv1alpha1 "github.com/etcdguardian/etcdguardian/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSnapshotEngine_TakeFullSnapshot(t *testing.T) {
	engine := NewSnapshotEngine(logr.Discard())

	backup := &etcdguardianv1alpha1.EtcdBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-backup",
			Namespace: "default",
		},
		Spec: etcdguardianv1alpha1.EtcdBackupSpec{
			BackupMode: etcdguardianv1alpha1.BackupModeFull,
		},
	}

	ctx := context.Background()
	snapshotPath, size, revision, err := engine.TakeFullSnapshot(ctx, backup)

	if err != nil {
		t.Fatalf("TakeFullSnapshot failed: %v", err)
	}

	if snapshotPath == "" {
		t.Error("Expected non-empty snapshot path")
	}

	if size <= 0 {
		t.Error("Expected positive snapshot size")
	}

	if revision <= 0 {
		t.Error("Expected positive revision number")
	}

	// Cleanup
	if snapshotPath != "" {
		os.Remove(snapshotPath)
	}
}

func TestSnapshotEngine_TakeIncrementalSnapshot(t *testing.T) {
	engine := NewSnapshotEngine(logr.Discard())

	backup := &etcdguardianv1alpha1.EtcdBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-incremental",
			Namespace: "default",
		},
		Spec: etcdguardianv1alpha1.EtcdBackupSpec{
			BackupMode: etcdguardianv1alpha1.BackupModeIncremental,
		},
	}

	ctx := context.Background()
	snapshotPath, size, revision, err := engine.TakeIncrementalSnapshot(ctx, backup)

	if err != nil {
		t.Fatalf("TakeIncrementalSnapshot failed: %v", err)
	}

	if snapshotPath == "" {
		t.Error("Expected non-empty snapshot path")
	}

	if size <= 0 {
		t.Error("Expected positive snapshot size")
	}

	if revision <= 0 {
		t.Error("Expected positive revision number")
	}

	// Cleanup
	if snapshotPath != "" {
		os.Remove(snapshotPath)
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

package snapshot

import (
	"context"
	"os"
	"testing"

	"github.com/go-logr/logr"
	etcdguardianv1alpha1 "github.com/etcdguardian/etcdguardian/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSnapshotEngine_TakeFullSnapshot(t *testing.T) {
	engine := NewSnapshotEngine(logr.Discard())

	backup := &etcdguardianv1alpha1.EtcdBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-backup",
			Namespace: "default",
		},
		Spec: etcdguardianv1alpha1.EtcdBackupSpec{
			BackupMode: etcdguardianv1alpha1.BackupModeFull,
		},
	}

	ctx := context.Background()
	snapshotPath, size, revision, err := engine.TakeFullSnapshot(ctx, backup)

	if err != nil {
		t.Fatalf("TakeFullSnapshot failed: %v", err)
	}

	if snapshotPath == "" {
		t.Error("Expected non-empty snapshot path")
	}

	if size <= 0 {
		t.Error("Expected positive snapshot size")
	}

	if revision <= 0 {
		t.Error("Expected positive revision number")
	}

	// Cleanup
	if snapshotPath != "" {
		os.Remove(snapshotPath)
	}
}

func TestSnapshotEngine_TakeIncrementalSnapshot(t *testing.T) {
	engine := NewSnapshotEngine(logr.Discard())

	backup := &etcdguardianv1alpha1.EtcdBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-incremental",
			Namespace: "default",
		},
		Spec: etcdguardianv1alpha1.EtcdBackupSpec{
			BackupMode: etcdguardianv1alpha1.BackupModeIncremental,
		},
	}

	ctx := context.Background()
	snapshotPath, size, revision, err := engine.TakeIncrementalSnapshot(ctx, backup)

	if err != nil {
		t.Fatalf("TakeIncrementalSnapshot failed: %v", err)
	}

	if snapshotPath == "" {
		t.Error("Expected non-empty snapshot path")
	}

	if size <= 0 {
		t.Error("Expected positive snapshot size")
	}

	if revision <= 0 {
		t.Error("Expected positive revision number")
	}

	// Cleanup
	if snapshotPath != "" {
		os.Remove(snapshotPath)
	}
}
