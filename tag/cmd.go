package tag

import (
	"bytes"
	"image"
	"io"

	"github.com/as/text"
	"github.com/as/ui/win"
	"golang.org/x/mobile/event/mouse"
)

func Pt(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func Visible(w *win.Win, q0, q1 int64) bool {
	if q0 < w.Origin() {
		return false
	}
	if q1 > w.Origin()+w.Frame.Nchars {
		return false
	}
	return true
}

func Paste(w text.Editor, e mouse.Event) (int64, int64) {
	n, _ := Clip.Read(ClipBuf)
	s := fromUTF16(ClipBuf[:n])
	q0, q1 := w.Dot()
	if q0 != q1 {
		w.Delete(q0, q1)
		q1 = q0
	}
	w.Insert(s, q0)
	w.Select(q0, q0+int64(len(s)))
	return w.Dot()
}

func Rdsel(w text.Editor) string {
	q0, q1 := w.Dot()
	return string(w.Bytes()[q0:q1])
}

func Snarf(w text.Editor, e mouse.Event) {
	n := copy(ClipBuf, toUTF16([]byte(Rdsel(w))))
	io.Copy(Clip, bytes.NewReader(ClipBuf[:n]))
	q0, q1 := w.Dot()
	w.Delete(q0, q1)
	w.Select(q0, q0)
}
