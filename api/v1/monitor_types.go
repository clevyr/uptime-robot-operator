/*
Copyright 2025.

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
	"fmt"
	"strings"
	"time"

	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MonitorSpec defines the desired state of Monitor.
type MonitorSpec struct {
	// Interval defines the reconcile interval.
	//+kubebuilder:default:="24h"
	Interval *metav1.Duration `json:"interval,omitempty"`

	// Prune enables garbage collection.
	//+kubebuilder:default:=true
	Prune bool `json:"prune,omitempty"`

	// Account references this object's Account. If not specified, the default will be used.
	Account corev1.LocalObjectReference `json:"account,omitempty"`

	// Monitor configures the Uptime Robot monitor.
	Monitor MonitorValues `json:"monitor"`

	//+kubebuilder:default:={{}}
	Contacts []MonitorContactRef `json:"contacts,omitempty"`

	// SourceRef optionally references the object that created this Monitor.
	SourceRef *corev1.TypedLocalObjectReference `json:"sourceRef,omitempty"`
}

// MonitorStatus defines the observed state of Monitor.
type MonitorStatus struct {
	Ready  bool                `json:"ready"`
	ID     string              `json:"id,omitempty"`
	Type   urtypes.MonitorType `json:"type,omitempty"`
	Status uint8               `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.monitor.status,statuspath=.status.status
//+kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
//+kubebuilder:printcolumn:name="Friendly Name",type="string",JSONPath=".spec.monitor.name"
//+kubebuilder:printcolumn:name="URL",type="string",JSONPath=".spec.monitor.url"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Monitor is the Schema for the monitors API.
type Monitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitorSpec   `json:"spec,omitempty"`
	Status MonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:generate=true
//+kubebuilder:validation:XValidation:rule="self.type != 'Keyword' || has(self.keyword)", message="Keyword config is required if type is Keyword"
//+kubebuilder:validation:XValidation:rule="self.type != 'Port' || has(self.port)", message="Port config is required if type is Port"

type MonitorValues struct {
	// Name sets the name that is shown in Uptime Robot.
	Name string `json:"name"`

	// URL is the URL or IP to monitor, including the scheme.
	URL string `json:"url"`

	// Type chooses the monitor type.
	//+kubebuilder:default:=HTTPS
	Type urtypes.MonitorType `json:"type,omitempty"`

	// Interval is the monitoring interval.
	//+kubebuilder:default:="60s"
	Interval *metav1.Duration `json:"interval,omitempty"`

	// Status toggles pause status for the monitor. 0 is paused, 1 is running.
	//+kubebuilder:default:=1
	Status uint8 `json:"status,omitempty"`

	// Timeout is the monitor timeout.
	//+kubebuilder:default:="30s"
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Method defines the HTTP verb to use.
	//+kubebuilder:default:="HEAD"
	Method urtypes.HTTPMethod `json:"method,omitempty"`

	// POST configures POST, PUT, PATCH, DELETE, and OPTIONS requests.
	POST *MonitorPOST `json:"post,omitempty"`

	// Keyword provides configuration for the Keyword monitor type.
	Keyword *MonitorKeyword `json:"keyword,omitempty"`

	// Port provides configuration for the Port monitor type.
	Port *MonitorPort `json:"port,omitempty"`

	// Auth enables monitor auth.
	Auth *MonitorAuth `json:"auth,omitempty"`
}

//+kubebuilder:object:generate=true

type MonitorKeyword struct {
	Type urtypes.KeywordType `json:"type"`

	//+kubebuilder:default:=false
	CaseSensitive *bool `json:"caseSensitive,omitempty"`

	Value string `json:"value"`
}

//+kubebuilder:validation:XValidation:rule="self.type != 'Custom' || has(self.number)", message="Number is required if type is Custom"
//+kubebuilder:validation:XValidation:rule="self.type == 'Custom' || !has(self.number)", message="Type must be Custom if Number is set"

type MonitorPort struct {
	Type urtypes.PortType `json:"type"`

	Number uint16 `json:"number,omitempty"`
}

type MonitorAuth struct {
	//+kubebuilder:default:="Basic"
	Type urtypes.MonitorAuthType `json:"type"`

	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	SecretName  string `json:"secretName,omitempty"`
	UsernameKey string `json:"usernameKey,omitempty"`
	PasswordKey string `json:"passwordKey,omitempty"`
}

type MonitorPOST struct {
	// Type defines the format of data to be sent with POST, PUT, PATCH, DELETE, and OPTIONS requests.
	//+kubebuilder:default:="KeyValue"
	Type urtypes.POSTType `json:"postType,omitempty"`

	// ContentType sets the Content-Type header for POST, PUT, PATCH, DELETE, and OPTIONS requests.
	//+kubebuilder:default:="text/html"
	ContentType urtypes.POSTContentType `json:"contentType,omitempty"`

	// Value is the JSON form of data to be sent with POST, PUT, PATCH, DELETE, and OPTIONS requests.
	Value string `json:"value,omitempty"`
}

// MonitorContactRef attaches alert contacts. If blank, the default will be used.
type MonitorContactRef struct {
	corev1.LocalObjectReference `json:",inline" mapstructure:",squash"`

	MonitorContactCommon `json:",inline" mapstructure:",squash"`
}

type MonitorContactCommon struct {
	// Threshold defines the number of minutes to wait to notify.
	//+kubebuilder:default:="1m"
	Threshold metav1.Duration `json:"threshold,omitempty"`

	// Recurrence defines the number of minutes between a repeat notification.
	// A value of 0, disables repeat notifications.
	Recurrence metav1.Duration `json:"recurrence,omitempty"`
}

type MonitorContact struct {
	ID string `json:"id"`

	MonitorContactCommon `json:",inline"`
}

func (m MonitorContact) String() string {
	return fmt.Sprintf(
		"%s_%d_%d",
		m.ID,
		int(m.MonitorContactCommon.Threshold.Round(time.Minute).Minutes()),
		int(m.MonitorContactCommon.Recurrence.Round(time.Minute).Minutes()),
	)
}

type MonitorContacts []MonitorContact

func (m MonitorContacts) String() string {
	results := make([]string, 0, len(m))
	for _, c := range m {
		results = append(results, c.String())
	}
	return strings.Join(results, "-")
}

//+kubebuilder:object:root=true

// MonitorList contains a list of Monitor.
type MonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitor{}, &MonitorList{})
}
