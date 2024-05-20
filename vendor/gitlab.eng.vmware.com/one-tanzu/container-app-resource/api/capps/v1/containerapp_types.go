/*
Copyright 2024.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=capp
// +kubebuilder:printcolumn:name=Name,JSONPath=.metadata.name,description=Name,type=string
// +kubebuilder:printcolumn:name=Description,JSONPath=.spec.description,description=Description,type=string
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// ContainerApp is the Schema for the containerapps API
type ContainerApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContainerAppSpec   `json:"spec,omitempty"`
	Status ContainerAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ContainerAppList contains a list of ContainerApp
type ContainerAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContainerApp `json:"items"`
}

// ContainerAppSpec defines the desired state of ContainerApp
type ContainerAppSpec struct {
	// Short description of the App
	Description string `json:"description,omitempty"`
	// Contact information about the App owners
	Contact map[string]string `json:"contact,omitempty"`

	// Contain details about the app contents: source repository, revision, build dates, etc.
	// Well known keys: "summary"
	Content map[ContentKey]string `json:"content,omitempty"`
	// Specifies build instructions for this app
	Build *Build `json:"build,omitempty"`
	// The image of the app
	Image string `json:"image,omitempty"`

	// Environment independent env variables (immutable)
	NonSecretEnv []NonSecretEnvVar `json:"nonSecretEnv,omitempty"`
	// Environment dependent env variables (mutable)
	SecretEnv []SecretEnvVar `json:"secretEnv,omitempty"`

	// Default number of app replicas
	Replicas *int `json:"replicas,omitempty"`

	// Wanted set of resources for each replica (cpu/mem/...)
	// Well known keys: "cpu", "memory"
	Resources map[ResourcesKey]string `json:"resources,omitempty"`

	// Determines the health of the App
	Probes *Probes `json:"probes,omitempty"`

	// The named network ports of the App
	Ports []Port `json:"ports,omitempty"`

	// The service bindings that the app is using
	ServiceBindings []ServiceBinding `json:"serviceBindings,omitempty"`

	RelatedRefs []RelatedRef `json:"relatedRefs,omitempty"`
}

type ContentKey string

// These are valid content types. The list is not exhaustive.
// The summary type is guaranteed to exist.
const (
	ContentKeySummary ContentKey = "summary"
)

type Build struct {
	// Path is relative to location of ContainerApp definition
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:example=".."
	Path string `json:"path,omitempty"`
	// Build artifact type (buildpack, Helm, Docker image)
	// +optional
	Buildpacks *Buildpacks `json:"buildpacks,omitempty"`
	// Build-dependent environment variables
	// +optional
	NonSecretEnv []NonSecretEnvVar `json:"nonSecretEnv,omitempty"`
}

type Buildpacks struct{}

type NonSecretEnvVar struct {
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type SecretEnvVar struct {
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Name         string        `json:"name"`
	SecretKeyRef *SecretKeyRef `json:"secretKeyRef,omitempty"`
}

type SecretKeyRef struct {
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Name string `json:"name"`
	Key  string `json:"key,omitempty"`
}

type ResourcesKey string

// Well known resources
const (
	ResourcesKeyCPU    ResourcesKey = "cpu"
	ResourcesKeyMemory ResourcesKey = "memory"
)

type Probes struct {
	// Periodic probe of the App liveness.
	// App will be restarted if the probe fails.
	Liveness  *Probe `json:"liveness,omitempty"`
	Readiness *Probe `json:"readiness,omitempty"`
	Statup    *Probe `json:"startup,omitempty"`
}

type Probe struct {
	// The action taken to determine the health of a container
	ProbeHandler `json:",inline"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
	// Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
	// value overrides the value provided by the pod spec.
	// Value must be non-negative integer. The value zero indicates stop immediately via
	// the kill signal (no opportunity to shut down).
	// This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
	// Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
}

// ProbeHandler defines a specific action that should be taken in a probe.
// One and only one of the fields must be specified.
type ProbeHandler struct {
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
type HTTPHeader struct {
	// The header field name.
	// This will be canonicalized upon output, so case-variant names will be understood as the same header.
	Name string `json:"name"`
	// The header field value
	Value string `json:"value"`
}

// HTTPGetAction describes an action based on HTTP Get requests.
type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `json:"path,omitempty"`
	// Name or number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	Port int `json:"port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional
	Host string `json:"host,omitempty"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +optional
	Scheme URIScheme `json:"scheme,omitempty"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty"`
}

// URIScheme identifies the scheme used for connection to a host for Get actions
// +enum
type URIScheme string

const (
	// URISchemeHTTP means that the scheme used will be http://
	URISchemeHTTP URIScheme = "HTTP"
	// URISchemeHTTPS means that the scheme used will be https://
	URISchemeHTTPS URIScheme = "HTTPS"
)

type Port struct {
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Name string `json:"name"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int `json:"port"`
}

type ServiceBinding struct {
	// +required
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +required
	// +kubebuilder:validation:Required
	Type string `json:"type"`
}

type RelatedRef struct {
	// +required
	// +kubebuilder:validation:Required
	//
	// Designates the type of reference the relatedRef refers to
	For      RelatedRefFor `json:"for"`
	APIGroup string        `json:"apiGroup,omitempty"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
	// The keyPath within the resource that depicts this ref.
	// The keyPath is in jsonpath format. In case Kind is a Secret, it is an established pattern to embed whole yaml/json files within a single secret key.
	// In this case, one can use the arrow -> operator to designate a jsonpath that selects a key within the embedded yaml structure.
	//
	// Here are a few examples:
	// ".data.replicas" -- selects the .data.replicas property within the secret
	// ".data.'values\.yaml'->[yaml].data.replicas" -- selects the .data.replicas field within the values.yaml file which is stored in the secret under .data.values.yaml path
	KeyPath       string `json:"keyPath,omitempty"`
	LabelSelector string `json:"labelSelector,omitempty"`
}

type RelatedRefFor string

const (
	RelatedRefFor_KubernetesListReplicas         RelatedRefFor = "kubernetes.list-replicas"
	RelatedRefFor_KubernetesSetSecretEnv         RelatedRefFor = "kubernetes.set-secret-env"
	RelatedRefFor_KubernetesScaleReplicas        RelatedRefFor = "kubernetes.scale-replicas"
	RelatedRefFor_KubernetesScaleResources       RelatedRefFor = "kubernetes.scale-resources"
	RelatedRefFor_KubernetesServiceBindingTarget RelatedRefFor = "kubernetes.service-binding-target"
	RelatedRefFor_KubernetesServiceTarget        RelatedRefFor = "kubernetes.service-target"
)

// ContainerAppStatus defines the observed state of ContainerApp
type ContainerAppStatus struct {
	// Conditions represent the latest available observations of a ContainerApp's current state.
	// +optional
	Conditions []ContainerAppCondition `json:"conditions,omitempty"`
	Replicas   []Replica               `json:"replicas,omitempty"`
}

type ContainerAppCondition struct {
	Type          ContainerAppConditionType `json:"type"`
	BaseCondition `json:",inline"`
}

// +kubebuilder:validation:Enum=Deploying;DeploySucceeded;DeployFailed
type ContainerAppConditionType string

const (
	ContainerAppDeploying       ContainerAppConditionType = "Deploying"
	ContainerAppDeploySucceeded ContainerAppConditionType = "DeploySucceeded"
	ContainerAppDeployFailed    ContainerAppConditionType = "DeployFailed"
)

type BaseCondition struct {
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

type Replica struct {
	Name string `json:"name,omitempty"`

	// Human-readable summary of the App contents
	// On Kubernetes, corresponds to the Pod's "containerapp.apps.tanzu.vmware.com/content-summary" annotation
	ContentSummary string `json:"contentSummary,omitempty"`

	// The availability target of this replica
	// On Kubernetes, corresponds to containing Namespace's "spaces.tanzu.vmware.com/availability-target" annotation
	AvailabilityTarget string `json:"availabilityTarget,omitempty"`

	Condition []ReplicaCondition `json:"conditions,omitempty"`
}

type ReplicaCondition struct {
	Type          ReplicaConditionType `json:"type"`
	BaseCondition `json:",inline"`
}

// +kubebuilder:validation:Enum=Initializing;Running;Succeeded;Failed
type ReplicaConditionType string

const (
	ReplicaInitializing ReplicaConditionType = "Initializing"
	ReplicaRunning      ReplicaConditionType = "Running"
	ReplicaSucceeded    ReplicaConditionType = "Succeeded"
	ReplicaFailed       ReplicaConditionType = "Failed"
)

func (a ContainerApp) NamespaceScoped() bool {
	return true
}

func init() {
	SchemeBuilder.Register(&ContainerApp{}, &ContainerAppList{})
}
