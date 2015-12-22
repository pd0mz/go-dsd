package dsd

// #cgo CFLAGS: -Iinclude -I/usr/local/include -L/usr/local/lib
// #cgo LDFLAGS: -lm -lcodec2 -L/usr/local/lib
/*
#include <inttypes.h>
#include <codec2/codec2.h>

typedef struct CODEC2 codec2_t;

*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	Codec2Mode3200 uint8 = iota
	Codec2Mode2400
	Codec2Mode1600
	Codec2Mode1400
	Codec2Mode1300
	Codec2Mode1200
	Codec2Mode700
	Codec2Mode700B
)

type Codec2 struct {
	codec2          *C.codec2_t
	samplesPerFrame int
	bitsPerFrame    int
}

func NewCodec2(mode uint8) *Codec2 {
	vs := &Codec2{}
	vs.codec2 = C.codec2_create(C.int(mode))
	vs.samplesPerFrame = int(C.codec2_samples_per_frame(vs.codec2))
	vs.bitsPerFrame = int(C.codec2_bits_per_frame(vs.codec2))
	return vs
}

func (vs *Codec2) Close() error {
	if vs.codec2 == nil {
		return ErrClosed
	}

	C.codec2_destroy(vs.codec2)
	vs.codec2 = nil

	return nil
}

func (vs *Codec2) Decode(bits []byte) ([]float32, error) {
	if len(bits) != vs.bitsPerFrame {
		return nil, fmt.Errorf("dsd: codec2 required %d bits per frame, got %d", vs.bitsPerFrame, len(bits))
	}
	var (
		u = make([]uint16, vs.samplesPerFrame)
	)
	C.codec2_decode(vs.codec2, (*C.short)(unsafe.Pointer(&u[0])), (*C.uchar)(unsafe.Pointer(&bits[0])))
	return u16stof32s(u), nil
}

func (vs *Codec2) Encode(samples []int32) ([]byte, error) {
	if len(samples) != vs.samplesPerFrame {
		return nil, fmt.Errorf("dsd: codec2 required %d samples per frame, got %d", vs.samplesPerFrame, len(samples))
	}
	var (
		s    = i32stou16s(samples)
		bits = make([]byte, vs.bitsPerFrame)
	)
	for i, sample := range samples {
		s[i] = i32tou16(sample)
	}
	C.codec2_encode(vs.codec2, (*C.uchar)(unsafe.Pointer(&bits[0])), (*C.short)(unsafe.Pointer(&s[0])))
	return bits, nil
}

var _ (VoiceStream) = (*Codec2)(nil)
