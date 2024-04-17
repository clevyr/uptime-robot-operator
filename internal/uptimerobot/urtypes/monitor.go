package urtypes

const (
	MonitorPaused = iota
	MonitorRunning
)

//go:generate enumer -type MonitorType -trimprefix Type -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=HTTPS;Keyword;Ping;Port;Heartbeat

type MonitorType uint8

const (
	TypeHTTPS MonitorType = iota + 1
	TypeKeyword
	TypePing
	TypePort
	TypeHeartbeat
)

//go:generate enumer -type MonitorAuthType -trimprefix Auth -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=Basic;Digest

type MonitorAuthType uint8

const (
	AuthBasic MonitorAuthType = iota + 1
	AuthDigest
)
