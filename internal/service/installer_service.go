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

// --- Types and Constants (保持不变) ---
type Step string
const ( StepDownload Step = "download"; StepInstall  Step = "install"; StepStart    Step = "start"; StepFinished Step = "finished"; StepError    Step = "error" )
type ProgressUpdate struct { /* ... 保持不变 ... */
	Step         Step   `json:"step"`
	Progress     int    `json:"progress"`
	StepProgress int    `json:"stepProgress"`
	Message      string `json:"message"`
	Error        string `json:"error,omitempty"`
	Done         bool   `json:"done"`
	RawLine      string `json:"rawLine,omitempty"`
}

// --- Interface Definition (修正 clientGone 类型) ---
type InstallerService interface {
	// **关键修正:** clientGone 的类型是 <-chan struct{}
	InstallMinikube(messageChan chan<- ProgressUpdate, clientGone <-chan struct{})
}

// --- Struct Definition (保持不变) ---
type installerService struct {
	cfg *configs.InstallerConfig
}

// --- Constructor (保持不变) ---
func NewInstallerService(cfg *configs.Config) InstallerService {
	return &installerService{ cfg: &cfg.Installer, }
}

// --- InstallMinikube Method (修正 clientGone 类型) ---
func (s *installerService) InstallMinikube(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}) { // **关键修正:** clientGone 类型
	defer close(messageChan)

	// ... (变量定义和 URL switch case 保持不变) ...
	var minikubeURL string
	var targetFileName string = "minikube-download"
	installTarget := s.cfg.MinikubePath
	simulatedInstallPath := "/usr/local/bin/minikube" // Default simulation path
	if runtime.GOOS == "windows" { simulatedInstallPath = "C:\\minikube\\minikube.exe" }
	osType := runtime.GOOS
	arch := runtime.GOARCH
	release := "latest"
	switch osType { /* ... url logic ... */
		case "linux":
			if arch == "amd64" { minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-linux-amd64", release); targetFileName = "minikube-linux-amd64"; } else if arch == "arm64" { minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-linux-arm64", release); targetFileName = "minikube-linux-arm64"; }
		case "darwin":
			if arch == "amd64" { minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-darwin-amd64", release); targetFileName = "minikube-darwin-amd64"; } else if arch == "arm64" { minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-darwin-arm64", release); targetFileName = "minikube-darwin-arm64"; }
		case "windows":
			minikubeURL = fmt.Sprintf("https://github.com/kubernetes/minikube/releases/%s/download/minikube-windows-amd64.exe", release); targetFileName = "minikube-windows-amd64.exe"; if installTarget == "" { installTarget = simulatedInstallPath }
	}
	if minikubeURL == "" { s.sendFinalUpdate(messageChan, StepError, 0, 0, fmt.Sprintf("不支持的 OS/Arch 组合: %s/%s", osType, arch), true, true); return }

	downloadPath := filepath.Join(s.cfg.DownloadDir, targetFileName)
	log.Printf("将下载到: %s", downloadPath)
	if err := os.MkdirAll(s.cfg.DownloadDir, 0755); err != nil { s.sendFinalUpdate(messageChan, StepError, 2, 0, fmt.Sprintf("无法创建下载目录 '%s': %v", s.cfg.DownloadDir, err), true, true); return }
	defer func() { /* ... cleanup logic ... */
		log.Printf("尝试清理下载的文件: %s", downloadPath)
		err := os.Remove(downloadPath)
		if err != nil && !os.IsNotExist(err) { log.Printf("警告: 清理下载文件 %s 失败: %v", downloadPath, err) } else if err == nil { log.Printf("成功清理下载的文件: %s", downloadPath) }
	}()

	// --- 调用步骤函数，传递正确的 clientGone 类型 ---
	if !s.executeDownloadStep(messageChan, clientGone, minikubeURL, downloadPath) { return }
	if !s.executeInstallStepSimulated(messageChan, clientGone, downloadPath, simulatedInstallPath) { return }
	s.executeMinikubeStartStep(messageChan, clientGone, installTarget)
}

