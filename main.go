package main

import (
	"fmt"
	"math"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mazznoer/csscolorparser"
)

type Xors struct {
	x int32
	y int32
	z int32
	w int32
}

type S struct {
	left  int
	right int
}

type Coordinate struct {
	x float64
	y float64
}

func (xors *Xors) randXors() float64 {
	t := xors.x ^ (xors.x << 11)
	xors.x = xors.y
	xors.y = xors.z
	xors.z = xors.w
	w2 := uint32(xors.w) >> 19
	t2 := uint32(t) >> 8
	xors.w = (xors.w ^ int32(w2)) ^ (t ^ int32(t2))
	return (float64(xors.w) / 4294967296.0) + 0.5
}

func (xors *Xors) fiXors() bool {
	if xors.randXors() < 0.5 {
		return true
	}
	return false
}

const (
	TypeSin  = 1
	TypeCos  = 2
	TypeNone = 0
)

func calc(_type int, num float64) float64 {
	if _type == TypeCos {
		return math.Cos(num)
	} else if _type == TypeSin {
		return math.Sin(num)
	} else {
		return 0.0
	}
}

func rg(ha *[]int, xors *Xors) float64 {
	ra := *ha
	c := xors.randXors()
	a := int32(float64(len(ra))*c) | 0

	b := ra[a]
	ra[a] = ra[len(ra)-1]
	*ha = ra[:len(ra)-1]
	return float64(b)
}

func visual_hash(text string, ctx *cairo.Context) {
	c := [4]int32{123456789, 362436069, 521288629, 0}
	for pos, char := range text {
		left := int32(char)
		right := int32((pos * 11) % 16)
		c[(pos+3)%4] ^= int32(left) << int32(right)
	}

	xors := &Xors{
		x: c[0],
		y: c[1],
		z: c[2],
		w: c[3],
	}

	var nn int
	for i := 0; i < 52; i++ {
		xors.randXors()
	}

	if xors.fiXors() {
		nn = 7
	} else {
		nn = 11
	}
	size := 4620
	ctx.SetOperator(cairo.OPERATOR_SOURCE)
	ctx.SetSourceRGBA(26, 32, 44, 0)
	ctx.Rectangle(0.0, 0.0, 400.0, 400.0)
	ctx.Fill()
	ctx.SetOperator(cairo.OPERATOR_LIGHTEN)

	h := [8]float64{0.0}
	h[2] = 0.3 + xors.randXors()*0.2
	h[3] = 0.1 + xors.randXors()*0.1
	h[5] = 1.0 + xors.randXors()*4.0
	h[6] = 1.0 + xors.randXors()
	h[7] = 1.0 + xors.randXors()
	h[0] = 0.4 + xors.randXors()*0.2
	for a := 2; a < 8; a++ {
		c := xors.fiXors()
		if c {
			h[a] *= -1.0
		}
	}
	ki := []int{1, 3, 5, 7, 9, 11}
	gu := []int{0, 0, 2, 4, 6, 8, 10}
	s := [8]S{{TypeNone, TypeNone}}
	q := [8]float64{0.0}
	var pr float64 = float64(int32(1.0+xors.randXors()*float64((nn-1)))|0) / float64(nn)
	for a := range [2]int{} {
		if xors.fiXors() {
			s[a] = S{TypeCos, TypeSin}
			q[a] = rg(&ki, xors) - pr
		} else {
			s[a] = S{TypeSin, TypeCos}
			q[a] = rg(&gu, xors) + pr
		}
	}

	for a := 2; a < 8; a++ {
		b := xors.fiXors()
		if len(ki) == 0 {
			b = false
		}
		if len(gu) == 0 {
			b = true
		}
		if b {
			q[a] = rg(&ki, xors)
		} else {
			q[a] = rg(&gu, xors)
		}
		if xors.fiXors() {
			q[a] *= -1.0
		}
		if a > 5 {
			b = !b
		}
		if b {
			s[a] = S{TypeCos, TypeNone}
		} else {
			s[a] = S{TypeSin, TypeNone}
		}

	}

	n := []float64{0.0, 0.0, 0.0}
	p := []Coordinate{}
	for a := range [3]int{} {
		if xors.fiXors() {
			n[a] = 1.0
		} else {
			n[a] = -1.0
		}
	}
	step := math.Pi * 2.0 / float64(size) * float64(nn)

	r := 0.0
	f := 0

	for f < size {
		c1 := calc(s[3].left, r*q[3])
		bf := calc(s[6].left, r*q[6]+c1*h[5]) * n[0]
		af := 1.0 + bf*h[0]
		df := calc(s[7].left, r*q[7])
		ef := -1.0 * df
		df *= (2.0 - af) * n[1]
		ef *= (2.0 - af) * n[2]
		c2 := calc(s[5].left, r*q[5])
		cf := calc(s[4].left, r*q[4]+c2*h[7]) / 4.0 * h[6] * (af - (1.0 - h[0]))
		xf := math.Sin(r*pr+cf)*af +
			calc(s[0].left, r*q[0])*h[2]*df +
			calc(s[1].left, r*q[1])*h[3]*ef
		yf := math.Cos(r*pr+cf)*af +
			calc(s[0].right, r*q[0])*h[2]*df +
			calc(s[1].right, r*q[1])*h[3]*ef
		p = append(p, Coordinate{xf*110.0 + 200.0, yf*110.0 + 200.0})
		r += step
		f += 1
	}
	ctx.NewPath()
	var hx int = 0
	for dx := 0; dx < 3; dx++ {
		gh := int(xors.randXors()*360.0) | 0
		hx += 1 + int(xors.randXors()*3.0) | 0
		ih := 50 + int(xors.randXors()*20.0) | 0
		for ah := 0; ah < len(p); ah++ {
			ctx.NewPath()
			ei := []Coordinate{}
			for bi := range [3]int{} {
				index := uint32((ah + bi*((dx+1)*int(hx)))) % uint32(len(p))
				ci := p[index]
				ei = append(ei, ci)
				ctx.LineTo(ci.x, ci.y)
			}
			fa := ei[0].x * (ei[1].y - ei[2].y)
			fa += ei[1].x * (ei[2].y - ei[0].y)
			fa += ei[2].x * (ei[0].y - ei[1].y)
			if fa > 45 && fa < 8000 {
				rgba, _ := csscolorparser.Parse(fmt.Sprint("hsla(", gh, ", ", ih, "%, 40%, ", 55.0/fa, ")"))
				ctx.SetSourceRGBA(rgba.R, rgba.G, rgba.B, rgba.A)
				ctx.Fill()
			}
		}
	}
}

