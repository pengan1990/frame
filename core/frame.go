package core

import (
	"context"
	"sync"

	"github.com/panjf2000/ants/v2"
)

const (
	buffSize = 128
)

type Frame struct {
	poolSize int
	reader   ReadPair
	writers  []WritePair
	chk      CheckPointer
	pool     *ants.Pool
	buffCh   chan interface{}
	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

type ReadPair interface {
	Reader
	BackPressure
}

type WritePair interface {
	Writer
	BackPressure
}

func NewFrame(reader ReadPair, ws []WritePair, chk CheckPointer, size int, ctx context.Context) (*Frame, error) {
	pool, err := ants.NewPool(size)
	if err != nil {
		return nil, err
	}
	subCtx, subCancel := context.WithCancel(ctx)

	ch := make(chan interface{}, buffSize)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	return &Frame{
		poolSize: size,
		reader:   reader,
		writers:  ws,
		chk:      chk,
		pool:     pool,
		buffCh:   ch,
		wg:       wg,
		ctx:      subCtx,
		cancel:   subCancel,
	}, nil
}

func (f *Frame) read() error {
	if f.reader.IsDone() {
		return nil
	}
	data, err := f.reader.Execute()
	if err != nil {
		return err
	}

	f.reader.Next(f.ctx)
	select {
	case f.buffCh <- data:
	case <-f.ctx.Done():
		return context.Canceled
	}
	return nil
}

func (f *Frame) write() {
	defer f.wg.Done()
	for {
		select {
		case <-f.ctx.Done():
			// TODO context done
			return
		case data, more := <-f.buffCh:
			if !more {
				// TODO warn for channel closed
				return
			}

			for _, w := range f.writers {
				if err := f.pool.Submit(func() {
					// backup pressure
					w.Next(f.ctx)

					// execute
					if err := w.Execute(data); err != nil {
						// TODO save error writer
					}
				}); err != nil {
					return
				}
			}

			// save point
			if err := f.chk.Save(data); err != nil {
				// TODO log for checkpoint error
				return
			}
		}
	}
}

func (f *Frame) Stop() {
	f.cancel()
}

func (f *Frame) Execute() error {
	go f.write()

	if err := f.read(); err != nil {
		return err
	}

	// close write channel buffer
	close(f.buffCh)

	// wait writer done means try best to write as much as possible
	f.wg.Wait()

	return nil
}
