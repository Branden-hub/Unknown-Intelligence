package main

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// SystemMetabolism represents the "physical" state of the AI.
type SystemMetabolism struct {
	Latency          time.Duration // Time to process a request
	MemorySaturation float64       // Percentage of allocated memory in use
	APICost          float64       // Simulated cost of API calls
}

// HomeostasisMonitor tracks the system's metabolic state.
type HomeostasisMonitor struct {
	metabolism SystemMetabolism
	mutex      sync.RWMutex
}

func NewHomeostasisMonitor() *HomeostasisMonitor {
	return &HomeostasisMonitor{}
}

// Monitor an asynchronous process that continuously updates the system's metabolic state.
func (hm *HomeostasisMonitor) Monitor() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			hm.mutex.Lock()

			// Simulate latency fluctuations
			hm.metabolism.Latency = time.Duration(100+rand.Intn(150)) * time.Millisecond

			// Get actual memory usage
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			// Using HeapAlloc as a proxy for memory saturation. A more complex calculation
			// would consider the total available memory.
			hm.metabolism.MemorySaturation = float64(m.HeapAlloc) / float64(m.Sys) * 100

			// Simulate API cost fluctuations
			hm.metabolism.APICost += rand.Float64() * 0.01

			hm.mutex.Unlock()
		}
	}()
}

// GetMetabolism safely returns the current metabolic state.
func (hm *HomeostasisMonitor) GetMetabolism() SystemMetabolism {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	return hm.metabolism
}
