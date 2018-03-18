package main

import (
	"image"
	"sync"
	"unicode/utf8"
)

type Row struct {
	lk  sync.Mutex
	r   image.Rectangle
	tag Text
	col []*Column
}

func (row *Row) Init(r image.Rectangle) *Row {
	if row == nil {
		row = &Row{}
	}
	display.ScreenImage.Draw(r, display.White, nil, image.ZP)
	row.col = []*Column{}
	row.r = r
	tagfile := NewTagFile()
	r1 := r
	r1.Max.Y = r1.Min.Y + tagfont.Height
	t := &row.tag
	t.Init(tagfile, r1, fontget(0, false, false, ""), tagcolors)
	t.what = Rowtag
	t.row = row
	t.w = nil
	t.col = nil
	r1.Min.Y = r1.Max.Y
	r1.Max.Y += display.ScaleSize(Border)
	display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
	t.Insert(0, []rune("Newcol Kill Putall Dump Exit"), true)
	t.SetSelect(t.file.b.nc(), t.file.b.nc())
	return row
}

func (r *Row) ncol() uint { return uint(len(r.col)) }

func (row *Row) Add(c *Column, x int) *Column {
	r := row.r
	var d *Column

	// Work out the geometry of the column.
	r.Min.Y = row.tag.fr.Rect.Max.Y + display.ScaleSize(Border)
	if x < r.Min.X && row.ncol() > 0 { // Take 40% of last column unless specified
		d = row.col[row.ncol()-1]
		x = d.r.Min.X + 3*d.r.Dx()/5
	}
	/* look for column we'll land on */
	var colidx int
	for colidx = 0; colidx < len(row.col); colidx++ {
		d = row.col[colidx]
		if x < d.r.Max.X {
			break
		}
	}
	if row.ncol() > 0 {
		if uint(colidx) < row.ncol() {
			colidx++ // Place new column after d
		}
		r = d.r
		if r.Dx() < 100 {
			return nil // Refuse columns too narrow
		}
		display.ScreenImage.Draw(r, display.White, nil, image.ZP)
		r1 := r
		r1.Max.X = min(x-display.ScaleSize(Border), r.Max.X-50)
		if r1.Dx() < 50 {
			r1.Max.X = r1.Min.X + 50
		}
		d.Resize(r1)
		r1.Min.X = r1.Max.X
		r1.Max.X = r1.Min.X + display.ScaleSize(Border)
		display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
		r.Min.X = r1.Max.X
	}
	if c == nil {
		c = &Column{}
		c.Init(r)
	} else {
		c.Resize(r)
	}
	c.row = row
	c.tag.row = row
	row.col = append(row.col, c)
	clearmouse()
	return c
}

func (r *Row) Resize(rect image.Rectangle) {
	or := row.r
	deltax := rect.Min.X - or.Min.X
	row.r = rect
	r1 := rect
	r1.Max.Y = r1.Min.Y + tagfont.Height
	row.tag.Resize(r1, true)
	r1.Min.Y = r1.Max.Y
	r1.Max.Y += display.ScaleSize(Border)
	display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
	rect.Min.Y = r1.Max.Y
	r1 = rect
	r1.Max.X = r1.Min.X
	for i := uint(0); i < row.ncol(); i++ {
		c := row.col[i]
		r1.Min.X = r1.Max.X
		/* the test should not be necessary, but guarantee we don't lose a pixel */
		if i == row.ncol()-1 {
			r1.Max.X = rect.Max.X
		} else {
			r1.Max.X = (c.r.Max.X-or.Min.X)*rect.Dx()/or.Dx() + deltax
		}
		if i > 0 {
			r2 := r1
			r2.Max.X = r2.Min.X + display.ScaleSize(Border)
			display.ScreenImage.Draw(r2, display.Black, nil, image.ZP)
			r1.Min.X = r2.Max.X
		}
		c.Resize(r1)
	}
}

func (r *Row) DragCol(c *Column, _ int) {
	Unimpl()
}

func (r *Row) Close(c *Column, dofree bool) {
	Unimpl()
}

func (r *Row) WhichCol(p image.Point) *Column {
	for i := uint(0); i < row.ncol(); i++ {
		c := row.col[i]
		if p.In(c.r) {
			return c
		}
	}
	return nil
}

func (r *Row) Which(p image.Point) *Text {
	if p.In(row.tag.all) {
		return &row.tag
	}
	c := row.WhichCol(p)
	if c != nil {
		return c.Which(p)
	}
	return nil
}

func (row *Row) Type(r rune, p image.Point) *Text {
	var (
		w *Window
		t *Text
	)

	if r == 0 {
		r = utf8.RuneError
	}

	clearmouse()
	row.lk.Lock()
	if bartflag {
		t = barttext
	} else {
		t = row.Which(p)
	}
	if t != nil && !(t.what == Tag && p.In(t.scrollr)) {
		w = t.w
		if w == nil {
			t.Type(r)
		} else {
			w.Lock('K')
			w.Type(t, r)
			/* Expand tag if necessary */
			if t.what == Tag {
				t.w.tagsafe = false
				if r == '\n' {
					t.w.tagexpand = true
				}
				w.Resize(w.r, true, true)
			}
			w.Unlock()
		}
	}
	row.lk.Unlock()
	return t
}

func (r *Row) Clean() int {
	Unimpl()
	return 0
}

func (r *Row) Dump(file string) {
	Unimpl()

}

func (r *Row) LoadFonts(file string) {
	Unimpl()

}

func (r *Row) Load(file string, initing bool) error {
	Unimpl()
	return nil
}

func (r *Row) AllWindows(f func(*Window, interface{}), arg interface{}) {
	for _, c := range r.col {
		for _, w := range c.w {
			f(w, arg)
		}
	}
}

func (r *Row) LookupWin(id int, dump bool) *Window {
	for _, c := range r.col {
		for _, w := range c.w {
			if dump && w.dumpid == id {
				return w
			}
			if !dump && w.id == id {
				return w
			}
		}
	}
	return nil
}
