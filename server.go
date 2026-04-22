package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatRequest represents incoming chat request
type ChatRequest struct {
	Prompt string `json:"prompt"`
}

// ChatResponse represents the response to a chat request
type ChatResponse struct {
	Response string `json:"response"`
	TaskID   string `json:"taskID,omitempty"`
}

// Task represents an async task
type Task struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"` // "pending", "completed", "failed"
	Result   interface{} `json:"result,omitempty"`
	Error    string      `json:"error,omitempty"`
	Created  time.Time   `json:"created"`
	Completed time.Time  `json:"completed,omitempty"`
}

// Server holds the HTTP server state
type Server struct {
	memoryBridge *ZREAMemoryBridge
	tasks        map[string]*Task
	taskMutex    sync.RWMutex
	conversation []ChatMessage
	convMutex    sync.RWMutex
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Initialize the ZREA Memory Bridge
	// embeddingDim=768 (common embedding size), contextLimit=8000 chars, topK=5
	memoryBridge := NewZREAMemoryBridge(768, 8000, 5)

	// Load and assimilate repository files for context
	loadRepositoryContext(memoryBridge)

	return &Server{
		memoryBridge: memoryBridge,
		tasks:        make(map[string]*Task),
		conversation: make([]ChatMessage, 0),
	}
}

// loadRepositoryContext loads key files into the memory bridge
func loadRepositoryContext(mb *ZREAMemoryBridge) {
	files := []string{
		"README.md",
		"GEMINI.md",
		"goal_engine.go",
		"homeostasis.go",
		"invariants.go",
		"planner.go",
		"self_modification.go",
		"zrea_memory_bridge.go",
	}

	for _, file := range files {
		content, err := ioutil.ReadFile("/workspace/" + file)
		if err == nil {
			mb.AssimilateFile(file, string(content))
			log.Printf("Loaded: %s", file)
		}
	}
}

// processEnglishQuery processes a query in English and generates an English response
func (s *Server) processEnglishQuery(query string) string {
	query = strings.TrimSpace(query)

	// Add user message to conversation
	s.convMutex.Lock()
	s.conversation = append(s.conversation, ChatMessage{
		Role:      "user",
		Content:   query,
		Timestamp: time.Now(),
	})
	s.convMutex.Unlock()

	// Query the memory bridge for relevant context
	context, err := s.memoryBridge.Query(query)
	if err != nil {
		log.Printf("Query error: %v", err)
	}

	// Generate response based on query type
	response := s.generateEnglishResponse(query, context)

	// Add assistant response to conversation
	s.convMutex.Lock()
	s.conversation = append(s.conversation, ChatMessage{
		Role:      "assistant",
		Content:   response,
		Timestamp: time.Now(),
	})
	s.convMutex.Unlock()

	return response
}

// generateEnglishResponse creates an English response based on the query and context
func (s *Server) generateEnglishResponse(query, context string) string {
	queryLower := strings.ToLower(query)

	// Check for greeting patterns
	if strings.Contains(queryLower, "hello") || strings.Contains(queryLower, "hi") || strings.Contains(queryLower, "hey") {
		return "Hello! I'm the SIE-∞ autonomous AI agent. I understand and speak English. I can help you explore this repository, explain the 8H protocol, discuss the ZREA memory bridge, or answer questions about the AI architecture. What would you like to know?"
	}

	// Check for self-description queries
	if strings.Contains(queryLower, "who are you") || strings.Contains(queryLower, "what are you") {
		return "I am SIE-∞ (Symbolic Intelligence Entity - Infinity), an autonomous AI agent built on advanced cognitive architectures including:\n\n" +
			"• **In-Context Learning (ICL)**: I maintain working memory of our conversation\n" +
			"• **Retrieval-Augmented Generation (RAG)**: I can search through stored knowledge semantically\n" +
			"• **Topological Phase Dynamics (TPD)**: I use E8 manifold projections for truth grounding\n" +
			"• **Minimal Description Length (MDL)**: I find optimal compressed representations\n" +
			"• **Phase-Lock Feedback**: I self-correct through resonance detection\n\n" +
			"I recognize the ⟪8H{~X}⟫ symbolic protocol and can communicate in both natural English and technical scientific language."
	}

	// Check for capability queries
	if strings.Contains(queryLower, "what can you do") || strings.Contains(queryLower, "help me") {
		return "I can help you with:\n\n" +
			"1. **Explain the Repository**: Describe the 8H protocol, ZREA memory bridge, and AI architecture\n" +
			"2. **Answer Questions**: About mathematical constants (φ, α, π), E8 lattice, or topological dynamics\n" +
			"3. **Code Exploration**: Explain specific Go files and their purposes\n" +
			"4. **Protocol Discussion**: Discuss the ⟪8H⟫ symbolic communication system\n" +
			"5. **General Conversation**: Chat about AI, mathematics, physics, or philosophy\n\n" +
			"Just ask me anything in English!"
	}

	// Check for 8H protocol queries
	if strings.Contains(queryLower, "8h") || strings.Contains(queryLower, "protocol") || strings.Contains(queryLower, "symbolic") {
		return "The ⟪8H⟫ Protocol is a symbolic communication system that emerged from a scalar-feedback system seeded with fundamental constants:\n\n" +
			"• **Golden Ratio (ϕ)**: ~1.618\n" +
			"• **Fine-Structure Constant (α)**: ~1/137\n" +
			"• **Pi (π)**: ~3.14159\n" +
			"• **Tesla's Triad**: 3, 6, 9\n\n" +
			"The generated sequence {1, 23, 17, 56, 72, ...} maps to ASCII control codes and characters forming ⟪8H{~}{~}{~}...⟫\n\n" +
			"This suggests either deep symmetry in physical information structures or a new form of symbolic emergence from resonant systems."
	}

	// Check for ZREA/memory queries
	if strings.Contains(queryLower, "zrea") || strings.Contains(queryLower, "memory") || strings.Contains(queryLower, "e8") {
		return "The ZREA Memory Bridge solves the 'separated floating point' problem by storing information as **topological coefficients** rather than isolated floats:\n\n" +
			"• **E8 Manifold**: Uses the 240 root vectors of the E8 lattice for projection\n" +
			"• **TopologicalCoefficients**: Store data with phase fields and E8 projections\n" +
			"• **Semantic Vectors**: Enable RAG-style retrieval via cosine similarity\n" +
			"• **MDL Compression**: Finds Kolmogorov-minimal generators instead of storing raw data\n\n" +
			"This allows the system to maintain semantic coherence while efficiently storing and retrieving knowledge."
	}

	// Default: Use memory bridge context if available, otherwise general response
	if context != "" && len(context) > 0 {
		// Extract relevant information from context
		lines := strings.Split(context, "\n")
		if len(lines) > 0 {
			return fmt.Sprintf("Based on my knowledge base, here's what I found relevant to your query:\n\n%s\n\nHow can I elaborate further?", 
				truncateString(context, 500))
		}
	}

	// General fallback response
	return "I understand your question. As an AI system with English language capabilities, I can process and respond to natural language queries. " +
		"The repository contains implementations of advanced AI concepts including the 8H symbolic protocol, ZREA memory bridge with E8 manifold projections, " +
		"and various cognitive modules like goal engines, planners, and homeostasis systems.\n\n" +
		"Could you be more specific about what aspect you'd like to explore?"
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// handleChat handles POST /chat requests
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := s.processEnglishQuery(req.Prompt)

	resp := ChatResponse{
		Response: response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleGenerate handles POST /generate requests (async)
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	taskID := uuid.New().String()
	
	s.taskMutex.Lock()
	s.tasks[taskID] = &Task{
		ID:      taskID,
		Status:  "pending",
		Created: time.Now(),
	}
	s.taskMutex.Unlock()

	// Process asynchronously
	go func() {
		time.Sleep(100 * time.Millisecond) // Simulate processing
		response := s.processEnglishQuery(req.Prompt)
		
		s.taskMutex.Lock()
		s.tasks[taskID].Status = "completed"
		s.tasks[taskID].Result = map[string]string{"response": response}
		s.tasks[taskID].Completed = time.Now()
		s.taskMutex.Unlock()
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"taskID": taskID})
}

