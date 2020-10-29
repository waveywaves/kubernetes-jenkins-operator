/*


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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	UnifiedPool BackupPoolType = "unified"
	OpenPool    BackupPoolType = "open"
)

// BackupPoolSpec defines the desired state of BackupPool
type BackupPoolSpec struct {
	// Type describes the type of backupPool does one want
	// "unified" creates a single PVC and all backups are stored in
	// the corresponding PV denoted by the Backup name
	// "open" creates a new PVC for each single Backup instance
	Type BackupPoolType `json:"type"`
}

type BackupPoolType string

type BackupPoolStatus struct {
	JenkinsRef            string `json:"jenkinsRef"`
	PersistentVolumeClaim string `json:"persistentVolumeClaimRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BackupPool is the Schema for the backuppools API
type BackupPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupPoolSpec   `json:"spec,omitempty"`
	Status BackupPoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackupPoolList contains a list of BackupPool
type BackupPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupPool{}, &BackupPoolList{})
}
