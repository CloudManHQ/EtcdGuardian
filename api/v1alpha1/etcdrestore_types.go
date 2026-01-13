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

// RestoreMode defines the restore mode
// +kubebuilder:validation:Enum=Full;Incremental;PointInTime
type RestoreMode string

const (
	RestoreModeFull        RestoreMode = "Full"
	RestoreModeIncremental RestoreMode = "Incremental"
	RestoreModePointInTime RestoreMode = "PointInTime"
)

// RestorePhase defines the phase of restore
// +kubebuilder:validation:Enum=Pending;Quiescing;Restoring;Validating;Completed;Failed
type RestorePhase string

const (
	RestorePhasePending    RestorePhase = "Pending"
	RestorePhaseQuiescing  RestorePhase = "Quiescing"
	RestorePhaseRestoring  RestorePhase = "Restoring"
	RestorePhaseValidating RestorePhase = "Validating"
	RestorePhaseCompleted  RestorePhase = "Completed"
	RestorePhaseFailed     RestorePhase = "Failed"
)

// EtcdRestoreSpec defines the desired state of EtcdRestore
type EtcdRestoreSpec struct {
	// BackupName is the name of the EtcdBackup resource to restore from
	// +kubebuilder:validation:Required
	BackupName string `json:"backupName"`

	// SnapshotLocation directly specifies the snapshot path (overrides BackupName)
	// +optional
	SnapshotLocation string `json:"snapshotLocation,omitempty"`

	// RestoreMode specifies the restore mode
	// +kubebuilder:validation:Required
	RestoreMode RestoreMode `json:"restoreMode"`

	// TargetRevision specifies the target etcd revision for point-in-time restore
	// +optional
	TargetRevision *int64 `json:"targetRevision,omitempty"`

	// EtcdCluster defines the target etcd cluster configuration
	// +kubebuilder:validation:Required
	EtcdCluster EtcdClusterConfig `json:"etcdCluster"`

	// QuiesceCluster indicates whether to quiesce the cluster before restore
	// +optional
	// +kubebuilder:default=true
	QuiesceCluster bool `json:"quiesceCluster,omitempty"`

	// PreRestoreHooks are hooks executed before restore
	// +optional
	PreRestoreHooks []Hook `json:"preRestoreHooks,omitempty"`

	// PostRestoreHooks are hooks executed after restore
	// +optional
	PostRestoreHooks []Hook `json:"postRestoreHooks,omitempty"`

	// VersionCompatibility defines version compatibility settings
	// +optional
	VersionCompatibility *VersionCompatibility `json:"versionCompatibility,omitempty"`

	// NamespaceFilter filters namespaces to restore (multi-tenant)
	// +optional
	NamespaceFilter []string `json:"namespaceFilter,omitempty"`
}

// EtcdClusterConfig defines etcd cluster configuration
type EtcdClusterConfig struct {
	// Endpoints is the list of etcd endpoints
	// +kubebuilder:validation:Required
	Endpoints []string `json:"endpoints"`

	// DataDir is the etcd data directory path
	// +kubebuilder:validation:Required
	DataDir string `json:"dataDir"`

	// Certificates for etcd connection
	// +optional
	Certificates *EtcdCertificates `json:"certificates,omitempty"`
}

// VersionCompatibility defines version compatibility settings
type VersionCompatibility struct {
	// Check indicates whether to check version compatibility
	// +optional
	// +kubebuilder:default=true
	Check bool `json:"check,omitempty"`

	// AllowMismatch allows version mismatches
	// +optional
	AllowMismatch bool `json:"allowMismatch,omitempty"`
}

// EtcdRestoreStatus defines the observed state of EtcdRestore
type EtcdRestoreStatus struct {
	// Phase represents the current phase of the restore
	// +optional
	Phase RestorePhase `json:"phase,omitempty"`

	// RestoredRevision is the etcd revision that was actually restored
	// +optional
	RestoredRevision int64 `json:"restoredRevision,omitempty"`

	// RestoredKeys is the number of keys restored
	// +optional
	RestoredKeys int `json:"restoredKeys,omitempty"`

	// Errors contains error messages
	// +optional
	Errors []string `json:"errors,omitempty"`

	// StartTime is when the restore started
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is when the restore completed
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Conditions represent the latest available observations of the restore's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=etcdrestore
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Mode",type=string,JSONPath=`.spec.restoreMode`
// +kubebuilder:printcolumn:name="Keys",type=integer,JSONPath=`.status.restoredKeys`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// EtcdRestore is the Schema for the etcdrestores API
type EtcdRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdRestoreSpec   `json:"spec,omitempty"`
	Status EtcdRestoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdRestoreList contains a list of EtcdRestore
type EtcdRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdRestore{}, &EtcdRestoreList{})
}
