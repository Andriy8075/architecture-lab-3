package painter

import (
	"golang.org/x/exp/shiny/screen"
	"image"
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
	done    chan struct{}
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})
	l.done = make(chan struct{}) // Ініціалізація каналу done

	l.mq = messageQueue{ops: make(chan Operation, 100)}

	go l.eventLoop()
	//go func() {
	//	for !l.stopReq || !l.mq.empty() {
	//		op := l.mq.pull()
	//		if op == nil {
	//			continue
	//		}
	//		update := op.Do(l.next)
	//		if update {
	//			l.Receiver.Update(l.next)
	//			l.next, l.prev = l.prev, l.next
	//		}
	//	}
	//	close(l.stop)
	//}()
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
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	//l.Post(OperationFunc(func(screen.Texture) {
	//	l.stopReq = true
	//}))
	//<-l.stop

	l.stopReq = true
	l.stop <- struct{}{}
	<-l.done
}

type messageQueue struct {
	//Ops []Operation
	//mu      sync.Mutex
	//blocked chan struct{}

	ops chan Operation
}

func (mq *messageQueue) push(op Operation) {
	//mq.mu.Lock()
	//defer mq.mu.Unlock()
	//mq.Ops = append(mq.Ops, op)
	//if mq.blocked != nil {
	//	close(mq.blocked)
	//	mq.blocked = nil
	//}

	mq.ops <- op
}

func (mq *messageQueue) pull() Operation {
	//mq.mu.Lock()
	//defer mq.mu.Unlock()
	//for len(mq.Ops) == 0 {
	//	mq.blocked = make(chan struct{})
	//	mq.mu.Unlock()
	//	<-mq.blocked
	//	mq.mu.Lock()
	//}
	//op := mq.Ops[0]
	//mq.Ops[0] = nil
	//mq.Ops = mq.Ops[1:]
	//return op

	return <-mq.ops
}

func (mq *messageQueue) empty() bool {
	//mq.mu.Lock()
	//defer mq.mu.Unlock()
	//return len(mq.Ops) == 0

	return len(mq.ops) == 0
}
