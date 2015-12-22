package dsd

func f32toi32(f float32) int32 {
	return int32(f * 2147483648.0)
}

func f32stoi32s(f []float32) []int32 {
	var i = make([]int32, len(f))
	for j, k := range f {
		i[j] = f32toi32(k)
	}
	return i
}

func i32tou16(i int32) uint16 {
	if i < 0 {
		return uint16(i>>16) ^ 0x8000
	}
	return uint16(i>>16) | 0x8000
}

func i32stou16s(i []int32) []uint16 {
	var u = make([]uint16, len(i))
	for j, k := range i {
		u[j] = i32tou16(k)
	}
	return u
}

func u16tof32(u uint16) float32 {
	return (float32(u) - 32768.0) / 32768.0
}

func u16stof32s(u []uint16) []float32 {
	var f = make([]float32, len(u))
	for j, k := range u {
		f[j] = u16tof32(k)
	}
	return f
}

func u16toi32(u uint16) int32 {
	if u < 32768 {
		return int32(u|0x8000) << 16
	}
	return int32(u^0x8000) << 16
}

func u16stoi32s(u []uint16) []int32 {
	var i = make([]int32, len(u))
	for j, k := range u {
		i[j] = u16toi32(k)
	}
	return i
}
