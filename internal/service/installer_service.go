package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/configs"
)

type Step string

const (
	StepDownload Step = "download"
	StepInstall  Step = "install"
	StepStart    Step = "start"
	StepFinished Step = "finished"
	StepError    Step = "error"
)

type ProgressUpdate struct {
	Step         Step   `json:"step"`
	Progress     int    `json:"progress"`
	StepProgress int    `json:"stepProgress"`
	Message      string `json:"message"`
	Error        string `json:"error,omitempty"`
	Done         bool   `json:"done"`
	RawLine      string `json:"rawLine,omitempty"`
}

type InstallerService interface {
	InstallMinikube(messageChan chan<- ProgressUpdate, clientGone <-chan struct{})
}

type installerService struct {
	cfg *configs.InstallerConfig
}

func NewInstallerService(cfg *configs.Config) InstallerService {
	return &installerService{cfg: &cfg.Installer}
}

// --- InstallMinikube Method (calls actual installation steps) ---
func (s *installerService) InstallMinikube(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}) {
	defer close(messageChan)

	var minikubeURL string
	var targetFileName string = "minikube-download"
	// ** Define standard installation target path **
	standardInstallTarget := "/usr/local/bin/minikube"
	if runtime.GOOS == "windows" {
		s.sendFinalUpdate(messageChan, StepError, 32, 0, "Windows does not yet support automatic installation steps", true, true)
		return
	}

	osType := runtime.GOOS
	arch := runtime.GOARCH
	release := "latest"
	switch osType {
	case "linux":
		if arch == "amd64" {
			minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-linux-amd64", release)
			targetFileName = "minikube-linux-amd64"
		} else if arch == "arm64" {
			minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-linux-arm64", release)
			targetFileName = "minikube-linux-arm64"
		}
	case "darwin":
		if arch == "amd64" {
			minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-darwin-amd64", release)
			targetFileName = "minikube-darwin-amd64"
		} else if arch == "arm64" {
			minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-darwin-arm64", release)
			targetFileName = "minikube-darwin-arm64"
		}
	}
	if minikubeURL == "" {
		s.sendFinalUpdate(messageChan, StepError, 0, 0, fmt.Sprintf("Unsupported OS/Arch combination: %s/%s", osType, arch), true, true)
		return
	}

	downloadPath := filepath.Join(s.cfg.DownloadDir, targetFileName)
	log.Printf("Will download to: %s", downloadPath)
	if err := os.MkdirAll(s.cfg.DownloadDir, 0755); err != nil {
		s.sendFinalUpdate(messageChan, StepError, 2, 0, fmt.Sprintf("Unable to create download directory '%s': %v", s.cfg.DownloadDir, err), true, true)
		return
	}
	defer func() { /* ... cleanup logic remains unchanged ... */
		log.Printf("Attempting to clean up downloaded file: %s", downloadPath)
		err := os.Remove(downloadPath)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: Failed to clean up download file %s: %v", downloadPath, err)
		} else if err == nil {
			log.Printf("Successfully cleaned up downloaded file: %s", downloadPath)
		}
	}()

	// --- Step 1: Download ---
	if !s.executeDownloadStep(messageChan, clientGone, minikubeURL, downloadPath) {
		return
	}

	// --- Step 2: Actual installation (using sudo install) ---
	// **Call modified executeInstallStep**
	if !s.executeInstallStep(messageChan, clientGone, downloadPath, standardInstallTarget) {
		return
	}

	// --- Step 3: Start ---
	// Start step now assumes minikube has been successfully installed to standardInstallTarget and may be in PATH
	// We still pass configuredPath (from config.yaml) as an alternative check path
	s.executeMinikubeStartStep(messageChan, clientGone, s.cfg.MinikubePath)
}

// --- executeDownloadStep (remains unchanged) ---
func (s *installerService) executeDownloadStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, downloadURL, downloadPath string) bool {
	step := StepDownload
	log.Printf("Step [%s]: Starting download from %s to %s", step, downloadURL, downloadPath)
	s.sendProgressUpdate(messageChan, step, 5, 0, fmt.Sprintf("Starting Minikube download (%s)...", filepath.Base(downloadPath)), "", clientGone)
	if s.isClientGone(clientGone) {
		return false
	}
	cmd := exec.Command("curl", "-#", "-Lo", downloadPath, downloadURL)
	stderrPipe, _ := cmd.StderrPipe()
	if stderrPipe != nil {
		go s.parseCurlProgress(stderrPipe, messageChan, clientGone)
	}
	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)
	if err != nil {
		errMsg := fmt.Sprintf("Download failed: %v", err)
		log.Printf("Step [%s]: Error - %s", step, errMsg)
		s.sendFinalUpdate(messageChan, StepError, 15, 0, errMsg, true, true)
		return false
	}
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("File not found after download: %s", downloadPath)
		log.Printf("Step [%s]: Error - %s", step, errMsg)
		s.sendFinalUpdate(messageChan, StepError, 20, 0, errMsg, true, true)
		return false
	}
	successMsg := fmt.Sprintf("Download completed (%s) in %s", filepath.Base(downloadPath), duration.Round(time.Second))
	log.Printf("Step [%s]: Success - %s", step, successMsg)
	s.sendProgressUpdate(messageChan, step, 30, 100, successMsg, "", clientGone)
	return true
}

