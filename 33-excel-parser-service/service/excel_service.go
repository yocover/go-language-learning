package service

import (
	"encoding/json"
	"example/use/pkg/configs"
	"example/use/pkg/consts"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.neuxnet.com/ai/global-kit/net/resty"
	"go.uber.org/zap"
)

type ExcelResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Sheets []Sheet `json:"sheets"`
	} `json:"data"`
}

type Sheet struct {
	SheetName string  `json:"sheet_name"`
	Tables    []Table `json:"tables"`
}

type Table struct {
	Columns []string `json:"columns"`
	Values  [][]any  `json:"values"`
}

type ExcelService struct {
	llmService *LLMService
	baseURL    string
	apiPath    string
}

func NewExcelService() *ExcelService {
	return &ExcelService{
		llmService: GetLLMService(),
		baseURL:    configs.Cfg.ExcelParseConfig.BaseURL,
		apiPath:    configs.Cfg.ExcelParseConfig.ApiPath,
	}
}

func (s *ExcelService) ProcessExcel(localFilePath string) error {

	// MacOS/Linux 输出: /path/to/file
	// 将多个路径元素连接成一个完整的路径
	// path1 := filepath.Join("/path", "to", "file")

	processPath := filepath.Join(getOutPutPath(), consts.AUTO_DIR)
	fmt.Printf("processPath: %s\n", processPath)

	// 1. 创建处理
	if err := os.MkdirAll(processPath, 0755); err != nil {
		return fmt.Errorf("创建管理目录失败: %v", err)
	}
	fmt.Printf("开始处理文件: %s\n", localFilePath)

	// 2. 解析Excel数据
	excelData, err := s.parseExcelFile(localFilePath)
	if err != nil {
		return fmt.Errorf("解析Excel文件失败: %v", err)
	}

	// // 正确处理 json.Marshal 的返回值
	// jsonData, err := json.Marshal(excelData)
	// if err != nil {
	// 	return fmt.Errorf("转换JSON数据失败: %v", err)
	// }

	// fmt.Printf("JSON value of excelData: %s\n", string(jsonData))

	// 3. 生成分析结果
	summary, err := s.analyzeExcelData(excelData, localFilePath)
	if err != nil {
		return fmt.Errorf("分析Excel数据失败: %v", err)
	}

	// 4. 保存分析结果
	fileId := s.getNowTime()
	if err := s.saveResults(processPath, fileId, summary); err != nil {
		return fmt.Errorf("保存分析结果失败: %v", err)
	}

	// 5. 上传obs
	uploadMarkdownToHwObs(processPath)

	// 6. 删除本地文件
	// deleteLocalFile(processPath)

	// fmt.Printf("分析结果: %s\n", summary)

	return nil
}

func (s *ExcelService) getNowTime() int64 {
	return time.Now().Unix()
}

func (s *ExcelService) saveResults(processPath string, fileId int64, content string) error {
	// 保存带页码版本
	allPath := filepath.Join(processPath, fmt.Sprintf("%d_all.md", fileId))
	if err := writeFile(allPath, fmt.Sprintf("<!-- page:1 -->\n%s", content)); err != nil {
		return fmt.Errorf("保存带页码版本失败: %v", err)
	}

	// 保存无页码版本
	noPagePath := filepath.Join(processPath, fmt.Sprintf("%d_all_nopage.md", fileId))
	if err := writeFile(noPagePath, content); err != nil {
		return fmt.Errorf("保存无页码版本失败: %v", err)
	}
	return nil
}

// 构建prompt
func (s *ExcelService) buildAnlysisPrompt(filePath string, sheetName string, table Table) string {
	var b strings.Builder

	// 基本信息
	fmt.Fprintf(&b, "File: %s\nSheet: %s\nColumns: %s\n\nSample Data:\n",
		filepath.Base(filePath),
		sheetName,
		strings.Join(table.Columns, ", "))

	// 样本数据（最多5行）
	maxRows := math.Min(float64(len(table.Values)), 5)
	for i := range int(maxRows) {
		b.WriteString(fmt.Sprintf("\nRow %d:", i+1))
		for j, val := range table.Values[i] {
			b.WriteString(fmt.Sprintf(" %s: %v,", table.Columns[j], val))
		}
	}

	b.WriteString("\n\nPlease analyze this data and provide key characteristics.")

	return b.String()
}

// 分析Excel 数据
func (s *ExcelService) analyzeExcelData(excelData *ExcelResponse, filePath string) (string, error) {
	sheet := excelData.Data.Sheets[0]
	table := sheet.Tables[0]
	prompt := s.buildAnlysisPrompt(filePath, sheet.SheetName, table)

	// 调用大模型进行数据分析
	temperature := float64(0)

	defer s.llmService.Stop()

	response, err := s.llmService.ChatCompletion(
		configs.Cfg.LLMConfig.Model,
		[]ChatMessage{
			{Role: "system", Content: "You are a prefessional data analyst."},
			{Role: "user", Content: prompt},
		},
		temperature,
		configs.Cfg.LLMConfig.MaxToken,
		&ResponseFormat{Type: "text"},
	)

	if err != nil {
		return "", fmt.Errorf("调用大模型失败: %v", err)
	}

	return response.Choices[0].Message.Content, nil
}

// parseExcelFile 解析Excel文件
func (s *ExcelService) parseExcelFile(filePath string) (*ExcelResponse, error) {
	// 直接调用服务的方法
	return s.getExcelService(filePath)
}

// getExcelService  获取excel 解析后的数据
func (s *ExcelService) getExcelService(filePath string) (*ExcelResponse, error) {
	url := fmt.Sprintf("%s%s", s.baseURL, s.apiPath)

	file, err := os.Open(filePath)
	if err != nil {
		zap.L().Error("打开文件失败", zap.Error(err))
		return nil, err
	}
	defer file.Close()

	// 获取实际的文件名
	fileName := filepath.Base(filePath)

	resp, err := resty.File(url, nil, nil, "file", fileName, file)
	if err != nil {
		zap.L().Error("请求服务失败", zap.Error(err))
		return nil, err
	}

	var response ExcelResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		zap.L().Error("解析服务返回数据失败", zap.Error(err))
		return nil, err
	}

	return &response, nil
}
