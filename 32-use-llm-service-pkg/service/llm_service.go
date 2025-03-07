package service

import (
	"encoding/json"
	"example/use/pkg/llm-service/configs"
	"fmt"
	"sync"

	"git.neuxnet.com/ai/global-kit/net/resty"
	"go.uber.org/zap"
)

type LLMService struct {
	baseurl     string
	apiKey      string
	logger      *zap.Logger
	sugar       *zap.SugaredLogger
	concurrency map[string]chan struct{} // 每个model 的并发通道
	mu          sync.RWMutex             // 读写锁
}

var (
	deepInfraService *LLMService
	llmOnce          sync.Once
)

func InitLLMService() {
	GetLLMService() // 调用GetLLMService()方法，触发初始化流程
	zap.L().Info("LLM service initialized!")
}

func GetLLMService() *LLMService {
	llmOnce.Do(func() {
		config := zap.NewProductionConfig()
		logger, _ := config.Build()
		sugar := logger.Sugar()

		deepInfraService = &LLMService{
			baseurl:     configs.Cfg.LLMConfig.Url,
			apiKey:      configs.Cfg.LLMConfig.ApiKey,
			logger:      logger,
			sugar:       sugar,
			concurrency: make(map[string]chan struct{}),
		}
	})

	for model, concurrency := range configs.Cfg.LLMConfig.ModelConcurrency {
		if concurrency <= 0 {
			concurrency = configs.Cfg.LLMConfig.DefaultConcurrency
		}
		deepInfraService.concurrency[model] = make(chan struct{}, concurrency)
	}
	return deepInfraService
}

func (s *LLMService) getModelConcurrencyChannel(model string) chan struct{} {
	s.mu.RLock()
	if ch, exists := s.concurrency[model]; exists {
		s.mu.RUnlock()
		return ch
	}
	s.mu.RUnlock()

	// 如果模型不存在，创建新的并发控制通道
	s.mu.Lock()
	defer s.mu.Unlock()

	// 双重检查
	if ch, exists := s.concurrency[model]; exists {
		return ch
	}

	// 使用默认并发度
	defaultConcurrency := configs.Cfg.LLMConfig.DefaultConcurrency

	ch := make(chan struct{}, defaultConcurrency)
	s.concurrency[model] = ch
	return ch
}

// 提供给外部的调用的方法
func (s *LLMService) ChatCompletion(model string, messages []ChatMessage, temperature float64, maxTokens int, responseFormat *ResponseFormat) (*DeepInfraResponse, error) {

	// 获取模型特定的并发控制通道
	concurrencyChannel := s.getModelConcurrencyChannel(model)

	// 获取并发许可
	concurrencyChannel <- struct{}{}

	defer func() {
		<-concurrencyChannel // 释放并发许可
	}()

	zap.L().Info("Sending chat completion request", zap.String("model", model), zap.Any("messages", messages))

	if responseFormat == nil {
		responseFormat = &ResponseFormat{Type: "text"}
	}

	requestBody := DeepInfraRequest{
		Model:          model,
		Messages:       messages,
		Temperature:    temperature,
		MaxTokens:      maxTokens,
		ResponseFormat: *responseFormat,
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", s.apiKey),
	}

	body, err := resty.PostWithTimeOut(s.baseurl+"/chat/completions", requestBody, headers, 240)
	if err != nil {

		zap.L().Error("Error sending response body", zap.Error(err))
		return nil, fmt.Errorf("error sending response body: %w", err)
	}

	var response DeepInfraResponse
	if err := json.Unmarshal(body, &response); err != nil {
		zap.L().Error("Error unmarshaling response body", zap.Error(err), zap.Any("response", response))
		return nil, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	if response.Error != nil {
		zap.L().Error("API returned error", zap.String("message", response.Error.Message), zap.String("type", response.Error.Type), zap.String("param", response.Error.Param), zap.String("code", response.Error.Code), zap.Any("response", response))
		return nil, fmt.Errorf("API returned error: %s, type: %s, param: %s, code: %s", response.Error.Message, response.Error.Type, response.Error.Param, response.Error.Code)
	}

	zap.L().Info("Received chat completion response", zap.String("response", string(body)), zap.Any("request", requestBody))
	return &response, nil
}

// 停止服务
func (s *LLMService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for model, ch := range s.concurrency {
		close(ch)
		delete(s.concurrency, model)
	}
	zap.L().Info("LLM service stopped!")
}

func StopLLMService() {

	fmt.Println("Stopping LLM service...")
	if deepInfraService != nil {
		deepInfraService.Stop()
	}
}
