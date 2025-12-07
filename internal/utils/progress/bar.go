package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/utils/terminal"
)

type ProgressBar struct {
	total       int64
	current     int64
	description string
	startTime   time.Time
	mu          sync.Mutex
	lastRefresh time.Time
}

func NewProgressBar(total int64, description string) *ProgressBar {
	return &ProgressBar{
		total:       total,
		description: description,
		startTime:   time.Now(),
	}
}

func (pb *ProgressBar) Increment(delta int64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	pb.current += delta
	if pb.current > pb.total {
		pb.current = pb.total
	}

	// Rate limit updates to avoid flickering (max 10fps)
	if time.Since(pb.lastRefresh) > 100*time.Millisecond || pb.current == pb.total {
		pb.render()
		pb.lastRefresh = time.Now()
	}
}

func (pb *ProgressBar) Finish() {
	pb.mu.Lock()
	pb.current = pb.total
	pb.render()
	pb.mu.Unlock()
	fmt.Println() // New line after completion
}

func (pb *ProgressBar) render() {
	if !terminal.IsTerminal() {
		return
	}

	percent := 0.0
	if pb.total > 0 {
		percent = float64(pb.current) / float64(pb.total) * 100
	}

	elapsed := time.Since(pb.startTime)
	rate := float64(pb.current) / elapsed.Seconds()
	if rate == 0 {
		rate = 1
	}
	remainingItems := float64(pb.total - pb.current)
	eta := time.Duration(remainingItems/rate) * time.Second

	// Format: Description [=======>    ] 45% (45/100) ETA: 2m30s
	width := 30
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}

	barStr := strings.Repeat("=", filled)
	if filled < width {
		barStr += ">"
		barStr += strings.Repeat(" ", width-filled-1)
	}

	// Use \r to overwrite line
	fmt.Printf("\r%s [%s] %.1f%% (%d/%d) ETA: %s",
		pb.description, barStr, percent, pb.current, pb.total, eta.Round(time.Second))
}
