package urtypes

//go:generate go run github.com/dmarkham/enumer -type HTTPMethod -trimprefix Method -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=HEAD;GET;POST;PUT;PATCH;DELETE;OPTIONS

type HTTPMethod uint8

const (
	MethodHEAD HTTPMethod = iota + 1
	MethodGET
	MethodPOST
	MethodPUT
	MethodPATCH
	MethodDELETE
	MethodOPTIONS
)
