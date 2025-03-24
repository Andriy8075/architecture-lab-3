package painter

import (
	"image"
	"image/color"
	"sync"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

type mockReceiver struct {
	lastTexture screen.Texture
	mu          sync.Mutex
}

func (m *mockReceiver) Update(t screen.Texture) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastTexture = t
}

type mockScreen struct{}

func (m *mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) { return nil, nil }
func (m *mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}
func (m *mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) { return nil, nil }

type mockTexture struct {
	screen.Texture
}

func (m *mockTexture) Release() {}

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		mr mockReceiver
	)

	l.Receiver = &mr

	scr := &mockScreen{}
	l.Start(scr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.Post(OperationFunc(func(t screen.Texture, s *TextureState) {
			s.Background = color.White
			s.UpdateNeeded = true
		}))
		l.Post(UpdateOp)
	}()

	wg.Wait()

	l.StopAndWait()

	if mr.lastTexture == nil {
		t.Error("Texture was not updated")
	}
}
