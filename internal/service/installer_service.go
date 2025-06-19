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

// --- InstallMinikube Method (调用实际安装步骤) ---
func (s *installerService) InstallMinikube(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}) {
	defer close(messageChan)

	var minikubeURL string
	var targetFileName string = "minikube-download"
	// ** 定义标准的安装目标路径 **
	standardInstallTarget := "/usr/local/bin/minikube"
	if runtime.GOOS == "windows" {
		s.sendFinalUpdate(messageChan, StepError, 32, 0, "Windows 尚不支持自动执行安装步骤", true, true)
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
		s.sendFinalUpdate(messageChan, StepError, 0, 0, fmt.Sprintf("不支持的 OS/Arch 组合: %s/%s", osType, arch), true, true)
		return
	}

	downloadPath := filepath.Join(s.cfg.DownloadDir, targetFileName)
	log.Printf("将下载到: %s", downloadPath)
	if err := os.MkdirAll(s.cfg.DownloadDir, 0755); err != nil {
		s.sendFinalUpdate(messageChan, StepError, 2, 0, fmt.Sprintf("无法创建下载目录 '%s': %v", s.cfg.DownloadDir, err), true, true)
		return
	}
	defer func() { /* ... cleanup logic 保持不变 ... */
		log.Printf("尝试清理下载的文件: %s", downloadPath)
		err := os.Remove(downloadPath)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("警告: 清理下载文件 %s 失败: %v", downloadPath, err)
		} else if err == nil {
			log.Printf("成功清理下载的文件: %s", downloadPath)
		}
	}()

	// --- 步骤 1: 下载 ---
	if !s.executeDownloadStep(messageChan, clientGone, minikubeURL, downloadPath) {
		return
	}

	// --- 步骤 2: 实际安装 (使用 sudo install) ---
	// **调用修改后的 executeInstallStep**
	if !s.executeInstallStep(messageChan, clientGone, downloadPath, standardInstallTarget) {
		return
	}

	// --- 步骤 3: 启动 ---
	// 启动步骤现在假设 minikube 已被成功安装到 standardInstallTarget 并可能位于 PATH 中
	// 我们仍然传递 configuredPath (来自 config.yaml) 作为备选检查路径
	s.executeMinikubeStartStep(messageChan, clientGone, s.cfg.MinikubePath)
}

// --- executeDownloadStep (保持不变) ---
func (s *installerService) executeDownloadStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, downloadURL, downloadPath string) bool {
	step := StepDownload
	log.Printf("步骤 [%s]: 开始从 %s 下载到 %s", step, downloadURL, downloadPath)
	s.sendProgressUpdate(messageChan, step, 5, 0, fmt.Sprintf("开始下载 Minikube (%s)...", filepath.Base(downloadPath)), "", clientGone)
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
		errMsg := fmt.Sprintf("下载失败: %v", err)
		log.Printf("步骤 [%s]: 错误 - %s", step, errMsg)
		s.sendFinalUpdate(messageChan, StepError, 15, 0, errMsg, true, true)
		return false
	}
	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("下载后文件未找到: %s", downloadPath)
		log.Printf("步骤 [%s]: 错误 - %s", step, errMsg)
		s.sendFinalUpdate(messageChan, StepError, 20, 0, errMsg, true, true)
		return false
	}
	successMsg := fmt.Sprintf("下载完成 (%s) in %s", filepath.Base(downloadPath), duration.Round(time.Second))
	log.Printf("步骤 [%s]: 成功 - %s", step, successMsg)
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
							s.sendProgressUpdate(messageChan, StepDownload, overallProgress, stepProgress, fmt.Sprintf("下载中... %.1f%%", percent), line, clientGone)
							lastOverallProgress = overallProgress
						}
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		log.Printf("解析 curl stderr 时出错: %v", err)
	}
}

