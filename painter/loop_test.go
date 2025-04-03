package painter

import (
	"errors"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw" // Add this import
	"image"
	"image/color"
	"sync"
	"testing"
	"time"
)

type mockReceiver struct {
	updateCount int
	lastTexture screen.Texture
}

func (m *mockReceiver) Update(t screen.Texture) {
	m.updateCount++
	m.lastTexture = t
}

type mockScreen struct{}

func (m *mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, errors.New("not implemented")
}

func (m *mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}

func (m *mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, errors.New("not implemented")
}

type mockTexture struct {
	screen.Texture
	fillCount int
}

func (m *mockTexture) Release() {}
func (m *mockTexture) Size() image.Point {
	return image.Pt(800, 800)
}
func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rect(0, 0, 800, 800)
}
func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.fillCount++
}

func TestLoop_Start(t *testing.T) {
	var (
		l  Loop
		mr mockReceiver
		ms mockScreen
	)

	l.Receiver = &mr

	l.Start(&ms)

	if l.next == nil || l.prev == nil {
		t.Error("expected textures to be initialized")
	}

	l.StopAndWait()
}

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		mr mockReceiver
		ms mockScreen
	)

	l.Receiver = &mr
	l.Start(&ms)

	opPosted := false
	testOp := OperationFunc(func(t screen.Texture) {
		opPosted = true
	})

	l.Post(testOp)
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))

	<-l.stop

	if !opPosted {
		t.Error("operation was not executed")
	}
}

func TestLoop_StopAndWait(t *testing.T) {
	var (
		l  Loop
		mr mockReceiver
		ms mockScreen
	)

	l.Receiver = &mr
	l.Start(&ms)

	stopCompleted := false
	go func() {
		l.StopAndWait()
		stopCompleted = true
	}()

	// Post a no-op to trigger the stop
	l.Post(OperationFunc(func(t screen.Texture) {}))
	time.Sleep(100 * time.Millisecond)

	if !stopCompleted {
		t.Error("StopAndWait didn't complete")
	}
}

func TestMessageQueue(t *testing.T) {
	mq := messageQueue{}

	// Test push and pull
	testOp := OperationFunc(func(t screen.Texture) {})
	mq.push(testOp)

	if mq.empty() {
		t.Error("queue should not be empty after push")
	}

	op := mq.pull()
	if op == nil {
		t.Error("expected to get operation from queue")
	}

	if !mq.empty() {
		t.Error("queue should be empty after pull")
	}

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			mq.push(testOp)
			wg.Done()
		}()
	}
	wg.Wait()

	if len(mq.Ops) != 10 {
		t.Errorf("expected 10 operations in queue, got %d", len(mq.Ops))
	}
}

func TestLoop_OperationOrder(t *testing.T) {
	var (
		l  Loop
		mr mockReceiver
		ms mockScreen
	)

	l.Receiver = &mr
	l.Start(&ms)

	var executionOrder []int

	// Post operations that record their execution order
	for i := 0; i < 5; i++ {
		j := i
		l.Post(OperationFunc(func(t screen.Texture) {
			executionOrder = append(executionOrder, j)
		}))
	}

	// Post stop operation
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))

	<-l.stop

	// Verify operations executed in order
	for i := 0; i < 5; i++ {
		if executionOrder[i] != i {
			t.Errorf("operations executed out of order, expected %d got %d", i, executionOrder[i])
		}
	}
}

func TestLoop_UpdatePropagation(t *testing.T) {
	var (
		l  Loop
		mr mockReceiver
		ms mockScreen
	)

	l.Receiver = &mr
	l.Start(&ms)

	// Post an operation that requests update
	l.Post(OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color.White, screen.Src)
	}))
	l.Post(UpdateOp)
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))

	<-l.stop

	if mr.updateCount != 1 {
		t.Errorf("expected 1 update, got %d", mr.updateCount)
	}
}
