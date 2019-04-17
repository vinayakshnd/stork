package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//ApplicationCloneResourceName is the name for the application clone resource
	ApplicationCloneResourceName = "applicationclone"
	//ApplicationCloneResourcePlural is the name in plural for the application clone resources
	ApplicationCloneResourcePlural = "applicationclones"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//ApplicationClone represents the cloning of application in different namespaces
type ApplicationClone struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            ApplicationCloneSpec   `json:"spec"`
	Status          ApplicationCloneStatus `json:"status,omitempty"`
}

//ApplicationCloneSpec defines the spec to create an application clone
type ApplicationCloneSpec struct {
	//NamespaceMapping to store the target:destination namespaces
	NamespaceMapping map[string]string `json:"namespace_mapping"`
	//Selectors for label on objects
	Selectors    map[string]string `json:"selectors"`
	PreExecRule  string            `json:"preExecRule"`
	PostExecRule string            `json:"postExecRule"`
}

//ApplicationCloneStatus defines the status of the clone
type ApplicationCloneStatus struct {
	//Status of the cloning process
	Status ApplicationCloneStatusType `json:"status"`
	//Stage of the cloning process
	Stage ApplicationCloneStageType `json:"stage"`
}

//ApplicationCloneStatusType defines status of the application being cloned
type ApplicationCloneStatusType string

const (
	//ApplicationCloneStatusInitial initial state when the cloning will start
	ApplicationCloneStatusInitial ApplicationCloneStatusType = ""
	//ApplicationCloneStatusPending when cloning is still pending
	ApplicationCloneStatusPending ApplicationCloneStatusType = "Pending"
	//ApplicationCloneStatusInProgress cloning in progress
	ApplicationCloneStatusInProgress ApplicationCloneStatusType = "InProgress"
	//ApplicationCloneStatusFailed when cloning has failed
	ApplicationCloneStatusFailed ApplicationCloneStatusType = "Failed"
	//ApplicationCloneStatusSuccess when cloning was a success
	ApplicationCloneStatusSuccess ApplicationCloneStatusType = "Success"
	//ApplicationCloneStatusPartialSuccess when cloning was only partially successful
	ApplicationCloneStatusPartialSuccess ApplicationCloneStatusType = "PartialSuccess"
)

//ApplicationCloneStageType defines the stage of the cloning process
type ApplicationCloneStageType string

const (
	//ApplicationCloneStageInitial when the cloning was started
	ApplicationCloneStageInitial ApplicationCloneStageType = ""
	//ApplicationCloneStagePreExecRule stage when pre-exec rules are being executed
	ApplicationCloneStagePreExecRule ApplicationCloneStageType = "PreExecRule"
	//ApplicationCloneStagePostExecRule stage when post-exec rules are being executed
	ApplicationCloneStagePostExecRule ApplicationCloneStageType = "PostExecRule"
	//ApplicationCloneStageVolumeClone stage where the volumes are being cloned
	ApplicationCloneStageVolumeClone ApplicationCloneStageType = "VolumeClone"
	//ApplicationCloneStageApplicationClone stage when applications are being cloned
	ApplicationCloneStageApplicationClone ApplicationCloneStageType = "ApplicationClone"
	//ApplicationCloneStageDone stage when the cloning is done
	ApplicationCloneStageDone ApplicationCloneStageType = "Done"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationCloneList is a list of ApplicationClones
type ApplicationCloneList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []ApplicationClone `json:"items"`
}
