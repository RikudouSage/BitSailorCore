package tfa

import "go.chrastecky.dev/bitsailor-core/bitwarden/dto"

type Request interface {
	SetProvider(kind *dto.TFAKind)
	SetRemember(remember *bool)
}
