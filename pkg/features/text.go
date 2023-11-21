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

const (
	// DefaultTextSize is used for all Text features unless configured
	// explicitly otherwise
	DefaultTextSize = 14.0 // units: points. So about 4.93mm
)

// Text describes a text feature
type Text struct {
	Origin geometry.Point
	Alignment
	Purpose
	Text string
	// Size somehow describes the size of the text. Specific units not defined
	// here but probably safest to use points.
	Size float64
	// Radians. 0 for normal orientation.
	Rotate float64
}

// TextOptionFunc functions mutate a Text structure
type TextOptionFunc func(*Text)

// WithAlignment is a Text option function that sets alignment for a text feature
func WithAlignment(align Alignment) TextOptionFunc {
	return func(t *Text) {
		t.Alignment = align
	}
}

// WithSize is a Text option function that sets size for a text feature
func WithSize(size float64) TextOptionFunc {
	return func(t *Text) {
		t.Size = size
	}
}

// WithRotation is a Text option function that configures rotation (in radians)
// for a text feature
func WithRotation(r float64) TextOptionFunc {
	return func(t *Text) {
		t.Rotate = r
	}
}

// NewText creates a new Text feature
func NewText(origin geometry.Point, text string, options ...TextOptionFunc) *Text {
	t := &Text{
		Origin: origin,
		Text:   text,
		Size:   14.0, // default of
	}
	for _, opt := range options {
		opt(t)
	}
	return t
}

// GetPurpose returns the intended purpose of this feature
func (t *Text) GetPurpose() Purpose {
	return t.Purpose
}

// SetPurpose sets the purpose for a text feature. This isn't going to make
// much sense set to Cutout, I think, but it satisfies the interface.
func (t *Text) SetPurpose(purpose Purpose) {
	t.Purpose = purpose
}

// String satisfies the Stringer interface to aid debug printing
func (t Text) String() string {
	return fmt.Sprintf("Text(x=%.2f, y=%.2f, size=%.2f, align=%s, purpose=%s, text=%q)",
		t.Origin.X, t.Origin.Y, t.Size, t.Alignment.String(), t.Purpose.String(), t.Text)
}
