// Package cubicnoise provides 1D & 2D random noise generator
// with bicubic interpolation.
//
// This is a golang port of the java implementation of
// https://github.com/jobtalle/CubicNoise
package cubicnoise

import "math"

const (
	rndA = 134775813
	rndB = 1103515245
)

// CubicNoise is a random noise generator with bicubic
// interpolation.
type CubicNoise struct {
	seed    int32
	octave  int32
	periodX int32
	periodY int32
}

// New creates a randon noise generator.
//
// Zero for periodX and periodY are replaced with
// math.MaxInt32, which means infinite period.
func New(seed, octave, periodX, periodY int32) *CubicNoise {
	if periodX == 0 {
		periodX = math.MaxInt32
	}
	if periodY == 0 {
		periodY = math.MaxInt32
	}
	return &CubicNoise{
		seed:    seed,
		octave:  octave,
		periodX: periodX,
		periodY: periodY,
	}
}

// Sample1D returns a 1D noise.
func (n *CubicNoise) Sample1D(x float64) float64 {
	xi := int32(x) / n.octave
	lerp := float64(x)/float64(n.octave) - float64(xi)

	return interpolate(
		randomize(n.seed, tile(xi-1, n.periodX), 0),
		randomize(n.seed, tile(xi, n.periodX), 0),
		randomize(n.seed, tile(xi+1, n.periodX), 0),
		randomize(n.seed, tile(xi+2, n.periodX), 0),
		lerp)*0.5 + 0.25
}

// Sample2D returns a 2D noise.
func (n *CubicNoise) Sample2D(x, y float64) float64 {
	xi := int32(x) / n.octave
	lerpX := float64(x)/float64(n.octave) - float64(xi)
	yi := int32(y) / n.octave
	lerpY := float64(y)/float64(n.octave) - float64(yi)

	var xSamples [4]float64
	for i := int32(0); i < 4; i++ {
		tileY := tile(yi-1+i, n.periodY)
		xSamples[i] = interpolate(
			randomize(n.seed, tile(xi-1, n.periodX), tileY),
			randomize(n.seed, tile(xi, n.periodX), tileY),
			randomize(n.seed, tile(xi+1, n.periodX), tileY),
			randomize(n.seed, tile(xi+2, n.periodX), tileY),
			lerpX)
	}
	return interpolate(xSamples[0], xSamples[1], xSamples[2], xSamples[3], lerpY)*0.5 + 0.25
}

func randomize(seed, x, y int32) float64 {
	return float64(((x^y)*rndA^(seed+x))*(((rndB*x)<<16)^(rndB*y)-rndA)) / float64(math.MaxInt32)
}
func tile(coordinate, period int32) int32 {
	return coordinate % period
}

func interpolate(a, b, c, d, x float64) float64 {
	p := (d - c) - (a - b)
	return x*x*x*p + x*x*((a-b)-p) + x*(c-a) + b
}