// --- **修改:** executeInstallStep (执行实际的 sudo install) ---
func (s *installerService) executeInstallStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, downloadedFile, installTarget string) bool {
	step := StepInstall
	log.Printf("步骤 [%s]: 尝试将 %s 安装到 %s (需要免密 sudo)", step, downloadedFile, installTarget)
	s.sendProgressUpdate(messageChan, step, 31, 10, fmt.Sprintf("准备执行安装命令 (sudo install %s %s)...", downloadedFile, installTarget), "", clientGone)

	// **安全警告**
	warningMsg := "警告：即将执行需要 sudo 权限的安装命令。请确保运行此服务的用户已被正确配置为可以免密码执行 'sudo install'。这存在安全风险！"
	log.Println(warningMsg)
	s.sendProgressUpdate(messageChan, step, 32, 20, warningMsg, warningMsg, clientGone) // 发送警告

	if s.isClientGone(clientGone) {
		return false
	}

	// --- 执行 sudo install 命令 ---
	cmd := exec.Command("sudo", "install", downloadedFile, installTarget)
	log.Printf("执行命令: %s", cmd.String())

	outputBytes, err := cmd.CombinedOutput() // 同时捕获 stdout 和 stderr
	output := string(outputBytes)
	if len(output) > 0 { // 只在有输出时记录
		log.Printf("sudo install output:\n%s", output)
		// 将 sudo 的输出也发送给前端日志
		s.sendProgressUpdate(messageChan, step, 35, 50, "安装命令输出:", output, clientGone)
	}

	if err != nil {
		errMsg := fmt.Sprintf("安装命令 'sudo install' 执行失败: %v", err)
		// 尝试从输出中解析更具体的错误
		if strings.Contains(output, "incorrect password attempt") || strings.Contains(output, "sudo: a password is required") {
			errMsg = "安装失败：执行 'sudo install' 需要密码或未配置免密 sudo。"
			log.Println("错误: sudo 需要密码。请配置免密 sudo。")
		} else if strings.Contains(output, "Permission denied") {
			errMsg = fmt.Sprintf("安装失败：权限被拒绝。无法写入目标目录 %s 或 sudo 配置不正确。", installTarget)
			log.Println("错误: 权限被拒绝。")
		} else if strings.Contains(output, "No such file or directory") && strings.Contains(output, downloadedFile) {
			errMsg = fmt.Sprintf("安装失败：源文件 '%s' 未找到 (可能下载失败或已被清理)。", downloadedFile)
			log.Printf("错误: 源文件 %s 未找到。", downloadedFile)
		} else if strings.Contains(output, "No such file or directory") && strings.Contains(output, filepath.Dir(installTarget)) {
			errMsg = fmt.Sprintf("安装失败：目标目录 '%s' 不存在。", filepath.Dir(installTarget))
			log.Printf("错误: 目标目录 %s 不存在。", filepath.Dir(installTarget))
		} else {
			log.Printf("错误: 'sudo install' 失败，错误: %v, 输出: %s", err, output)
		}
		s.sendFinalUpdate(messageChan, StepError, 38, 80, errMsg, true, true) // 在失败时更新进度为接近完成安装步骤
		return false
	}

	// 安装命令成功执行
	successMsg := fmt.Sprintf("成功将 Minikube 安装到 %s", installTarget)
	log.Printf("步骤 [%s]: %s", step, successMsg)
	s.sendProgressUpdate(messageChan, step, 40, 100, successMsg, "", clientGone) // Install step complete
	return true
}

