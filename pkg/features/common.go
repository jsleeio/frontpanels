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

// Package features encapsulate information about features on a panel, such as
// drill holes (Circles), legend text (Text), and so on.
package features

import "fmt"

// Purpose is intended to convey the application for the feature, eg. to
// differentiate decorative circles in a panel silkscreen vs. drill holes
type Purpose int

// Cutout et al specify purposes for panel features; this will be necessary
// in modules using Features in order to render them in the correct places.
// For example, in a Gerber file, a circle could be silkscreen markings or
// a drill hole
const (
	// Marking features are intended to be used to create aesthetic or legend
	// features on a panel. This is intentionally the first item in order to
	// make it the zero-value/default
	Marking Purpose = iota // this MUST be the first item
	// Cutout features are intended to be used to create a hole/void in a
	// panel
	Cutout // this MUST be the last item
)

// String satisfies the Stringer interface to aid debug printing
func (p Purpose) String() string {
	switch p {
	case Marking:
		return "marking"
	case Cutout:
		return "cutout"
	}
	panic(fmt.Sprintf("invalid Purpose value (valid range is %d..%d): %d",
		int(Marking), int(Cutout), int(p)))
}

// Feature interface. Intentionally small.
type Feature interface {
	GetPurpose() Purpose
	SetPurpose(Purpose)
}

// Alignment specifies an alignment relative to a feature, typically the
// feature origin. Most likely to be used with Text features.
type Alignment int

// TopLeft et al specify alignments relative to a defined point, eg.
// how Text is positioned relative to its origin.
const (
	TopLeft Alignment = iota
	TopCentre
	TopRight
	CentreLeft
	Centre
	CentreRight
	BottomLeft
	BottomCentre
	BottomRight
)

// String satisfies the Stringer interface to aid debug printing
func (a Alignment) String() string {
	switch a {
	case TopLeft:
		return "top-left"
	case TopCentre:
		return "top-centre"
	case TopRight:
		return "top-right"
	case CentreLeft:
		return "centre-left"
	case Centre:
		return "centre"
	case CentreRight:
		return "centre-right"
	case BottomLeft:
		return "bottom-left"
	case BottomCentre:
		return "bottom-centre"
	case BottomRight:
		return "bottom-right"
	}
	panic(fmt.Sprintf("invalid Alignment value (valid range is %d..%d): %d",
		int(TopLeft), int(BottomRight), int(a)))
}
