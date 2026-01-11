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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackupMode defines the backup mode
// +kubebuilder:validation:Enum=Full;Incremental
type BackupMode string

const (
	BackupModeFull        BackupMode = "Full"
	BackupModeIncremental BackupMode = "Incremental"
)

// BackupPhase defines the phase of backup
// +kubebuilder:validation:Enum=Pending;Validating;Preparing;Snapshotting;Uploading;Validating_Snapshot;TriggeringVelero;Completed;Failed
type BackupPhase string

const (
	BackupPhasePending             BackupPhase = "Pending"
	BackupPhaseValidating          BackupPhase = "Validating"
	BackupPhasePreparing           BackupPhase = "Preparing"
	BackupPhaseSnapshotting        BackupPhase = "Snapshotting"
	BackupPhaseUploading           BackupPhase = "Uploading"
	BackupPhaseValidatingSnapshot  BackupPhase = "Validating_Snapshot"
	BackupPhaseTriggeringVelero    BackupPhase = "TriggeringVelero"
	BackupPhaseCompleted           BackupPhase = "Completed"
	BackupPhaseFailed              BackupPhase = "Failed"
)

// StorageProvider defines the storage provider type
// +kubebuilder:validation:Enum=S3;OSS;GCS;Azure
type StorageProvider string

const (
	StorageProviderS3    StorageProvider = "S3"
	StorageProviderOSS   StorageProvider = "OSS"
	StorageProviderGCS   StorageProvider = "GCS"
	StorageProviderAzure StorageProvider = "Azure"
)

// EtcdBackupSpec defines the desired state of EtcdBackup
type EtcdBackupSpec struct {
	// Schedule defines the cron schedule for periodic backups (optional)
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// BackupMode specifies whether this is a full or incremental backup
	// +kubebuilder:validation:Required
	BackupMode BackupMode `json:"backupMode"`

	// EtcdEndpoints is the list of etcd endpoints (optional, auto-discovery if empty)
	// +optional
	EtcdEndpoints []string `json:"etcdEndpoints,omitempty"`

	// EtcdCertificates contains etcd TLS certificate references
	// +optional
	EtcdCertificates *EtcdCertificates `json:"etcdCertificates,omitempty"`

	// StorageLocation defines where to store the backup
	// +kubebuilder:validation:Required
	StorageLocation StorageLocation `json:"storageLocation"`

	// Encryption configuration for backup encryption
	// +optional
	Encryption *EncryptionConfig `json:"encryption,omitempty"`

	// RetentionPolicy defines how long to keep backups
	// +optional
	RetentionPolicy *RetentionPolicy `json:"retentionPolicy,omitempty"`

	// Validation settings for backup validation
	// +optional
	Validation *ValidationConfig `json:"validation,omitempty"`

	// VeleroIntegration for Velero backup integration
	// +optional
	VeleroIntegration *VeleroIntegration `json:"veleroIntegration,omitempty"`

	// NamespaceSelector for multi-tenant backup filtering
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// Hooks for pre/post backup scripts
	// +optional
	Hooks *BackupHooks `json:"hooks,omitempty"`
}

// EtcdCertificates defines TLS certificates for etcd connection
type EtcdCertificates struct {
	// CA certificate secret reference
	// +optional
	CA string `json:"ca,omitempty"`

	// Client certificate secret reference
	// +optional
	Cert string `json:"cert,omitempty"`

	// Client key secret reference
	// +optional
	Key string `json:"key,omitempty"`
}