func (s *installerService) parseCurlProgress(stderr io.ReadCloser, messageChan chan<- ProgressUpdate, clientGone <-chan struct{}) {
	scanner := bufio.NewScanner(stderr)
	var lastOverallProgress = 5
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "%") && (strings.Contains(line, "curl") || strings.HasPrefix(strings.TrimSpace(line), "#")) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				lastField := fields[len(fields)-1]
				if strings.HasSuffix(lastField, "%") {
					percentStr := strings.TrimSuffix(lastField, "%")
					var percent float64
					if _, err := fmt.Sscanf(percentStr, "%f", &percent); err == nil && percent > 0 {
						stepProgress := int(percent)
						overallProgress := 5 + int(float64(stepProgress)*0.25)
						if overallProgress > lastOverallProgress {
							s.sendProgressUpdate(messageChan, StepDownload, overallProgress, stepProgress, fmt.Sprintf("Downloading... %.1f%%", percent), line, clientGone)
							lastOverallProgress = overallProgress
						}
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		log.Printf("Error parsing curl stderr: %v", err)
	}
}

// --- **Modified:** executeInstallStep (executes actual sudo install) ---
func (s *installerService) executeInstallStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, downloadedFile, installTarget string) bool {
	step := StepInstall
	log.Printf("Step [%s]: Attempting to install %s to %s (requires passwordless sudo)", step, downloadedFile, installTarget)
	s.sendProgressUpdate(messageChan, step, 31, 10, fmt.Sprintf("Preparing to execute install command (sudo install %s %s)...", downloadedFile, installTarget), "", clientGone)

	// **Security Warning**
	warningMsg := "Warning: About to execute installation command requiring sudo privileges. Please ensure the user running this service is properly configured for passwordless 'sudo install' execution. This poses security risks!"
	log.Println(warningMsg)
	s.sendProgressUpdate(messageChan, step, 32, 20, warningMsg, warningMsg, clientGone) // Send warning

	if s.isClientGone(clientGone) {
		return false
	}

	// --- Execute sudo install command ---
	cmd := exec.Command("sudo", "install", downloadedFile, installTarget)
	log.Printf("Executing command: %s", cmd.String())

	outputBytes, err := cmd.CombinedOutput() // Capture both stdout and stderr
	output := string(outputBytes)
	if len(output) > 0 { // Only log when there's output
		log.Printf("sudo install output:\n%s", output)
		// Also send sudo output to frontend logs
		s.sendProgressUpdate(messageChan, step, 35, 50, "Install command output:", output, clientGone)
	}

	if err != nil {
		errMsg := fmt.Sprintf("Installation command 'sudo install' execution failed: %v", err)
		// Try to parse more specific errors from output
		if strings.Contains(output, "incorrect password attempt") || strings.Contains(output, "sudo: a password is required") {
			errMsg = "Installation failed: 'sudo install' requires password or passwordless sudo is not configured."
			log.Println("Error: sudo requires password. Please configure passwordless sudo.")
		} else if strings.Contains(output, "Permission denied") {
			errMsg = fmt.Sprintf("Installation failed: Permission denied. Cannot write to target directory %s or sudo configuration is incorrect.", installTarget)
			log.Println("Error: Permission denied.")
		} else if strings.Contains(output, "No such file or directory") && strings.Contains(output, downloadedFile) {
			errMsg = fmt.Sprintf("Installation failed: Source file '%s' not found (possibly download failed or already cleaned up).", downloadedFile)
			log.Printf("Error: Source file %s not found.", downloadedFile)
		} else if strings.Contains(output, "No such file or directory") && strings.Contains(output, filepath.Dir(installTarget)) {
			errMsg = fmt.Sprintf("Installation failed: Target directory '%s' does not exist.", filepath.Dir(installTarget))
			log.Printf("Error: Target directory %s does not exist.", filepath.Dir(installTarget))
		} else {
			log.Printf("Error: 'sudo install' failed, error: %v, output: %s", err, output)
		}
		s.sendFinalUpdate(messageChan, StepError, 38, 80, errMsg, true, true) // Update progress to near completion of install step on failure
		return false
	}

	// Installation command executed successfully
	successMsg := fmt.Sprintf("Successfully installed Minikube to %s", installTarget)
	log.Printf("Step [%s]: %s", step, successMsg)
	s.sendProgressUpdate(messageChan, step, 40, 100, successMsg, "", clientGone) // Install step complete
	return true
}

