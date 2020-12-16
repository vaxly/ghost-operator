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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BlogSpec defines the desired state of Blog
// +k8s:openapi-gen=true
type BlogSpec struct {
	// Ghost deployment repicas
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Ghost container image, by default using latest ghost image from docker hub registry.
	// NOTE: This operator only support ghost image from docker official image. https://hub.docker.com/_/ghost/
	// +optional
	Image string `json:"image,omitempty"`
	// Ghost configuration. This field will be written as ghost configuration. Saved in configmap and mounted
	// in /etc/ghost/config/config.json and symlinked to /var/lib/ghost/config.production.json
	Config GhostConfig `json:"config"`
	// +optional
	Persistent GhostPersistent `json:"persistent,omitempty"`
	// +optional
	Ingress GhostIngress `json:"ingress,omitempty"`

	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
}

// BlogStatus defines the observed state of Blog
// +k8s:openapi-gen=true
type BlogStatus struct {
	Replicas int32  `json:"replicas,omitempty"`
	DBHost   string `json:"dbHost,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Blog is the Schema for the blogs API
type Blog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlogSpec   `json:"spec,omitempty"`
	Status BlogStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BlogList contains a list of Blog
type BlogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Blog `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Blog{}, &BlogList{})
}
