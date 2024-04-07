package uptimerobot

//go:generate enumer -type MonitorType -json -trimprefix Type

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
