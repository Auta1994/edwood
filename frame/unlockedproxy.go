package frame

import (
	"image"

	"log"
)

// SelectScrollUpdater is a struct subset of frame.Frame that contains only those
// methods needed to implement a frame.Select callback.
type  SelectScrollUpdater interface {
	GetFrameFillStatus() FrameFillStatus
	Charofpt(pt image.Point) int
	DefaultFontHeight() int
	Delete(int,  int) int 
	Insert( []rune,  int) bool
	IsLastLineFull() bool
	Rect() image.Rectangle
}

type selectscrollupdaterimpl  frameimpl

func (up *selectscrollupdaterimpl) GetFrameFillStatus() FrameFillStatus {
	log.Println("selectscrollupdaterimpl.GetFrameFillStatus")
	f := (*frameimpl)(up)
	return FrameFillStatus{
			Nchars:   f.nchars,
			Nlines:   f.nlines,
			Maxlines: f.maxlines,
		}
}

func (up *selectscrollupdaterimpl) Charofpt(pt image.Point) int {
	log.Println("selectscrollupdaterimpl.Charofpt")
	f := (*frameimpl)(up)
	return 	f.charofptimpl(pt)
}

func (up *selectscrollupdaterimpl) DefaultFontHeight() int {
	log.Println("selectscrollupdaterimpl.DefaultFontHeight")
	f := (*frameimpl)(up)
	return f.defaultfontheight
}

func (up *selectscrollupdaterimpl) Delete(p0, p1 int) int {
	log.Println("selectscrollupdaterimpl.Delete")
	f := (*frameimpl)(up)
	return f.deleteimpl(p0, p1)
}


func (up *selectscrollupdaterimpl) Insert(r []rune, p0 int) bool {
	log.Println("selectscrollupdaterimpl.Insert")
	f := (*frameimpl)(up)
	return f.insertimpl(r, p0)
}

func (up *selectscrollupdaterimpl) IsLastLineFull() bool {
	log.Println("selectscrollupdaterimpl.IsLastLineFull")
	f := (*frameimpl)(up)
	return f.lastlinefull
}


func (up *selectscrollupdaterimpl) Rect() image.Rectangle {
	log.Println("selectscrollupdaterimpl.Rect")
	f := (*frameimpl)(up)
	return f.rect
}