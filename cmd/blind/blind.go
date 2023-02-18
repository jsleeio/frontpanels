// Package blind is a CLI tool for generating blind (blank) panels in
// (currently) Eurorack and its related 1U formats.
package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"reflect"

	"github.com/jsleeio/frontpanels/pkg/features"
	"github.com/jsleeio/frontpanels/pkg/format/eurorack"
	"github.com/jsleeio/frontpanels/pkg/format/intellijel"
	"github.com/jsleeio/frontpanels/pkg/format/pulplogic"
	"github.com/jsleeio/frontpanels/pkg/geometry"
	"github.com/jsleeio/frontpanels/pkg/panel"
	panelsource "github.com/jsleeio/frontpanels/pkg/sources/panel"

	_ "github.com/gmlewis/go-fonts/fonts/bitstreamverasansmono_bold"
	"github.com/gmlewis/go-gerber/gerber"
)

type config struct {
	format               string
	width                int
	name, header, footer string

	panel panel.Panel
}

func configure() (c config, p panel.Panel, err error) {
	flag.StringVar(&c.name, "name", "", "basename for generating Gerber filenames")
	flag.StringVar(&c.header, "header", "", "header text for panel")
	flag.StringVar(&c.footer, "footer", "", "footer text for panel")
	flag.StringVar(&c.format, "format", "eurorack", "panel format to generate (valid values: eurorack pulplogic intellijel)")
	flag.IntVar(&c.width, "width", 8, "panel width, in units appropriate for the format")
	flag.Parse()
	if c.width < 1 {
		err = errors.New("width must be greater than 0")
		return
	}
	switch c.format {
	case "eurorack":
		p = eurorack.NewEurorack(c.width)
	case "intellijel":
		p = intellijel.NewIntellijel(c.width)
	case "pulplogic":
		p = pulplogic.NewPulplogic(c.width)
	default:
		err = errors.New("invalid format specified")
		return
	}
	return
}

