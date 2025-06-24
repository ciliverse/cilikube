package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

type InstallerHandler struct {
	installerService service.InstallerService
}

func NewInstallerHandler(is service.InstallerService) *InstallerHandler {
	return &InstallerHandler{
		installerService: is,
	}
}

// StreamMinikubeInstallation handles the SSE request.
func (h *InstallerHandler) StreamMinikubeInstallation(c *gin.Context) {
	// 设置 SSE 头部
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// CORS 由中间件处理

	// 刷新头部
	c.Writer.Flush()

	// 创建通道
	messageChan := make(chan service.ProgressUpdate) // 使用 service.ProgressUpdate

	// 获取客户端断开连接的通知 channel
	// 使用 c.Request.Context().Done() 是更现代和推荐的方式
	clientGone := c.Request.Context().Done() // 类型是 <-chan struct{}

	log.Println("SSE: 连接已建立，启动安装服务 Goroutine。")
	// 在新 Goroutine 中启动服务
	go h.installerService.InstallMinikube(messageChan, clientGone)
	// 将 clientGone (<-chan struct{}) 传递给服务

	log.Println("SSE: Handler 开始监听服务消息并推送到客户端...")
	// 在当前 Goroutine 中处理流，直到结束或出错
	err := h.streamUpdatesToClient(c, messageChan, clientGone) // 将 clientGone (<-chan struct{}) 传递给辅助函数
	if err != nil {
		log.Printf("SSE: 流处理错误: %v", err)
	}
	log.Println("SSE: Handler 流处理结束。")
}

// streamUpdatesToClient 辅助函数，处理从 service 发来的消息并推送到 client
func (h *InstallerHandler) streamUpdatesToClient(c *gin.Context, messageChan <-chan service.ProgressUpdate, clientGone <-chan struct{}) error {
	defer log.Println("SSE: streamUpdatesToClient 循环结束。")
	for {
		select {
		case <-clientGone: // 监听 Context.Done() channel
			log.Println("SSE: 客户端断开连接 (Context Done)。")
			return nil // 客户端断开，正常退出
		case update, ok := <-messageChan:
			if !ok {
				log.Println("SSE: 服务通道已关闭。")
				return nil // 服务完成或出错，正常退出循环
			}

			// 收到更新，准备发送
			log.Printf("SSE: 从服务收到更新: Step=%s, Progress=%d, Done=%t", update.Step, update.Progress, update.Done)

			jsonData, err := json.Marshal(update)
			if err != nil {
				log.Printf("SSE: 序列化服务更新失败: %v", err)
				// 尝试通知客户端
				_, writeErr := fmt.Fprintf(c.Writer, "event: error\ndata: {\"error\": \"Internal server error marshalling update: %v\"}\n\n", err)
				if writeErr != nil {
					log.Printf("SSE: 向客户端写入序列化错误信息时失败: %v", writeErr)
					return writeErr // 返回写入错误
				}
				c.Writer.Flush()
				continue // 继续监听下一条消息
			}

			// 发送数据
			_, writeErr := fmt.Fprintf(c.Writer, "event: message\ndata: %s\n\n", string(jsonData))
			if writeErr != nil {
				log.Printf("SSE: 向客户端写入数据失败: %v", writeErr)
				return writeErr // 返回写入错误
			}

			// 刷新确保发送
			if f, ok := c.Writer.(http.Flusher); ok {
				f.Flush()
			} else {
				log.Println("SSE: 警告 - ResponseWriter 不支持 Flushing。")
			}

			// 如果是最后一条消息，则退出
			if update.Done {
				log.Println("SSE: 已发送最终更新，正常退出流处理。")
				return nil
			}
		}
	}
}
