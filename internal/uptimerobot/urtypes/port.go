package urtypes

//go:generate go run github.com/dmarkham/enumer -type PortType -trimprefix Port -json -text

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
