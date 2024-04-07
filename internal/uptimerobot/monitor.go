package uptimerobot

import (
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

	// Timeout is the monitor timeout.
	//+kubebuilder:default:="30s"
	Timeout metav1.Duration `json:"timeout,omitempty"`

	// HTTPMethod defines the HTTP verb to use.
	//+kubebuilder:default:="GET"
	HTTPMethod urtypes.HTTPMethod `json:"httpMethod,omitempty"`
}