// StorageLocation defines the storage backend configuration
type StorageLocation struct {
	// Provider specifies the storage provider (S3, OSS, GCS, Azure)
	// +kubebuilder:validation:Required
	Provider StorageProvider `json:"provider"`

	// Bucket name
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`

	// Prefix for storage path
	// +optional
	Prefix string `json:"prefix,omitempty"`

	// Region of the storage
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	// Endpoint for custom storage endpoint (e.g., MinIO)
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// CredentialsSecret is the name of the secret containing credentials
	// +kubebuilder:validation:Required
	CredentialsSecret string `json:"credentialsSecret"`
}

// EncryptionConfig defines encryption settings
type EncryptionConfig struct {
	// Enabled specifies whether encryption is enabled
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// KMSKeyID for KMS-based encryption
	// +optional
	KMSKeyID string `json:"kmsKeyID,omitempty"`

	// EncryptionSecret for client-side encryption
	// +optional
	EncryptionSecret string `json:"encryptionSecret,omitempty"`
}

// RetentionPolicy defines backup retention rules
type RetentionPolicy struct {
	// MaxBackups is the maximum number of backups to retain
	// +optional
	MaxBackups *int `json:"maxBackups,omitempty"`

	// MaxAge is the maximum age of backups to retain (e.g., "720h")
	// +optional
	MaxAge *metav1.Duration `json:"maxAge,omitempty"`
}

// ValidationConfig defines validation settings
type ValidationConfig struct {
	// Enabled specifies whether automatic validation is enabled
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// ConsistencyCheck enables hash consistency checking
	// +optional
	ConsistencyCheck bool `json:"consistencyCheck,omitempty"`
}

// VeleroIntegration defines Velero integration settings
type VeleroIntegration struct {
	// Enabled specifies whether to trigger Velero backup
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// BackupName is the name of the associated Velero backup
	// +optional
	BackupName string `json:"backupName,omitempty"`
}

// BackupHooks defines pre and post backup hooks
type BackupHooks struct {
	// PreBackup hooks executed before backup
	// +optional
	PreBackup []Hook `json:"preBackup,omitempty"`

	// PostBackup hooks executed after backup
	// +optional
	PostBackup []Hook `json:"postBackup,omitempty"`
}

// Hook defines a single hook configuration
type Hook struct {
	// Name of the hook
	Name string `json:"name"`

	// Type of hook (Exec, HTTP, etc.)
	Type string `json:"type"`

	// Exec defines command execution hook
	// +optional
	Exec *ExecHook `json:"exec,omitempty"`
}

// ExecHook defines command execution hook
type ExecHook struct {
	// Command to execute
	Command []string `json:"command"`

	// Timeout for command execution
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// EtcdBackupStatus defines the observed state of EtcdBackup
type EtcdBackupStatus struct {
	// Phase represents the current phase of the backup
	// +optional
	Phase BackupPhase `json:"phase,omitempty"`

	// SnapshotSize is the size of the snapshot in bytes
	// +optional
	SnapshotSize int64 `json:"snapshotSize,omitempty"`

	// SnapshotLocation is the full path where the snapshot is stored
	// +optional
	SnapshotLocation string `json:"snapshotLocation,omitempty"`

	// EtcdRevision is the etcd revision at the time of backup
	// +optional
	EtcdRevision int64 `json:"etcdRevision,omitempty"`

	// ValidationResult contains validation results
	// +optional
	ValidationResult *ValidationResult `json:"validationResult,omitempty"`

	// StartTime is when the backup started
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is when the backup completed
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// VeleroBackupName is the name of the associated Velero backup
	// +optional
	VeleroBackupName string `json:"veleroBackupName,omitempty"`

	// VeleroBackupUID is the UID of the associated Velero backup
	// +optional
	VeleroBackupUID string `json:"veleroBackupUID,omitempty"`

	// Conditions represent the latest available observations of the backup's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	// +optional
	Message string `json:"message,omitempty"`
}

// ValidationResult contains the results of backup validation
type ValidationResult struct {
	// Valid indicates whether the backup passed validation
	Valid bool `json:"valid"`

	// Hash is the computed hash of the backup
	// +optional
	Hash string `json:"hash,omitempty"`

	// Message provides additional validation information
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=etcdbkp
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Mode",type=string,JSONPath=`.spec.backupMode`
// +kubebuilder:printcolumn:name="Size",type=string,JSONPath=`.status.snapshotSize`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// EtcdBackup is the Schema for the etcdbackups API
type EtcdBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdBackupSpec   `json:"spec,omitempty"`
	Status EtcdBackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdBackupList contains a list of EtcdBackup
type EtcdBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdBackup{}, &EtcdBackupList{})
}
