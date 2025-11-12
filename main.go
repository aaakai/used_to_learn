package main

import (
	"awesomeProject/voice"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化语音管理器
	vm := voice.NewVoiceManager()
	if err := vm.Init(); err != nil {
		panic(err)
	}
	defer vm.Close()

	// 创建HTTP服务器
	router := gin.Default()

	// 注册语音路由
	voiceHandler := voice.NewVoiceHandler(vm)
	voiceHandler.RegisterRoutes(router)

	// 启动服务器
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
