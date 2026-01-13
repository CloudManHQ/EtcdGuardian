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

// EtcdBackupScheduleSpec defines the desired state of EtcdBackupSchedule
type EtcdBackupScheduleSpec struct {
	// Schedule defines the cron schedule
	// +kubebuilder:validation:Required
	Schedule string `json:"schedule"`

	// BackupTemplate is the template for creating backups
	// +kubebuilder:validation:Required
	BackupTemplate EtcdBackupTemplateSpec `json:"backupTemplate"`

	// Suspend suspends the schedule
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// AIOptimization enables AI-driven schedule optimization
	// +optional
	AIOptimization *AIOptimizationConfig `json:"aiOptimization,omitempty"`
}

// EtcdBackupTemplateSpec defines the backup template
type EtcdBackupTemplateSpec struct {
	// Metadata for generated backups
	// +optional
	Metadata metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the backup specification
	// +kubebuilder:validation:Required
	Spec EtcdBackupSpec `json:"spec"`
}

// AIOptimizationConfig defines AI optimization settings
type AIOptimizationConfig struct {
	// Enabled enables AI optimization
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// MinInterval is the minimum backup interval
	// +optional
	MinInterval *metav1.Duration `json:"minInterval,omitempty"`

	// MaxInterval is the maximum backup interval
	// +optional
	MaxInterval *metav1.Duration `json:"maxInterval,omitempty"`
}

// BackupRef references a backup
type BackupRef struct {
	// Name of the backup
	Name string `json:"name"`

	// UID of the backup
	UID string `json:"uid"`

	// CreationTimestamp of the backup
	CreationTimestamp metav1.Time `json:"creationTimestamp"`

	// Status of the backup
	Status BackupPhase `json:"status"`
}

// EtcdBackupScheduleStatus defines the observed state of EtcdBackupSchedule
type EtcdBackupScheduleStatus struct {
	// LastBackupTime is the time of the last backup
	// +optional
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`

	// NextBackupTime is the time of the next scheduled backup (AI-optimized)
	// +optional
	NextBackupTime *metav1.Time `json:"nextBackupTime,omitempty"`

	// BackupHistory contains recent backup references
	// +optional
	BackupHistory []BackupRef `json:"backupHistory,omitempty"`

	// Conditions represent the latest available observations of the schedule's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=etcdbkpsched
// +kubebuilder:printcolumn:name="Schedule",type=string,JSONPath=`.spec.schedule`
// +kubebuilder:printcolumn:name="Suspend",type=boolean,JSONPath=`.spec.suspend`
// +kubebuilder:printcolumn:name="LastBackup",type="date",JSONPath=".status.lastBackupTime"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// EtcdBackupSchedule is the Schema for the etcdbackupschedules API
type EtcdBackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdBackupScheduleSpec   `json:"spec,omitempty"`
	Status EtcdBackupScheduleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdBackupScheduleList contains a list of EtcdBackupSchedule
type EtcdBackupScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdBackupSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdBackupSchedule{}, &EtcdBackupScheduleList{})
}
