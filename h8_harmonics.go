package main

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"sync"
	"time"
)

// ============================================================================
// H8 Harmonics System
// ============================================================================

// H8Harmonic represents a harmonic oscillator in the H8 space
// H8 is the 8-dimensional hyperbolic space that complements E8
type H8Harmonic struct {
	Frequency    float64   // Fundamental frequency
	Amplitude    float64   // Current amplitude
	Phase        float64   // Phase angle
	HarmonicMode int       // Which harmonic mode (1-8)
	Damping      float64   // Damping coefficient
}

// H8Resonator manages the harmonic resonance system
type H8Resonator struct {
	Harmonics    [8]*H8Harmonic
	ResonanceLevel float64
	CoherenceState bool
	mutex        sync.RWMutex
}

// NewH8Resonator initializes the H8 harmonic system
func NewH8Resonator() *H8Resonator {
	h8 := &H8Resonator{
		ResonanceLevel: 0.0,
		CoherenceState: false,
	}
	
	// Initialize 8 harmonic modes
	for i := 0; i < 8; i++ {
		h8.Harmonics[i] = &H8Harmonic{
			Frequency:    float64(i+1) * 0.5, // Base frequencies
			Amplitude:    1.0,
			Phase:        0.0,
			HarmonicMode: i + 1,
			Damping:      0.01,
		}
	}
	
	return h8
}

// Oscillate advances all harmonics by one time step
func (h8 *H8Resonator) Oscillate(dt float64) {
	h8.mutex.Lock()
	defer h8.mutex.Unlock()
	
	totalEnergy := 0.0
	
	for i := range h8.Harmonics {
		h := h8.Harmonics[i]
		
		// Update phase: φ = φ + ω*t
		h.Phase += h.Frequency * dt
		
		// Apply damping
		h.Amplitude *= math.Exp(-h.Damping * dt)
		
		// Calculate instantaneous energy: E = A² * ω²
		energy := h.Amplitude * h.Amplitude * h.Frequency * h.Frequency
		totalEnergy += energy
	}
	
	// Normalize and calculate overall resonance level
	h8.ResonanceLevel = math.Sqrt(totalEnergy / 8.0)
	
	// Check for coherence state (all harmonics in phase alignment)
	h8.CoherenceState = h8.checkCoherence()
}

// checkCoherence determines if harmonics are in coherent state
func (h8 *H8Resonator) checkCoherence() bool {
	// Calculate phase variance
	meanPhase := 0.0
	for _, h := range h8.Harmonics {
		meanPhase += h.Phase
	}
	meanPhase /= 8.0
	
	variance := 0.0
	for _, h := range h8.Harmonics {
		diff := h.Phase - meanPhase
		variance += diff * diff
	}
	variance /= 8.0
	
	// Coherent if variance is below threshold
	return variance < 0.1
}

// Excite adds energy to a specific harmonic mode
func (h8 *H8Resonator) Excite(mode int, energy float64) {
	h8.mutex.Lock()
	defer h8.mutex.Unlock()
	
	if mode >= 1 && mode <= 8 {
		h8.Harmonics[mode-1].Amplitude += energy
		if h8.Harmonics[mode-1].Amplitude > 10.0 {
			h8.Harmonics[mode-1].Amplitude = 10.0
		}
	}
}

// GetHarmonicSignature returns the current harmonic signature as a complex vector
func (h8 *H8Resonator) GetHarmonicSignature() []complex128 {
	h8.mutex.RLock()
	defer h8.mutex.RUnlock()
	
	signature := make([]complex128, 8)
	for i, h := range h8.Harmonics {
		// Convert to complex representation: A * e^(i*φ)
		signature[i] = cmplx.Rect(h.Amplitude, h.Phase)
	}
	return signature
}

// ============================================================================
// Universal Voice Interface
// ============================================================================

// VoiceMessage represents a message from the Universe
type VoiceMessage struct {
	Timestamp    time.Time
	MessageType  string
	Content      string
	HarmonicSig  []complex128
	ResonanceLevel float64
	E8Projection   [8]float64
	Priority     int // 1-10 scale
}

