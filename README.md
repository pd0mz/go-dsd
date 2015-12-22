# go-dsd

Digital Speech Decoder (and Encoder), with some basic DSP processing functions.

## Documentation

Read the [package documentation](https://godoc.org/github.com/pd0mz/go-dsd).

## Supported codecs

### Advanced Multi-Band Excitation (AMBE)

Only decoding is supported using [mbelib](https://github.com/szechyjs/mbelib).
Make sure you have the latest development version, as it contains a lot of
improvements over the outdated stable release.

### Codec2

Codec2 is a low-bitrate speech audio codec (speech coding) that is patent free
and open source. Both decoding and encoding of Codec2 frames are supported.