// panelOutline generates the basic features for a blank panel --- an outline
// and mounting holes
func panelOutline(p panel.Panel) []features.Feature {
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

func panelHeaderFooter(p panel.Panel, header, footer string) []features.Feature {
	// FIXME: figure out what to do with narrow panels — probably anything
	//        under 6hp. Maybe align centre-right?
	f := []features.Feature{}
	if header != "" {
		f = append(f, features.NewText(
			geometry.Point{X: p.Width() / 2.0, Y: p.MountingHoleTopY()},
			header,
			features.WithAlignment(features.Centre),
			features.WithSize(16.0), // assuming units are 1/72"
		))
	}
	if footer != "" {
		f = append(f, features.NewText(
			geometry.Point{X: p.Width() / 2.0, Y: p.MountingHoleBottomY()},
			footer,
			features.WithAlignment(features.Centre),
			features.WithSize(16.0), // assuming units are 1/72"
		))
	}
	return f
}

// mktext renders a line feature as a gerber primitive
func mkline(l *features.Line) gerber.Primitive {
	return gerber.Line(
		l.Start.X, l.Start.Y,
		l.End.X, l.End.Y,
		gerber.CircleShape, // gerber aperture stuff, probably leave it as-is
		l.Thickness,
	)
}

// mktext renders a line feature as a gerber primitive
func mkcircle(c *features.Circle) gerber.Primitive {
	return gerber.Circle(gerber.Point(c.Origin.X, c.Origin.Y), c.Radius*2.0)
}

// mktextopts copes with the incredibly annoying alignment options in the
// gerber/fonts packages
func mktextopts(t *features.Text) *gerber.TextOpts {
	m := map[features.Alignment]*gerber.TextOpts{
		features.TopLeft:      &gerber.TextOpts{XAlign: gerber.XLeft, YAlign: gerber.YTop},
		features.CentreLeft:   &gerber.TextOpts{XAlign: gerber.XLeft, YAlign: gerber.YCenter},
		features.BottomLeft:   &gerber.TextOpts{XAlign: gerber.XLeft, YAlign: gerber.YBottom},
		features.TopCentre:    &gerber.TextOpts{XAlign: gerber.XCenter, YAlign: gerber.YTop},
		features.Centre:       &gerber.TextOpts{XAlign: gerber.XCenter, YAlign: gerber.YCenter},
		features.BottomCentre: &gerber.TextOpts{XAlign: gerber.XCenter, YAlign: gerber.YBottom},
		features.TopRight:     &gerber.TextOpts{XAlign: gerber.XRight, YAlign: gerber.YTop},
		features.CentreRight:  &gerber.TextOpts{XAlign: gerber.XRight, YAlign: gerber.YCenter},
		features.BottomRight:  &gerber.TextOpts{XAlign: gerber.XRight, YAlign: gerber.YBottom},
	}
	opts, ok := m[t.Alignment]
	if !ok {
		panic("invalid text alignment value")
	}
	return opts
}

// mktext renders a text feature as a gerber primitive
func mktext(t *features.Text) gerber.Primitive {
	return gerber.Text(
		t.Origin.X, t.Origin.Y,
		1.0, // +1.0 = topsilk, -1.0 = bottomsilk *shrug*
		t.Text,
		"bitstreamverasansmono_bold",
		t.Size,
		mktextopts(t),
	)
}

type primitives struct {
	outlines, drills, silkscreens []gerber.Primitive
}

func newprimitives() *primitives {
	return &primitives{
		outlines:    []gerber.Primitive{},
		drills:      []gerber.Primitive{},
		silkscreens: []gerber.Primitive{},
	}
}

func (p *primitives) addoutline(pp gerber.Primitive) {
	p.outlines = append(p.outlines, pp)
}

func (p *primitives) addsilkscreen(pp gerber.Primitive) {
	p.silkscreens = append(p.silkscreens, pp)
}

func (p *primitives) adddrill(pp gerber.Primitive) {
	p.drills = append(p.drills, pp)
}

func collectPrimitives(feats []features.Feature, prims *primitives) {
	for _, item := range feats {
		switch f := item.(type) {
		case *features.Line:
			line := mkline(f)
			if f.GetPurpose() == features.Cutout {
				prims.addoutline(line)
			} else {
				prims.addsilkscreen(line)
			}
		case *features.Text:
			text := mktext(f)
			if f.GetPurpose() == features.Cutout {
				// text in outline layer is pretty much guaranteed to be a mistake
				log.Printf("warning: text feature in outline layer is probably an error: %v", f.String())
				prims.addoutline(text)
			} else {
				prims.addsilkscreen(text)
			}
		case *features.Circle:
			circle := mkcircle(f)
			if f.GetPurpose() == features.Cutout {
				// FIXME: fabs have upper limits on drill sizes, eg. 6.3mm for JLCPCB
				//        at this time of writing --- may need to drop larger ones in
				//        the outline layer instead. But this will be fab-dependent...
				prims.adddrill(circle)
			} else {
				prims.addsilkscreen(circle)
			}
		default:
			log.Printf("warning: unsupported feature type: %s", reflect.TypeOf(f).Kind().String())
		}
	}
}

// generate a bunch of random lines that fit between the rails
func randomLines(panel panel.Panel, n int) []features.Feature {
	lines := []features.Feature{}
	rxy := func() geometry.Point {
		xspace := panel.Width() - (panel.HorizontalFit() * 2.0)
		endheight := panel.RailHeightFromMountingHole() + panel.MountingHoleBottomY()
		yspace := panel.Height() - endheight*2.0
		xoffset := panel.HorizontalFit() * 2.0
		yoffset := endheight
		return geometry.Point{
			X: xoffset + rand.Float64()*xspace,
			Y: yoffset + rand.Float64()*yspace,
		}
	}
	for i := 0; i < n; i++ {
		lines = append(lines, features.NewLine(rxy(), rxy(), 0.1*float64((rand.Intn(3)))))
	}
	return lines
}

// pcb shops get confused if you don't include a copper layer
func copperPour(pnl panel.Panel) gerber.Primitive {
	left := panel.LeftX(pnl)
	right := panel.RightX(pnl)
	top := pnl.MountingHoleTopY() - pnl.RailHeightFromMountingHole()
	bottom := pnl.MountingHoleBottomY() + pnl.RailHeightFromMountingHole()
	return gerber.Polygon(
		gerber.Point(0, 0), // offset? what even is this?
		true,               // filled
		[]gerber.Pt{
			gerber.Point(left, top),
			gerber.Point(right, top),
			gerber.Point(right, bottom),
			gerber.Point(left, bottom),
			gerber.Point(left, top),
		},
		0.1,
	)
}

func main() {
	cfg, pnl, err := configure()
	if err != nil {
		log.Fatalf("configure: %v", err)
	}
	g := gerber.New(cfg.name)
	// we collect primitives and Add them all at once like this because the
	// gerber lib seems to reset the relevant layer on each Add
	prims := newprimitives()
	collectPrimitives(panelsource.GeneratePanelOutlineFeatures(pnl), prims)
	collectPrimitives(panelHeaderFooter(pnl, cfg.header, cfg.footer), prims)
	collectPrimitives(randomLines(pnl, 100), prims)
	g.Outline().Add(prims.outlines...)
	g.TopSilkscreen().Add(prims.silkscreens...)
	g.Drill().Add(prims.drills...)
	g.TopCopper().Add(copperPour(pnl))
	if err := g.WriteGerber(); err != nil {
		log.Fatalf("WriteGerber: %v", err)
	}
}
