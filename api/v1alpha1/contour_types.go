// Copyright Project Contour Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// OwningContourNameLabel is the owner reference label used for a Contour
	// created by the operator. The value should be the name of the contour.
	OwningContourNameLabel = "contour.operator.projectcontour.io/owning-contour-name"

	// OwningContourNsLabel is the owner reference label used for a Contour
	// created by the operator. The value should be the namespace of the contour.
	OwningContourNsLabel = "contour.operator.projectcontour.io/owning-contour-namespace"

	// ContourFinalizer is the name of the finalizer used for a Contour.
	ContourFinalizer = "contour.operator.projectcontour.io/finalizer"
)

// +kubebuilder:object:root=true

// Contour is the Schema for the contours API.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].status`
// +kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].reason`
type Contour struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of Contour.
	Spec ContourSpec `json:"spec,omitempty"`
	// Status defines the observed state of Contour.
	Status ContourStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContourList contains a list of Contour.
type ContourList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Contour `json:"items"`
}

// ContourSpec defines the desired state of Contour.
type ContourSpec struct {
	// Replicas is the desired number of Contour replicas. If unset,
	// defaults to 2.
	//
	// +kubebuilder:default=2
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas,omitempty"`

	// Namespace defines the schema of a Contour namespace. See each field for
	// additional details. Namespace name should be the same namespace as the
	// Gateway when GatewayClassRef is set.
	//
	// TODO [danehans]: Ignore Namespace when GatewayClassRef is set.
	// xref: https://github.com/projectcontour/contour-operator/issues/212
	//
	// +kubebuilder:default={name: "projectcontour", removeOnDeletion: false}
	Namespace NamespaceSpec `json:"namespace,omitempty"`

	// NetworkPublishing defines the schema for publishing Contour to a network.
	//
	// See each field for additional details.
	//
	// +kubebuilder:default={envoy: {type: LoadBalancerService, containerPorts: {{name: http, portNumber: 8080}, {name: https, portNumber: 8443}}}}
	NetworkPublishing NetworkPublishing `json:"networkPublishing,omitempty"`

	// GatewayClassRef is a reference to a GatewayClass name used for
	// managing a Contour.
	//
	// +kubebuilder:validation:MaxLength=253
	// +optional
	GatewayClassRef *string `json:"gatewayClassRef,omitempty"`

	// IngressClassName is the name of the IngressClass used by Contour. If unset,
	// Contour will process all ingress objects without an ingress class annotation
	// or ingress objects with an annotation matching ingress-class=contour. When
	// specified, Contour will only process ingress objects that match the provided
	// class.
	//
	// For additional IngressClass details, refer to:
	//   https://projectcontour.io/docs/main/config/annotations/#ingress-class
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +optional
	IngressClassName *string `json:"ingressClassName,omitempty"`
}

// NamespaceSpec defines the schema of a Contour namespace.
type NamespaceSpec struct {
	// Name is the name of the namespace to run Contour and dependent
	// resources. If unset, defaults to "projectcontour".
	//
	// +kubebuilder:default=projectcontour
	Name string `json:"name,omitempty"`

	// RemoveOnDeletion will remove the namespace when the Contour is
	// deleted. If set to True, deletion will not occur if any of the
	// following conditions exist:
	//
	// 1. The Contour namespace is "default", "kube-system" or the
	//    contour-operator's namespace.
	//
	// 2. Another Contour exists in the namespace.
	//
	// 3. The namespace does not contain the Contour owning label.
	//
	// +kubebuilder:default=false
	RemoveOnDeletion bool `json:"removeOnDeletion,omitempty"`
}

// NetworkPublishing defines the schema for publishing Contour to a network.
type NetworkPublishing struct {
	// Envoy provides the schema for publishing the network endpoints of Envoy.
	//
	// If unset, defaults to:
	//   type: LoadBalancerService
	//   containerPorts:
	//   - name: http
	//     portNumber: 8080
	//   - name: https
	//     portNumber: 8443
	//
	// +kubebuilder:default={type: LoadBalancerService, loadBalancer: {scope: External, providerParameters: {type: AWS}}, containerPorts: {{name: http, portNumber: 8080}, {name: https, portNumber: 8443}}}
	Envoy EnvoyNetworkPublishing `json:"envoy,omitempty"`
}

