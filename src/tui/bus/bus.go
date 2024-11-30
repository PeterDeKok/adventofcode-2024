package bus

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

type MsgBus struct {
	ch  chan tea.Msg
	ctx context.Context
}

func CreateMsgBus(ctx context.Context) *MsgBus {
	ch := make(chan tea.Msg)

	context.AfterFunc(ctx, func() {
		close(ch)
	})

	b := &MsgBus{
		ch:  ch,
		ctx: ctx,
	}

	return b
}

func (b *MsgBus) Relay(p *tea.Program) {
	for {
		select {
		case msg, ok := <-b.ch:
			if !ok {
				return
			}
			p.Send(msg)
		}
	}
}

func (b *MsgBus) Send(msg tea.Msg) {
	go func() {
		select {
		case <-b.ctx.Done():
		case b.ch <- msg:
		}
	}()
}
