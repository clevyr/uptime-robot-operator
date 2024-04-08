package urtypes

//go:generate enumer -type PortType -json -trimprefix Port

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=HTTP;FTP;SMTP;POP3;IMAP;Custom

type PortType uint8

const (
	PortHTTP PortType = iota + 1
	PortFTP
	PortSMTP
	PortPOP3
	PortIMAP
	PortCustom
)
