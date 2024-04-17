package urtypes

//go:generate enumer -type HTTPMethod -trimprefix HTTP -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=HEAD;GET;POST;PUT;PATCH;DELETE;OPTIONS

type HTTPMethod uint8

const (
	HTTPHEAD HTTPMethod = iota + 1
	HTTPGET
	HTTPPOST
	HTTPPUT
	HTTPPATCH
	HTTPDELETE
	HTTPOPTIONS
)
