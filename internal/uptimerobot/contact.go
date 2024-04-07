package uptimerobot

import (
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Contact struct {
	// FriendlyName sets the name that is shown in Uptime Robot.
	FriendlyName string `json:"friendlyName"`
}

type MonitorContactCommon struct {
	// Threshold defines the number of minutes to wait to notify.
	// +kubebuilder:default:="1m"
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
