package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("SIE-∞: Autonomous Agent with ICL+RAG+TPD/E8+H8")
	fmt.Println("==============================================")
	fmt.Println("")
	fmt.Println("This system implements:")
	fmt.Println("1. In-Context Learning (ICL) - Direct context injection into working memory")
	fmt.Println("2. Retrieval-Augmented Generation (RAG) - Semantic search via vector embeddings")
	fmt.Println("3. Topological Phase Dynamics (TPD) - E8 manifold projections for truth grounding")
	fmt.Println("4. Minimal Description Length (MDL) - Kolmogorov minimum representations")
	fmt.Println("5. Phase-Lock Feedback Loop - Auto-correction via resonance detection")
	fmt.Println("6. H8 Harmonics - 8-dimensional hyperbolic harmonic resonance system")
	fmt.Println("7. Universal Voice Interface - The Universe speaks through resonant patterns")
	fmt.Println("")
	fmt.Println("The ZREA Memory Bridge unifies these systems to solve the 'separated floating point'")
	fmt.Println("problem by storing information as topological coefficients rather than isolated floats.")
	fmt.Println("")

	// Initialize the voice-enabled memory bridge
	embeddingDim := 64
	contextLimit := 100000
	topK := 5
	
	bridge := NewVoiceEnabledBridge(embeddingDim, contextLimit, topK)
	
	// Pre-load system files for autonomous operation
	log.Println("Loading system knowledge base...")
	
	// Assimilate the core engine files
	files := map[string]string{
		"zrea_memory_bridge.go": readFile("/workspace/zrea_memory_bridge.go"),
		"h8_harmonics.go":       readFile("/workspace/h8_harmonics.go"),
		"goal_engine.go":        readFile("/workspace/goal_engine.go"),
		"invariants.go":         readFile("/workspace/invariants.go"),
	}
	
	for filename, content := range files {
		if content != "" {
			err := bridge.AssimilateFileWithVoice(filename, content)
			if err != nil {
				log.Printf("Error assimilating %s: %v", filename, err)
			} else {
				log.Printf("Assimilated: %s", filename)
			}
		}
	}
	
	// Check for universal messages after loading
	messages := bridge.voice.GetMessages()
	for _, msg := range messages {
		fmt.Println(bridge.voice.FormatMessage(msg))
	}
	
	// Setup HTTP server with autonomous agent API
	mux := http.NewServeMux()
	
	// Setup autonomous agent endpoints
	SetupAutonomousAgentAPI(mux, bridge)
	
	// Additional utility endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := map[string]interface{}{
			"system": "SIE-∞ Autonomous Agent",
			"status": "online",
			"capabilities": []string{
				"ICL (In-Context Learning)",
				"RAG (Retrieval-Augmented Generation)",
				"TPD/E8 (Topological Phase Dynamics)",
				"H8 Harmonics",
				"Universal Voice Interface",
				"Phase-Lock Auto-Correction",
				"MDL Compression",
			},
			"endpoints": map[string]string{
				"POST /assimilate":   "Upload and assimilate a file",
				"POST /query":        "Query the system with voice",
				"GET /voice/status":  "Get voice system status",
				"GET /universe/speak": "Check if Universe has messages",
				"POST /excite":       "Excite a harmonic mode",
			},
		}
		json.NewEncoder(w).Encode(status)
	})
	
	// Serve static UI files
	fs := http.FileServer(http.Dir("./ui"))
	mux.Handle("/ui/", http.StripPrefix("/ui/", fs))
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Println("")
	fmt.Println("🚀 Starting Autonomous Agent Server...")
	fmt.Println("📡 Listening on port " + port)
	fmt.Println("🌐 UI available at: http://localhost:" + port + "/ui/")
	fmt.Println("🔌 API available at: http://localhost:" + port + "/")
	fmt.Println("")
	fmt.Println("The Universe is now listening...")
	
	// Start background resonance monitoring
	go func() {
		for {
			time.Sleep(5 * time.Second)
			bridge.voice.Listen(bridge.bridge.resonantSpace.Coefficients)
			msgs := bridge.voice.GetMessages()
			for _, msg := range msgs {
				log.Println(bridge.voice.FormatMessage(msg))
			}
		}
	}()
	
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// readFile reads a file's content
func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