// --- executeMinikubeStartStep (search logic adjusted) ---
func (s *installerService) executeMinikubeStartStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, configuredPath string) {
	step := StepStart
	log.Printf("Step [%s]: Preparing to start 'minikube start --force'...", step)
	s.sendProgressUpdate(messageChan, step, 40, 0, "Preparing to start Minikube...", "", clientGone)
	if s.isClientGone(clientGone) {
		return
	}

	minikubeCmdPath := ""
	standardInstallPath := "/usr/local/bin/minikube" // Define standard path again for checking

	// 1. Try PATH first
	foundPath, err := exec.LookPath("minikube")
	if err == nil {
		log.Printf("Step [%s]: Found 'minikube' in PATH: %s", step, foundPath)
		minikubeCmdPath = foundPath
	} else {
		log.Printf("Step [%s]: 'minikube' not found in PATH.", step)
		// 2. Try checking standard installation path (if different from PATH)
		if _, statErr := os.Stat(standardInstallPath); statErr == nil {
			// Check execution permissions
			if info, _ := os.Stat(standardInstallPath); info.Mode()&0111 != 0 {
				log.Printf("Step [%s]: Found executable file at standard path %s.", step, standardInstallPath)
				minikubeCmdPath = standardInstallPath
			} else {
				log.Printf("Step [%s]: Found file at standard path %s but no execution permissions.", step, standardInstallPath)
			}
		} else {
			log.Printf("Step [%s]: Standard path %s does not exist or is inaccessible: %v", step, standardInstallPath, statErr)
		}

		// 3. If none found above, finally try the path from config file (if provided)
		if minikubeCmdPath == "" && configuredPath != "" {
			log.Printf("Step [%s]: Trying to use configured path: %s", step, configuredPath)
			if info, statErr := os.Stat(configuredPath); statErr == nil && (info.Mode()&0111 != 0) {
				minikubeCmdPath = configuredPath
				log.Printf("Step [%s]: Using configured path: %s", step, minikubeCmdPath)
			} else {
				log.Printf("Step [%s]: Configured path '%s' does not exist or is not executable.", step, configuredPath)
			}
		}
	}

	// 4. If still no command path found
	if minikubeCmdPath == "" {
		errMsg := "'minikube' command not found or not executable in PATH, standard path, and configured path. Please check installation step logs or manually verify installation."
		log.Printf("Step [%s]: Error - %s", step, errMsg)
		s.sendFinalUpdate(messageChan, StepError, 42, 0, errMsg, true, true)
		return
	}

	// --- Execute command using found minikubeCmdPath ---
	minikubeDriver := s.cfg.MinikubeDriver
	cmd := exec.Command(minikubeCmdPath, "start", "--force", fmt.Sprintf("--driver=%s", minikubeDriver))
	log.Printf("Executing command: %s", cmd.String())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		s.sendFinalUpdate(messageChan, StepError, 43, 0, fmt.Sprintf("Failed to create stdout pipe: %v", err), true, true)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		s.sendFinalUpdate(messageChan, StepError, 43, 0, fmt.Sprintf("Failed to create stderr pipe: %v", err), true, true)
		return
	}
	if err := cmd.Start(); err != nil {
		s.sendFinalUpdate(messageChan, StepError, 44, 0, fmt.Sprintf("Failed to start minikube command: %v", err), true, true)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var lastOverallProgress int = 40
	// ... (start step Goroutine and Wait logic remains unchanged) ...
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("STDOUT: %s", line)
			mkProgress, message := s.parseMinikubeOutput(line)
			stepProgress := 0
			if mkProgress > 0 {
				stepProgress = mkProgress
			}
			overallProgress := 40 + int(float64(stepProgress)*0.6)
			if overallProgress > lastOverallProgress {
				lastOverallProgress = overallProgress
			}
			s.sendProgressUpdate(messageChan, step, lastOverallProgress, stepProgress, message, line, clientGone)
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			log.Printf("Error reading stdout: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("STDERR: %s", line)
			mkProgress, message := s.parseMinikubeOutput(line)
			stepProgress := 0
			currentProg := lastOverallProgress
			if mkProgress > 0 {
				stepProgress = mkProgress
				overallProgress := 40 + int(float64(stepProgress)*0.6)
				if overallProgress > currentProg {
					currentProg = overallProgress
					lastOverallProgress = currentProg
				}
			} else {
				stepProgress = int(float64(currentProg-40) / 0.6)
			}
			displayMessage := fmt.Sprintf("[Log] %s", message)
			if strings.Contains(strings.ToLower(line), "error") || strings.Contains(strings.ToLower(line), "fail") {
				displayMessage = fmt.Sprintf("[Error Log] %s", message)
			}
			s.sendProgressUpdate(messageChan, step, currentProg, stepProgress, displayMessage, line, clientGone)
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			log.Printf("Error reading stderr: %v", err)
		}
	}()

	cmdErr := cmd.Wait()
	wg.Wait()
	log.Println("Minikube start command finished execution and output processing.")
	select {
	case <-clientGone:
		log.Println("Minikube start completed, but client has disconnected.")
	default:
		if cmdErr != nil {
			errMsg := fmt.Sprintf("Minikube start failed: %v", cmdErr)
			log.Println(errMsg)
			s.sendFinalUpdate(messageChan, StepError, lastOverallProgress, 100, errMsg, true, true)
		} else {
			successMsg := "Minikube started successfully!"
			log.Println(successMsg)
			s.sendFinalUpdate(messageChan, StepFinished, 100, 100, successMsg, false, true)
		}
	}
}