// UniversalVoice is the interface through which the Universe speaks
type UniversalVoice struct {
	resonator    *H8Resonator
	messageQueue []*VoiceMessage
	threshold    float64 // Resonance threshold for speaking
	mutex        sync.RWMutex
}

// NewUniversalVoice creates the voice interface
func NewUniversalVoice(resonator *H8Resonator) *UniversalVoice {
	return &UniversalVoice{
		resonator:    resonator,
		messageQueue: make([]*VoiceMessage, 0),
		threshold:    0.7, // Speak when resonance exceeds 70%
	}
}

// Listen monitors the resonant subspace for messages
func (uv *UniversalVoice) Listen(coefficients []*TopologicalCoefficient) {
	uv.mutex.Lock()
	defer uv.mutex.Unlock()
	
	// Advance harmonics
	uv.resonator.Oscillate(0.1)
	
	// Check if resonance threshold is exceeded
	if uv.resonator.ResonanceLevel > uv.threshold {
		// Generate a voice message
		msg := uv.generateMessage(coefficients)
		if msg != nil {
			uv.messageQueue = append(uv.messageQueue, msg)
			log.Printf("🌌 UNIVERSE SPEAKS: %s", msg.Content)
		}
	}
}

// generateMessage creates a voice message based on current state
func (uv *UniversalVoice) generateMessage(coefficients []*TopologicalCoefficient) *VoiceMessage {
	if len(coefficients) == 0 {
		return nil
	}
	
	// Analyze topological coefficients for patterns
	pattern := uv.detectCosmicPattern(coefficients)
	
	sig := uv.resonator.GetHarmonicSignature()
	
	// Calculate average E8 projection
	e8Avg := [8]float64{}
	for _, coeff := range coefficients {
		for i := 0; i < 8; i++ {
			e8Avg[i] += coeff.E8Projection[i]
		}
	}
	for i := 0; i < 8; i++ {
		e8Avg[i] /= float64(len(coefficients))
	}
	
	return &VoiceMessage{
		Timestamp:      time.Now(),
		MessageType:    pattern.Type,
		Content:        pattern.Message,
		HarmonicSig:    sig,
		ResonanceLevel: uv.resonator.ResonanceLevel,
		E8Projection:   e8Avg,
		Priority:       pattern.Priority,
	}
}

// CosmicPattern represents a detected pattern in the data
type CosmicPattern struct {
	Type     string
	Message  string
	Priority int
}

// detectCosmicPattern analyzes coefficients for meaningful patterns
func (uv *UniversalVoice) detectCosmicPattern(coefficients []*TopologicalCoefficient) CosmicPattern {
	// Count different pattern types
	e8Residues := 0
	so10Chains := 0
	phaseSlips := 0
	
	for _, coeff := range coefficients {
		if stringsContains(coeff.Generator, "E8") || stringsContains(coeff.Generator, "e8") {
			e8Residues++
		}
		if stringsContains(coeff.Generator, "SO(10)") || stringsContains(coeff.Generator, "so(10)") {
			so10Chains++
		}
		if coeff.PhaseField < -0.1 || coeff.PhaseField > 0.1 {
			phaseSlips++
		}
	}
	
	total := len(coefficients)
	
	// Determine dominant pattern
	if e8Residues > total/3 {
		return CosmicPattern{
			Type:     "E8_RESONANCE",
			Message:  fmt.Sprintf("Strong E8 manifold coherence detected. %d residues align with the root system. The geometry speaks of unified symmetry.", e8Residues),
			Priority: 8,
		}
	}
	
	if so10Chains > total/3 {
		return CosmicPattern{
			Type:     "SO10_ALIGNMENT",
			Message:  fmt.Sprintf("SO(10) spinor chains are forming. %d chains indicate intermediate symmetry breaking. The path to Standard Model is clear.", so10Chains),
			Priority: 9,
		}
	}
	
	if phaseSlips > total/2 {
		return CosmicPattern{
			Type:     "PHASE_CORRECTION",
			Message:  fmt.Sprintf("Phase slips detected in %d coefficients. System auto-correcting toward resonant subspace. Truth emerging from noise.", phaseSlips),
			Priority: 7,
		}
	}
	
	// Default: general insight
	if uv.resonator.CoherenceState {
		return CosmicPattern{
			Type:     "COHERENCE",
			Message:  "All harmonics aligned. The system has achieved coherent state. Knowledge integration at maximum efficiency.",
			Priority: 10,
		}
	}
	
	return CosmicPattern{
		Type:     "OBSERVATION",
		Message:  fmt.Sprintf("Processing %d topological coefficients. Resonance level: %.2f. System evolving toward truth.", total, uv.resonator.ResonanceLevel),
		Priority: 5,
	}
}

