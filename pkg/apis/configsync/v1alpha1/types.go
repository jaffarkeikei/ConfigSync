package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=configsyncs,scope=Namespaced,shortName=cs

// ConfigSync is the Schema for the configsyncs API
type ConfigSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSyncSpec   `json:"spec,omitempty"`
	Status ConfigSyncStatus `json:"status,omitempty"`
}

// ConfigSyncSpec defines the desired state of ConfigSync
type ConfigSyncSpec struct {
	// GitRepository is the URL of the Git repository containing configurations
	GitRepository string `json:"gitRepository"`

	// Branch is the Git branch to use
	Branch string `json:"branch,omitempty"`

	// Path is the path within the Git repository where configurations are stored
	Path string `json:"path"`

	// Environment specifies which environment this ConfigSync instance manages
	// +kubebuilder:validation:Enum=development;staging;production
	Environment string `json:"environment"`

	// SyncInterval specifies how often to check for changes in Git
	// +kubebuilder:default="5m"
	SyncInterval string `json:"syncInterval,omitempty"`

	// AutoApprove specifies whether to automatically apply changes
	// +kubebuilder:default=false
	AutoApprove bool `json:"autoApprove,omitempty"`

	// DriftDetection enables detection and remediation of configuration drift
	// +kubebuilder:default=true
	DriftDetection bool `json:"driftDetection,omitempty"`
}

// ConfigSyncStatus defines the observed state of ConfigSync
type ConfigSyncStatus struct {
	// LastSyncTime is the time of the last successful sync
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`

	// LastCommitID is the Git commit ID that was last synced
	LastCommitID string `json:"lastCommitID,omitempty"`

	// Conditions represent the latest available observations of ConfigSync's state
	Conditions []ConfigSyncCondition `json:"conditions,omitempty"`
}

// ConfigSyncConditionType is a valid value for ConfigSyncCondition.Type
type ConfigSyncConditionType string

const (
	// ConfigSyncReady means the ConfigSync has been successfully configured and is ready to sync
	ConfigSyncReady ConfigSyncConditionType = "Ready"

	// ConfigSyncSynced means the ConfigSync has successfully synced with the Git repository
	ConfigSyncSynced ConfigSyncConditionType = "Synced"

	// ConfigSyncError means the ConfigSync encountered an error
	ConfigSyncError ConfigSyncConditionType = "Error"

	// ConfigSyncDrifted means configuration drift was detected
	ConfigSyncDrifted ConfigSyncConditionType = "Drifted"
)

// ConfigSyncCondition contains details about the state of the ConfigSync at a certain point in time
type ConfigSyncCondition struct {
	// Type of ConfigSync condition
	Type ConfigSyncConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown
	Status metav1.ConditionStatus `json:"status"`

	// LastTransitionTime is the last time the condition transitioned from one status to another
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason for the condition's last transition
	Reason string `json:"reason,omitempty"`

	// Message contains human-readable message about condition
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true

// ConfigSyncList contains a list of ConfigSync resources
type ConfigSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigSync `json:"items"`
}
