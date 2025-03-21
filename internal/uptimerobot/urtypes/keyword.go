package urtypes

//go:generate go run github.com/dmarkham/enumer -type KeywordType -trimprefix Keyword -json -text

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=Exists;NotExists

type KeywordType uint8

const (
	KeywordExists KeywordType = iota + 1
	KeywordNotExists
)
