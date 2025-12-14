package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// DecisionCard is the structured proposal format for the AI.
type DecisionCard struct {
	ProposalID         string
	TargetModule       string
	Rationale          string
	ActionCodeDiff     string
	PredictedEpsilonGain float64
	PredictedIGain     float64
	CalculatedRiskScore  float64
	RiskAdjustedReward   float64
}

// PlannerReasonor generates plans and proposals.
type PlannerReasoner struct {
	goalEngine   *GoalEngine
	memory       *MemoryConsolidator // Link to long-term memory
}

func NewPlannerReasoner(ge *GoalEngine, mem *MemoryConsolidator) *PlannerReasoner {
	return &PlannerReasoner{goalEngine: ge, memory: mem}
}

// generateProposal simulates the generation of a self-improvement proposal.
func (pr *PlannerReasoner) generateProposal(anomaly string) DecisionCard {
	// In a real implementation, this would be a complex reasoning process.
	// Here, we simulate the generation of a proposal based on a detected anomaly.
	targetModule := "HarmonicFoldingEngine"
	action := "Refactor compression layer to use a more efficient algorithm."

	// Check against Avoidance Rules from past failures.
	for _, rule := range pr.memory.GetAvoidanceRules() {
		if strings.Contains(rule, targetModule) {
			// If a rule exists for this module, generate a different proposal.
			targetModule = "MemoryCoreSystem"
			action = "Optimize data indexing for faster retrieval."
		}
	}

	// Simulate the "Simulation Chamber"
	predictedEpsilonGain := rand.Float64() * 0.1 // Predict a gain of 0-10%
	predictedIGain := rand.Float64() * 0.05 // Predict a gain of 0-5%
	riskScore := rand.Float64() * 0.25      // Predict a risk of 0-25%

	rar := pr.goalEngine.CalculateRiskAdjustedReward(predictedEpsilonGain, predictedIGain, riskScore)

	return DecisionCard{
		ProposalID:         uuid.New().String(),
		TargetModule:       targetModule,
		Rationale:          fmt.Sprintf("Detected anomaly: %s. A refactoring is proposed.", anomaly),
		ActionCodeDiff:     action,
		PredictedEpsilonGain: predictedEpsilonGain,
		PredictedIGain:     predictedIGain,
		CalculatedRiskScore:  riskScore,
		RiskAdjustedReward:   rar,
	}
}