// --- Step Execution Functions (修正 clientGone 类型) ---
// **关键修正:** 所有接收 clientGone 的函数参数类型改为 <-chan struct{}
func (s *installerService) executeDownloadStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, downloadURL, downloadPath string) bool {
	// ... (函数体基本不变, 除了调用 isClientGone) ...
	step := StepDownload
	log.Printf("步骤 [%s]: 开始从 %s 下载到 %s", step, downloadURL, downloadPath)
	s.sendProgressUpdate(messageChan, step, 5, 0, fmt.Sprintf("开始下载 Minikube (%s)...", filepath.Base(downloadPath)), "", clientGone)
	if s.isClientGone(clientGone) { return false } // 使用修正后的 isClientGone
	cmd := exec.Command("curl", "-#", "-Lo", downloadPath, downloadURL)
	stderrPipe, _ := cmd.StderrPipe()
	if stderrPipe != nil { go s.parseCurlProgress(stderrPipe, messageChan, clientGone) } // 传递 clientGone
	startTime := time.Now(); err := cmd.Run(); duration := time.Since(startTime)
	if err != nil { errMsg := fmt.Sprintf("下载失败: %v", err); log.Printf("步骤 [%s]: 错误 - %s", step, errMsg); s.sendFinalUpdate(messageChan, StepError, 15, 0, errMsg, true, true); return false }
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) { errMsg := fmt.Sprintf("下载后文件未找到: %s", downloadPath); log.Printf("步骤 [%s]: 错误 - %s", step, errMsg); s.sendFinalUpdate(messageChan, StepError, 20, 0, errMsg, true, true); return false }
	successMsg := fmt.Sprintf("下载完成 (%s) in %s", filepath.Base(downloadPath), duration.Round(time.Second)); log.Printf("步骤 [%s]: 成功 - %s", step, successMsg); s.sendProgressUpdate(messageChan, step, 30, 100, successMsg, "", clientGone)
	return true
}

func (s *installerService) parseCurlProgress(stderr io.ReadCloser, messageChan chan<- ProgressUpdate, clientGone <-chan struct{}) { // **关键修正:** clientGone 类型
	// ... (函数体保持不变, 内部调用 sendProgressUpdate 会使用正确的 clientGone) ...
	scanner := bufio.NewScanner(stderr)
	var lastOverallProgress = 5
	for scanner.Scan() {
		// ... parsing logic ...
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
                             s.sendProgressUpdate(messageChan, StepDownload, overallProgress, stepProgress, fmt.Sprintf("下载中... %.1f%%", percent), line, clientGone) // 传递 clientGone
                             lastOverallProgress = overallProgress
                         }
                     }
                 }
             }
        }
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) { log.Printf("解析 curl stderr 时出错: %v", err) }
}

func (s *installerService) executeInstallStepSimulated(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, downloadedFile, simulatedInstallPath string) bool { // **关键修正:** clientGone 类型
	// ... (函数体基本不变, 除了调用 isClientGone 和 sendProgressUpdate) ...
	step := StepInstall
	log.Printf("步骤 [%s]: 模拟将 %s 安装到 %s", step, downloadedFile, simulatedInstallPath)
	warningMsg := fmt.Sprintf("注意：此步骤为模拟! 你需要手动或通过其他方式执行类似 'sudo install %s %s' 的操作，或确保 Minikube 在 PATH 中。", downloadedFile, simulatedInstallPath)
	log.Printf("步骤 [%s]: %s", step, warningMsg)
	s.sendProgressUpdate(messageChan, step, 35, 50, warningMsg, warningMsg, clientGone) // 传递 clientGone
	if s.isClientGone(clientGone) { return false }
	time.Sleep(1 * time.Second)
	minikubePathToCheck := "minikube"; if s.cfg.MinikubePath != "" { minikubePathToCheck = s.cfg.MinikubePath }
	foundPath, err := exec.LookPath(minikubePathToCheck)
	if err == nil { log.Printf("步骤 [%s]: 检测到现有/可用的 minikube 可执行文件: %s", step, foundPath); s.sendProgressUpdate(messageChan, step, 38, 80, fmt.Sprintf("检测到可用的 Minikube: %s", foundPath), "", clientGone); } else { log.Printf("步骤 [%s]: 未在 PATH 或配置路径 (%s) 中检测到 minikube。", step, minikubePathToCheck); s.sendProgressUpdate(messageChan, step, 38, 80, "未检测到 Minikube (将在下一步尝试启动)", "", clientGone); }
	successMsg := "安装步骤模拟完成。"; log.Printf("步骤 [%s]: 模拟完成。", step); s.sendProgressUpdate(messageChan, step, 40, 100, successMsg, "", clientGone)
	return true
}

