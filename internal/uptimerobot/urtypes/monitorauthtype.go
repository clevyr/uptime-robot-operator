package urtypes

//go:generate enumer -type MonitorAuthType -json -trimprefix Auth

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=Basic;Digest

type MonitorAuthType uint8

const (
	AuthBasic MonitorAuthType = iota + 1
	AuthDigest
)
