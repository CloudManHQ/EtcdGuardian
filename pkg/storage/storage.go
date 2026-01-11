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

package storage

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	etcdguardianv1alpha1 "github.com/etcdguardian/etcdguardian/api/v1alpha1"
)

// Storage defines the interface for storage backends
type Storage interface {
	// Upload uploads a snapshot to storage
	Upload(ctx context.Context, localPath string, backup *etcdguardianv1alpha1.EtcdBackup) (string, error)
	
	// Download downloads a snapshot from storage
	Download(ctx context.Context, remotePath, localPath string) error
	
	// List lists snapshots in storage
	List(ctx context.Context, prefix string) ([]SnapshotMetadata, error)
	
	// Delete deletes a snapshot from storage
	Delete(ctx context.Context, remotePath string) error
	
	// GetMetadata gets snapshot metadata
	GetMetadata(ctx context.Context, remotePath string) (*SnapshotMetadata, error)
}

// SnapshotMetadata contains snapshot metadata
type SnapshotMetadata struct {
	Name              string
	Path              string
	Size              int64
	CreationTimestamp int64
	EtcdVersion       string
}

// NewStorage creates a new storage backend based on the provider
func NewStorage(provider etcdguardianv1alpha1.StorageProvider, location etcdguardianv1alpha1.StorageLocation, k8sClient client.Client, namespace string) (Storage, error) {
	switch provider {
	case etcdguardianv1alpha1.StorageProviderS3:
		return NewS3Storage(location, k8sClient, namespace)
	case etcdguardianv1alpha1.StorageProviderOSS:
		return NewOSSStorage(location, k8sClient, namespace)
	case etcdguardianv1alpha1.StorageProviderGCS:
		return nil, fmt.Errorf("GCS storage not yet implemented")
	case etcdguardianv1alpha1.StorageProviderAzure:
		return nil, fmt.Errorf("Azure storage not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", provider)
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

package storage

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	etcdguardianv1alpha1 "github.com/etcdguardian/etcdguardian/api/v1alpha1"
)

// Storage defines the interface for storage backends
type Storage interface {
	// Upload uploads a snapshot to storage
	Upload(ctx context.Context, localPath string, backup *etcdguardianv1alpha1.EtcdBackup) (string, error)
	
	// Download downloads a snapshot from storage
	Download(ctx context.Context, remotePath, localPath string) error
	
	// List lists snapshots in storage
	List(ctx context.Context, prefix string) ([]SnapshotMetadata, error)
	
	// Delete deletes a snapshot from storage
	Delete(ctx context.Context, remotePath string) error
	
	// GetMetadata gets snapshot metadata
	GetMetadata(ctx context.Context, remotePath string) (*SnapshotMetadata, error)
}

// SnapshotMetadata contains snapshot metadata
type SnapshotMetadata struct {
	Name              string
	Path              string
	Size              int64
	CreationTimestamp int64
	EtcdVersion       string
}

// NewStorage creates a new storage backend based on the provider
func NewStorage(provider etcdguardianv1alpha1.StorageProvider, location etcdguardianv1alpha1.StorageLocation, k8sClient client.Client, namespace string) (Storage, error) {
	switch provider {
	case etcdguardianv1alpha1.StorageProviderS3:
		return NewS3Storage(location, k8sClient, namespace)
	case etcdguardianv1alpha1.StorageProviderOSS:
		return NewOSSStorage(location, k8sClient, namespace)
	case etcdguardianv1alpha1.StorageProviderGCS:
		return nil, fmt.Errorf("GCS storage not yet implemented")
	case etcdguardianv1alpha1.StorageProviderAzure:
		return nil, fmt.Errorf("Azure storage not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", provider)
	}
}
