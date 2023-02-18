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

// Package panel faciliates generation of basic panel features (outline and
// mounting holes)
package panel

import (
	"github.com/jsleeio/frontpanels/pkg/features"
	"github.com/jsleeio/frontpanels/pkg/panel"
)

// GeneratePanelOutlineFeatures generates the basic features for a blank panel:
// an outline and some mounting holes
func GeneratePanelOutlineFeatures(p panel.Panel) []features.Feature {
	top := features.NewLine(panel.TopLeft(p), panel.TopRight(p), 0.1)
	top.SetPurpose(features.Cutout)
	bottom := features.NewLine(panel.BottomLeft(p), panel.BottomRight(p), 0.1)
	bottom.SetPurpose(features.Cutout)
	left := features.NewLine(panel.TopLeft(p), panel.BottomLeft(p), 0.1)
	left.SetPurpose(features.Cutout)
	right := features.NewLine(panel.TopRight(p), panel.BottomRight(p), 0.1)
	right.SetPurpose(features.Cutout)
	f := []features.Feature{top, bottom, left, right}
	for _, centre := range p.MountingHoles() {
		hole := features.NewCircle(centre, p.MountingHoleDiameter()/2.0)
		hole.SetPurpose(features.Cutout)
		f = append(f, hole)
	}
	return f
}
