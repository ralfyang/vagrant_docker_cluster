package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type CommandRequest struct {
	Command string `json:"command"`
	Arg     string `json:"arg,omitempty"`
}

type CommandResponse struct {
	Status []VMStatus `json:"status"`
	Error  string     `json:"error,omitempty"`
}

type VMStatus struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	Provider string `json:"provider"`
}

// RemoveANSI removes ANSI escape codes from a string
func RemoveANSI(str string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}

func parseStatus(output string) []VMStatus {
	lines := regexp.MustCompile(`\r?\n`).Split(output, -1)
	var statuses []VMStatus
	for _, line := range lines {
		parts := regexp.MustCompile(`\s+`).Split(line, 4)
		if len(parts) >= 4 {
			statuses = append(statuses, VMStatus{Name: parts[0], State: parts[2], Provider: parts[3]})
		}
	}
	return statuses
}

func executeCommand(command string, arg string) (string, string) {
	cmdStr := fmt.Sprintf("./ctl_api.sh %s", command)
	if arg != "" {
		cmdStr = fmt.Sprintf("./ctl_api.sh %s %s", command, arg)
	}

	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Dir = "/vm/vagrant_docker_cluster" // 스크립트가 있는 디렉토리로 설정

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return RemoveANSI(out.String()), RemoveANSI(stderr.String())
	}
	return RemoveANSI(out.String()), ""
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received command: %s, Arg: %s", req.Command, req.Arg)
	output, stderr := executeCommand(req.Command, req.Arg)

	var resp CommandResponse
	if req.Command == "status" {
		resp.Status = parseStatus(output)
	} else {
		resp.Status = []VMStatus{{Name: req.Arg, State: output}}
	}
	if stderr != "" {
		resp.Error = stderr
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/execute", commandHandler).Methods("POST")

	// CORS 설정 추가
	corsObj := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsObj, headersOk, methodsOk)(r)))
}

