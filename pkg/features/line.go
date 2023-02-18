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
// drill holes (Circle), legend text (Text), and so on.
package features

import (
	"fmt"

	"github.com/jsleeio/frontpanels/pkg/geometry"
)

// Line describes a line feature
type Line struct {
	Start, End geometry.Point
	Thickness  float64
	Purpose
}

// NewLine initializes a new Line object
func NewLine(start, end geometry.Point, thickness float64) *Line {
	if thickness < 0.0 {
		panic("line thickness must be a positive value")
	}
	return &Line{Start: start, End: end, Thickness: thickness}
}

// GetPurpose returns the intended purpose of this feature
func (l *Line) GetPurpose() Purpose {
	return l.Purpose
}

// SetPurpose sets the purpose for a line feature
func (l *Line) SetPurpose(purpose Purpose) {
	l.Purpose = purpose
}

// String satisfies the Stringer interface to aid debug printing
func (l *Line) String() string {
	return fmt.Sprintf("Line(x1=%.2f, y1=%.2f, x2=%.2f, y2=%.2f, thickness=%.2f, purpose=%s)",
		l.Start.X, l.Start.Y, l.End.X, l.End.Y, l.Thickness, l.Purpose.String())
}