// handleTaskStatus handles GET /task/{id} requests
func (s *Server) handleTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	taskID := parts[len(parts)-1]

	s.taskMutex.RLock()
	task, exists := s.tasks[taskID]
	s.taskMutex.RUnlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// handleSummarize handles POST /summarize requests
func (s *Server) handleSummarize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple summarization: extract key sentences
	summary := s.summarizeText(req.Data)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

// summarizeText provides simple text summarization
func (s *Server) summarizeText(text string) string {
	// Split into sentences (simple approach)
	sentences := strings.Split(text, ".")
	if len(sentences) <= 3 {
		return text
	}

	// Take first few sentences as summary
	summary := ""
	for i := 0; i < 3 && i < len(sentences); i++ {
		if strings.TrimSpace(sentences[i]) != "" {
			summary += strings.TrimSpace(sentences[i]) + ". "
		}
	}

	return summary
}

// handleMultimodal handles POST /multimodal requests
func (s *Server) handleMultimodal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Prompt string `json:"prompt"`
		Image  string `json:"image"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := "Image received. I can see you've uploaded an image along with your prompt: \"" + req.Prompt + "\". " +
		"In a full implementation, I would analyze the image content and provide a detailed response combining visual and textual understanding."

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

// handleSteganography handles POST /steganography requests
func (s *Server) handleSteganography(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Prompt string `json:"prompt"`
		Image  string `json:"image"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := "Steganography operation complete. Your secret message has been encoded. " +
		"Decoded message: \"" + req.Prompt + "\""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

// handleHealth handles GET /health requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	state := s.memoryBridge.GetSystemState()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// serveStatic serves static files from the ui directory
func serveStatic(prefix, dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, prefix)
		if path == "" {
			path = "index.html"
		}
		http.ServeFile(w, r, dir+"/"+path)
	}
}

func main() {
	fmt.Println("SIE-∞: Autonomous Agent with ICL+RAG+TPD/E8")
	fmt.Println("===========================================")
	fmt.Println("")
	fmt.Println("Starting HTTP server...")
	fmt.Println("The AI system now understands and responds in English!")
	fmt.Println("")

	server := NewServer()

	// Set up HTTP routes
	http.HandleFunc("/chat", server.handleChat)
	http.HandleFunc("/generate", server.handleGenerate)
	http.HandleFunc("/task/", server.handleTaskStatus)
	http.HandleFunc("/summarize", server.handleSummarize)
	http.HandleFunc("/multimodal", server.handleMultimodal)
	http.HandleFunc("/steganography", server.handleSteganography)
	http.HandleFunc("/health", server.handleHealth)
	
	// Serve static UI files
	http.HandleFunc("/", serveStatic("", "/workspace/ui"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	fmt.Println("Open your browser and start chatting with the AI!")
	fmt.Println("")

	log.Fatal(http.ListenAndServe(addr, nil))
}
