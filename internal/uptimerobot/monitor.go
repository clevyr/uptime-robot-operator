package uptimerobot

import (
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:generate=true
//+kubebuilder:validation:XValidation:rule="self.type != 'Keyword' || has(self.keyword)", message="Keyword config is required if type is Keyword"

type Monitor struct {
	// FriendlyName sets the name that is shown in Uptime Robot.
	FriendlyName string `json:"friendlyName"`

	// URL is the URL or IP to monitor, including the scheme.
	URL string `json:"url"`

	// Type chooses the monitor type.
	//+kubebuilder:default:=HTTPS
	Type urtypes.MonitorType `json:"type,omitempty"`

	// Interval is the monitoring interval.
	//+kubebuilder:default:="60s"
	Interval metav1.Duration `json:"interval,omitempty"`

	// Status toggles pause status for the monitor. 0 is paused, 1 is running.
	//+kubebuilder:default:=1
	Status uint8 `json:"status,omitempty"`

	// Timeout is the monitor timeout.
	//+kubebuilder:default:="30s"
	Timeout metav1.Duration `json:"timeout,omitempty"`

	// HTTPMethod defines the HTTP verb to use.
	//+kubebuilder:default:="HEAD"
	HTTPMethod urtypes.HTTPMethod `json:"httpMethod,omitempty"`

	// Keyword provides configuration for the Keyword monitor type.
	Keyword *MonitorKeyword `json:"keyword,omitempty"`
}

type MonitorKeyword struct {
	Type urtypes.KeywordType `json:"type"`

	//+kubebuilder:default:=false
	CaseSensitive bool `json:"caseSensitive,omitempty"`

	Value string `json:"value"`
}
