package uptimerobot

//go:generate enumer -type MonitorType -json -trimprefix Type -transform lower

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=https;keyword;ping;port;heartbeat

type MonitorType uint8

const (
	TypeHTTPS MonitorType = iota + 1
	TypeKeyword
	TypePing
	TypePort
	TypeHeartbeat
)
