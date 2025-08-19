package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bryanagamk/todo-app-postgresql/internal/events"
)

type LoggerPool struct {
	n        int
	filePath string

	once sync.Once
	fh   *os.File
	mu   sync.Mutex // serialize writes to file
}

func NewLoggerPool(n int, filePath string) *LoggerPool {
	if n <= 0 {
		n = 2
	}
	return &LoggerPool{n: n, filePath: filePath}
}

func (p *LoggerPool) openFile() error {
	var err error
	p.once.Do(func() {
		p.fh, err = os.OpenFile(p.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	})
	return err
}

// Run: start N workers yang membaca dari bus.Subscribe().
// Setiap event akan ditulis JSON satu baris ke file.
func (p *LoggerPool) Run(ctx context.Context, bus *events.Bus) {
	for i := 0; i < p.n; i++ {
		bus.AddConsumer()
		go func(workerID int) {
			defer bus.DoneConsumer()

			_ = p.openFile()
			if p.fh == nil {
				return
			}

			for {
				select {
				case <-ctx.Done():
					return
				case ev, ok := <-bus.Subscribe():
					if !ok {
						return // bus closed
					}
					// Simulasi kerja: encode ke JSON dan tulis
					row := map[string]interface{}{
						"ts":     time.Now().Format(time.RFC3339Nano),
						"type":   ev.Type,
						"event":  ev, // akan berisi payload
						"worker": workerID,
					}
					b, err := json.Marshal(row)
					if err != nil {
						continue
					}
					// Pastikan write ke file tidak race
					p.mu.Lock()
					fmt.Fprintln(p.fh, string(b))
					p.mu.Unlock()
				}
			}
		}(i + 1)
	}
}

func (p *LoggerPool) Close() error {
	if p.fh != nil {
		return p.fh.Close()
	}
	return nil
}
