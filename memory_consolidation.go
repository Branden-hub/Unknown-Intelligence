package main

import (
	"log"
	"sync"
	"time"
)

// MemoryConsolidator handles the "dreaming" cycle.
type MemoryConsolidator struct {
	// Rules learned from past failures.
	AvoidanceRules []string
	mutex          sync.RWMutex
}

func NewMemoryConsolidator() *MemoryConsolidator {
	return &MemoryConsolidator{
		AvoidanceRules: make([]string, 0),
	}
}

// DreamCycle simulates the process of learning from past failures.
func (mc *MemoryConsolidator) DreamCycle(failedProposals []DecisionCard) {
	go func() {
		log.Println("Dream Cycle Initiated: Consolidating memories...")
		time.Sleep(10 * time.Second) // Simulate time-intensive analysis

		mc.mutex.Lock()
		defer mc.mutex.Unlock()

		for _, proposal := range failedProposals {
			// A real implementation would involve complex analysis to generalize failure modes.
			// Here, we create a simple rule based on the failed proposal's target module.
			rule := "Avoid modifications to " + proposal.TargetModule + " that resulted in low RAR."
			mc.AvoidanceRules = append(mc.AvoidanceRules, rule)
			log.Printf("New Avoidance Rule Learned: %s", rule)
		}

		log.Println("Dream Cycle Complete.")
	}()
}

// GetAvoidanceRules safely returns the current set of avoidance rules.
func (mc *MemoryConsolidator) GetAvoidanceRules() []string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.AvoidanceRules
}
