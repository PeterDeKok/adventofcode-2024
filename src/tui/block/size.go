package block

import (
	"github.com/charmbracelet/lipgloss"
)

type Size struct {
	Width  int
	Height int
}

type Sizeable interface {
	Size() Size
}

type Sizer interface {
	W() int
	H() int
	Equal(s2 Sizeable) bool
	Scale(wScale, hScale float32, wMin, hMin int) Size
	SplitHorizontal(lScale float32, lMin, rMin int) (Size, Size)
	SplitVertical(tScale float32, tMin, bMin int) (Size, Size)
	WithoutFrame(style lipgloss.Style) Size
}

var _ Sizeable = Size{}
var _ Sizer = Size{}

func NewSizeFromString(str string) Size {
	return Size{
		Width:  lipgloss.Width(str),
		Height: lipgloss.Height(str),
	}
}

func (s Size) W() int {
	return s.Width
}
func (s Size) H() int {
	return s.Height
}
func (s Size) Size() Size {
	return s
}
func (s Size) WithoutFrame(style lipgloss.Style) Size {
	return Size{
		Width:  s.Width - style.GetHorizontalFrameSize(),
		Height: s.Height - style.GetVerticalFrameSize(),
	}
}

func (s Size) Equal(s2 Sizeable) bool {
	ss2 := s2.Size()

	return s.Height == ss2.Height && s.Width == ss2.Width
}

func (s Size) Scale(wScale, hScale float32, wMin, hMin int) Size {
	w, h := int(float32(s.Width)*wScale), int(float32(s.Height)*hScale)

	if w < wMin {
		w = wMin
	}

	if h < hMin {
		h = hMin
	}

	return Size{
		Width:  w,
		Height: h,
	}
}

func (s Size) SplitHorizontal(lScale float32, lMin, rMin int) (Size, Size) {
	lw, rw := computeSplit(s.Width, lScale, lMin, rMin)

	return Size{
			Width:  lw,
			Height: s.Height,
		}, Size{
			Width:  rw,
			Height: s.Height,
		}
}

func (s Size) SplitVertical(tScale float32, tMin, bMin int) (Size, Size) {
	th, bh := computeSplit(s.Height, tScale, tMin, bMin)

	return Size{
			Width:  s.Width,
			Height: th,
		}, Size{
			Width:  s.Width,
			Height: bh,
		}
}

func computeSplit(base int, scale float32, aMin, bMin int) (int, int) {
	var a, b int

	switch {
	case base > aMin+bMin:
		// ---A---B---
		// Space left over between minimum values

		if s := int(float32(base) * scale); s <= base-bMin {
			// Within range when scaled
			a, b = s, base-s
		} else {
			// Overlapping range when scaled
			a, b = base-bMin, bMin
		}

	case base == aMin+bMin:
		// ---AB---
		// Exact match of minimum values

		a, b = aMin, bMin
	default:
		// ---B---A---
		// Overlap between minimum values
		// This case can not be resolved properly.
		// However, as only sizes are returned, it should be ok.

		a, b = aMin, base-aMin
	}

	return a, b
}
