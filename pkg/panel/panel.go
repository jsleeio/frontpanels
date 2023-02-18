// Copyright 2019 John Slee <jslee@jslee.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Package panel provides a common interface for panel descriptions, and basic
// implementations of some required functions.
package panel

import (
	"github.com/jsleeio/frontpanels/pkg/geometry"
)

// Panel types encapsulate physical characteristics of a panel.
//
// All coordinates, distances and sizes are indicated in millimetres
type Panel interface {
	// MountingHoles returns a list of geometry.Points indicating mounting hole locations
	MountingHoles() []geometry.Point

	// MountingHoleDiameter returns the appropriate mounting hole diameter for
	// the panel format
	MountingHoleDiameter() float64

	// Height returns the Y-dimension size for the panel, eg. 128.5mm for panels
	// in the Eurorack system. This does NOT include tolerance adjustments
	Height() float64

	// Height returns the X-dimension size for the panel. This does not include
	// tolerance adjustments
	Width() float64

	// HorizontalFit returns the panel tolerance amount in the horizontal axis.
	// When creating panel outlines, this tolerance amount should be added to
	// the left-edge X coordinate, and subtracted from the right-edge coordinate,
	// resulting in the panel being slightly narrower than the "correct" width.
	//
	// As this is intended to improve panel fit in a system with panels of varying
	// tolerances, this adjustment should only be applied to the left and right
	// edges of the outline, and NOT to the X coordinates of any other features
	// of the panel. (and especially not the mounting holes!)
	HorizontalFit() float64

	// CornerRadius returns the radius of the panel corners, if applicable. A
	// zero value will result in no corner segments being generated, and so the
	// board outline will consist of four straight lines.
	CornerRadius() float64

	// RailHeightFromMountingHole indicates how far up (from centre of bottom
	// mounting hole) or down (from centre of top mounting hole) the mounting
	// rail extends. This can be used to define KeepOut areas on the panel
	//
	// eg. For all Eurorack-related formats this is likely to be around 5mm,
	// though the exact figure will differ with rail type.
	//
	// This is primarily used to determine how much empty space there is between
	// the mounting rails, so best to err on the side of larger than smaller
	RailHeightFromMountingHole() float64

	// MountingHoleTopY returns the Y coordinate for the top row of mounting
	// holes
	MountingHoleTopY() float64

	// MountingHoleBottomY returns the Y coordinate for the bottom row of
	// mounting holes
	MountingHoleBottomY() float64

	// HeaderLocation returns the location of the header text
	HeaderLocation() geometry.Point

	// FooterLocation returns the location of the footer text
	FooterLocation() geometry.Point
}

// The following functions are probably appropriate for many front panel types,
// but not all, and so are provided here to be used as required.

// LeftX returns the left edge coordinate of a panel, adjusted for horizontal
// fit
func LeftX(spec Panel) float64 {
	return spec.HorizontalFit() / 2
}

// RightX returns the right edge coordinate of a panel, adjusted for horizontal
// fit
func RightX(spec Panel) float64 {
	return spec.Width() - spec.HorizontalFit()/2
}

// TopY returns the top edge coordinate of a panel
func TopY(spec Panel) float64 {
	return spec.Height()
}

// BottomY returns the bottom edge coordinate of a panel
func BottomY(spec Panel) float64 {
	return 0
}

// TopLeft returns the top-left corner coordinate of a panel, adjusted for
// horizontal fit
func TopLeft(spec Panel) geometry.Point {
	return geometry.Point{X: LeftX(spec), Y: TopY(spec)}
}

// TopRight returns the top-right corner coordinate of a panel, adjusted for
// horizontal fit
func TopRight(spec Panel) geometry.Point {
	return geometry.Point{X: RightX(spec), Y: TopY(spec)}
}

// BottomLeft returns the bottom-left corner coordinate of a panel, adjusted
// for horizontal fit
func BottomLeft(spec Panel) geometry.Point {
	return geometry.Point{X: LeftX(spec), Y: BottomY(spec)}
}

// BottomRight returns the bottom-right corner coordinate of a panel, adjusted
// for horizontal fit
func BottomRight(spec Panel) geometry.Point {
	return geometry.Point{X: RightX(spec), Y: BottomY(spec)}
}
