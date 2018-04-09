package img

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"

	"github.com/as/shiny/screen"
	"github.com/as/text"
	"github.com/as/ui"
)

type Node struct {
	Sp, size, pad image.Point
	dirty         bool
}

func (n Node) Size() image.Point {
	return n.Size()
}
func (n Node) Pad() image.Point {
	return n.Sp.Add(n.pad)
}

type Img struct {
	Node
	*ui.Dev
	img image.Image
	b   screen.Buffer
	ScrollBar
	org int64
	text.Editor
}

type ScrollBar struct {
	bar     image.Rectangle
	Scrollr image.Rectangle
}

var testimage, _ = ioutil.ReadFile(`C:\d\suns.png`)

func New(dev *ui.Dev, sp, size, pad image.Point, img image.Image) *Img {
	ed, _ := text.Open(text.NewBuffer())
	b := dev.NewBuffer(size)
	if img == nil {
		img, _, _ = image.Decode(bytes.NewReader(testimage))
	}

	w := &Img{
		img:    img,
		Node:   Node{Sp: sp, size: size, pad: pad},
		Dev:    dev,
		b:      b,
		Editor: ed,
	}
	w.init()
	return w
}

func (w *Img) Mark() { w.dirty = true }

func (w *Img) init() {
	w.Blank()
	w.Fill()
	q0, q1 := w.Dot()
	w.Select(q0, q1)
	w.Mark()
}

func (w *Img) Blank() {
	if w.b == nil {
		return
	}
	r := w.b.RGBA().Bounds()
	draw.Draw(w.b.RGBA(), r, image.Black, image.ZP, draw.Src)
	if w.Sp.Y > 0 {
		r.Min.Y--
	}
	w.Mark()
	//	w.drawsb()
}

func (n *Img) Bounds() image.Rectangle { return image.Rectangle{n.Sp, n.Size()} }
func (w *Img) Buffer() screen.Buffer   { return w.b }
func (w *Img) Bytes() []byte           { return w.Editor.Bytes() }
func (w *Img) Dirty() bool             { return w.dirty }
func (w *Img) Len() int64              { return w.Editor.Len() }
func (w Img) Loc() image.Rectangle     { return image.Rectangle{w.Sp, w.Sp.Add(w.size)} }
func (w *Img) Move(sp image.Point)     { w.Sp = sp }
func (w *Img) Origin() int64           { return w.org }
func (w *Img) Refresh() {
	w.Upload()
	w.Window().Upload(w.Sp, w.b, w.b.Bounds())
	w.dirty = false
}
func (w *Img) Upload() {
	if !w.dirty {
		return
	}
	r := w.img.Bounds()
	b := w.b
	draw.Draw(b.RGBA(), b.RGBA().Bounds(), w.img, w.img.Bounds().Min, draw.Src)
	w.Window().Upload(r.Min, w.b, r)
	w.dirty = false
}
func (w *Img) Resize(size image.Point) {
	if size.Y < 100 {
		size.Y = 100
	}
	b := w.NewBuffer(size)
	w.size = size
	w.b.Release()
	w.b = b
	draw.Draw(b.RGBA(), b.RGBA().Bounds(), w.img, w.img.Bounds().Min, draw.Src)
	w.Refresh()
}
func (w *Img) Fill()                           {}
func (w *Img) Clicksb(pt image.Point, dir int) {}
func (w *Img) Scroll(dl int)                   {}
func (w *Img) SetOrigin(org int64, exact bool) {}
