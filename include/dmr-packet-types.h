#ifndef DMRPACKET_TYPES_H_
#define DMRPACKET_TYPES_H_

#include "base/types.h"

typedef struct {
	flag_t bits[98+10+48+10+98]; // See DMR AI spec. page 85.
} dmrpacket_payload_bits_t;

typedef struct {
	flag_t bits[98*2];
} dmrpacket_payload_info_bits_t;

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

typedef struct {
	uint8_t bytes[sizeof(dmrpacket_payload_voice_bits_t)/8];
} dmrpacket_payload_voice_bytes_t;

#endif
