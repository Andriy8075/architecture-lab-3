package ui

import (
	"golang.org/x/exp/shiny/imageutil"
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz size.Event
	T  image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.T.X = 400
	pw.T.Y = 400
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)
		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true
		}
	}
	return false
}

//func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
//	switch e := e.(type) {
//	case size.Event:
//		pw.sz = e
//	case mouse.Event:
//		if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
//			pw.drawTAt(float64(e.X), float64(e.Y)) // Конвертація float32 до float64
//		}
//	case error:
//		log.Printf("ERROR: %s", e)
//	case paint.Event:
//		if t != nil {
//			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
//		} else {
//			pw.drawDefaultUI()
//		}
//		pw.w.Publish()
//	}
//}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
				pw.T = image.Point{
					X: int(e.X),
					Y: int(e.Y),
				}
				pw.w.Send(paint.Event{})
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawTAt(x, y float64) {
	size := 200
	thickness := 50

	// Horizontal part
	rect1 := image.Rect(
		int(x)-size/2, int(y)-thickness/2,
		int(x)+size/2, int(y)+thickness/2,
	)
	pw.w.Fill(rect1, color.RGBA{B: 0xff, A: 0xff}, draw.Src)

	// Vertical part
	rect2 := image.Rect(
		int(x)-thickness/2, int(y)-thickness/2,
		int(x)+thickness/2, int(y)+size/2,
	)
	pw.w.Fill(rect2, color.RGBA{B: 0xff, A: 0xff}, draw.Src)

	pw.w.Publish()
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src) // Фон.

	pw.drawTAt(float64(pw.T.X), float64(pw.T.Y))

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}
