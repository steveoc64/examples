package fractal

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/effect"
)

type fractal struct {
	currIterations          uint
	currScale, currX, currY float64

	window       fyne.Window
	canvas       fyne.CanvasObject
	backingStore draw.Image
	mode         rune
	dirty        bool
}

func (f *fractal) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	rect := image.Rect(0, 0, size.Width, size.Height)
	f.backingStore = image.NewRGBA(rect)
	f.dirty = true
	f.canvas.Resize(size)
}

func (f *fractal) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(320, 240)
}

func (f *fractal) refresh() {
	if f.currScale >= 1.0 {
		f.currIterations = 100
	} else {
		f.currIterations = uint(100 * (1 + math.Pow((math.Log10(1/f.currScale)), 1.25)))
	}

	f.window.Canvas().Refresh(f.canvas)
}

func (f *fractal) scaleChannel(c float64, start, end uint32) uint8 {
	if end >= start {
		return (uint8)(c*float64(uint8(end-start))) + uint8(start)
	}

	return (uint8)((1-c)*float64(uint8(start-end))) + uint8(end)
}

func (f *fractal) scaleColor(c float64, start, end color.Color) color.Color {
	r1, g1, b1, _ := start.RGBA()
	r2, g2, b2, _ := end.RGBA()
	return color.RGBA{f.scaleChannel(c, r1, r2), f.scaleChannel(c, g1, g2), f.scaleChannel(c, b1, b2), 0xff}
}

func (f *fractal) mandelbrot(px, py, w, h int) color.Color {
	drawScale := 3.5 * f.currScale
	aspect := (float64(h) / float64(w))
	cRe := ((float64(px)/float64(w))-0.5)*drawScale + f.currX
	cIm := ((float64(py)/float64(w))-(0.5*aspect))*drawScale - f.currY

	var i uint
	var x, y, xsq, ysq float64

	for i = 0; i < f.currIterations && (xsq+ysq <= 4); i++ {
		xNew := float64(xsq-ysq) + cRe
		y = 2*x*y + cIm
		x = xNew

		xsq = x * x
		ysq = y * y
	}

	if i == f.currIterations {
		return theme.BackgroundColor()
	}

	mu := (float64(i) / float64(f.currIterations))
	c := math.Sin((mu / 2) * math.Pi)

	return f.scaleColor(c, theme.PrimaryColor(), theme.TextColor())
}

func (f *fractal) fractalRune(r rune) {
	switch r {
	case '+':
		f.currScale /= 1.1
		f.dirty = true
	case '-':
		f.currScale *= 1.1
		f.dirty = true
	case '1', '2', '3', '4', '5', '6':
		f.mode = r
	case ' ':
		f.mode = ' '
		f.dirty = true
	}

	f.refresh()
}

func (f *fractal) fractalKey(ev *fyne.KeyEvent) {
	delta := f.currScale * 0.2
	if ev.Name == fyne.KeyUp {
		f.currY -= delta
		f.dirty = true
	} else if ev.Name == fyne.KeyDown {
		f.currY += delta
		f.dirty = true
	} else if ev.Name == fyne.KeyLeft {
		f.currX += delta
		f.dirty = true
	} else if ev.Name == fyne.KeyRight {
		f.currX -= delta
		f.dirty = true
	}
	if f.dirty {
		f.refresh()

	}
}

func (f *fractal) render(w, h int) image.Image {
	if f.dirty {
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				f.backingStore.Set(x, y, f.mandelbrot(x, y, w, h))
			}
		}
	}
	f.dirty = false
	switch f.mode {
	case '1':
		f.backingStore = adjust.Hue(f.backingStore, 8)
	case '2':
		//f.backingStore = blur.Box(f.backingStore, 3.0)
		f.backingStore = effect.Sharpen(f.backingStore)
	case '3':
		f.backingStore = effect.Dilate(f.backingStore, 3)
	case '4':
		f.backingStore = effect.Emboss(f.backingStore)
	case '5':
		f.backingStore = effect.Erode(f.backingStore, 3)
	case '6':
		f.backingStore = effect.Sobel(f.backingStore)
	}
	return f.backingStore
}

// Show loads a Mandelbrot fractal example window for the specified app context
func Show(app fyne.App) {
	window := app.NewWindow("Fractal")
	window.SetPadded(false)
	fractal := &fractal{window: window}
	//fractal.canvas = canvas.NewRasterWithPixels(fractal.mandelbrot)
	fractal.canvas = canvas.NewRaster(fractal.render)

	fractal.currIterations = 100
	fractal.currScale = 1.0
	fractal.currX = -0.75
	fractal.currY = 0.0

	window.SetContent(fyne.NewContainerWithLayout(fractal, fractal.canvas))
	window.Canvas().SetOnTypedRune(fractal.fractalRune)
	window.Canvas().SetOnTypedKey(fractal.fractalKey)
	window.Show()
}
