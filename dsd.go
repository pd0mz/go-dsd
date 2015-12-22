// Package dsd implements Digital Speec Decoders (and Encoders)
package dsd

import "errors"

var (
	ErrClosed         = errors.New("dsd: voice stream not open")
	ErrNotImplemented = errors.New("dsd: not implemented")
)

type Decoder interface {
	Decode(bits []byte) ([]float32, error)
}

type Encoder interface {
	Encode([]int32) ([]byte, error)
}

type VoiceStream interface {
	Decoder
	Encoder
	Close() error
}
