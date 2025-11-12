package voice

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
# 发起呼叫
curl -X POST http://localhost:8080/call -d '{"number":"123","timeout":30}'

# 接听
curl -X POST http://localhost:8080/answer

# 拒绝
curl -X POST http://localhost:8080/reject

# 查看状态
curl http://localhost:8080/status
*/

type VoiceHandler struct {
	manager *VoiceManager
}

func NewVoiceHandler(manager *VoiceManager) *VoiceHandler {
	return &VoiceHandler{manager: manager}
}

func (h *VoiceHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/call", h.call)
	router.POST("/answer", h.answer)
	router.POST("/reject", h.reject)
	router.GET("/status", h.status)
}

func (h *VoiceHandler) call(c *gin.Context) {
	var req struct {
		Number  string `json:"number"`
		Timeout int    `json:"timeout"` // 秒
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manager.Call(req.Number, time.Duration(req.Timeout)*time.Second); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": h.manager.GetState()})
}

func (h *VoiceHandler) answer(c *gin.Context) {
	if err := h.manager.Answer(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": h.manager.GetState()})
}

func (h *VoiceHandler) reject(c *gin.Context) {
	if err := h.manager.Reject(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": h.manager.GetState()})
}

func (h *VoiceHandler) status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": h.manager.GetState()})
}
