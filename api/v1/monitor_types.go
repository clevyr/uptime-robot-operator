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
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot"
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MonitorContactRef attaches alert contacts. If blank, the default will be used.
type MonitorContactRef struct {
	corev1.LocalObjectReference `json:",inline"`

	uptimerobot.MonitorContactCommon `json:",inline"`
}

// MonitorSpec defines the desired state of Monitor
type MonitorSpec struct {
	// Interval defines the reconcile interval.
	//+kubebuilder:default:="24h"
	Interval metav1.Duration `json:"interval,omitempty"`

	// Prune enables garbage collection.
	//+kubebuilder:default:=true
	Prune bool `json:"prune,omitempty"`

	// Monitor configures the Uptime Robot monitor.
	Monitor uptimerobot.Monitor `json:"monitor"`

	// +kubebuilder:default:={{}}
	Contacts []MonitorContactRef `json:"contacts,omitempty"`
}

// MonitorStatus defines the observed state of Monitor
type MonitorStatus struct {
	Created bool                `json:"created"`
	ID      string              `json:"id,omitempty"`
	Type    urtypes.MonitorType `json:"type,omitempty"`
	Status  uint8               `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.monitor.status,statuspath=.status.status
//+kubebuilder:printcolumn:name="Created",type="boolean",JSONPath=".status.created"
//+kubebuilder:printcolumn:name="Friendly Name",type="string",priority=1,JSONPath=".spec.monitor.friendlyName"
//+kubebuilder:printcolumn:name="URL",type="string",priority=1,JSONPath=".spec.monitor.url"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Monitor is the Schema for the monitors API
type Monitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitorSpec   `json:"spec,omitempty"`
	Status MonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MonitorList contains a list of Monitor
type MonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitor{}, &MonitorList{})
}
