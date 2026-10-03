package dto

type TFAKind int

const (
	TFAKindAuthenticator TFAKind = iota
	TFAKindEmail
)
