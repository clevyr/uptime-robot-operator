package urtypes

//go:generate go run github.com/dmarkham/enumer -type POSTType -trimprefix Type -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=KeyValue;RawData

type POSTType uint8

const (
	TypeKeyValue POSTType = iota + 1
	TypeRawData
)

//go:generate go run github.com/dmarkham/enumer -type POSTContentType -linecomment -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=text/html;application/json

type POSTContentType uint8

const (
	ContentTypeHTML POSTContentType = iota // text/html
	ContentTypeJSON                        // application/json
)
