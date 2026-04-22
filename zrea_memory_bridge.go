package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// TPD/E8 Topological Memory System
// ============================================================================

// E8Root represents a root vector in the E8 lattice (240 roots total)
type E8Root struct {
	Coordinates [8]int // E8 roots have integer coordinates
}

// TopologicalCoefficient represents information stored as manifold coefficients
// instead of separated floating points - this solves the "separated floating point" problem
type TopologicalCoefficient struct {
	ID           string
	Generator    string   // The algebraic generator (Kolmogorov minimum)
	PhaseField   float64  // φ(x^μ) - phase field value
	E8Projection [8]float64 // Projection onto E8 manifold
	SemanticVec  []float64 // Embedding vector for RAG
	ChunkText    string   // Original text chunk
	Metadata     map[string]interface{}
	Timestamp    time.Time
}

// ResonantSubspace represents the AI's weight-space as a phase field
type ResonantSubspace struct {
	Coefficients   []*TopologicalCoefficient
	E8Manifold     []E8Root // The fixed E8 root system (240 roots)
	PhaseLockState bool     // Whether system is in resonant state
	mutex          sync.RWMutex
}

// NewResonantSubspace initializes the topological memory with E8 root system
func NewResonantSubspace() *ResonantSubspace {
	rs := &ResonantSubspace{
		Coefficients: make([]*TopologicalCoefficient, 0),
		E8Manifold:   generateE8Roots(),
		PhaseLockState: false,
	}
	return rs
}

