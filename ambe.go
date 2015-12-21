package dsd

// #cgo CFLAGS: -Iinclude
// #cgo LDFLAGS: -lm -lmbe
/*
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <mbelib.h>

#define VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT	160

typedef uint8_t flag_t;

typedef struct {
	flag_t bits[72];
} dmrpacket_payload_ambe_frame_bits_t;

typedef union {
	struct {
		flag_t bits[108*2];
	} raw;
	struct {
		dmrpacket_payload_ambe_frame_bits_t frames[3];
	} ambe_frames;
} dmrpacket_payload_voice_bits_t;

static uint8_t voicestreams_decode_deinterleave_matrix_w[36] = {
	0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
	0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
	0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
	0x00, 0x01, 0x00, 0x01, 0x00, 0x02,
	0x00, 0x02, 0x00, 0x02, 0x00, 0x02,
	0x00, 0x02, 0x00, 0x02, 0x00, 0x02
};

static uint8_t voicestreams_decode_deinterleave_matrix_x[36] = {
	0x17, 0x0a, 0x16, 0x09, 0x15, 0x08,
	0x14, 0x07, 0x13, 0x06, 0x12, 0x05,
	0x11, 0x04, 0x10, 0x03, 0x0f, 0x02,
	0x0e, 0x01, 0x0d, 0x00, 0x0c, 0x0a,
	0x0b, 0x09, 0x0a, 0x08, 0x09, 0x07,
	0x08, 0x06, 0x07, 0x05, 0x06, 0x04
};

static uint8_t voicestreams_decode_deinterleave_matrix_y[36] = {
	0x00, 0x02, 0x00, 0x02, 0x00, 0x02,
	0x00, 0x02, 0x00, 0x03, 0x00, 0x03,
	0x01, 0x03, 0x01, 0x03, 0x01, 0x03,
	0x01, 0x03, 0x01, 0x03, 0x01, 0x03,
	0x01, 0x03, 0x01, 0x03, 0x01, 0x03,
	0x01, 0x03, 0x01, 0x03, 0x01, 0x03
};

static uint8_t voicestreams_decode_deinterleave_matrix_z[36] = {
	0x05, 0x03, 0x04, 0x02, 0x03, 0x01,
	0x02, 0x00, 0x01, 0x0d, 0x00, 0x0c,
	0x16, 0x0b, 0x15, 0x0a, 0x14, 0x09,
	0x13, 0x08, 0x12, 0x07, 0x11, 0x06,
	0x10, 0x05, 0x0f, 0x04, 0x0e, 0x03,
	0x0d, 0x02, 0x0c, 0x01, 0x0b, 0x00
};

typedef struct voice_stream_s {
	mbe_parms cur_mp;
	mbe_parms prev_mp;
	mbe_parms prev_mp_enhanced;
	uint8_t decodequality;
} voicestream_t;

typedef struct {
	float samples[VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT];
} voicestreams_decoded_frame_t;

voicestream_t* init_voicestream(uint8_t decodequality) {
	voicestream_t *voicestream = malloc(sizeof(voicestream_t));
	voicestream->decodequality = decodequality;
	mbe_initMbeParms(
		&voicestream->cur_mp,
		&voicestream->prev_mp,
		&voicestream->prev_mp_enhanced);
	return voicestream;
}

void free_voicestream(voicestream_t *voicestream)
{
	free(voicestream);
}

float voicestreams_decoded_frame_sample(
	voicestreams_decoded_frame_t *decoded_frame,
	int sample)
{
	return decoded_frame->samples[sample];
}

voicestreams_decoded_frame_t *voicestreams_decode_ambe_frame(voicestream_t *voicestream, dmrpacket_payload_ambe_frame_bits_t *ambe_frame_bits) {
	static voicestreams_decoded_frame_t decoded_frame;
	char deinterleaved_ambe_frame_bits[4][24];
	uint8_t j;
	uint8_t *w, *x, *y, *z;
	int errs, errs2;
	char err_str[64];
	char ambe_d[49];

	// Deinterleaving
	w = voicestreams_decode_deinterleave_matrix_w;
	x = voicestreams_decode_deinterleave_matrix_x;
	y = voicestreams_decode_deinterleave_matrix_y;
	z = voicestreams_decode_deinterleave_matrix_z;

	for (j = 0; j < sizeof(dmrpacket_payload_ambe_frame_bits_t); j += 2) {
		deinterleaved_ambe_frame_bits[*w][*x] = ambe_frame_bits->bits[j];
		deinterleaved_ambe_frame_bits[*y][*z] = ambe_frame_bits->bits[j+1];
		w++;
		x++;
		y++;
		z++;
	}

	mbe_processAmbe3600x2450Framef(decoded_frame.samples, &errs, &errs2, err_str, deinterleaved_ambe_frame_bits, ambe_d, &voicestream->cur_mp, &voicestream->prev_mp, &voicestream->prev_mp_enhanced, voicestream->decodequality);

	if (errs2 > 0)
		fprintf(stderr, "dsd/ambe: mbelib decoding errors: %u %s\n", errs2, err_str);

	for (j = 0; j < VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT; j++)
		decoded_frame.samples[j] /= 32767.0;

	return &decoded_frame;
}

void voicestream_decode(
	voicestream_t *voicestream,
	char **bits,
	float **samplesptr) {

	int i, j;
	float *samples = (float *)samplesptr;
	dmrpacket_payload_voice_bits_t *voice_bits = (dmrpacket_payload_voice_bits_t *) bits;
	//fprintf(stderr, "dsd/ambe: decoding %d bits\n", sizeof(voice_bits->raw.bits));
	voicestreams_decoded_frame_t *decoded_frame;

	for (i = 0; i < 3; i++) {
		//fprintf(stderr, "dsd/ambe: decode frame %d\n", i);
		decoded_frame = voicestreams_decode_ambe_frame(
			voicestream,
			&voice_bits->ambe_frames.frames[i]);

		for (j = 0; j < VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT; j++)
			samples[j+(i*VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT)] = decoded_frame->samples[j];
	}
}


*/
import "C"
import (
	"fmt"
	"unsafe"
)

type AMBEVoiceDecodeQuality uint8

type AMBEVoiceStream struct {
	voicestream *C.voicestream_t
}

func NewAMBEVoiceStream(qual AMBEVoiceDecodeQuality) *AMBEVoiceStream {
	vs := &AMBEVoiceStream{}
	vs.voicestream = C.init_voicestream(C.uint8_t(qual))
	return vs
}

func (vs *AMBEVoiceStream) Decode(bits []byte) ([]float32, error) {
	if len(bits) != 216 {
		return nil, fmt.Errorf("need 216 bits, got %d", len(bits))
	}
	samples := make([]float32,
		C.VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT*3,
		C.VOICESTREAMS_DECODED_AMBE_FRAME_SAMPLES_COUNT*3)
	C.voicestream_decode(
		vs.voicestream,
		(**C.char)(unsafe.Pointer(&bits[0])),
		(**C.float)(unsafe.Pointer(&samples[0])))
	return samples, nil
}

func (vs *AMBEVoiceStream) Close() {
	C.free_voicestream(vs.voicestream)
	vs.voicestream = nil
}