func (s *installerService) executeMinikubeStartStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, configuredPath string) { // **关键修正:** clientGone 类型
	// ... (函数体基本不变, 除了调用 isClientGone 和 sendProgressUpdate) ...
	step := StepStart
	log.Printf("步骤 [%s]: 准备启动 'minikube start --force'...", step)
	s.sendProgressUpdate(messageChan, step, 40, 0, "准备启动 Minikube...", "", clientGone) // 传递 clientGone
	if s.isClientGone(clientGone) { return }
	minikubeCmdPath := "minikube"; foundPath, err := exec.LookPath("minikube")
	if err != nil { log.Printf("步骤 [%s]: 'minikube' 未在 PATH 中找到。", step); if configuredPath != "" { log.Printf("步骤 [%s]: 尝试使用配置的路径: %s", step, configuredPath); if info, statErr := os.Stat(configuredPath); statErr == nil && (info.Mode()&0111 != 0) { minikubeCmdPath = configuredPath; log.Printf("步骤 [%s]: 使用配置的路径: %s", step, minikubeCmdPath); } else { errMsg := fmt.Sprintf("'minikube' 命令未在 PATH 中找到，并且配置的路径 '%s' 不存在或不可执行。请确保 Minikube 已正确安装或配置路径正确。", configuredPath); log.Printf("步骤 [%s]: 错误 - %s", step, errMsg); s.sendFinalUpdate(messageChan, StepError, 42, 0, errMsg, true, true); return; } } else { errMsg := "'minikube' 命令未在 PATH 中找到，并且未配置特定的 minikube 路径。请确保 Minikube 已安装并在 PATH 中。"; log.Printf("步骤 [%s]: 错误 - %s", step, errMsg); s.sendFinalUpdate(messageChan, StepError, 42, 0, errMsg, true, true); return; } } else { log.Printf("步骤 [%s]: 在 PATH 中找到 'minikube': %s", step, foundPath); minikubeCmdPath = foundPath; }

	minikubeDriver := s.cfg.MinikubeDriver; cmd := exec.Command(minikubeCmdPath, "start", "--force", fmt.Sprintf("--driver=%s", minikubeDriver)); log.Printf("执行命令: %s", cmd.String())
	stdoutPipe, err := cmd.StdoutPipe(); if err != nil { s.sendFinalUpdate(messageChan, StepError, 43, 0, fmt.Sprintf("创建 stdout pipe 失败: %v", err), true, true); return }
	stderrPipe, err := cmd.StderrPipe(); if err != nil { s.sendFinalUpdate(messageChan, StepError, 43, 0, fmt.Sprintf("创建 stderr pipe 失败: %v", err), true, true); return }
	if err := cmd.Start(); err != nil { s.sendFinalUpdate(messageChan, StepError, 44, 0, fmt.Sprintf("启动 minikube 命令失败: %v", err), true, true); return }

	var wg sync.WaitGroup; wg.Add(2); var lastOverallProgress int = 40
	go func() { defer wg.Done(); scanner := bufio.NewScanner(stdoutPipe); for scanner.Scan() { line := scanner.Text(); log.Printf("STDOUT: %s", line); mkProgress, message := s.parseMinikubeOutput(line); stepProgress := 0; if mkProgress > 0 { stepProgress = mkProgress }; overallProgress := 40 + int(float64(stepProgress)*0.6); if overallProgress > lastOverallProgress { lastOverallProgress = overallProgress }; s.sendProgressUpdate(messageChan, step, lastOverallProgress, stepProgress, message, line, clientGone); }; if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) { log.Printf("读取 stdout 时出错: %v", err) } }()
	go func() { defer wg.Done(); scanner := bufio.NewScanner(stderrPipe); for scanner.Scan() { line := scanner.Text(); log.Printf("STDERR: %s", line); mkProgress, message := s.parseMinikubeOutput(line); stepProgress := 0; currentProg := lastOverallProgress; if mkProgress > 0 { stepProgress = mkProgress; overallProgress := 40 + int(float64(stepProgress)*0.6); if overallProgress > currentProg { currentProg = overallProgress; lastOverallProgress = currentProg; } } else { stepProgress = int(float64(currentProg-40) / 0.6); }; displayMessage := fmt.Sprintf("[Log] %s", message); if strings.Contains(strings.ToLower(line), "error") || strings.Contains(strings.ToLower(line), "fail") { displayMessage = fmt.Sprintf("[错误日志] %s", message) }; s.sendProgressUpdate(messageChan, step, currentProg, stepProgress, displayMessage, line, clientGone); }; if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) { log.Printf("读取 stderr 时出错: %v", err) } }()

	cmdErr := cmd.Wait(); wg.Wait(); log.Println("Minikube start command finished execution and output processing.")
	select { case <-clientGone: log.Println("Minikube start 完成, 但客户端已断开连接."); default: if cmdErr != nil { errMsg := fmt.Sprintf("Minikube start 失败: %v", cmdErr); log.Println(errMsg); s.sendFinalUpdate(messageChan, StepError, lastOverallProgress, 100, errMsg, true, true); } else { successMsg := "Minikube 启动成功!"; log.Println(successMsg); s.sendFinalUpdate(messageChan, StepFinished, 100, 100, successMsg, false, true); } }
}

