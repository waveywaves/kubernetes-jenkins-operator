package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SeedJobSpec defines the desired state of SeedJob
type SeedJobSpec struct {
	// ID is the unique seed job name
	ID string `json:"id,omitempty"`

	// CredentialID is the Kubernetes secret name which stores repository access credentials
	CredentialID string `json:"credentialID,omitempty"`

	// Description is the description of the seed job
	// +optional
	Description string `json:"description,omitempty"`

	// Targets is the repository path where are seed job definitions
	Targets string `json:"targets,omitempty"`

	// RepositoryBranch is the repository branch where are seed job definitions
	RepositoryBranch string `json:"repositoryBranch,omitempty"`

	// RepositoryURL is the repository access URL. Can be SSH or HTTPS.
	RepositoryURL string `json:"repositoryUrl,omitempty"`

	// JenkinsRef refers to the JenkinsInstance where this SeedJob is created
	JenkinsRef string `json:"jenkinsRef,omitempty"`

	// JenkinsCredentialType is the https://jenkinsci.github.io/kubernetes-credentials-provider-plugin/ credential type
	// +optional
	JenkinsCredentialType JenkinsCredentialType `json:"credentialType,omitempty"`

	// BitbucketPushTrigger is used for Bitbucket web hooks
	// +optional
	BitbucketPushTrigger bool `json:"bitbucketPushTrigger"`

	// GitHubPushTrigger is used for GitHub web hooks
	// +optional
	GitHubPushTrigger bool `json:"githubPushTrigger"`

	// BuildPeriodically is setting for scheduled trigger
	// +optional
	BuildPeriodically string `json:"buildPeriodically"`

	// PollSCM is setting for polling changes in SCM
	// +optional
	PollSCM string `json:"pollSCM"`

	// IgnoreMissingFiles is setting for Job DSL API plugin to ignore files that miss
	// +optional
	IgnoreMissingFiles bool `json:"ignoreMissingFiles"`

	// AdditionalClasspath is setting for Job DSL API plugin to set Additional Classpath
	// +optional
	AdditionalClasspath string `json:"additionalClasspath"`

	// FailOnMissingPlugin is setting for Job DSL API plugin that fails job if required plugin is missing
	// +optional
	FailOnMissingPlugin bool `json:"failOnMissingPlugin"`

	// UnstableOnDeprecation is setting for Job DSL API plugin that sets build status as unstable if build using deprecated features
	// +optional
	UnstableOnDeprecation bool `json:"unstableOnDeprecation"`
}

// SeedJobStatus defines the observed state of SeedJob
type SeedJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SeedJob is the Schema for the seedjobs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=seedjobs,scope=Namespaced
type SeedJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SeedJobSpec   `json:"spec,omitempty"`
	Status SeedJobStatus `json:"status,omitempty"`
}

// JenkinsCredentialType defines type of Jenkins credential used to seed job mechanism
type JenkinsCredentialType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SeedJobList contains a list of SeedJob
type SeedJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SeedJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SeedJob{}, &SeedJobList{})
}

