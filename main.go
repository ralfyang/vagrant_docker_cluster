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

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type CommandRequest struct {
	Command  string `json:"command"`
	Arg      string `json:"arg"`
	Password string `json:"password"`
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

type MemoryInfo struct {
	AvailableMemory string `json:"availableMemory"`
	TotalMemory     string `json:"totalMemory"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var logClients = make(map[*websocket.Conn]bool)

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

func getPassword() string {
	return getEnv("PASSWORD", "defaultpassword")
}

func authenticate(password string) bool {
	return password == getPassword()
}

func executeCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !authenticate(req.Password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	go func() {
		scanner := bufio.NewScanner(outputPipe)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("Output: %s\n", line)
			for client := range logClients {
				err := client.WriteMessage(websocket.TextMessage, []byte(line))
				if err != nil {
					log.Printf("Error writing to websocket client: %s\n", err)
					client.Close()
					delete(logClients, client)
				}
			}
		}
	}()

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

	resp, err := json.Marshal(map[string]interface{}{"status": vmStatuses})
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

func getMemoryInfo(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("free", "-h")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to get memory info", http.StatusInternalServerError)
		return
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		http.Error(w, "Failed to parse memory info", http.StatusInternalServerError)
		return
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 7 {
		http.Error(w, "Failed to parse memory info", http.StatusInternalServerError)
		return
	}

	totalMemory := fields[1]
	availableMemory := fields[6]

	resp, err := json.Marshal(MemoryInfo{
		AvailableMemory: availableMemory,
		TotalMemory:     totalMemory,
	})
	if err != nil {
		http.Error(w, "Failed to marshal memory info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %s\n", err)
		return
	}
	logClients[ws] = true

	defer func() {
		delete(logClients, ws)
		ws.Close()
	}()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %s\n", err)
			break
			break
		}
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	config, err := os.ReadFile("Vagrantfile")
	if err != nil {
		http.Error(w, "Failed to read config", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{"config": string(config)})
	if err != nil {
		http.Error(w, "Failed to marshal config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func main() {
	loadEnv()
	http.HandleFunc("/execute", executeCommand)
	http.HandleFunc("/ipinfo", getIPInfo)
	http.HandleFunc("/memoryinfo", getMemoryInfo)
	http.HandleFunc("/logs", logHandler)
	http.HandleFunc("/config", configHandler)

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