// --- parseMinikubeOutput (保持不变) ---
func (s *installerService) parseMinikubeOutput(line string) (progress int, message string) { /* ... 保持不变 ... */
	lineLower := strings.ToLower(line); message = line
	if strings.Contains(line, "minikube v") { return 5, "Initializing..." }; if strings.Contains(lineLower, "using the") && strings.Contains(lineLower, "driver") { return 10, line }; if strings.Contains(lineLower, "starting control plane node") { return 15, line }; if strings.Contains(lineLower, "creating") && (strings.Contains(lineLower, "container") || strings.Contains(lineLower, "vm")) { return 20, line }; if strings.Contains(lineLower, "preparing kubernetes") { return 30, line }; if strings.Contains(lineLower, "pulling base image") { return 35, line }; if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "kubelet") { return 40, "Downloading Kubelet" }; if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "kubeadm") { return 45, "Downloading Kubeadm" }; if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "kubectl") { return 50, "Downloading Kubectl" }; if strings.Contains(lineLower, "downloading") && strings.Contains(lineLower, "cni") { return 55, "Downloading CNI plugins" }; if strings.Contains(lineLower, "downloading") { return 60, line }; if strings.Contains(lineLower, "verifying kubernetes components") { return 65, line }; if strings.Contains(lineLower, "generating certificates") { return 70, line }; if strings.Contains(lineLower, "booting up control plane") { return 75, line }; if strings.Contains(lineLower, "configuring") || strings.Contains(lineLower, "waiting for") { return 80, line }; if strings.Contains(lineLower, "setting up kubeconfig") { return 85, line }; if strings.Contains(lineLower, "enabling addons") { return 90, line }; if strings.Contains(lineLower, "kubectl is now configured") { return 95, line }; if strings.Contains(lineLower, "done!") || strings.Contains(lineLower, "successfully") { return 98, line }; return -1, message
}

// --- Helper Methods (修正 clientGone 类型, send 使用非阻塞) ---
// **关键修正:** isClientGone, sendProgressUpdate, sendFinalUpdate 的 clientGone 参数类型改为 <-chan struct{}
func (s *installerService) isClientGone(clientGone <-chan struct{}) bool {
	select {
	case <-clientGone: log.Println("SSE Service: 检测到客户端断开。"); return true
	default: return false
	}
}

func (s *installerService) sendProgressUpdate(messageChan chan<- ProgressUpdate, step Step, overallProgress, stepProgress int, message string, rawLine string, clientGone <-chan struct{}) {
	if s.isClientGone(clientGone) { log.Println("SSE Service: 客户端已断开，不发送进度更新。"); return }
	update := ProgressUpdate{ Step: step, Progress: overallProgress, StepProgress: stepProgress, Message: message, Done: false, RawLine: rawLine, }
	select {
	case messageChan <- update: // Sent successfully
	default: log.Printf("警告: SSE 消息通道阻塞或前端未接收，跳过更新: Step=%s, Progress=%d", step, overallProgress)
	}
}

func (s *installerService) sendFinalUpdate(messageChan chan<- ProgressUpdate, step Step, overallProgress, stepProgress int, message string, isError bool, done bool) {
	log.Printf("尝试发送最终更新: Step=%s, Progress=%d, Error=%t, Done=%t, Message=%s", step, overallProgress, isError, done, message)
	update := ProgressUpdate{ Step: step, Progress: overallProgress, StepProgress: stepProgress, Message: message, Done: done, }
	if isError { update.Error = message }
	select {
	case messageChan <- update: log.Println("最终更新已发送到通道。")
	case <-time.After(500 * time.Millisecond): log.Println("警告: 最终 SSE 更新发送超时 (通道阻塞或前端未接收)。")
	}
}