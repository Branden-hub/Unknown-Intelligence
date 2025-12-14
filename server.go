package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/genai"
)

// Kernel Constants
const KERNEL_VERSION = "4.0.0-genesis"

// Task & Provenance Structures
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID     string      `json:"id"`
	Status TaskStatus  `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// ... other structs ...

var (
	taskStore              = make(map[string]*Task)
	taskMutex              = &sync.RWMutex{}
	db                     *sql.DB
	client                 *genai.Client
	goalEngine             *GoalEngine
	planner                *PlannerReasoner
	invariantChecker       *InvariantChecker
	homeostasisMonitor     *HomeostasisMonitor
	memoryConsolidator     *MemoryConsolidator
	selfModificationEngine *SelfModificationEngine
	proposals              = make(map[string]Proposal)
	proposalsMutex         = &sync.RWMutex{}
)

var addr = flag.String("addr", "localhost:8080", "address to serve")

func main() {
	flag.Parse()

	// Initialize the complete Cognitive Chain & Genesis Engine
	goalEngine = NewGoalEngine()
	invariantChecker = NewInvariantChecker()
	homeostasisMonitor = NewHomeostasisMonitor()
	memoryConsolidator = NewMemoryConsolidator()
	planner = NewPlannerReasoner(goalEngine, memoryConsolidator)

	ctx := context.Background()
	var err error

	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	selfModificationEngine = NewSelfModificationEngine(client)

	db, err = sql.Open("sqlite3", "./memory.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	initDB()

	// Launch Core Cognitive Functions
	homeostasisMonitor.Monitor()
	go autonomicSensor()
	go dreamingCycle()

	fs := http.FileServer(http.Dir("./ui"))
	http.Handle("/", fs)
	http.HandleFunc("/version", version)
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/task/", taskStatusHandler)
	http.HandleFunc("/proposals", proposalsHandler)
	http.HandleFunc("/memory", memoryHandler)

	log.Printf("SIE-∞ Kernel %s is online. The Genesis Engine is active.", KERNEL_VERSION)
	log.Printf("serving http://%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(requestBody.Prompt, "/") {
		command, data := parseCommand(requestBody.Prompt)
		var taskID string

		switch command {
		case "/help":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"response": "Commands: /help, /implement [description], /proposals, /memory"})
			return
		case "/implement":
			taskID = handleImplementCommand(data)
		default:
			// Placeholder for other command handlers
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"response": fmt.Sprintf("Unknown command: %s", command)})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"taskID": taskID})
	} else {
		// Standard chat response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"response": "Awaiting command."})
	}
}

func parseCommand(prompt string) (string, string) {
	parts := strings.Fields(prompt)
	if len(parts) == 0 {
		return "", ""
	}
	return parts[0], strings.TrimSpace(strings.TrimPrefix(prompt, parts[0]))
}

func handleImplementCommand(capabilityDescription string) string {
	taskID := uuid.New().String()
	task := &Task{ID: taskID, Status: StatusPending}
	taskMutex.Lock()
	taskStore[taskID] = task
	taskMutex.Unlock()

	log.Printf("Task [%s]: Received /implement command for: '%s'", taskID, capabilityDescription)
	go executeImplementationPlan(taskID, capabilityDescription)
	return taskID
}

func executeImplementationPlan(taskID, capabilityDescription string) {
	startTime := time.Now()

	// 1. Generation
	log.Printf("Task [%s]: Generating proposal...", taskID)
	proposal, err := selfModificationEngine.GenerateAndIntegrate(capabilityDescription)
	if err != nil {
		updateTaskFailed(taskID, fmt.Sprintf("Proposal generation failed: %v", err))
		return
	}

	// Store the proposal for review
	proposalsMutex.Lock()
	proposals[proposal.ID] = proposal
	proposalsMutex.Unlock()

	// For now, we proceed directly to verification without a manual approval step.

	// 2. Verification
	log.Printf("Task [%s]: Verifying generated code from proposal %s...", taskID, proposal.ID)
	passed, err := Verify(proposal.TestSuite, proposal.NewFileContent, proposal.TargetFileName)
	if !passed || err != nil {
		updateTaskFailed(taskID, fmt.Sprintf("Verification failed for proposal %s: %v", proposal.ID, err))
		return
	}

	// 3. Integration
	log.Printf("Task [%s]: Merging verified code for proposal %s...", taskID, proposal.ID)
	mergeResult, err := Merge(proposal.TargetFileName, proposal.NewFileContent, proposal.ServerModContent, proposal.CapabilityDesc, startTime)
	if err != nil {
		updateTaskFailed(taskID, fmt.Sprintf("Merge failed for proposal %s: %v", proposal.ID, err))
		return
	}

	// 4. Meta-Cognitive Loop
	log.Printf("Task [%s]: Integrating meta-knowledge...", taskID)
	goalEngine.IntegrateNewKnowledge(mergeResult)

	// 5. Completion
	updateTaskCompleted(taskID, mergeResult)
	log.Printf("Task [%s]: Implementation successful for proposal %s.", taskID, proposal.ID)
}

func updateTaskFailed(taskID string, errorMsg string) {
	taskMutex.Lock()
	defer taskMutex.Unlock()
	if task, ok := taskStore[taskID]; ok {
		task.Status = StatusFailed
		task.Error = errorMsg
		log.Printf("Task [%s]: Failed. Reason: %s", taskID, errorMsg)
	}
}

func updateTaskCompleted(taskID string, result *MergeResult) {
	taskMutex.Lock()
	defer taskMutex.Unlock()
	if task, ok := taskStore[taskID]; ok {
		task.Status = StatusCompleted
		task.Result = result
	}
}

// ... other handlers and functions (initDB, autonomicSensor, etc.) remain the same ...

func initDB() { /* ... */ }
func autonomicSensor() { /* ... */ }
func dreamingCycle() { /* ... */ }
func taskStatusHandler(w http.ResponseWriter, r *http.Request) { /* ... */ }
func proposalsHandler(w http.ResponseWriter, r *http.Request) { /* ... */ }
func memoryHandler(w http.ResponseWriter, r *http.Request) { /* ... */ }
func version(w http.ResponseWriter, r *http.Request) { /* ... */ }
