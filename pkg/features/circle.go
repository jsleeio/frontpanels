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

import (
	"fmt"

	"github.com/jsleeio/frontpanels/pkg/geometry"
)

// Circle describes a circular feature
type Circle struct {
	Origin geometry.Point
	Radius float64
	Purpose
}

// NewCircle initializes a new Circle object
func NewCircle(origin geometry.Point, radius float64) *Circle {
	if radius < 0.0 {
		panic("circle radius must be a positive value")
	}
	return &Circle{Origin: origin, Radius: radius}
}

// GetPurpose returns the intended purpose of this feature
func (c *Circle) GetPurpose() Purpose {
	return c.Purpose
}

// SetPurpose sets the purpose for a circle feature
func (c *Circle) SetPurpose(purpose Purpose) {
	c.Purpose = purpose
}

// String satisfies the Stringer interface to aid debug printing
func (c *Circle) String() string {
	return fmt.Sprintf("Circle(x=%.2f, y=%.2f, r=%.2f, purpose=%s)",
		c.Origin.X, c.Origin.Y, c.Radius, c.Purpose.String())
}
