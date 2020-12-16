// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Blog) DeepCopyInto(out *Blog) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Blog.
func (in *Blog) DeepCopy() *Blog {
	if in == nil {
		return nil
	}
	out := new(Blog)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Blog) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlogList) DeepCopyInto(out *BlogList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Blog, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlogList.
func (in *BlogList) DeepCopy() *BlogList {
	if in == nil {
		return nil
	}
	out := new(BlogList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BlogList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlogSpec) DeepCopyInto(out *BlogSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.Config.DeepCopyInto(&out.Config)
	in.Persistent.DeepCopyInto(&out.Persistent)
	in.Ingress.DeepCopyInto(&out.Ingress)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlogSpec.
func (in *BlogSpec) DeepCopy() *BlogSpec {
	if in == nil {
		return nil
	}
	out := new(BlogSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlogStatus) DeepCopyInto(out *BlogStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlogStatus.
func (in *BlogStatus) DeepCopy() *BlogStatus {
	if in == nil {
		return nil
	}
	out := new(BlogStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostConfig) DeepCopyInto(out *GhostConfig) {
	*out = *in
	out.Mail = in.Mail
	in.Logging.DeepCopyInto(&out.Logging)
	out.Server = in.Server
	out.Database = in.Database
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostConfig.
func (in *GhostConfig) DeepCopy() *GhostConfig {
	if in == nil {
		return nil
	}
	out := new(GhostConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostDatabase) DeepCopyInto(out *GhostDatabase) {
	*out = *in
	out.Connection = in.Connection
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostDatabase.
func (in *GhostDatabase) DeepCopy() *GhostDatabase {
	if in == nil {
		return nil
	}
	out := new(GhostDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostDatabaseConnection) DeepCopyInto(out *GhostDatabaseConnection) {
	*out = *in
	out.Port = in.Port
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostDatabaseConnection.
func (in *GhostDatabaseConnection) DeepCopy() *GhostDatabaseConnection {
	if in == nil {
		return nil
	}
	out := new(GhostDatabaseConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostIngress) DeepCopyInto(out *GhostIngress) {
	*out = *in
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.TLS = in.TLS
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostIngress.
func (in *GhostIngress) DeepCopy() *GhostIngress {
	if in == nil {
		return nil
	}
	out := new(GhostIngress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostIngressTLS) DeepCopyInto(out *GhostIngressTLS) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostIngressTLS.
func (in *GhostIngressTLS) DeepCopy() *GhostIngressTLS {
	if in == nil {
		return nil
	}
	out := new(GhostIngressTLS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostPersistent) DeepCopyInto(out *GhostPersistent) {
	*out = *in
	if in.StorageClass != nil {
		in, out := &in.StorageClass, &out.StorageClass
		*out = new(string)
		**out = **in
	}
	out.Size = in.Size.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostPersistent.
func (in *GhostPersistent) DeepCopy() *GhostPersistent {
	if in == nil {
		return nil
	}
	out := new(GhostPersistent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GhostServer) DeepCopyInto(out *GhostServer) {
	*out = *in
	out.Port = in.Port
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GhostServer.
func (in *GhostServer) DeepCopy() *GhostServer {
	if in == nil {
		return nil
	}
	out := new(GhostServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Logging) DeepCopyInto(out *Logging) {
	*out = *in
	if in.Transports != nil {
		in, out := &in.Transports, &out.Transports
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Logging.
func (in *Logging) DeepCopy() *Logging {
	if in == nil {
		return nil
	}
	out := new(Logging)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mail) DeepCopyInto(out *Mail) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mail.
func (in *Mail) DeepCopy() *Mail {
	if in == nil {
		return nil
	}
	out := new(Mail)
	in.DeepCopyInto(out)
	return out
}
