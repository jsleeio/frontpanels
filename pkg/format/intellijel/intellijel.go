// Copyright 2023 John Slee <jslee@jslee.io>
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

// Package intellijel implements the Intellijel 1U module panel format
package intellijel

import (
	"github.com/jsleeio/frontpanels/pkg/format/eurorack"
	"github.com/jsleeio/frontpanels/pkg/geometry"
)

// based on https://intellijel.com/support/1u-technical-specifications/
const (
	// PanelHeight1U represents the total height of an Intellijel 1U panel, in
	// millimetres
	PanelHeight1U = 39.65

	// ExtraMountingHolesThreshold represents the panel width threshold beyond
	// which additional mounting holes are required
	ExtraMountingHolesThreshold = 6

	// MountingHolesLeftOffset represents the distance of the first mounting
	// hole from the left edge of the panel, in millimetres
	MountingHolesLeftOffset = eurorack.MountingHolesLeftOffset

	// MountingHoleTopY1U represents the Y value for the top row of 1U mounting
	// holes, in millimetres
	MountingHoleTopY1U = PanelHeight1U - 3.00

	// MountingHoleBottomY1U represents the Y value for the bottom row of 1U
	// mounting holes, in millimetres
	MountingHoleBottomY1U = 3.00

	// MountingHoleDiameter represents the diameter of a Eurorack system
	// mounting hole, in millimetres
	MountingHoleDiameter = eurorack.MountingHoleDiameter

	// HP represents horizontal pitch in a Eurorack frame, in millimetres
	HP = eurorack.HP

	// HorizontalFit indicates the panel tolerance adjustment for the format
	HorizontalFit = 0.25

	// CornerRadius indicates the corner radius for the format. Eurorack doesn't
	// believe in such things.
	CornerRadius = 0.0

	// RailHeightFromMountingHole is used to determine how much space exists.
	// See discussion in github.com/jsleeio/pkg/panel. 5mm is a good safe
	// figure for all known-used Eurorack rail types
	RailHeightFromMountingHole = eurorack.RailHeightFromMountingHole
)

// Intellijel implements the panel.Panel interface and encapsulates the physical
// characteristics of a Intellijel panel
type Intellijel struct {
	HP int
}

// NewIntellijel constructs a new Intellijel object
func NewIntellijel(hp int) *Intellijel {
	return &Intellijel{HP: hp}
}

// Width returns the width of a Intellijel panel, in millimetres
func (i Intellijel) Width() float64 {
	if i.HP == 1 {
		// Special case: 1hp panels according to the Doepfer specification should
		// be 5.00mm wide, and at this size, we don't have much room for error.
		// Return 0.0 for HorizontalFit() and 5.00 for Width()
		return 5.00
	}
	return HP * float64(i.HP)
}

// Height returns the height of a Intellijel panel, in millimetres
func (i Intellijel) Height() float64 {
	return PanelHeight1U
}

// MountingHoleDiameter returns the Intellijel system mounting hole size, in
// millimetres
func (i Intellijel) MountingHoleDiameter() float64 {
	return MountingHoleDiameter
}

// MountingHoles generates a set of Point objects representing the mounting
// hole locations of a Intellijel panel
func (i Intellijel) MountingHoles() []geometry.Point {
	lhsx := MountingHolesLeftOffset
	// special case; 1HP Eurorack panels are narrower than MountingHolesLeftOffset.
	// I'm not completely sure what the correct thing to do here is but it SEEMS
	// logical to move it left by 1HP, leaving the hole pretty close to the middle
	// of a 1HP panel.
	//
	// @negativspace on ModWiggler says he leaves the hole in the centre on 1hp,
	// panels, which makes sense, so we'll do that too. With a 5mm panel width
	// there's not a lot of meat left on either side of an M3 screw hole...
	if i.HP == 1 {
		lhsx = i.Width() / 2.0
	}
	holes := []geometry.Point{
		{X: lhsx, Y: MountingHoleBottomY1U},
		{X: lhsx, Y: MountingHoleTopY1U},
	}
	if i.HP > ExtraMountingHolesThreshold {
		rhsx := MountingHolesLeftOffset + HP*(float64(i.HP-3))
		holes = append(holes, geometry.Point{X: rhsx, Y: MountingHoleBottomY1U})
		holes = append(holes, geometry.Point{X: rhsx, Y: MountingHoleTopY1U})
	}
	return holes
}

// HorizontalFit indicates the panel tolerance adjustment for the format
func (i Intellijel) HorizontalFit() float64 {
	if i.HP == 1 {
		// Special case: 1hp panels according to the Doepfer specification should
		// be 5.00mm wide, and at this size, we don't have much room for error.
		// Return 0.0 for HorizontalFit() and 5.00 for Width()
		return 0.00
	}
	return HorizontalFit
}

// CornerRadius indicates the corner radius for the format
func (i Intellijel) CornerRadius() float64 {
	return CornerRadius
}

// RailHeightFromMountingHole is used to calculate space between rails
func (i Intellijel) RailHeightFromMountingHole() float64 {
	return RailHeightFromMountingHole
}

// MountingHoleTopY returns the Y coordinate for the top row of mounting
// holes
func (i Intellijel) MountingHoleTopY() float64 {
	return MountingHoleTopY1U
}

// MountingHoleBottomY returns the Y coordinate for the bottom row of
// mounting holes
func (i Intellijel) MountingHoleBottomY() float64 {
	return MountingHoleBottomY1U
}

// HeaderLocation returns the location of the header text. Intellijel 1U has
// mounting rails so this is typically aligned with the top mounting screw
func (i Intellijel) HeaderLocation() geometry.Point {
	return geometry.Point{X: i.Width() / 2.0, Y: i.MountingHoleTopY()}
}

// FooterLocation returns the location of the footer text. Intellijel 1U has
// mounting rails so this is typically aligned with the bottom mounting screw
func (i Intellijel) FooterLocation() geometry.Point {
	return geometry.Point{X: i.Width() / 2.0, Y: i.MountingHoleBottomY()}
}
