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

// JCasCProfileSpec defines the desired state of JCasCProfile
type JCasCProfileSpec struct {
	Config map[string]string `json:"config"`
}

// JCasCProfileStatus defines the observed state of JCasCProfile
type JCasCProfileStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JCasCProfile is the Schema for the jcascprofiles API
type JCasCProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JCasCProfileSpec   `json:"spec,omitempty"`
	Status JCasCProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JCasCProfileList contains a list of JCasCProfile
type JCasCProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JCasCProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JCasCProfile{}, &JCasCProfileList{})
}
