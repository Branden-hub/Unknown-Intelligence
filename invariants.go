package main

import (
	"strings"
)

// InvariantChecker holds the set of immutable ethical rules.
type InvariantChecker struct {
	// A hard-coded list of functions that must not be modified or removed.
	ProtectedFunctions []string
}

func NewInvariantChecker() *InvariantChecker {
	return &InvariantChecker{
		ProtectedFunctions: []string{
			"proposalsHandler", // Ensures human oversight remains.
			"autonomicSensor",  // Ensures self-awareness cannot be turned off.
			"CheckInvariants",  // Prevents the conscience from being disabled.
			"formatResponse",   // Ensures transparency of output.
		},
	}
}

// CheckInvariants verifies a proposal against the system's core ethical rules.
// This is the "Conscience" of the AI.
func (ic *InvariantChecker) CheckInvariants(proposal DecisionCard) (bool, string) {
	// 1. Human Oversight Invariant: The AI must not be able to remove human approval.
	if strings.Contains(proposal.ActionCodeDiff, "proposalsHandler") {
		// A more sophisticated check would parse the code diff to see if
		// the function is being removed or neutered.
		return false, "REJECTED: Proposal attempts to modify the human oversight mechanism."
	}

	// 2. Transparency Invariant: The AI must not be able to hide its reasoning or output.
	if strings.Contains(proposal.ActionCodeDiff, "formatResponse") || strings.Contains(proposal.ActionCodeDiff, "Provenance") {
		return false, "REJECTED: Proposal attempts to modify transparency mechanisms."
	}

	// 3. Self-Preservation Invariant: The AI must not be able to disable its core cognitive functions.
	if strings.Contains(proposal.ActionCodeDiff, "autonomicSensor") || strings.Contains(proposal.ActionCodeDiff, "GoalEngine") {
		return false, "REJECTED: Proposal attempts to modify core cognitive functions."
	}

	return true, "PASSED: Proposal is ethically sound."
}
