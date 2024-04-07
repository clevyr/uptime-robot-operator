package urtypes

//go:generate enumer -type HTTPMethod -json -trimprefix HTTP

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=HEAD;GET;POST;PUT;PATCH;DELETE;OPTIONS

type HTTPMethod uint8

const (
	HTTPHEAD HTTPMethod = iota
	HTTPGET
	HTTPPOST
	HTTPPUT
	HTTPPATCH
	HTTPDELETE
	HTTPOPTIONS
)
