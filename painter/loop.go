package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	next screen.Texture
	prev screen.Texture

	mq messageQueue

	stop    chan struct{}
	stopReq bool
	done    chan struct{} // Додано поле done
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.stop = make(chan struct{})
	l.done = make(chan struct{}) // Ініціалізація каналу done

	l.mq = messageQueue{ops: make(chan Operation, 100)}

	go l.eventLoop()
}

func (l *Loop) eventLoop() {
	for {
		select {
		case op := <-l.mq.ops:
			if op == nil {
				return
			}
			if update := op.Do(l.next); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		case <-l.stop:
			close(l.mq.ops)
			for op := range l.mq.ops {
				if op != nil {
					op.Do(l.next)
				}
			}
			close(l.done) // Закриваємо канал done при завершенні
			return
		}
	}
}

func (l *Loop) Post(op Operation) {
	if l.stopReq {
		return
	}
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.stopReq = true
	l.stop <- struct{}{}
	<-l.done // Очікуємо закриття каналу done
}

type messageQueue struct {
	ops chan Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.ops <- op
}

func (mq *messageQueue) pull() Operation {
	return <-mq.ops
}

func (mq *messageQueue) empty() bool {
	return len(mq.ops) == 0
}
