package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type CommandRequest struct {
	Command string `json:"command"`
	Arg     string `json:"arg"`
}

type VMStatus struct {
	Name     string `json:"Name"`
	State    string `json:"State"`
	Provider string `json:"Provider"`
}

func executeCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Received command: %s, Arg: %s\n", req.Command, req.Arg)

	var cmd *exec.Cmd
	switch req.Command {
	case "status":
		cmd = exec.Command("vagrant", "status", "--machine-readable")
	case "start":
		cmd = exec.Command("vagrant", "up", req.Arg)
	case "stop":
		cmd = exec.Command("vagrant", "halt", req.Arg)
	case "reload":
		cmd = exec.Command("vagrant", "reload", req.Arg)
	case "reboot":
		cmd = exec.Command("vagrant", "halt", req.Arg, "&&", "vagrant", "up", req.Arg)
	case "remove":
		cmd = exec.Command("vagrant", "destroy", "-f", req.Arg)
	}

	if cmd == nil {
		http.Error(w, "Unknown command", http.StatusBadRequest)
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command error: %s\n", err)
		http.Error(w, "Command execution failed", http.StatusInternalServerError)
		return
	}

	log.Printf("Command output: %s\n", string(output))

	if req.Command == "status" {
		lines := strings.Split(string(output), "\n")
		var statuses []VMStatus
		for _, line := range lines {
			fields := strings.Split(line, ",")
			if len(fields) >= 4 && fields[2] == "state" {
				status := VMStatus{
					Name:     fields[1],
					State:    fields[3],
					Provider: fields[2],
				}
				statuses = append(statuses, status)
			}
		}

		resp, err := json.Marshal(map[string]interface{}{"status": statuses})
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Command executed successfully"))
	}
}

func main() {
	http.HandleFunc("/execute", executeCommand)

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

