# Go 语言 flag 库使用示例

本目录包含了 Go 语言标准库 `flag` 包的各种使用示例，帮助您理解如何在 Go 程序中处理命令行参数。

## 示例列表

1. **基本用法** (`basic/`)
   - 展示了 flag 库的基本用法
   - 演示了如何定义和解析不同类型的命令行参数

2. **变量方式定义** (`var/`)
   - 展示了使用变量方式定义命令行参数
   - 使用 `flag.TypeVar()` 函数将参数绑定到变量

3. **短标志和长标志** (`shorthand/`)
   - 展示了如何同时支持短标志和长标志
   - 例如同时支持 `-n` 和 `-name`

4. **自定义 flag 类型** (`custom/`)
   - 展示了如何创建自定义的 flag 类型
   - 包含字符串列表和枚举类型的示例

5. **子命令和 FlagSet** (`subcommand/`)
   - 展示了如何使用 `flag.FlagSet` 实现子命令
   - 类似 `git commit`、`git push` 这样的命令结构

6. **自定义帮助信息** (`usage/`)
   - 展示了如何自定义 flag 的帮助信息
   - 通过修改 `flag.Usage` 函数实现

## 运行示例

每个示例都可以单独编译和运行。例如：

```bash
# 编译并运行基本示例
cd basic
go build -o basic
./basic -name=John -age=30 -married

# 运行带有子命令的示例
cd ../subcommand
go build -o subcommand
./subcommand server -port=9090 -mode=production
./subcommand client -server=example.com:9090 -debug
```

## flag 库的主要特点

1. **简单易用**：标准库提供的 flag 包使用简单，适合大多数命令行程序
2. **类型安全**：自动进行类型转换，避免手动解析参数的麻烦
3. **自定义类型**：可以通过实现 `flag.Value` 接口支持自定义类型
4. **子命令支持**：通过 `flag.FlagSet` 可以实现子命令功能
5. **自定义帮助**：可以自定义帮助信息的格式和内容

## 常见用法

- 使用 `-flag=value` 或 `-flag value` 形式传递参数
- 布尔类型参数可以使用 `-flag` 形式（默认为 true）
- 使用 `-h` 或 `-help` 查看帮助信息
- 使用 `flag.Parse()` 解析命令行参数
- 使用 `flag.Args()` 获取未解析的参数 