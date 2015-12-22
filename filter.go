package dsd

type Filter interface {
	Step(in float32) float32
}

type FIR struct {
	state []float32
	coefs []float32
}

func NewFIR(coefs []float32) *FIR {
	c := make([]float32, len(coefs))
	s := make([]float32, len(coefs))
	for i, v := range coefs {
		c[i] = v
	}
	return &FIR{coefs: c, state: s}
}

func (f *FIR) Step(in float32) float32 {
	var result float32
	for i := 0; i < len(f.state)-1; i++ {
		f.state[i] = f.state[i+1]
		result += f.state[i+1] * f.coefs[i]
	}

	f.state[len(f.state)-1] = in
	result += f.coefs[len(f.state)-1] * in
	return result
}

type IIR struct {
	A     []float32
	B     []float32
	state []float32
}

func NewIIR(A, B []float32) *IIR {
	l := 0
	if len(A)+1 > len(B) {
		l = len(A) + 1
	} else {
		l = len(B)
	}

	a := make([]float32, l-1)
	b := make([]float32, l)

	for i, v := range A {
		a[i] = v
	}

	for i, v := range B {
		b[i] = v
	}

	s := make([]float32, l)
	return &IIR{A: a, B: b, state: s}
}

func (f *IIR) Step(in float32) float32 {
	var result float32
	for i := 0; i < len(f.state)-1; i++ {
		f.state[i+1] = f.state[i]
		result += -1 * f.state[i] * f.A[i]
	}
	f.state[0] = result
	result = 0
	for i, v := range f.B {
		result += f.state[i] * v
	}
	return result
}