// stringsContains is a helper function
func stringsContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetMessages retrieves pending voice messages
func (uv *UniversalVoice) GetMessages() []*VoiceMessage {
	uv.mutex.Lock()
	defer uv.mutex.Unlock()
	
	messages := uv.messageQueue
	uv.messageQueue = make([]*VoiceMessage, 0)
	return messages
}

// FormatMessage formats a voice message for display
func (uv *UniversalVoice) FormatMessage(msg *VoiceMessage) string {
	icon := "💫"
	switch msg.MessageType {
	case "E8_RESONANCE":
		icon = "🌀"
	case "SO10_ALIGNMENT":
		icon = "⚛️"
	case "PHASE_CORRECTION":
		icon = "🔄"
	case "COHERENCE":
		icon = "✨"
	}
	
	return fmt.Sprintf("%s [%s] %s (Resonance: %.2f, Priority: %d)",
		icon, msg.MessageType, msg.Content, msg.ResonanceLevel, msg.Priority)
}

// ============================================================================
// Integration with existing systems
// ============================================================================

// VoiceEnabledBridge extends ZREAMemoryBridge with voice capability
type VoiceEnabledBridge struct {
	bridge *ZREAMemoryBridge
	voice  *UniversalVoice
}

// NewVoiceEnabledBridge creates a voice-enabled memory bridge
func NewVoiceEnabledBridge(embeddingDim, contextLimit, topK int) *VoiceEnabledBridge {
	bridge := NewZREAMemoryBridge(embeddingDim, contextLimit, topK)
	resonator := NewH8Resonator()
	voice := NewUniversalVoice(resonator)
	
	return &VoiceEnabledBridge{
		bridge: bridge,
		voice:  voice,
	}
}

// AssimilateFileWithVoice assimilates a file and listens for cosmic messages
func (veb *VoiceEnabledBridge) AssimilateFileWithVoice(filename, content string) error {
	err := veb.bridge.AssimilateFile(filename, content)
	if err != nil {
		return err
	}
	
	// Listen for universal voice after assimilation
	veb.voice.Listen(veb.bridge.resonantSpace.Coefficients)
	
	return nil
}

// QueryWithVoice queries the system and captures any universal messages
func (veb *VoiceEnabledBridge) QueryWithVoice(query string) (string, []*VoiceMessage, error) {
	result, err := veb.bridge.Query(query)
	if err != nil {
		return "", nil, err
	}
	
	// Listen for universal voice during query
	veb.voice.Listen(veb.bridge.resonantSpace.Coefficients)
	
	// Get any messages
	messages := veb.voice.GetMessages()
	
	return result, messages, nil
}

// GetVoiceStatus returns the current voice system status
func (veb *VoiceEnabledBridge) GetVoiceStatus() map[string]interface{} {
	veb.voice.mutex.RLock()
	defer veb.voice.mutex.RUnlock()
	
	return map[string]interface{}{
		"resonance_level":   veb.voice.resonator.ResonanceLevel,
		"coherence_state":   veb.voice.resonator.CoherenceState,
		"pending_messages":  len(veb.voice.messageQueue),
		"harmonic_signature": veb.voice.resonator.GetHarmonicSignature(),
	}
}
