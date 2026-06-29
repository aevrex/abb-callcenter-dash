package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	router *mux.Router
}

type PageData struct {
	Title       string
	Agents      []AgentData
	Queues      []QueueData
	StateCounts map[string]map[string]map[string]int // dept → team → state → count
}

type QueueData struct {
	Label   string `json:"label"`
	Count   int    `json:"count"`
	Longest string `json:"longest"`
}

type AgentData struct {
	Name         []string `json:"name"`
	TeamName     string   `json:"teamName"`
	StatusChange string   `json:"last_status_change"`
	State        string   `json:"raw_status"`
}

func NewApp() *App {
	app := &App{
		router: mux.NewRouter(),
	}

	app.routes()

	return app
}

func (app *App) routes() {
	app.router.HandleFunc("/", app.handleHome).Methods("GET")
	app.router.HandleFunc("/TL", app.handleTL).Methods("GET")
	app.router.HandleFunc("/queues", app.handleQueues).Methods("GET")
	app.router.HandleFunc("/tvqueues", app.handleTLQueues).Methods("GET")
	app.router.HandleFunc("/agents", app.handleAgents).Methods("GET")
}

func (app *App) Run() {
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", app.router))
}

// Handlers
func (app *App) handleHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, "TV.html", PageData{
		Title: "Dashboard",
	})
}

func (app *App) handleTL(w http.ResponseWriter, r *http.Request) {
	app.render(w, "TL.html", PageData{
		Title: "TV View",
	})
}

func (app *App) handleQueues(w http.ResponseWriter, r *http.Request) {
	queueURL  := "http://dashboardbeo01.prd.aussiebb.io/api/v1/queues/"
	csURL     := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rcs"
	salesURL  := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rsales"
	becoURL   := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=beco"
	activationsURL   := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rca"


	queues, err := fetchQueueData(queueURL)
	if err != nil {
		http.Error(w, "Failed to fetch queue data", http.StatusInternalServerError)
		log.Println("failed to fetch queue data:", err)
		return
	}

	csAgents, err := fetchAgentData(csURL)
	if err != nil {
		http.Error(w, "Failed to fetch CS agent data", http.StatusInternalServerError)
		log.Println("failed to fetch CS agent data:", err)
		return
	}

	salesAgents, err := fetchAgentData(salesURL)
	if err != nil {
		http.Error(w, "Failed to fetch sales agent data", http.StatusInternalServerError)
		log.Println("failed to fetch sales agent data:", err)
		return
	}

	becoAgents, err := fetchAgentData(becoURL)
	if err != nil {
		http.Error(w, "Failed to fetch sales agent data", http.StatusInternalServerError)
		log.Println("failed to fetch sales agent data:", err)
		return
	}

	activationsAgents, err := fetchAgentData(activationsURL)
	if err != nil {
		http.Error(w, "Failed to fetch sales agent data", http.StatusInternalServerError)
		log.Println("failed to fetch sales agent data:", err)
		return
	}

	app.renderPartial(w, "queues.html", PageData{
		Queues: queues,
		StateCounts: map[string]map[string]map[string]int{
			"Res CS":    countByState(csAgents),
			"Res Sales": countByState(salesAgents),
			"BECO Assurance": countByState(becoAgents),
			"Res Activations": countByState(activationsAgents),
		},
	})
}

func (app *App) handleTLQueues(w http.ResponseWriter, r *http.Request) {
	queueURL := "http://dashboardbeo01.prd.aussiebb.io/api/v1/queues/"
	csURL    := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rcs"
	salesURL := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rsales"
	becoURL   := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=beco"
	activationsURL   := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rca"

	queues, err := fetchQueueData(queueURL)
	if err != nil {
		http.Error(w, "Failed to fetch queue data", http.StatusInternalServerError)
		log.Println("failed to fetch queue data:", err)
		return
	}

	csAgents, err := fetchAgentData(csURL)
	if err != nil {
		http.Error(w, "Failed to fetch CS agent data", http.StatusInternalServerError)
		log.Println("failed to fetch CS agent data:", err)
		return
	}

	salesAgents, err := fetchAgentData(salesURL)
	if err != nil {
		http.Error(w, "Failed to fetch sales agent data", http.StatusInternalServerError)
		log.Println("failed to fetch sales agent data:", err)
		return
	}

	becoAgents, err := fetchAgentData(becoURL)
	if err != nil {
		http.Error(w, "Failed to fetch sales agent data", http.StatusInternalServerError)
		log.Println("failed to fetch sales agent data:", err)
		return
	}

	activationsAgents, err := fetchAgentData(activationsURL)
	if err != nil {
		http.Error(w, "Failed to fetch sales agent data", http.StatusInternalServerError)
		log.Println("failed to fetch sales agent data:", err)
		return
	}

	app.renderPartial(w, "tvQueues.html", PageData{
		Queues: queues,
		StateCounts: map[string]map[string]map[string]int{
			"Res CS":    countByState(csAgents),
			"Res Sales": countByState(salesAgents),
			"BECO Assurance": countByState(becoAgents),
			"Res Activations": countByState(activationsAgents),
		},
	})
}

func (app *App) handleAgents(w http.ResponseWriter, r *http.Request) {
	agentsURL := "http://dashboardbeo01.prd.aussiebb.io/api/v1/agents?team=rcs"

	agents, err := fetchAgentData(agentsURL)
	if err != nil {
		http.Error(w, "Failed to fetch agent data", http.StatusInternalServerError)
		log.Println("failed to fetch agent data:", err)
		return
	}

	app.renderPartial(w, "agents.html", PageData{
		Agents: agents,
	})
}

// Rendering htmx content
func (app *App) render(w http.ResponseWriter, page string, data PageData) {
	files := []string{
		"templates/index.html",
		"templates/" + page,
	}

	views, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		log.Println("template parse error:", err)
		return
	}

	err = views.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
		log.Println("template render error:", err)
		return
	}
}

func (app *App) renderPartial(w http.ResponseWriter, page string, data PageData) {
	views, err := template.ParseFiles("templates/partials/" + page)
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = views.ExecuteTemplate(w, page, data)
	if err != nil {
		http.Error(w, "Template render error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

// Data fetches
func fetchQueueData(url string) ([]QueueData, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", response.Status)
	}

	var queueData []QueueData

	err = json.NewDecoder(response.Body).Decode(&queueData)
	if err != nil {
		return nil, fmt.Errorf("JSON decoding failed: %w", err)
	}

	return queueData, nil
}

func fetchAgentData(url string) ([]AgentData, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", response.Status)
	}

	var agentData []AgentData

	err = json.NewDecoder(response.Body).Decode(&agentData)
	if err != nil {
		return nil, fmt.Errorf("JSON decoding failed: %w", err)
	}

	return agentData, nil
}

func countByState(agents []AgentData) map[string]map[string]int {
	counts := make(map[string]map[string]int)
	for _, agent := range agents {
		if counts[agent.TeamName] == nil {
			counts[agent.TeamName] = make(map[string]int)
		}
		counts[agent.TeamName][agent.State]++
	}
	return counts
}

func main() {
	app := NewApp()
	app.Run()
}
