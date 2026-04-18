package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/tunsuy/synapse/internal/initializer"
	"github.com/tunsuy/synapse/pkg/extension"
)

// version 由构建时注入
var version = "dev"

func main() {
	app := &cli.App{
		Name:    "synapse",
		Usage:   "Personal Knowledge Hub — 个人知识中枢",
		Version: version,
		Description: `Synapse 从各种 AI 助手对话中自动沉淀、整理、反哺知识，
让你的每一次 AI 对话都成为知识复利。

架构：扩展点模型（Extension Point Model）
  - Source:    数据源（原始内容从哪来）
  - Processor: 处理引擎（原始内容 → 结构化知识）
  - Store:     存储底座（知识文件的持久化）
  - Indexer:   检索引擎（知识库检索）
  - Consumer:  消费端（知识输出为各种形式）
  - Auditor:   质量审计（知识库健康检查）`,
		Commands: []*cli.Command{
			initCommand(),
			pluginCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// initCommand 返回 init 子命令
func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "初始化 knowhub 知识库仓库",
		Description: `创建 knowhub 仓库的标准目录结构，包括：
  - .synapse/schema.yaml  知识规范
  - .synapse/config.yaml  扩展点配置
  - profile/me.md         用户画像
  - topics/               主题知识
  - entities/             实体页
  - concepts/             概念页
  - inbox/                增量缓冲区
  - journal/              时间线
  - graph/                知识图谱`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "知识库根目录路径（默认为当前目录）",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "知识库拥有者名称",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "覆盖已有文件",
			},
		},
		Action: func(c *cli.Context) error {
			return initializer.Init(initializer.Options{
				Path:  c.String("path"),
				Name:  c.String("name"),
				Force: c.Bool("force"),
			})
		},
	}
}

// pluginCommand 返回 plugin 子命令
func pluginCommand() *cli.Command {
	return &cli.Command{
		Name:    "plugin",
		Aliases: []string{"plugins"},
		Usage:   "管理扩展点插件",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "列出所有已注册的插件",
				Action: func(c *cli.Context) error {
					return pluginListAction()
				},
			},
		},
	}
}

// pluginListAction 列出所有已注册的插件
func pluginListAction() error {
	printSection("Sources", extension.ListSources())
	printSection("Processors", extension.ListProcessors())
	printSection("Stores", extension.ListStores())
	printSection("Indexers", extension.ListIndexers())
	printSection("Consumers", extension.ListConsumers())
	printSection("Auditors", extension.ListAuditors())
	return nil
}

// printSection 打印一个扩展点分类
func printSection(title string, items []string) {
	fmt.Printf("\n%s:\n", title)
	if len(items) == 0 {
		fmt.Println("  (none)")
		return
	}
	for _, item := range items {
		label := "built-in"
		fmt.Printf("  ✅ %-25s (%s)\n", item, label)
	}
}

// init 引入所有内置扩展点实现（通过 plugins.go）
// 以下是 main 包中的初始化说明，实际的 blank import 在 plugins.go 中
func init() {
	// 确保版本信息的格式
	if version == "" {
		version = "dev"
	}
	_ = strings.TrimSpace(version)
}