func main() {
	gtk.Init(nil)
	var surface *cairo.Surface = nil
	var fileName string = "visual_hash.png"
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("GOTK3 VISUAL HASH")
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
	vbox.SetMarginStart(16)
	vbox.SetMarginEnd(16)
	vbox.SetMarginTop(16)
	vbox.SetMarginBottom(16)
	textName, _ := gtk.EntryNew()
	btnGenerate, _ := gtk.ButtonNewWithLabel("GENERATE")
	btnSave, _ := gtk.ButtonNewWithLabel("SAVE IMAGE")
	canvas, _ := gtk.DrawingAreaNew()
	canvas.SetSizeRequest(400, 400)

	vbox.Add(textName)
	vbox.Add(btnGenerate)
	vbox.Add(canvas)
	vbox.Add(btnSave)

	win.Add(vbox)
	win.ShowAll()

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	btnGenerate.Connect("clicked", func(btn *gtk.Button) {
		canvas.QueueDrawArea(0, 0, 400, 400)
	})

	btnSave.Connect("clicked", func(btn *gtk.Button) {
		if surface == nil {
			return
		}
		x := canvas.GetAllocation().GetX()
		y := canvas.GetAllocation().GetY()

		pixbuf, err := gdk.PixbufGetFromSurface(surface, x, y, 400, 400)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		dialog, _ := gtk.FileChooserNativeDialogNew("Save File", win, gtk.FILE_CHOOSER_ACTION_SAVE, "Save", "Cancel")
		filter, _ := gtk.FileFilterNew()
		filter.AddMimeType("image/png")
		dialog.AddFilter(filter)
		dialog.SetFilename(fileName)
		response := dialog.Run()
		if response == int(gtk.RESPONSE_ACCEPT) || response == int(gtk.RESPONSE_YES) {
			pixbuf.SavePNG(dialog.GetFilename(), 0)
		}
	})

	canvas.Connect("draw", func(widget *gtk.DrawingArea, cr *cairo.Context) {
		name, err := textName.GetText()
		if err == nil {
			visual_hash(name, cr)
			if cr != nil {
				surface = cr.GetTarget()
				fileName = name + ".png"
			}
		}
	})

	gtk.Main()
}
