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
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Refer to // https://ghost.org/docs/concepts/config

// GhostDatabaseConnection defines ghost database connection.
type GhostDatabaseConnection struct {
	// sqlite filename.
	// +optional
	Filename string `json:"filename,omitempty"`
	// mysql host
	// +optional
	Host string `json:"host,omitempty"`
	// mysql port
	// +optional
	Port intstr.IntOrString `json:"port,omitempty"`
	// mysql database user
	// +optional
	User string `json:"user,omitempty"`
	// mysql database password of user
	// +optional
	Password string `json:"password,omitempty"`
	// mysql database name
	// +optional
	Database string `json:"database,omitempty"`
}

type GhostServer struct {
	Host string             `json:"host"`
	Port intstr.IntOrString `json:"port"`
}

// GhostDatabase defines ghost database config.
type GhostDatabase struct {
	// Client is ghost database client.
	// +kubebuilder:validation:Enum=sqlite3;mysql
	Client string `json:"client"`
	// +optional
	Connection GhostDatabaseConnection `json:"connection"`
}

// Mail defines ghost mail config.
type Mail struct {
	Transport string `json:"transport"`
}

// Logging defines ghost logging type config.
type Logging struct {
	Transports []string `json:"transports"`
}

// GhostConfig defines related ghost configuration based on https://ghost.org/docs/concepts/config
type GhostConfig struct {
	Mail     Mail          `json:"mail"`
	URL      string        `json:"url"`
	Logging  Logging       `json:"logging"`
	Server   GhostServer   `json:"server"`
	Database GhostDatabase `json:"database"`
}

// GhostPersistent defines peristent volume
type GhostPersistent struct {
	Enabled bool `json:"enabled"`
	// If not defined, default will be used
	StorageClass *string `json:"storageClass,omitempty"`
	// size of storage
	Size resource.Quantity `json:"size"`
}

// GhostIngressTLS defines ingress tls
type GhostIngressTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName"`
}

// GhostIngress defines ingress
type GhostIngress struct {
	Enabled bool `json:"enabled"`
	// +optional
	// +listType=set
	Hosts []string `json:"hosts,omitempty"`
	// +optional
	TLS GhostIngressTLS `json:"tls,omitempty"`
	// Additional annotations passed to ".metadata.annotations" in networking.k8s.io/ingress object.
	// This is useful for configuring ingress through annotation field like: ingress-class, static-ip, etc
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}
