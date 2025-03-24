package main

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"net/http"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
	"github.com/roman-mazur/architecture-lab-3/ui"
)

func main() {
	var (
		pv ui.Visualizer

		state  painter.State
		opLoop painter.Loop
		parser lang.Parser
	)

	pv.Title = "Simple painter"
	pv.OnScreenReady = func(s screen.Screen) {
		state.Texture, _ = s.NewTexture(image.Pt(800, 800))
		opLoop.Start(s)
	}

	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
