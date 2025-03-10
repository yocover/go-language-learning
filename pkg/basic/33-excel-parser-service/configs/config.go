package configs

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var Cfg AppConfig
var once sync.Once

func InitConfig() {
	once.Do(func() {

		// TODO: 从配置文件中读取配置信息，并赋值给 cfg
		data, err := os.ReadFile("util/dev.yml")
		if err != nil {
			log.Fatal(err)
		}
		if err := yaml.Unmarshal(data, &Cfg); err != nil {
			log.Fatal(err)
		}
	})
}

type AppConfig struct {
	Log              ZapLogConfig     `yaml:"log"`
	LLMConfig        LLMConfig        `yaml:"llmConfig"`
	BusinessConfig   BusinessConfig   `yaml:"businessConfig"`
	ExcelParseConfig ExcelParseConfig `yaml:"excelParseConfig"`
}

type BusinessConfig struct {
	ObsBucket           string `yaml:"obsBucket"`
	ObsExpire           int64  `yaml:"obsExpire"`
	UploadPath          string `yaml:"uploadPath"`
	LocalFileRoot       string `yaml:"localFileRoot"`
	LocalOutputPath     string `yaml:"localOutputPath"`
	CondaPath           string `yaml:"condaPath"`
	CondaEnv            string `yaml:"condaEnv"`
	CondaShellPath      string `yaml:"condaShellPath"`
	SofficePath         string `yaml:"sofficePath"`
	OcrmypdfPath        string `yaml:"ocrmypdfPath"`
	OcrmypdfLanguage    string `yaml:"ocrmypdfLanguage"`
	PdfConvertWorkers   int    `yaml:"pdfConvertWorkers"`
	PdfConvertQueueSize int    `yaml:"pdfConvertQueueSize"`
	CudaVisibleDevices  int    `yaml:"cudaVisibleDevices"`
	FilePageSaveTimeout int    `yaml:"filePageSaveTimeout"`
}

type ZapLogConfig struct {
	FilePath   string `yaml:"file_path" json:"file_path"`
	Level      string `yaml:"level" json:"level"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
}

type LLMConfig struct {
	Model              string         `yaml:"model"`
	Url                string         `yaml:"url"`
	ApiKey             string         `yaml:"api_key"`
	ContextMaxToken    int            `yaml:"context_max_tokens"`
	MaxToken           int            `yaml:"max_token"`
	FixTitlePrompt     string         `yaml:"fix_title_prompt"`
	ModelConcurrency   map[string]int `yaml:"model_concurrency"`   // 每个模型的并发配置
	DefaultConcurrency int            `yaml:"default_concurrency"` // 默认的并发配置
}

type ExcelParseConfig struct {
	BaseURL string `yaml:"baseUrl"`
	ApiPath string `yaml:"apiPath"`
}
