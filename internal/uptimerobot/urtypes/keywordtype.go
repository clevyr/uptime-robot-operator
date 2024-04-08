package urtypes

//go:generate enumer -type KeywordType -json -trimprefix Keyword

//+kubebuilder:validation:Type:=string
//+kubebuilder:validation:Enum:=Exists;NotExists

type KeywordType uint8

const (
	KeywordExists KeywordType = iota + 1
	KeywordNotExists
)
