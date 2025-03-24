package painter

import (
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type State struct {
	bgColor color.Color
	bgRect  *BgRect
	figures []*TFigure
	Texture screen.Texture // Змінили на експортоване поле (велика літера)
}

func (s *State) UpdateTexture() {
	// Очищаємо текстуру
	s.Texture.Fill(s.Texture.Bounds(), s.bgColor, screen.Src)

	// Малюємо bgRect якщо він є
	if s.bgRect != nil {
		s.bgRect.Do(s.Texture)
	}

	// Малюємо всі фігури
	for _, fig := range s.figures {
		fig.Do(s.Texture)
	}
}

func (s *State) SetBgRect(rect *BgRect) {
	s.bgRect = rect
}

func (s *State) AddFigure(fig *TFigure) {
	s.figures = append(s.figures, fig)
}

func (s *State) MoveFigures(x, y float32) {
	for _, fig := range s.figures {
		fig.X = x
		fig.Y = y
	}
}

func (s *State) Reset() {
	s.bgColor = color.Black
	s.bgRect = nil
	s.figures = nil
}
