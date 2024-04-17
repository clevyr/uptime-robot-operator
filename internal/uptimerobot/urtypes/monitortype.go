package urtypes

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
