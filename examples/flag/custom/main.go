package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// StringList 是一个自定义的字符串列表类型
type StringList []string

// String 实现 flag.Value 接口的 String 方法
func (s *StringList) String() string {
	return strings.Join(*s, ",")
}

// Set 实现 flag.Value 接口的 Set 方法
func (s *StringList) Set(value string) error {
	if value == "" {
		return errors.New("值不能为空")
	}
	*s = append(*s, value)
	return nil
}

// 定义一个枚举类型
type LogLevel string

const (
	Debug   LogLevel = "debug"
	Info    LogLevel = "info"
	Warning LogLevel = "warning"
	Error   LogLevel = "error"
)

// String 实现 flag.Value 接口的 String 方法
func (l *LogLevel) String() string {
	return string(*l)
}

// Set 实现 flag.Value 接口的 Set 方法
func (l *LogLevel) Set(value string) error {
	switch LogLevel(value) {
	case Debug, Info, Warning, Error:
		*l = LogLevel(value)
		return nil
	default:
		return fmt.Errorf("无效的日志级别: %s, 必须是 debug/info/warning/error 之一", value)
	}
}

func main() {
	// 自定义字符串列表类型
	var tags StringList
	flag.Var(&tags, "tag", "添加标签 (可多次使用)")

	// 自定义枚举类型
	level := LogLevel(Info) // 默认值为 info
	flag.Var(&level, "level", "设置日志级别 (debug/info/warning/error)")

	// 解析命令行参数
	flag.Parse()

	// 输出结果
	fmt.Println("标签列表:", tags)
	fmt.Println("日志级别:", level)

	// 使用示例
	fmt.Println("\n使用示例:")
	fmt.Println("  ./custom -tag=golang -tag=programming")
	fmt.Println("  ./custom -level=debug")
}