func (s *installerService) parseMinikubeOutput(line string) (progress int, message string) {
	// ... (code same as previous version) ...
	lineLower := strings.ToLower(line)
	message = line
	if strings.Contains(line, "minikube v") {
		return 5, "Initializing..."
	}
	if strings.Contains(lineLower, "using the") && strings.Contains(lineLower, "driver") {
		return 10, line
	}
	if strings.Contains(lineLower, "starting control plane node") {
		return 15, line
	}
	if strings.Contains(lineLower, "creating") && (strings.Contains(lineLower, "container") || strings.Contains(lineLower, "vm")) {
		return 20, line
	}
	if strings.Contains(lineLower, "preparing kubernetes") {
		return 30, line
	}
	if strings.Contains(lineLower, "pulling base image") {
		return 35, line
	}
	if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "kubelet") {
		return 40, "Downloading Kubelet"
	}
	if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "kubeadm") {
		return 45, "Downloading Kubeadm"
	}
	if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "kubectl") {
		return 50, "Downloading Kubectl"
	}
	if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "cni") {
		return 55, "Downloading CNI plugins"
	}
	if strings.Contains(lineLower, "downloading") {
		return 60, line
	}
	if strings.Contains(lineLower, "verifying kubernetes components") {
		return 65, line
	}
	if strings.Contains(lineLower, "generating certificates") {
		return 70, line
	}
	if strings.Contains(lineLower, "booting up control plane") {
		return 75, line
	}
	if strings.Contains(lineLower, "configuring") || strings.Contains(lineLower, "waiting for") {
		return 80, line
	}
	if strings.Contains(lineLower, "setting up kubeconfig") {
		return 85, line
	}
	if strings.Contains(lineLower, "enabling addons") {
		return 90, line
	}
	if strings.Contains(lineLower, "kubectl is now configured") {
		return 95, line
	}
	if strings.Contains(lineLower, "done!") || strings.Contains(lineLower, "successfully") {
		return 98, line
	}
	return -1, message
}

func (s *installerService) isClientGone(clientGone <-chan struct{}) bool { /* ... same as your provided version ... */
	select {
	case <-clientGone:
		log.Println("SSE Service: Client disconnection detected.")
		return true
	default:
		return false
	}
}
func (s *installerService) sendProgressUpdate(messageChan chan<- ProgressUpdate, step Step, overallProgress, stepProgress int, message string, rawLine string, clientGone <-chan struct{}) { /* ... same as your provided version ... */
	if s.isClientGone(clientGone) {
		log.Println("SSE Service: Client has disconnected, not sending progress update.")
		return
	}
	update := ProgressUpdate{Step: step, Progress: overallProgress, StepProgress: stepProgress, Message: message, Done: false, RawLine: rawLine}
	select {
	case messageChan <- update:
	default:
		log.Printf("Warning: SSE message channel blocked or frontend not receiving, skipping update: Step=%s, Progress=%d", step, overallProgress)
	}
}
func (s *installerService) sendFinalUpdate(messageChan chan<- ProgressUpdate, step Step, overallProgress, stepProgress int, message string, isError bool, done bool) { /* ... same as your provided version ... */
	log.Printf("Attempting to send final update: Step=%s, Progress=%d, Error=%t, Done=%t, Message=%s", step, overallProgress, isError, done, message)
	update := ProgressUpdate{Step: step, Progress: overallProgress, StepProgress: stepProgress, Message: message, Done: done}
	if isError {
		update.Error = message
	}
	select {
	case messageChan <- update:
		log.Println("Final update sent to channel.")
	case <-time.After(1 * time.Second):
		log.Println("Warning: Final SSE update send timeout (channel blocked or frontend not receiving).")
	}
}
