package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
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
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
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

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	case size.Event:
		pw.sz = e
	case error:
		log.Printf("ERROR: %s", e)
	case paint.Event:
		pw.drawUI()
	}
}

func (pw *Visualizer) drawUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src) // Фон білий
	pw.drawT()
	pw.w.Publish()
}

func (pw *Visualizer) drawT() {
	size := 200 // Загальний розмір фігури
	thickness := 50
	x, y := 400, 400 // Центр фігури

	// Горизонтальна частина "Т"
	rect1 := image.Rect(x-size/2, y-thickness/2, x+size/2, y+thickness/2)
	pw.w.Fill(rect1, color.RGBA{0, 0, 255, 255}, draw.Src)

	// Вертикальна частина "Т"
	rect2 := image.Rect(x-thickness/2, y-thickness/2, x+thickness/2, y+size/2)
	pw.w.Fill(rect2, color.RGBA{0, 0, 255, 255}, draw.Src)
}