// EnvoyNetworkPublishing defines the schema to publish Envoy to a network.
// +union
type EnvoyNetworkPublishing struct {
	// Type is the type of publishing strategy to use. Valid values are:
	//
	// * LoadBalancerService
	//
	// In this configuration, network endpoints for Envoy use container networking.
	// A Kubernetes LoadBalancer Service is created to publish Envoy network
	// endpoints. The Service uses port 80 to publish Envoy's HTTP network endpoint
	// and port 443 to publish Envoy's HTTPS network endpoint.
	//
	// See: https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer
	//
	// * NodePortService
	//
	// Publishes Envoy network endpoints using a Kubernetes NodePort Service.
	//
	// In this configuration, Envoy network endpoints use container networking. A Kubernetes
	// NodePort Service is created to publish the network endpoints.
	//
	// See: https://kubernetes.io/docs/concepts/services-networking/service/#nodeport
	//
	// * ClusterIPService
	//
	// Publishes Envoy network endpoints using a Kubernetes ClusterIP Service.
	//
	// In this configuration, Envoy network endpoints use container networking. A Kubernetes
	// ClusterIP Service is created to publish the network endpoints.
	//
	// See: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
	//
	// +unionDiscriminator
	// +kubebuilder:default=LoadBalancerService
	Type NetworkPublishingType `json:"type,omitempty"`

	// LoadBalancer holds parameters for the load balancer. Present only if type is
	// LoadBalancerService.
	//
	// If unspecified, defaults to an external Classic AWS ELB.
	//
	// +kubebuilder:default={scope: External, providerParameters: {type: AWS}}
	LoadBalancer LoadBalancerStrategy `json:"loadBalancer,omitempty"`

	// NodePorts is a list of network ports to expose on each node's IP at a static
	// port number using a NodePort Service. Present only if type is NodePortService.
	// A ClusterIP Service, which the NodePort Service routes to, is automatically
	// created. You'll be able to contact the NodePort Service, from outside the
	// cluster, by requesting <NodeIP>:<NodePort>.
	//
	// If type is NodePortService and nodePorts is unspecified, two nodeports will be
	// created, one named "http" and the other named "https", with port numbers auto
	// assigned by Kubernetes API server. For additional information on the NodePort
	// Service, see:
	//
	//  https://kubernetes.io/docs/concepts/services-networking/service/#nodeport
	//
	// Names and port numbers must be unique in the list. Two ports must be specified,
	// one named "http" for Envoy's insecure service and one named "https" for Envoy's
	// secure service.
	//
	// +kubebuilder:validation:MinItems=2
	// +kubebuilder:validation:MaxItems=2
	// +optional
	NodePorts []NodePort `json:"nodePorts,omitempty"`

	// ContainerPorts is a list of container ports to expose from the Envoy container(s).
	// Exposing a port here gives the system additional information about the network
	// connections the Envoy container uses, but is primarily informational. Not specifying
	// a port here DOES NOT prevent that port from being exposed by Envoy. Any port which is
	// listening on the default "0.0.0.0" address inside the Envoy container will be accessible
	// from the network. Names and port numbers must be unique in the list container ports. Two
	// ports must be specified, one named "http" for Envoy's insecure service and one named
	// "https" for Envoy's secure service.
	//
	// TODO [danehans]: Update minItems to 1, requiring only https when the following issue
	// is fixed: https://github.com/projectcontour/contour/issues/2577.
	//
	// TODO [danehans]: Increase maxItems when https://github.com/projectcontour/contour/pull/3263
	// is implemented.
	//
	// +kubebuilder:validation:MinItems=2
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:default={{name: http, portNumber: 8080}, {name: https, portNumber: 8443}}
	ContainerPorts []ContainerPort `json:"containerPorts,omitempty"`
}

// NetworkPublishingType is a way to publish network endpoints.
// +kubebuilder:validation:Enum=LoadBalancerService;NodePortService;ClusterIPService
type NetworkPublishingType string

const (
	// LoadBalancerServicePublishingType publishes a network endpoint using a Kubernetes
	// LoadBalancer Service.
	LoadBalancerServicePublishingType NetworkPublishingType = "LoadBalancerService"

	// NodePortServicePublishingType publishes a network endpoint using a Kubernetes
	// NodePort Service.
	NodePortServicePublishingType NetworkPublishingType = "NodePortService"

	// ClusterIPServicePublishingType publishes a network endpoint using a Kubernetes
	// ClusterIP Service.
	ClusterIPServicePublishingType NetworkPublishingType = "ClusterIPService"
)

// LoadBalancerStrategy holds parameters for a load balancer.
type LoadBalancerStrategy struct {
	// Scope indicates the scope at which the load balancer is exposed.
	// Possible values are "External" and "Internal".
	//
	// +kubebuilder:default=External
	Scope LoadBalancerScope `json:"scope,omitempty"`

	// ProviderParameters contains load balancer information specific to
	// the underlying infrastructure provider.
	//
	// +kubebuilder:default={type: "AWS"}
	ProviderParameters ProviderLoadBalancerParameters `json:"providerParameters,omitempty"`
}

// LoadBalancerScope is the scope at which a load balancer is exposed.
// +kubebuilder:validation:Enum=Internal;External
type LoadBalancerScope string

