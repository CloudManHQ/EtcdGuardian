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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	etcdguardianv1alpha1 "github.com/etcdguardian/etcdguardian/api/v1alpha1"
	"github.com/etcdguardian/etcdguardian/pkg/snapshot"
	"github.com/etcdguardian/etcdguardian/pkg/storage"
	"github.com/etcdguardian/etcdguardian/pkg/validation"
)

const (
	backupFinalizer = "etcdguardian.io/finalizer"
)

// EtcdBackupReconciler reconciles a EtcdBackup object
type EtcdBackupReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=etcdguardian.io,resources=etcdbackups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=etcdguardian.io,resources=etcdbackups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=etcdguardian.io,resources=etcdbackups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *EtcdBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", req.NamespacedName)

	// Fetch the EtcdBackup instance
	backup := &etcdguardianv1alpha1.EtcdBackup{}
	err := r.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("EtcdBackup resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get EtcdBackup")
		return ctrl.Result{}, err
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(backup, backupFinalizer) {
		controllerutil.AddFinalizer(backup, backupFinalizer)
		if err := r.Update(ctx, backup); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Handle deletion
	if !backup.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, backup)
	}

	// Check if backup is already completed or failed
	if backup.Status.Phase == etcdguardianv1alpha1.BackupPhaseCompleted ||
		backup.Status.Phase == etcdguardianv1alpha1.BackupPhaseFailed {
		return ctrl.Result{}, nil
	}

	// Initialize status if it's pending
	if backup.Status.Phase == "" {
		backup.Status.Phase = etcdguardianv1alpha1.BackupPhasePending
		backup.Status.StartTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, backup); err != nil {
			log.Error(err, "Failed to update backup status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Execute backup phases
	switch backup.Status.Phase {
	case etcdguardianv1alpha1.BackupPhasePending:
		return r.validateConfig(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseValidating:
		return r.prepareBackup(ctx, backup)
	case etcdguardianv1alpha1.BackupPhasePreparing:
		return r.takeSnapshot(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseSnapshotting:
		return r.uploadSnapshot(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseUploading:
		return r.validateSnapshot(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseValidatingSnapshot:
		return r.triggerVelero(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseTriggeringVelero:
		return r.completeBackup(ctx, backup)
	}

	return ctrl.Result{}, nil
}

// validateConfig validates the backup configuration
func (r *EtcdBackupReconciler) validateConfig(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Validating backup configuration")

	// Validate storage location
	if backup.Spec.StorageLocation.Bucket == "" {
		return r.updateStatusFailed(ctx, backup, "Storage bucket is required")
	}

	// Validate credentials secret exists
	secret := &client.ObjectKey{
		Name:      backup.Spec.StorageLocation.CredentialsSecret,
		Namespace: backup.Namespace,
	}
	if err := r.Get(ctx, *secret, &client.Object{}); err != nil {
		if errors.IsNotFound(err) {
			return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Credentials secret %s not found", backup.Spec.StorageLocation.CredentialsSecret))
		}
		return ctrl.Result{}, err
	}

	// Move to next phase
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseValidating
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// prepareBackup prepares the backup environment
func (r *EtcdBackupReconciler) prepareBackup(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Preparing backup")

	// TODO: Execute pre-backup hooks if defined

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhasePreparing
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// takeSnapshot performs the etcd snapshot
func (r *EtcdBackupReconciler) takeSnapshot(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Taking etcd snapshot")

	// Create snapshot engine
	snapshotEngine := snapshot.NewSnapshotEngine(log)

	// Perform snapshot based on backup mode
	var snapshotPath string
	var snapshotSize int64
	var etcdRevision int64
	var err error

	if backup.Spec.BackupMode == etcdguardianv1alpha1.BackupModeFull {
		snapshotPath, snapshotSize, etcdRevision, err = snapshotEngine.TakeFullSnapshot(ctx, backup)
	} else {
		snapshotPath, snapshotSize, etcdRevision, err = snapshotEngine.TakeIncrementalSnapshot(ctx, backup)
	}

	if err != nil {
		return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to take snapshot: %v", err))
	}

	// Update status with snapshot info
	backup.Status.SnapshotSize = snapshotSize
	backup.Status.EtcdRevision = etcdRevision
	backup.Status.SnapshotLocation = snapshotPath
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseSnapshotting
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// uploadSnapshot uploads the snapshot to storage
func (r *EtcdBackupReconciler) uploadSnapshot(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Uploading snapshot to storage")

	// Create storage backend
	storageBackend, err := storage.NewStorage(backup.Spec.StorageLocation.Provider, backup.Spec.StorageLocation, r.Client, backup.Namespace)
	if err != nil {
		return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to create storage backend: %v", err))
	}

	// Upload snapshot
	location, err := storageBackend.Upload(ctx, backup.Status.SnapshotLocation, backup)
	if err != nil {
		return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to upload snapshot: %v", err))
	}

	backup.Status.SnapshotLocation = location
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseUploading
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// validateSnapshot validates the uploaded snapshot
func (r *EtcdBackupReconciler) validateSnapshot(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Validating snapshot")

	if backup.Spec.Validation != nil && backup.Spec.Validation.Enabled {
		validator := validation.NewValidator(log)
		result, err := validator.ValidateSnapshot(ctx, backup.Status.SnapshotLocation)
		if err != nil {
			return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to validate snapshot: %v", err))
		}

		backup.Status.ValidationResult = &etcdguardianv1alpha1.ValidationResult{
			Valid:   result.Valid,
			Hash:    result.Hash,
			Message: result.Message,
		}

		if !result.Valid {
			return r.updateStatusFailed(ctx, backup, "Snapshot validation failed")
		}
	}

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseValidatingSnapshot
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// triggerVelero triggers Velero backup if enabled
func (r *EtcdBackupReconciler) triggerVelero(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))

	if backup.Spec.VeleroIntegration != nil && backup.Spec.VeleroIntegration.Enabled {
		log.Info("Triggering Velero backup")
		// TODO: Implement Velero integration
		// For now, just mark as completed
	}

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseTriggeringVelero
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// completeBackup marks the backup as completed
func (r *EtcdBackupReconciler) completeBackup(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Backup completed successfully")

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseCompleted
	backup.Status.CompletionTime = &metav1.Time{Time: time.Now()}
	backup.Status.Message = "Backup completed successfully"

	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// updateStatusFailed updates the backup status to failed
func (r *EtcdBackupReconciler) updateStatusFailed(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup, message string) (ctrl.Result, error) {
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseFailed
	backup.Status.CompletionTime = &metav1.Time{Time: time.Now()}
	backup.Status.Message = message

	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, fmt.Errorf("%s", message)
}

// handleDeletion handles the deletion of a backup
func (r *EtcdBackupReconciler) handleDeletion(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))

	if controllerutil.ContainsFinalizer(backup, backupFinalizer) {
		// TODO: Clean up snapshot from storage if needed

		log.Info("Removing finalizer")
		controllerutil.RemoveFinalizer(backup, backupFinalizer)
		if err := r.Update(ctx, backup); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EtcdBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&etcdguardianv1alpha1.EtcdBackup{}).
		Complete(r)
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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	etcdguardianv1alpha1 "github.com/etcdguardian/etcdguardian/api/v1alpha1"
	"github.com/etcdguardian/etcdguardian/pkg/snapshot"
	"github.com/etcdguardian/etcdguardian/pkg/storage"
	"github.com/etcdguardian/etcdguardian/pkg/validation"
)

const (
	backupFinalizer = "etcdguardian.io/finalizer"
)

// EtcdBackupReconciler reconciles a EtcdBackup object
type EtcdBackupReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=etcdguardian.io,resources=etcdbackups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=etcdguardian.io,resources=etcdbackups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=etcdguardian.io,resources=etcdbackups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *EtcdBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", req.NamespacedName)

	// Fetch the EtcdBackup instance
	backup := &etcdguardianv1alpha1.EtcdBackup{}
	err := r.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("EtcdBackup resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get EtcdBackup")
		return ctrl.Result{}, err
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(backup, backupFinalizer) {
		controllerutil.AddFinalizer(backup, backupFinalizer)
		if err := r.Update(ctx, backup); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Handle deletion
	if !backup.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, backup)
	}

	// Check if backup is already completed or failed
	if backup.Status.Phase == etcdguardianv1alpha1.BackupPhaseCompleted ||
		backup.Status.Phase == etcdguardianv1alpha1.BackupPhaseFailed {
		return ctrl.Result{}, nil
	}

	// Initialize status if it's pending
	if backup.Status.Phase == "" {
		backup.Status.Phase = etcdguardianv1alpha1.BackupPhasePending
		backup.Status.StartTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, backup); err != nil {
			log.Error(err, "Failed to update backup status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Execute backup phases
	switch backup.Status.Phase {
	case etcdguardianv1alpha1.BackupPhasePending:
		return r.validateConfig(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseValidating:
		return r.prepareBackup(ctx, backup)
	case etcdguardianv1alpha1.BackupPhasePreparing:
		return r.takeSnapshot(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseSnapshotting:
		return r.uploadSnapshot(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseUploading:
		return r.validateSnapshot(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseValidatingSnapshot:
		return r.triggerVelero(ctx, backup)
	case etcdguardianv1alpha1.BackupPhaseTriggeringVelero:
		return r.completeBackup(ctx, backup)
	}

	return ctrl.Result{}, nil
}

// validateConfig validates the backup configuration
func (r *EtcdBackupReconciler) validateConfig(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Validating backup configuration")

	// Validate storage location
	if backup.Spec.StorageLocation.Bucket == "" {
		return r.updateStatusFailed(ctx, backup, "Storage bucket is required")
	}

	// Validate credentials secret exists
	secret := &client.ObjectKey{
		Name:      backup.Spec.StorageLocation.CredentialsSecret,
		Namespace: backup.Namespace,
	}
	if err := r.Get(ctx, *secret, &client.Object{}); err != nil {
		if errors.IsNotFound(err) {
			return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Credentials secret %s not found", backup.Spec.StorageLocation.CredentialsSecret))
		}
		return ctrl.Result{}, err
	}

	// Move to next phase
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseValidating
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// prepareBackup prepares the backup environment
func (r *EtcdBackupReconciler) prepareBackup(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Preparing backup")

	// TODO: Execute pre-backup hooks if defined

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhasePreparing
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// takeSnapshot performs the etcd snapshot
func (r *EtcdBackupReconciler) takeSnapshot(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Taking etcd snapshot")

	// Create snapshot engine
	snapshotEngine := snapshot.NewSnapshotEngine(log)

	// Perform snapshot based on backup mode
	var snapshotPath string
	var snapshotSize int64
	var etcdRevision int64
	var err error

	if backup.Spec.BackupMode == etcdguardianv1alpha1.BackupModeFull {
		snapshotPath, snapshotSize, etcdRevision, err = snapshotEngine.TakeFullSnapshot(ctx, backup)
	} else {
		snapshotPath, snapshotSize, etcdRevision, err = snapshotEngine.TakeIncrementalSnapshot(ctx, backup)
	}

	if err != nil {
		return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to take snapshot: %v", err))
	}

	// Update status with snapshot info
	backup.Status.SnapshotSize = snapshotSize
	backup.Status.EtcdRevision = etcdRevision
	backup.Status.SnapshotLocation = snapshotPath
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseSnapshotting
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// uploadSnapshot uploads the snapshot to storage
func (r *EtcdBackupReconciler) uploadSnapshot(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Uploading snapshot to storage")

	// Create storage backend
	storageBackend, err := storage.NewStorage(backup.Spec.StorageLocation.Provider, backup.Spec.StorageLocation, r.Client, backup.Namespace)
	if err != nil {
		return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to create storage backend: %v", err))
	}

	// Upload snapshot
	location, err := storageBackend.Upload(ctx, backup.Status.SnapshotLocation, backup)
	if err != nil {
		return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to upload snapshot: %v", err))
	}

	backup.Status.SnapshotLocation = location
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseUploading
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// validateSnapshot validates the uploaded snapshot
func (r *EtcdBackupReconciler) validateSnapshot(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Validating snapshot")

	if backup.Spec.Validation != nil && backup.Spec.Validation.Enabled {
		validator := validation.NewValidator(log)
		result, err := validator.ValidateSnapshot(ctx, backup.Status.SnapshotLocation)
		if err != nil {
			return r.updateStatusFailed(ctx, backup, fmt.Sprintf("Failed to validate snapshot: %v", err))
		}

		backup.Status.ValidationResult = &etcdguardianv1alpha1.ValidationResult{
			Valid:   result.Valid,
			Hash:    result.Hash,
			Message: result.Message,
		}

		if !result.Valid {
			return r.updateStatusFailed(ctx, backup, "Snapshot validation failed")
		}
	}

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseValidatingSnapshot
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// triggerVelero triggers Velero backup if enabled
func (r *EtcdBackupReconciler) triggerVelero(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))

	if backup.Spec.VeleroIntegration != nil && backup.Spec.VeleroIntegration.Enabled {
		log.Info("Triggering Velero backup")
		// TODO: Implement Velero integration
		// For now, just mark as completed
	}

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseTriggeringVelero
	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// completeBackup marks the backup as completed
func (r *EtcdBackupReconciler) completeBackup(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))
	log.Info("Backup completed successfully")

	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseCompleted
	backup.Status.CompletionTime = &metav1.Time{Time: time.Now()}
	backup.Status.Message = "Backup completed successfully"

	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// updateStatusFailed updates the backup status to failed
func (r *EtcdBackupReconciler) updateStatusFailed(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup, message string) (ctrl.Result, error) {
	backup.Status.Phase = etcdguardianv1alpha1.BackupPhaseFailed
	backup.Status.CompletionTime = &metav1.Time{Time: time.Now()}
	backup.Status.Message = message

	if err := r.Status().Update(ctx, backup); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, fmt.Errorf("%s", message)
}

// handleDeletion handles the deletion of a backup
func (r *EtcdBackupReconciler) handleDeletion(ctx context.Context, backup *etcdguardianv1alpha1.EtcdBackup) (ctrl.Result, error) {
	log := r.Log.WithValues("etcdbackup", client.ObjectKeyFromObject(backup))

	if controllerutil.ContainsFinalizer(backup, backupFinalizer) {
		// TODO: Clean up snapshot from storage if needed

		log.Info("Removing finalizer")
		controllerutil.RemoveFinalizer(backup, backupFinalizer)
		if err := r.Update(ctx, backup); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EtcdBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&etcdguardianv1alpha1.EtcdBackup{}).
		Complete(r)
}