// --- executeMinikubeStartStep (查找逻辑调整) ---
func (s *installerService) executeMinikubeStartStep(messageChan chan<- ProgressUpdate, clientGone <-chan struct{}, configuredPath string) {
	step := StepStart
	log.Printf("步骤 [%s]: 准备启动 'minikube start --force'...", step)
	s.sendProgressUpdate(messageChan, step, 40, 0, "准备启动 Minikube...", "", clientGone)
	if s.isClientGone(clientGone) {
		return
	}

	minikubeCmdPath := ""
	standardInstallPath := "/usr/local/bin/minikube" // 再次定义标准路径以供检查

	// 1. 优先尝试 PATH
	foundPath, err := exec.LookPath("minikube")
	if err == nil {
		log.Printf("步骤 [%s]: 在 PATH 中找到 'minikube': %s", step, foundPath)
		minikubeCmdPath = foundPath
	} else {
		log.Printf("步骤 [%s]: 'minikube' 未在 PATH 中找到。", step)
		// 2. 尝试检查标准安装路径 (如果与 PATH 不同)
		if _, statErr := os.Stat(standardInstallPath); statErr == nil {
			// 检查执行权限
			if info, _ := os.Stat(standardInstallPath); info.Mode()&0111 != 0 {
				log.Printf("步骤 [%s]: 在标准路径 %s 中找到可执行文件。", step, standardInstallPath)
				minikubeCmdPath = standardInstallPath
			} else {
				log.Printf("步骤 [%s]: 在标准路径 %s 中找到文件但无执行权限。", step, standardInstallPath)
			}
		} else {
			log.Printf("步骤 [%s]: 标准路径 %s 不存在或无法访问: %v", step, standardInstallPath, statErr)
		}

		// 3. 如果上面都没找到，最后尝试配置文件中的路径 (如果提供了)
		if minikubeCmdPath == "" && configuredPath != "" {
			log.Printf("步骤 [%s]: 尝试使用配置的路径: %s", step, configuredPath)
			if info, statErr := os.Stat(configuredPath); statErr == nil && (info.Mode()&0111 != 0) {
				minikubeCmdPath = configuredPath
				log.Printf("步骤 [%s]: 使用配置的路径: %s", step, minikubeCmdPath)
			} else {
				log.Printf("步骤 [%s]: 配置的路径 '%s' 不存在或不可执行。", step, configuredPath)
			}
		}
	}

	// 4. 如果最终还是没找到命令路径
	if minikubeCmdPath == "" {
		errMsg := "'minikube' 命令在 PATH、标准路径和配置路径中均未找到或不可执行。请检查安装步骤日志或手动验证安装。"
		log.Printf("步骤 [%s]: 错误 - %s", step, errMsg)
		s.sendFinalUpdate(messageChan, StepError, 42, 0, errMsg, true, true)
		return
	}

	// --- 使用找到的 minikubeCmdPath 执行命令 ---
	minikubeDriver := s.cfg.MinikubeDriver
	cmd := exec.Command(minikubeCmdPath, "start", "--force", fmt.Sprintf("--driver=%s", minikubeDriver))
	log.Printf("执行命令: %s", cmd.String())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		s.sendFinalUpdate(messageChan, StepError, 43, 0, fmt.Sprintf("创建 stdout pipe 失败: %v", err), true, true)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		s.sendFinalUpdate(messageChan, StepError, 43, 0, fmt.Sprintf("创建 stderr pipe 失败: %v", err), true, true)
		return
	}
	if err := cmd.Start(); err != nil {
		s.sendFinalUpdate(messageChan, StepError, 44, 0, fmt.Sprintf("启动 minikube 命令失败: %v", err), true, true)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var lastOverallProgress int = 40
	// ... (启动步骤的 Goroutine 和 Wait 逻辑保持不变) ...
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
			log.Printf("读取 stdout 时出错: %v", err)
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
				displayMessage = fmt.Sprintf("[错误日志] %s", message)
			}
			s.sendProgressUpdate(messageChan, step, currentProg, stepProgress, displayMessage, line, clientGone)
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			log.Printf("读取 stderr 时出错: %v", err)
		}
	}()

	cmdErr := cmd.Wait()
	wg.Wait()
	log.Println("Minikube start command finished execution and output processing.")
	select {
	case <-clientGone:
		log.Println("Minikube start 完成, 但客户端已断开连接.")
	default:
		if cmdErr != nil {
			errMsg := fmt.Sprintf("Minikube start 失败: %v", cmdErr)
			log.Println(errMsg)
			s.sendFinalUpdate(messageChan, StepError, lastOverallProgress, 100, errMsg, true, true)
		} else {
			successMsg := "Minikube 启动成功!"
			log.Println(successMsg)
			s.sendFinalUpdate(messageChan, StepFinished, 100, 100, successMsg, false, true)
		}
	}
}

func (s *installerService) parseMinikubeOutput(line string) (progress int, message string) {
	// ... (代码与之前版本相同) ...
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

func (s *installerService) isClientGone(clientGone <-chan struct{}) bool { /* ... 与你提供的版本相同 ... */
	select {
	case <-clientGone:
		log.Println("SSE Service: 检测到客户端断开。")
		return true
	default:
		return false
	}
}
func (s *installerService) sendProgressUpdate(messageChan chan<- ProgressUpdate, step Step, overallProgress, stepProgress int, message string, rawLine string, clientGone <-chan struct{}) { /* ... 与你提供的版本相同 ... */
	if s.isClientGone(clientGone) {
		log.Println("SSE Service: 客户端已断开，不发送进度更新。")
		return
	}
	update := ProgressUpdate{Step: step, Progress: overallProgress, StepProgress: stepProgress, Message: message, Done: false, RawLine: rawLine}
	select {
	case messageChan <- update:
	default:
		log.Printf("警告: SSE 消息通道阻塞或前端未接收，跳过更新: Step=%s, Progress=%d", step, overallProgress)
	}
}
func (s *installerService) sendFinalUpdate(messageChan chan<- ProgressUpdate, step Step, overallProgress, stepProgress int, message string, isError bool, done bool) { /* ... 与你提供的版本相同 ... */
	log.Printf("尝试发送最终更新: Step=%s, Progress=%d, Error=%t, Done=%t, Message=%s", step, overallProgress, isError, done, message)
	update := ProgressUpdate{Step: step, Progress: overallProgress, StepProgress: stepProgress, Message: message, Done: done}
	if isError {
		update.Error = message
	}
	select {
	case messageChan <- update:
		log.Println("最终更新已发送到通道。")
	case <-time.After(1 * time.Second):
		log.Println("警告: 最终 SSE 更新发送超时 (通道阻塞或前端未接收)。")
	}
}
