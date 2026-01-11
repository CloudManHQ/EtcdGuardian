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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	// BackupDuration tracks backup duration in seconds
	BackupDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "etcdguardian_backup_duration_seconds",
			Help:    "Duration of etcd backup operations in seconds",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1s to ~17min
		},
		[]string{"backup_mode", "status"},
	)

	// BackupSize tracks backup size in bytes
	BackupSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcdguardian_backup_size_bytes",
			Help: "Size of etcd backup in bytes",
		},
		[]string{"backup_name", "backup_mode"},
	)

	// BackupTotal tracks total number of backups
	BackupTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcdguardian_backup_total",
			Help: "Total number of etcd backups",
		},
		[]string{"status"},
	)

	// EtcdDBSize tracks etcd database size
	EtcdDBSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcdguardian_etcd_db_size_bytes",
			Help: "Size of etcd database in bytes",
		},
		[]string{"endpoint"},
	)

	// EtcdRevision tracks current etcd revision
	EtcdRevision = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcdguardian_etcd_revision",
			Help: "Current etcd revision number",
		},
		[]string{"endpoint"},
	)

	// ValidationFailures tracks validation failures
	ValidationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcdguardian_validation_failures_total",
			Help: "Total number of validation failures",
		},
		[]string{"reason"},
	)

	// RestoreDuration tracks restore duration in seconds
	RestoreDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "etcdguardian_restore_duration_seconds",
			Help:    "Duration of etcd restore operations in seconds",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10),
		},
		[]string{"restore_mode"},
	)

	// RestoreTotal tracks total number of restores
	RestoreTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcdguardian_restore_total",
			Help: "Total number of etcd restores",
		},
		[]string{"status"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		BackupDuration,
		BackupSize,
		BackupTotal,
		EtcdDBSize,
		EtcdRevision,
		ValidationFailures,
		RestoreDuration,
		RestoreTotal,
	)
}

// RecordBackupDuration records the duration of a backup operation
func RecordBackupDuration(mode, status string, duration float64) {
	BackupDuration.WithLabelValues(mode, status).Observe(duration)
}

// RecordBackupSize records the size of a backup
func RecordBackupSize(name, mode string, size int64) {
	BackupSize.WithLabelValues(name, mode).Set(float64(size))
}

// IncBackupTotal increments the total backup counter
func IncBackupTotal(status string) {
	BackupTotal.WithLabelValues(status).Inc()
}

// SetEtcdDBSize sets the etcd database size
func SetEtcdDBSize(endpoint string, size int64) {
	EtcdDBSize.WithLabelValues(endpoint).Set(float64(size))
}

// SetEtcdRevision sets the current etcd revision
func SetEtcdRevision(endpoint string, revision int64) {
	EtcdRevision.WithLabelValues(endpoint).Set(float64(revision))
}

// IncValidationFailures increments the validation failure counter
func IncValidationFailures(reason string) {
	ValidationFailures.WithLabelValues(reason).Inc()
}

// RecordRestoreDuration records the duration of a restore operation
func RecordRestoreDuration(mode string, duration float64) {
	RestoreDuration.WithLabelValues(mode).Observe(duration)
}

// IncRestoreTotal increments the total restore counter
func IncRestoreTotal(status string) {
	RestoreTotal.WithLabelValues(status).Inc()
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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	// BackupDuration tracks backup duration in seconds
	BackupDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "etcdguardian_backup_duration_seconds",
			Help:    "Duration of etcd backup operations in seconds",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1s to ~17min
		},
		[]string{"backup_mode", "status"},
	)

	// BackupSize tracks backup size in bytes
	BackupSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcdguardian_backup_size_bytes",
			Help: "Size of etcd backup in bytes",
		},
		[]string{"backup_name", "backup_mode"},
	)

	// BackupTotal tracks total number of backups
	BackupTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcdguardian_backup_total",
			Help: "Total number of etcd backups",
		},
		[]string{"status"},
	)

	// EtcdDBSize tracks etcd database size
	EtcdDBSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcdguardian_etcd_db_size_bytes",
			Help: "Size of etcd database in bytes",
		},
		[]string{"endpoint"},
	)

	// EtcdRevision tracks current etcd revision
	EtcdRevision = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcdguardian_etcd_revision",
			Help: "Current etcd revision number",
		},
		[]string{"endpoint"},
	)

	// ValidationFailures tracks validation failures
	ValidationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcdguardian_validation_failures_total",
			Help: "Total number of validation failures",
		},
		[]string{"reason"},
	)

	// RestoreDuration tracks restore duration in seconds
	RestoreDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "etcdguardian_restore_duration_seconds",
			Help:    "Duration of etcd restore operations in seconds",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10),
		},
		[]string{"restore_mode"},
	)

	// RestoreTotal tracks total number of restores
	RestoreTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcdguardian_restore_total",
			Help: "Total number of etcd restores",
		},
		[]string{"status"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		BackupDuration,
		BackupSize,
		BackupTotal,
		EtcdDBSize,
		EtcdRevision,
		ValidationFailures,
		RestoreDuration,
		RestoreTotal,
	)
}

// RecordBackupDuration records the duration of a backup operation
func RecordBackupDuration(mode, status string, duration float64) {
	BackupDuration.WithLabelValues(mode, status).Observe(duration)
}

// RecordBackupSize records the size of a backup
func RecordBackupSize(name, mode string, size int64) {
	BackupSize.WithLabelValues(name, mode).Set(float64(size))
}

// IncBackupTotal increments the total backup counter
func IncBackupTotal(status string) {
	BackupTotal.WithLabelValues(status).Inc()
}

// SetEtcdDBSize sets the etcd database size
func SetEtcdDBSize(endpoint string, size int64) {
	EtcdDBSize.WithLabelValues(endpoint).Set(float64(size))
}

// SetEtcdRevision sets the current etcd revision
func SetEtcdRevision(endpoint string, revision int64) {
	EtcdRevision.WithLabelValues(endpoint).Set(float64(revision))
}

// IncValidationFailures increments the validation failure counter
func IncValidationFailures(reason string) {
	ValidationFailures.WithLabelValues(reason).Inc()
}

// RecordRestoreDuration records the duration of a restore operation
func RecordRestoreDuration(mode string, duration float64) {
	RestoreDuration.WithLabelValues(mode).Observe(duration)
}

// IncRestoreTotal increments the total restore counter
func IncRestoreTotal(status string) {
	RestoreTotal.WithLabelValues(status).Inc()
}