var (
	// InternalLoadBalancer is a load balancer that is exposed only on the
	// cluster's private network.
	InternalLoadBalancer LoadBalancerScope = "Internal"

	// ExternalLoadBalancer is a load balancer that is exposed on the
	// cluster's public network (which is typically on the Internet).
	ExternalLoadBalancer LoadBalancerScope = "External"
)

// ProviderLoadBalancerParameters holds desired load balancer information
// specific to the underlying infrastructure provider.
//
// +union
type ProviderLoadBalancerParameters struct {
	// Type is the underlying infrastructure provider for the load balancer.
	// Allowed values are "AWS", "Azure", and "GCP".
	//
	// +unionDiscriminator
	// +kubebuilder:default=AWS
	Type LoadBalancerProviderType `json:"type,omitempty"`

	// AWS provides configuration settings that are specific to AWS
	// load balancers.
	//
	// If empty, defaults will be applied. See specific aws fields for
	// details about their defaults.
	//
	// +optional
	AWS *AWSLoadBalancerParameters `json:"aws,omitempty"`
}

// LoadBalancerProviderType is the underlying infrastructure provider for the
// load balancer. Allowed values are "AWS", "Azure", and "GCP".
//
// +kubebuilder:validation:Enum=AWS;Azure;GCP
type LoadBalancerProviderType string

const (
	AWSLoadBalancerProvider   LoadBalancerProviderType = "AWS"
	AzureLoadBalancerProvider LoadBalancerProviderType = "Azure"
	GCPLoadBalancerProvider   LoadBalancerProviderType = "GCP"
)

// AWSLoadBalancerParameters provides configuration settings that are specific to
// AWS load balancers.
type AWSLoadBalancerParameters struct {
	// Type is the type of AWS load balancer to manage.
	//
	// Valid values are:
	//
	// * "Classic": A Classic load balancer makes routing decisions at either the
	//   transport layer (TCP/SSL) or the application layer (HTTP/HTTPS). See
	//   the following for additional details:
	//
	//     https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#clb
	//
	// * "NLB": A Network load balancer makes routing decisions at the transport
	//   layer (TCP/SSL). See the following for additional details:
	//
	//     https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#nlb
	//
	// If unset, defaults to "Classic".
	//
	// +kubebuilder:default=Classic
	Type AWSLoadBalancerType `json:"type,omitempty"`
}

// AWSLoadBalancerType is the type of AWS load balancer to manage.
// +kubebuilder:validation:Enum=Classic;NLB
type AWSLoadBalancerType string

const (
	AWSClassicLoadBalancer AWSLoadBalancerType = "Classic"
	AWSNetworkLoadBalancer AWSLoadBalancerType = "NLB"
)

// NodePort is the schema to specify a network port for a NodePort Service.
type NodePort struct {
	// Name is an IANA_SVC_NAME within the NodePort Service.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`

	// PortNumber is the network port number to expose for the NodePort Service.
	// If unspecified, a port number will be assigned from the the cluster's
	// nodeport service range, i.e. --service-node-port-range flag
	// (default: 30000-32767).
	//
	// If specified, the number must:
	//
	// 1. Not be used by another NodePort Service.
	// 2. Be within the cluster's nodeport service range, i.e. --service-node-port-range
	//    flag (default: 30000-32767).
	// 3. Be a valid network port number, i.e. greater than 0 and less than 65536.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	PortNumber *int32 `json:"portNumber,omitempty"`
}

// ContainerPort is the schema to specify a network port for a container.
// A container port gives the system additional information about network
// connections a container uses, but is primarily informational.
type ContainerPort struct {
	// Name is an IANA_SVC_NAME within the pod.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`

	// PortNumber is the network port number to expose on the envoy pod.
	// The number must be greater than 0 and less than 65536.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	PortNumber int32 `json:"portNumber"`
}

const (
	// ContourAvailableConditionType indicates that the contour is running
	// and available.
	ContourAvailableConditionType = "Available"
)

// ContourStatus defines the observed state of Contour.
type ContourStatus struct {
	// AvailableContours is the number of observed available replicas
	// according to the Contour deployment. The deployment and its pods
	// will reside in the namespace specified by spec.namespace.name of
	// the contour.
	AvailableContours int32 `json:"availableContours"`

	// AvailableEnvoys is the number of observed available pods from
	// the Envoy daemonset. The daemonset and its pods will reside in the
	// namespace specified by spec.namespace.name of the contour.
	AvailableEnvoys int32 `json:"availableEnvoys"`

	// Conditions represent the observations of a contour's current state.
	// Known condition types are "Available". Reference the condition type
	// for additional details.
	//
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Contour{}, &ContourList{})
}
