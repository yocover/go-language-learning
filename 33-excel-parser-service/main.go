package main

import (
	"example/use/pkg/configs"
	"example/use/pkg/service"
)

func main() {
	configs.InitConfig()
	service.InitLLMService()

	// llmService := service.GetLLMService()
	// defer service.StopLLMService()
	// chatWithLLM(llmService)
	service.ParseFeature("files/Diagram_en.xlsx")
}
