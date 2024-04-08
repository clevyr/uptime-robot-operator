package uptimerobot

import (
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:generate=true
//+kubebuilder:validation:XValidation:rule="self.type != 'Keyword' || has(self.keyword)", message="Keyword config is required if type is Keyword"
//+kubebuilder:validation:XValidation:rule="self.type != 'Port' || has(self.port)", message="Port config is required if type is Port"

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

	// Port provides configuration for the Port monitor type.
	Port *MonitorPort `json:"port,omitempty"`

	// Auth enables monitor auth.
	Auth *MonitorAuth `json:"auth,omitempty"`
}

type MonitorKeyword struct {
	Type urtypes.KeywordType `json:"type"`

	//+kubebuilder:default:=false
	CaseSensitive bool `json:"caseSensitive,omitempty"`

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