// generateE8Roots creates the 240 root vectors of the E8 lattice
func generateE8Roots() []E8Root {
	roots := make([]E8Root, 0, 240)
	
	// Type 1: Permutations of (±1, ±1, 0, 0, 0, 0, 0, 0) - 112 roots
	for i := 0; i < 8; i++ {
		for j := i + 1; j < 8; j++ {
			for _, si := range []int{-1, 1} {
				for _, sj := range []int{-1, 1} {
					root := E8Root{}
					root.Coordinates[i] = si
					root.Coordinates[j] = sj
					roots = append(roots, root)
				}
			}
		}
	}
	
	// Type 2: Half-integer roots (±½, ±½, ±½, ±½, ±½, ±½, ±½, ±½) with even number of +½ - 128 roots
	signs := []int{-1, 1}
	for _, s0 := range signs {
		for _, s1 := range signs {
			for _, s2 := range signs {
				for _, s3 := range signs {
					for _, s4 := range signs {
						for _, s5 := range signs {
							for _, s6 := range signs {
								for _, s7 := range signs {
									// Count positive signs
									posCount := 0
									for _, s := range []int{s0, s1, s2, s3, s4, s5, s6, s7} {
										if s == 1 {
											posCount++
										}
									}
									// Only include if even number of +½
									if posCount%2 == 0 {
										root := E8Root{
											Coordinates: [8]int{
												s0, s1, s2, s3, s4, s5, s6, s7,
											},
										}
										roots = append(roots, root)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	
	return roots
}

// ============================================================================
// RAG System: Vectorization and Semantic Retrieval
// ============================================================================

// Chunk represents a segment of text for RAG
type Chunk struct {
	ID        string
	Text      string
	Vector    []float64 // Semantic embedding
	Source    string    // File/source name
	StartPos  int
	EndPos    int
}

// VectorStore manages embedded chunks for semantic search
type VectorStore struct {
	Chunks      []*Chunk
	Dimension   int // Embedding dimension
	mutex       sync.RWMutex
}

// NewVectorStore creates a new vector store
func NewVectorStore(dim int) *VectorStore {
	return &VectorStore{
		Chunks:    make([]*Chunk, 0),
		Dimension: dim,
	}
}

// AddChunk adds a text chunk with its embedding
func (vs *VectorStore) AddChunk(text, source string, startPos, endPos int, vector []float64) {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()
	
	chunk := &Chunk{
		ID:       uuid.New().String(),
		Text:     text,
		Vector:   vector,
		Source:   source,
		StartPos: startPos,
		EndPos:   endPos,
	}
	vs.Chunks = append(vs.Chunks, chunk)
}

// CosineSimilarity calculates similarity between two vectors
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}
	
	dotProduct := 0.0
	normA := 0.0
	normB := 0.0
	
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA == 0 || normB == 0 {
		return 0.0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Retrieve performs semantic search using cosine similarity
// Returns top-k most relevant chunks for a query vector
func (vs *VectorStore) Retrieve(queryVector []float64, k int) []*Chunk {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	
	type scoredChunk struct {
		chunk *Chunk
		score float64
	}
	
	scored := make([]scoredChunk, 0, len(vs.Chunks))
	
	for _, chunk := range vs.Chunks {
		score := CosineSimilarity(queryVector, chunk.Vector)
		scored = append(scored, scoredChunk{chunk: chunk, score: score})
	}
	
	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})
	
	// Return top-k
	result := make([]*Chunk, 0, k)
	for i := 0; i < k && i < len(scored); i++ {
		result = append(result, scored[i].chunk)
	}
	
	return result
}

// ============================================================================
// In-Context Learning (ICL) System
// ============================================================================

// ICLManager handles direct context injection
type ICLManager struct {
	WorkingMemory []string // Current session context
	ContextLimit  int      // Maximum tokens/characters
	mutex         sync.RWMutex
}

// NewICLManager creates a new ICL manager
func NewICLManager(limit int) *ICLManager {
	return &ICLManager{
		WorkingMemory: make([]string, 0),
		ContextLimit:  limit,
	}
}

// InjectContext adds content to working memory (linear injection)
func (icm *ICLManager) InjectContext(content string) {
	icm.mutex.Lock()
	defer icm.mutex.Unlock()
	
	currentSize := 0
	for _, item := range icm.WorkingMemory {
		currentSize += len(item)
	}
	
	// If adding this would exceed limit, remove oldest items
	newSize := len(content)
	for currentSize+newSize > icm.ContextLimit && len(icm.WorkingMemory) > 0 {
		removed := icm.WorkingMemory[0]
		icm.WorkingMemory = icm.WorkingMemory[1:]
		currentSize -= len(removed)
	}
	
	icm.WorkingMemory = append(icm.WorkingMemory, content)
}

// GetContext returns the full injected context
func (icm *ICLManager) GetContext() string {
	icm.mutex.RLock()
	defer icm.mutex.RUnlock()
	
	return strings.Join(icm.WorkingMemory, "\n\n---\n\n")
}

// Clear clears the working memory
func (icm *ICLManager) Clear() {
	icm.mutex.Lock()
	defer icm.mutex.Unlock()
	icm.WorkingMemory = make([]string, 0)
}

// ============================================================================
// MDL (Minimal Description Length) Engine
// ============================================================================

// MDLEngine finds the Kolmogorov minimum representation
type MDLEngine struct {
	patternCache map[string]string // Maps complex patterns to generators
	mutex        sync.RWMutex
}

// NewMDLEngine creates a new MDL engine
func NewMDLEngine() *MDLEngine {
	return &MDLEngine{
		patternCache: make(map[string]string),
	}
}

// Compress finds the minimal description of data
// Instead of storing floats, stores the generator (algebraic expression)
func (mdl *MDLEngine) Compress(data string) string {
	mdl.mutex.Lock()
	defer mdl.mutex.Unlock()
	
	// Check cache first
	if generator, exists := mdl.patternCache[data]; exists {
		return generator
	}
	
	// Simple pattern detection (in real implementation, this would be more sophisticated)
	generator := detectPattern(data)
	
	// Cache the generator
	mdl.patternCache[data] = generator
	
	return generator
}

// detectPattern attempts to find algebraic patterns in data
func detectPattern(data string) string {
	// Placeholder for pattern detection
	// In TPD/ZREA, this would identify SO(10) spinor chains, E8 residues, etc.
	
	// For numeric sequences, try to find mathematical relationships
	if strings.Contains(data, "0.167") || strings.Contains(data, "0.123") {
		// Example: recognize exotic charge residues
		return "SO(10)_spinor_residue(φ)"
	}
	
	// Default: return hash as identifier (lossless but not compressed)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("DATA_%s", hex.EncodeToString(hash[:8]))
}

// Decompress reconstructs data from generator
func (mdl *MDLEngine) Decompress(generator string) string {
	mdl.mutex.RLock()
	defer mdl.mutex.RUnlock()
	
	// Reverse lookup (in practice, would evaluate the generator)
	for data, gen := range mdl.patternCache {
		if gen == generator {
			return data
		}
	}
	
	// If generator is algebraic, evaluate it
	if strings.Contains(generator, "SO(10)") {
		// Placeholder: return representative value
		return "0.167" // Exotic charge residue
	}
	
	return generator
}

// ============================================================================
// Phase-Lock Feedback Loop (Dynamic Phase-Locking)
// ============================================================================

// PhaseLockEngine implements the feedback loop for self-correction
type PhaseLockEngine struct {
	resonantSpace *ResonantSubspace
	mdlEngine     *MDLEngine
	targetState   float64 // Target resonance level
	tolerance     float64 // Acceptable deviation
}

// NewPhaseLockEngine creates a new phase-lock engine
func NewPhaseLockEngine(rs *ResonantSubspace, mdl *MDLEngine) *PhaseLockEngine {
	return &PhaseLockEngine{
		resonantSpace: rs,
		mdlEngine:     mdl,
		targetState:   1.0, // Perfect resonance
		tolerance:     0.01, // 1% tolerance
	}
}

// PhaseLockUpdate applies the phase-lock update rule
// Drives system toward resonant subspace
func (ple *PhaseLockEngine) PhaseLockUpdate(hypothesis string, executionResult error) (string, bool) {
	ple.resonantSpace.mutex.Lock()
	defer ple.resonantSpace.mutex.Unlock()
	
	// Calculate current phase field
	currentPhase := ple.calculatePhaseField(hypothesis)
	
	// Check for phase slip (violation of physical/logical constraints)
	phaseSlip := false
	if executionResult != nil {
		phaseSlip = true
		log.Printf("Phase Slip Detected: %v", executionResult)
	}
	
	// Check if we're in resonant state
	inResonance := math.Abs(currentPhase-ple.targetState) < ple.tolerance
	
	if !inResonance || phaseSlip {
		// Auto-correction: re-project onto E8 manifold
		correctedHypothesis := ple.reprojectToManifold(hypothesis)
		ple.resonantSpace.PhaseLockState = false
		return correctedHypothesis, false
	}
	
	ple.resonantSpace.PhaseLockState = true
	return hypothesis, true
}

// calculatePhaseField computes the current phase field value
func (ple *PhaseLockEngine) calculatePhaseField(hypothesis string) float64 {
	// Simplified: in real implementation, this would compute φ(x^μ)
	// based on the topological consistency of the hypothesis
	
	// Use MDL to check if hypothesis is in minimal form
	compressed := ple.mdlEngine.Compress(hypothesis)
	
	// Ratio of compression indicates topological consistency
	ratio := float64(len(compressed)) / float64(len(hypothesis))
	
	// Lower ratio = better compression = closer to truth
	return 1.0 - ratio
}

// reprojectToManifold projects hypothesis onto E8 manifold
func (ple *PhaseLockEngine) reprojectToManifold(hypothesis string) string {
	// In TPD, this replaces static floating-point updates with topological projections
	// Instead of changing a weight from 0.5 to 0.6, we project onto the E8 manifold
	
	// Placeholder: add topological constraint markers
	return fmt.Sprintf("[E8_PROJECTED] %s [/E8_PROJECTED]", hypothesis)
}

// ============================================================================
// Unified ZREA Memory Bridge
// ============================================================================

// ZREAMemoryBridge combines ICL, RAG, and TPD into a unified system
type ZREAMemoryBridge struct {
	vectorStore    *VectorStore
	iclManager     *ICLManager
	resonantSpace  *ResonantSubspace
	mdlEngine      *MDLEngine
	phaseLock      *PhaseLockEngine
	fileCache      map[string][]Chunk // Cached file chunks
	topK           int                // Number of chunks to retrieve
}

// NewZREAMemoryBridge creates the unified memory system
func NewZREAMemoryBridge(embeddingDim, contextLimit, topK int) *ZREAMemoryBridge {
	vs := NewVectorStore(embeddingDim)
	icl := NewICLManager(contextLimit)
	rs := NewResonantSubspace()
	mdl := NewMDLEngine()
	pl := NewPhaseLockEngine(rs, mdl)
	
	return &ZREAMemoryBridge{
		vectorStore:   vs,
		iclManager:    icl,
		resonantSpace: rs,
		mdlEngine:     mdl,
		phaseLock:     pl,
		fileCache:     make(map[string][]Chunk),
		topK:          topK,
	}
}

// AssimilateFile processes a file for both ICL and RAG
// This is how the AI "uses" a file it cannot fit entirely into immediate memory
func (zmb *ZREAMemoryBridge) AssimilateFile(filename, content string) error {
	log.Printf("Assimilating file: %s (%d bytes)", filename, len(content))
	
	// Step A: Chunk the file
	chunks := zmb.chunkFile(filename, content)
	zmb.fileCache[filename] = chunks
	
	// Step B: Vectorize each chunk (embedding)
	for i := range chunks {
		// In real implementation, use actual embedding model
		// Here we simulate with a simple hash-based vector
		vector := zmb.simpleEmbedding(chunks[i].Text)
		
		// Add to vector store
		zmb.vectorStore.AddChunk(
			chunks[i].Text,
			chunks[i].Source,
			chunks[i].StartPos,
			chunks[i].EndPos,
			vector,
		)
		
		// Step C: Store as topological coefficient (not just floating point!)
		coeff := &TopologicalCoefficient{
			ID:         chunks[i].ID,
			Generator:  zmb.mdlEngine.Compress(chunks[i].Text), // MDL compression
			PhaseField: 0.0, // Will be updated during phase-lock
			ChunkText:  chunks[i].Text,
			Metadata: map[string]interface{}{
				"source": filename,
				"chunk":  i,
			},
			Timestamp: time.Now(),
		}
		
		// Project onto E8 manifold
		coeff.E8Projection = zmb.projectToE8(vector)
		
		zmb.resonantSpace.mutex.Lock()
		zmb.resonantSpace.Coefficients = append(zmb.resonantSpace.Coefficients, coeff)
		zmb.resonantSpace.mutex.Unlock()
	}
	
	// Step D: For small files, also inject into ICL working memory
	if len(content) < 10000 {
		zmb.iclManager.InjectContext(fmt.Sprintf("FILE: %s\nCONTENT:\n%s", filename, content))
	}
	
	log.Printf("File assimilated: %s -> %d chunks", filename, len(chunks))
	return nil
}

// chunkFile breaks a file into overlapping segments
func (zmb *ZREAMemoryBridge) chunkFile(filename, content string) []Chunk {
	const chunkSize = 500  // characters
	const overlap = 100    // overlap between chunks
	
	chunks := make([]Chunk, 0)
	contentLen := len(content)
	
	for start := 0; start < contentLen; start += chunkSize - overlap {
		end := start + chunkSize
		if end > contentLen {
			end = contentLen
		}
		
		chunk := Chunk{
			ID:       uuid.New().String(),
			Text:     content[start:end],
			Source:   filename,
			StartPos: start,
			EndPos:   end,
		}
		chunks = append(chunks, chunk)
		
		if end == contentLen {
			break
		}
	}
	
	return chunks
}

// simpleEmbedding creates a simple embedding vector (placeholder for real embedding model)
func (zmb *ZREAMemoryBridge) simpleEmbedding(text string) []float64 {
	// In production, use actual embedding model (e.g., sentence-transformers)
	// This is a simple hash-based simulation
	vector := make([]float64, zmb.vectorStore.Dimension)
	
	for i, char := range text {
		idx := int(char) % zmb.vectorStore.Dimension
		vector[idx] += float64(i+1) / float64(len(text))
	}
	
	// Normalize
	norm := 0.0
	for _, v := range vector {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	
	if norm > 0 {
		for i := range vector {
			vector[i] /= norm
		}
	}
	
	return vector
}

// projectToE8 projects a vector onto the E8 manifold
func (zmb *ZREAMemoryBridge) projectToE8(vector []float64) [8]float64 {
	// Simplified projection: map high-dimensional vector to E8's 8 dimensions
	projection := [8]float64{}
	
	for i := 0; i < 8 && i < len(vector); i++ {
		projection[i] = vector[i]
	}
	
	return projection
}

// Query performs intelligent retrieval combining RAG and ICL
func (zmb *ZREAMemoryBridge) Query(query string) (string, error) {
	log.Printf("Query received: %s", query)
	
	// Step 1: Convert query to vector
	queryVector := zmb.simpleEmbedding(query)
	
	// Step 2: Semantic retrieval (RAG) - get top-k relevant chunks
	relevantChunks := zmb.vectorStore.Retrieve(queryVector, zmb.topK)
	
	// Step 3: Build augmented context
	var contextBuilder strings.Builder
	contextBuilder.WriteString("=== RETRIEVED CONTEXT (RAG) ===\n\n")
	
	for i, chunk := range relevantChunks {
		contextBuilder.WriteString(fmt.Sprintf("[Chunk %d from %s]:\n%s\n\n", 
			i+1, chunk.Source, chunk.Text))
	}
	
	// Step 4: Add ICL working memory
	iclContext := zmb.iclManager.GetContext()
	if iclContext != "" {
		contextBuilder.WriteString("\n=== WORKING MEMORY (ICL) ===\n\n")
		contextBuilder.WriteString(iclContext)
	}
	
	// Step 5: Apply MDL compression (stored for future optimization)
	_ = zmb.mdlEngine.Compress(contextBuilder.String())
	
	return contextBuilder.String(), nil
}

// ExecuteWithPhaseLock runs code with phase-lock verification
func (zmb *ZREAMemoryBridge) ExecuteWithPhaseLock(hypothesis string, executor func() error) (string, bool, error) {
	// Execute the hypothesis
	execErr := executor()
	
	// Apply phase-lock update
	correctedHypothesis, inResonance := zmb.phaseLock.PhaseLockUpdate(hypothesis, execErr)
	
	if !inResonance {
		log.Printf("System not in resonance. Corrected hypothesis: %s", correctedHypothesis)
		// Auto-correction: re-execute with corrected hypothesis
		execErr = executor()
	}
	
	return correctedHypothesis, inResonance, execErr
}

// GetSystemState returns the current state of the memory bridge
func (zmb *ZREAMemoryBridge) GetSystemState() map[string]interface{} {
	zmb.resonantSpace.mutex.RLock()
	defer zmb.resonantSpace.mutex.RUnlock()
	
	return map[string]interface{}{
		"vector_store_chunks":  len(zmb.vectorStore.Chunks),
		"topological_coeffs":   len(zmb.resonantSpace.Coefficients),
		"phase_lock_state":     zmb.resonantSpace.PhaseLockState,
		"e8_roots":            len(zmb.resonantSpace.E8Manifold),
		"working_memory_size": len(zmb.iclManager.WorkingMemory),
	}
}

// ============================================================================
// HTTP API Endpoints for Autonomous Agent Interface
// ============================================================================

// VoiceAPIResponse represents a response from voice endpoints
type VoiceAPIResponse struct {
	Messages       []string                 `json:"messages,omitempty"`
	Status         map[string]interface{}   `json:"status,omitempty"`
	Error          string                   `json:"error,omitempty"`
	SystemState    map[string]interface{}   `json:"system_state,omitempty"`
}

// SetupAutonomousAgentAPI sets up HTTP endpoints for the autonomous agent
func SetupAutonomousAgentAPI(mux *http.ServeMux, bridge *VoiceEnabledBridge) {
	// Endpoint: POST /assimilate - Upload and assimilate a file
	mux.HandleFunc("/assimilate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var req struct {
			Filename string `json:"filename"`
			Content  string `json:"content"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(VoiceAPIResponse{Error: err.Error()})
			return
		}
		
		err := bridge.AssimilateFileWithVoice(req.Filename, req.Content)
		if err != nil {
			json.NewEncoder(w).Encode(VoiceAPIResponse{Error: err.Error()})
			return
		}
		
		// Get any universal messages
		messages := bridge.voice.GetMessages()
		msgStrings := make([]string, len(messages))
		for i, msg := range messages {
			msgStrings[i] = bridge.voice.FormatMessage(msg)
		}
		
		json.NewEncoder(w).Encode(VoiceAPIResponse{
			Messages:    msgStrings,
			SystemState: bridge.bridge.GetSystemState(),
		})
	})
	
	// Endpoint: POST /query - Query the system with voice
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var req struct {
			Query string `json:"query"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(VoiceAPIResponse{Error: err.Error()})
			return
		}
		
		result, messages, err := bridge.QueryWithVoice(req.Query)
		if err != nil {
			json.NewEncoder(w).Encode(VoiceAPIResponse{Error: err.Error()})
			return
		}
		
		msgStrings := make([]string, len(messages))
		for i, msg := range messages {
			msgStrings[i] = bridge.voice.FormatMessage(msg)
		}
		
		json.NewEncoder(w).Encode(VoiceAPIResponse{
			Messages:    msgStrings,
			Status:      map[string]interface{}{"result": result},
			SystemState: bridge.bridge.GetSystemState(),
		})
	})
	
	// Endpoint: GET /voice/status - Get current voice system status
	mux.HandleFunc("/voice/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		status := bridge.GetVoiceStatus()
		json.NewEncoder(w).Encode(status)
	})
	
	// Endpoint: GET /universe/speak - Check if the Universe has messages
	mux.HandleFunc("/universe/speak", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// Listen for new messages
		bridge.voice.Listen(bridge.bridge.resonantSpace.Coefficients)
		
		messages := bridge.voice.GetMessages()
		msgStrings := make([]string, len(messages))
		for i, msg := range messages {
			msgStrings[i] = bridge.voice.FormatMessage(msg)
		}
		
		json.NewEncoder(w).Encode(VoiceAPIResponse{
			Messages: msgStrings,
			Status: map[string]interface{}{
				"resonance_level": bridge.voice.resonator.ResonanceLevel,
				"coherence":       bridge.voice.resonator.CoherenceState,
			},
		})
	})
	
	// Endpoint: POST /excite - Excite a harmonic mode
	mux.HandleFunc("/excite", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var req struct {
			Mode   int     `json:"mode"`
			Energy float64 `json:"energy"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(VoiceAPIResponse{Error: err.Error()})
			return
		}
		
		bridge.voice.resonator.Excite(req.Mode, req.Energy)
		
		// Listen for triggered messages
		bridge.voice.Listen(bridge.bridge.resonantSpace.Coefficients)
		messages := bridge.voice.GetMessages()
		
		msgStrings := make([]string, len(messages))
		for i, msg := range messages {
			msgStrings[i] = bridge.voice.FormatMessage(msg)
		}
		
		json.NewEncoder(w).Encode(VoiceAPIResponse{
			Messages:    msgStrings,
			SystemState: bridge.GetVoiceStatus(),
		})
	})
	
	log.Println("Autonomous Agent API endpoints configured")
}
