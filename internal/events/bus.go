package events

import (
	"context"
	"sync"
	"time"
)

type Bus struct {
	ch      chan Event
	wg      sync.WaitGroup
	closed  bool
	mu      sync.Mutex
	bufSize int
}

func NewBus(buffer int) *Bus {
	if buffer <= 0 {
		buffer = 128
	}
	return &Bus{
		ch:      make(chan Event, buffer),
		bufSize: buffer,
	}
}

// Publish non-blocking jika ada buffer; blocking jika penuh.
// Hormati context supaya tidak deadlock saat shutdown.
func (b *Bus) Publish(ctx context.Context, e Event) error {
	b.mu.Lock()
	closed := b.closed
	b.mu.Unlock()
	if closed {
		return nil // ignore publishes after close
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.ch <- e:
		return nil
	}
}

// Subscribe mengembalikan channel hanya-baca untuk workers.
func (b *Bus) Subscribe() <-chan Event {
	return b.ch
}

// Close menutup bus setelah memastikan worker selesai (di sisi worker pakai wg).
func (b *Bus) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true
	close(b.ch)
	b.mu.Unlock()

	// Tunggu consumer yang terdaftar via Add/Done selesai.
	b.wg.Wait()
}

// AddConsumer dipanggil worker saat start; Done saat selesai.
func (b *Bus) AddConsumer()  { b.wg.Add(1) }
func (b *Bus) DoneConsumer() { b.wg.Done() }

// Helper buat bikin Event cepat
func NewTodoCreated(id int64, title string) Event {
	return Event{
		Type:      EventTodoCreated,
		Timestamp: time.Now().UTC(),
		Payload:   TodoCreatedPayload{ID: id, Title: title},
	}
}
