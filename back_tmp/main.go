package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

type CommandRequest struct {
	Command string `json:"command"`
	Arg     string `json:"arg"`
}

type VMStatus struct {
	Name     string `json:"Name"`
	State    string `json:"State"`
	IP       string `json:"IP"`
	Port     string `json:"Port"`
	Provider string `json:"Provider"`
}

type IPInfo struct {
	PublicIP  string `json:"publicIP"`
	PrivateIP string `json:"privateIP"`
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file")
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func getPublicIP() string {
	return getEnv("Public_IP", "Unavailable")
}

func getPrivateIP() string {
	return getEnv("Private_IP", "Unavailable")
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
		cmd = exec.Command("vagrant", "halt", req.Arg)
		if err := cmd.Run(); err != nil {
			log.Printf("Command error: %s\n", err)
			http.Error(w, "Command execution failed", http.StatusInternalServerError)
			return
		}
		cmd = exec.Command("vagrant", "up", req.Arg)
	case "remove":
		cmd = exec.Command("vagrant", "destroy", "-f", req.Arg)
	default:
		http.Error(w, "Unknown command", http.StatusBadRequest)
		return
	}

	if cmd == nil {
		http.Error(w, "Unknown command", http.StatusBadRequest)
		return
	}

	outputPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error getting stdout pipe: %s\n", err)
		http.Error(w, "Error executing command", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %s\n", err)
		http.Error(w, "Error executing command", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(outputPipe)
	var output []string
	for scanner.Scan() {
		line := scanner.Text()
		output = append(output, line)
		log.Printf("Output: %s\n", line)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command error: %s\n", err)
		http.Error(w, "Command execution failed", http.StatusInternalServerError)
		return
	}

	var vmStatuses []VMStatus
	if req.Command == "status" {
		vmStatuses, err = parseVagrantStatus()
		if err != nil {
			log.Printf("Error parsing vagrant status: %s\n", err)
			http.Error(w, "Error parsing status", http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(map[string]interface{}{"status": vmStatuses, "output": strings.Join(output, "\n")})
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func parseVagrantStatus() ([]VMStatus, error) {
	cmd := exec.Command("vagrant", "status", "--machine-readable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var statuses []VMStatus
	ipMap, portMap, err := getVMIpPortMapping()
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) >= 4 && fields[2] == "state" {
			ip := ipMap[fields[1]]
			port := portMap[fields[1]]
			status := VMStatus{
				Name:     fields[1],
				State:    fields[3],
				IP:       ip,
				Port:     port,
				Provider: fields[2],
			}
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}

func getVMIpPortMapping() (map[string]string, map[string]string, error) {
	file, err := os.Open("Vagrantfile")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	ipMap := make(map[string]string)
	portMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	ipRegex := regexp.MustCompile(`ip:\s*"([^"]+)"`)
	portRegex := regexp.MustCompile(`host:\s*"([^"]+)"`)
	nameRegex := regexp.MustCompile(`node_id\s*=\s*"([^"]+)"`)

	var currentName string
	for scanner.Scan() {
		line := scanner.Text()
		if matches := nameRegex.FindStringSubmatch(line); matches != nil {
			currentName = matches[1]
		} else if matches := ipRegex.FindStringSubmatch(line); matches != nil {
			ipMap[currentName] = matches[1]
		} else if matches := portRegex.FindStringSubmatch(line); matches != nil {
			portMap[currentName] = matches[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return ipMap, portMap, nil
}

func getIPInfo(w http.ResponseWriter, r *http.Request) {
	publicIP := getPublicIP()
	privateIP := getPrivateIP()

	resp, err := json.Marshal(IPInfo{
		PublicIP:  publicIP,
		PrivateIP: privateIP,
	})
	if err != nil {
		http.Error(w, "Failed to marshal IP info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func main() {
	loadEnv()
	http.HandleFunc("/execute", executeCommand)
	http.HandleFunc("/ipinfo", getIPInfo)

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
