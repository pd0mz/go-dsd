// Package dsd implements Digital Speec Decoders (and Encoders)
package dsd

import "errors"

var (
	ErrClosed         = errors.New("dsd: voice stream not open")
	ErrNotImplemented = errors.New("dsd: not implemented")
)

// Decoder transforms encoded bits to float32 samples.
type Decoder interface {
	Decode(bits []byte) ([]float32, error)
}

// Encoder transforms float32 samples to encoded bits.
type Encoder interface {
	Encode([]float32) ([]byte, error)
}

// VoiceStream is a complete codec with a decode and an encoder.
type VoiceStream interface {
	Decoder
	Encoder
	Close() error
}
