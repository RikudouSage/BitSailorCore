package crypto

type KDFType int

const (
	KDFTypeSHA256 KDFType = iota
	KDFTypeArgon2ID
)
