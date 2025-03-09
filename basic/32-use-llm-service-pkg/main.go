package main

import (
	"bufio"
	"example/use/pkg/llm-service/configs"
	"example/use/pkg/llm-service/service"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func chatWithLLM(llmService *service.LLMService) {
	reader := bufio.NewReader(os.Stdin)
	messages := []service.ChatMessage{
		{Role: "system", Content: "你是一个有帮助的助手"},
	}

	fmt.Println("开始和 AI 对话（输入 'quit' 退出）:")
	for {
		fmt.Print("\n你: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取输入错误: %v", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "quit" {
			fmt.Println("对话结束")
			break
		}

		messages = append(messages, service.ChatMessage{Role: "user", Content: input})

		// 记录开始时间
		startTime := time.Now()

		response, err := llmService.ChatCompletion(
			configs.Cfg.LLMConfig.Model,
			messages,
			0,
			configs.Cfg.LLMConfig.MaxToken,
			&service.ResponseFormat{Type: "text"},
		)

		// 计算响应时间
		duration := time.Since(startTime)

		if err != nil {
			log.Printf("LLM 响应错误: %v", err)
			continue
		}

		aiResponse := response.Choices[0].Message.Content
		fmt.Printf("\nAI: %s\n", aiResponse)
		fmt.Printf("\n[响应时间: %.2f 秒]\n", duration.Seconds())

		messages = append(messages, service.ChatMessage{Role: "assistant", Content: aiResponse})
	}
}

func main() {
	configs.InitConfig()
	service.InitLLMService()

	llmService := service.GetLLMService()
	defer service.StopLLMService()

	chatWithLLM(llmService)
}
