package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/tunsuy/synapse/internal/config"
	"github.com/tunsuy/synapse/internal/engine"
	"github.com/tunsuy/synapse/internal/schema"
	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
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
		Before: func(c *cli.Context) error {
			// 自动确保全局配置文件存在
			cfgPath, created, err := config.EnsureGlobalConfig()
			if err != nil {
				// 非致命错误，仅告警
				log.Printf("WARN: ensure global config: %v", err)
				return nil
			}
			if created {
				fmt.Printf("📝 Created global config template: %s\n", cfgPath)
				fmt.Println("   Please edit this file to configure your store and extensions.")
				fmt.Println("   Then run 'synapse check' to verify your configuration.")
				fmt.Println()
			}
			return nil
		},
		Commands: []*cli.Command{
			initCommand(),
			checkCommand(),
			collectCommand(),
			searchCommand(),
			auditCommand(),
			installCommand(),
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
		Usage: "初始化知识库",
		Description: `根据全局配置（~/.synapse/config.yaml）中指定的 Store 后端初始化知识库。

不同 Store 的初始化行为不同：
  - local-store:  在本地创建知识库目录结构
  - github-store: 在 GitHub 仓库中创建知识库骨架

初始化前会自动检查：
  - 全局配置是否存在且有效
  - Store 中是否已有知识库（防止覆盖）

使用 --force 可强制重新初始化。

示例：
  synapse init                     # 使用全局配置初始化
  synapse init --name "张三"       # 指定知识库拥有者名称
  synapse init --config my.yaml    # 使用指定配置文件
  synapse init --force             # 强制重新初始化`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "配置文件路径（默认 ~/.synapse/config.yaml）",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "知识库拥有者名称",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "强制覆盖已有知识库",
			},
		},
		Action: initAction,
	}
}

// initAction 执行 init 命令
func initAction(c *cli.Context) error {
	// 1. 确定配置文件路径
	cfgPath := c.String("config")
	if cfgPath == "" {
		var err error
		cfgPath, err = config.GlobalConfigPath()
		if err != nil {
			return fmt.Errorf("get global config path: %w", err)
		}
	}

	// 2. 加载配置
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config %s: %w\nPlease edit your config file or run 'synapse check' to diagnose", cfgPath, err)
	}

	// 3. 创建 Store 实例
	store, err := extension.GetStore(cfg.Synapse.Store.Name, cfg.Synapse.Store.Config)
	if err != nil {
		return fmt.Errorf("create store %q: %w", cfg.Synapse.Store.Name, err)
	}

	ctx := context.Background()

	// 4. 检查是否已初始化
	if !c.Bool("force") {
		initialized, err := store.Initialized(ctx)
		if err != nil {
			return fmt.Errorf("check store status: %w", err)
		}
		if initialized {
			fmt.Printf("⚠️  Knowledge base already initialized in %s store.\n", store.Name())
			fmt.Println("   Use --force to reinitialize (existing data will NOT be deleted).")
			return nil
		}
	}

	// 5. 准备 schema 数据
	s := schema.Default()
	schemaData, err := s.MarshalYAML()
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}

	// 6. 执行初始化
	fmt.Printf("🧠 Initializing knowledge base via %s...\n\n", store.Name())

	opts := extension.InitOptions{
		Name:       c.String("name"),
		SchemaData: schemaData,
		SchemaObj:  s,
		Force:      c.Bool("force"),
	}

	if err := store.Init(ctx, opts); err != nil {
		return fmt.Errorf("init knowledge base: %w", err)
	}

	fmt.Println("\n✅ Knowledge base initialized successfully!")
	fmt.Println()
	fmt.Println("📁 Knowledge base structure:")
	fmt.Println("   ├── .synapse/schema.yaml  # 知识规范")
	for i, pt := range s.PageTypes {
		prefix := "├──"
		if i == len(s.PageTypes)-1 {
			prefix = "└──"
		}
		emoji := ""
		if pt.Emoji != "" {
			emoji = pt.Emoji + " "
		}
		fmt.Printf("   %s %s%s\n", prefix, emoji, pt.Description)
	}
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Println("   1. Edit profile/me.md to describe yourself")
	fmt.Println("   2. Start accumulating knowledge with AI assistants!")
	fmt.Printf("   3. Run 'synapse install <assistant>' to install AI Skill\n")

	return nil
}

// checkCommand 返回 check 子命令
func checkCommand() *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "检查配置文件是否完整有效",
		Description: `检查全局配置文件（~/.synapse/config.yaml）的完整性和有效性。

检查项目：
  - 配置文件是否存在
  - 必填字段是否完整（store.name 等）
  - Store 扩展点是否已注册
  - Source、Processor 等扩展点是否已注册
  - 环境变量占位符是否已替换

示例：
  synapse check
  synapse check --config /path/to/config.yaml`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "配置文件路径（默认 ~/.synapse/config.yaml）",
			},
		},
		Action: checkAction,
	}
}

// checkAction 执行 check 命令
func checkAction(c *cli.Context) error {
	cfgPath := c.String("config")
	if cfgPath == "" {
		var err error
		cfgPath, err = config.GlobalConfigPath()
		if err != nil {
			return fmt.Errorf("get global config path: %w", err)
		}
	}

	fmt.Println("🔍 Checking Synapse configuration...")
	fmt.Printf("   Config: %s\n\n", cfgPath)

	issues := 0
	warnings := 0

	// 1. 检查配置文件是否存在
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		fmt.Println("   ❌ Config file not found")
		fmt.Printf("      Run any synapse command to auto-create the template,\n")
		fmt.Printf("      or manually create: %s\n", cfgPath)
		return fmt.Errorf("config file not found: %s", cfgPath)
	}
	fmt.Println("   ✅ Config file exists")

	// 2. 尝试加载配置
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Printf("   ❌ Config parse error: %v\n", err)
		return fmt.Errorf("invalid config: %w", err)
	}
	fmt.Println("   ✅ Config file is valid YAML")

	// 3. 检查版本
	if cfg.Synapse.Version == "" {
		fmt.Println("   ❌ Missing synapse.version")
		issues++
	} else {
		fmt.Printf("   ✅ Version: %s\n", cfg.Synapse.Version)
	}

	// 4. 检查 Store 配置
	storeName := cfg.Synapse.Store.Name
	if storeName == "" {
		fmt.Println("   ❌ Missing synapse.store.name (required)")
		issues++
	} else {
		fmt.Printf("   ✅ Store: %s\n", storeName)

		// 检查 Store 是否已注册
		storeNames := extension.ListStores()
		found := false
		for _, n := range storeNames {
			if n == storeName {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("   ❌ Store %q is not registered (available: %v)\n", storeName, storeNames)
			issues++
		} else {
			fmt.Printf("   ✅ Store %q is registered\n", storeName)
		}

		// 检查 Store 配置中是否有未替换的环境变量占位符
		for key, val := range cfg.Synapse.Store.Config {
			if s, ok := val.(string); ok && strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
				envName := s[2 : len(s)-1]
				if os.Getenv(envName) == "" {
					fmt.Printf("   ⚠️  Store config %q = %q — environment variable %s is not set\n", key, s, envName)
					warnings++
				} else {
					fmt.Printf("   ✅ Store config %q — env %s is set\n", key, envName)
				}
			}
		}
	}

	// 5. 检查 Source 配置
	if len(cfg.Synapse.Sources) == 0 {
		fmt.Println("   ⚠️  No sources configured (collect command won't work)")
		warnings++
	} else {
		for _, src := range cfg.Synapse.Sources {
			if src.IsEnabled() {
				srcNames := extension.ListSources()
				found := false
				for _, n := range srcNames {
					if n == src.Name {
						found = true
						break
					}
				}
				if found {
					fmt.Printf("   ✅ Source: %s (registered)\n", src.Name)
				} else {
					fmt.Printf("   ⚠️  Source %q is not registered (available: %v)\n", src.Name, srcNames)
					warnings++
				}
			}
		}
	}

	// 6. 检查 Processor 配置
	if cfg.Synapse.Processor == nil {
		fmt.Println("   ⚠️  No processor configured (collect command won't work)")
		warnings++
	} else {
		procNames := extension.ListProcessors()
		found := false
		for _, n := range procNames {
			if n == cfg.Synapse.Processor.Name {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("   ✅ Processor: %s (registered)\n", cfg.Synapse.Processor.Name)
		} else {
			fmt.Printf("   ⚠️  Processor %q is not registered (available: %v)\n", cfg.Synapse.Processor.Name, procNames)
			warnings++
		}
	}

	// 7. 汇总
	fmt.Println()
	if issues > 0 {
		fmt.Printf("❌ Found %d error(s) and %d warning(s). Please fix the errors before using synapse.\n", issues, warnings)
		return fmt.Errorf("configuration has %d error(s)", issues)
	}
	if warnings > 0 {
		fmt.Printf("⚠️  Found %d warning(s). Configuration is usable but some features may not work.\n", warnings)
	} else {
		fmt.Println("✅ Configuration is valid! You can now run 'synapse init' to initialize your knowledge base.")
	}

	return nil
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

// collectCommand 返回 collect 子命令
func collectCommand() *cli.Command {
	return &cli.Command{
		Name:  "collect",
		Usage: "采集知识：将 AI 对话内容通过管道传入知识库",
		Description: `运行完整的 collect 管道：Source.Fetch → Processor.Process → Store.Write

支持两种输入方式：
  1. 通过 --content 参数直接传入内容
  2. 通过 stdin 管道输入内容（如 echo "content" | synapse collect）

示例：
  synapse collect --content "Go的并发模型基于goroutine" --title "Go并发" --topics "Go,并发"
  echo "学习笔记..." | synapse collect --title "今日学习" --topics "Go"
  cat notes.md | synapse collect --topics "分布式系统" --entities "Raft"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "配置文件路径（默认 .synapse/config.yaml）",
				Value:   ".synapse/config.yaml",
			},
			&cli.StringFlag{
				Name:  "content",
				Usage: "原始内容（也可通过 stdin 传入）",
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "内容标题",
			},
			&cli.StringFlag{
				Name:  "topics",
				Usage: "建议主题（逗号分隔）",
			},
			&cli.StringFlag{
				Name:  "entities",
				Usage: "建议实体（逗号分隔）",
			},
			&cli.StringFlag{
				Name:  "concepts",
				Usage: "建议概念（逗号分隔）",
			},
			&cli.StringFlag{
				Name:  "key-points",
				Usage: "关键知识点（逗号分隔）",
			},
			&cli.StringFlag{
				Name:  "source",
				Usage: "数据来源标识",
				Value: "codebuddy",
			},
			&cli.StringFlag{
				Name:  "session-id",
				Usage: "会话标识",
			},
		},
		Action: collectAction,
	}
}

// collectAction 执行 collect 命令
func collectAction(c *cli.Context) error {
	content := c.String("content")

	// 如果没有通过 --content 传入，尝试从 stdin 读取
	if content == "" {
		if hasStdinInput() {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			content = strings.TrimSpace(string(data))
		}
	}

	if content == "" {
		return fmt.Errorf("no content provided; use --content flag or pipe content via stdin")
	}

	// 构建 FetchOptions.Config
	fetchConfig := map[string]any{
		"content": content,
	}

	if title := c.String("title"); title != "" {
		fetchConfig["title"] = title
	}
	if topics := c.String("topics"); topics != "" {
		fetchConfig["suggested_topics"] = topics
	}
	if entities := c.String("entities"); entities != "" {
		fetchConfig["suggested_entities"] = entities
	}
	if concepts := c.String("concepts"); concepts != "" {
		fetchConfig["suggested_concepts"] = concepts
	}
	if keyPoints := c.String("key-points"); keyPoints != "" {
		fetchConfig["key_points"] = keyPoints
	}
	if source := c.String("source"); source != "" {
		fetchConfig["source"] = source
	}
	if sessionID := c.String("session-id"); sessionID != "" {
		fetchConfig["session_id"] = sessionID
	}

	// 创建引擎
	eng, err := engine.New(c.String("config"))
	if err != nil {
		return fmt.Errorf("create engine: %w", err)
	}

	ctx := context.Background()
	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: fetchConfig,
		},
	}

	fmt.Println("🧠 Collecting knowledge...")
	if err := eng.Collect(ctx, opts); err != nil {
		return fmt.Errorf("collect: %w", err)
	}

	fmt.Println("✅ Knowledge collected successfully!")
	return nil
}

// searchCommand 返回 search 子命令
func searchCommand() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "搜索知识库",
		Description: `在知识库中搜索匹配的知识文件

M2 阶段使用基于文件遍历 + 文本匹配的简单搜索。
后续 M3/M4 将通过 Indexer 扩展点支持向量检索等高级搜索。

示例：
  synapse search goroutine
  synapse search --type topic "并发模型"
  synapse search --limit 5 golang`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "配置文件路径（默认 .synapse/config.yaml）",
				Value:   ".synapse/config.yaml",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "按页面类型过滤（topic/entity/concept/inbox/journal）",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "最大返回数量",
				Value:   20,
			},
		},
		Action: searchAction,
	}
}

// searchAction 执行 search 命令
func searchAction(c *cli.Context) error {
	query := strings.Join(c.Args().Slice(), " ")
	if query == "" {
		return fmt.Errorf("search query is required; usage: synapse search <query>")
	}

	// 创建引擎
	eng, err := engine.New(c.String("config"))
	if err != nil {
		return fmt.Errorf("create engine: %w", err)
	}

	ctx := context.Background()
	store := eng.Store()

	// M2 简单搜索：遍历所有目录，文本匹配
	dirs := []string{"topics", "entities", "concepts", "inbox", "journal", "profile"}
	if t := c.String("type"); t != "" {
		dirs = typeToDir(t)
	}

	limit := c.Int("limit")
	queryLower := strings.ToLower(query)
	var results []searchHit
	total := 0

	for _, dir := range dirs {
		files, err := store.List(ctx, dir, model.ListOptions{Recursive: true})
		if err != nil {
			continue
		}

		for _, fi := range files {
			kf, err := store.Read(ctx, fi.Path)
			if err != nil {
				continue
			}

			score := matchScore(kf, queryLower)
			if score > 0 {
				total++
				if len(results) < limit {
					results = append(results, searchHit{
						file:  kf,
						info:  fi,
						score: score,
					})
				}
			}
		}
	}

	// 输出结果
	fmt.Printf("🔍 Search results for %q (%d found)\n\n", query, total)

	if len(results) == 0 {
		fmt.Println("  No matching knowledge files found.")
		return nil
	}

	for i, hit := range results {
		typeIcon := pageTypeIcon(string(hit.file.Frontmatter.Type))
		fmt.Printf("  %d. %s [%s] %s\n", i+1, typeIcon, hit.file.Frontmatter.Type, hit.file.Frontmatter.Title)
		fmt.Printf("     📁 %s\n", hit.info.Path)

		if len(hit.file.Frontmatter.Tags) > 0 {
			fmt.Printf("     🏷️  %s\n", strings.Join(hit.file.Frontmatter.Tags, ", "))
		}
		fmt.Println()
	}

	if total > limit {
		fmt.Printf("  ... and %d more (use --limit to show more)\n", total-limit)
	}

	return nil
}

// auditCommand 返回 audit 子命令
func auditCommand() *cli.Command {
	return &cli.Command{
		Name:  "audit",
		Usage: "审计知识库健康状态",
		Description: `执行知识库健康检查，报告问题和统计信息

M2 阶段的基础审计包括：
  - Frontmatter 完整性检查
  - 双向链接有效性检查（断链检测）
  - 孤儿页面检测
  - 知识库统计

示例：
  synapse audit
  synapse audit --config /path/to/.synapse/config.yaml`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "配置文件路径（默认 .synapse/config.yaml）",
				Value:   ".synapse/config.yaml",
			},
		},
		Action: auditAction,
	}
}

// auditAction 执行 audit 命令
func auditAction(c *cli.Context) error {
	// 创建引擎
	eng, err := engine.New(c.String("config"))
	if err != nil {
		return fmt.Errorf("create engine: %w", err)
	}

	ctx := context.Background()
	store := eng.Store()

	fmt.Println("🔍 Auditing knowledge base...")
	fmt.Println()

	// 收集所有知识文件
	dirs := []string{"profile", "topics", "entities", "concepts", "inbox", "journal"}
	var allFiles []model.KnowledgeFile
	filesByType := make(map[model.PageType]int)
	allPaths := make(map[string]bool)

	for _, dir := range dirs {
		files, err := store.List(ctx, dir, model.ListOptions{Recursive: true})
		if err != nil {
			continue
		}
		for _, fi := range files {
			kf, err := store.Read(ctx, fi.Path)
			if err != nil {
				continue
			}
			allFiles = append(allFiles, kf)
			filesByType[kf.Frontmatter.Type]++
			// 提取文件名（不含扩展名）作为可链接的标识
			base := strings.TrimSuffix(filepath.Base(fi.Path), filepath.Ext(fi.Path))
			allPaths[base] = true
			// 标题也作为可链接的标识
			allPaths[strings.ToLower(kf.Frontmatter.Title)] = true
		}
	}

	// 审计
	var issues []auditIssue
	totalLinks := 0
	brokenLinks := 0
	orphanCount := 0

	for _, kf := range allFiles {
		// 1. 检查 Frontmatter 完整性
		if kf.Frontmatter.Title == "" {
			issues = append(issues, auditIssue{
				severity: "error",
				path:     kf.Path,
				message:  "missing title in frontmatter",
			})
		}
		if kf.Frontmatter.Type == "" {
			issues = append(issues, auditIssue{
				severity: "error",
				path:     kf.Path,
				message:  "missing type in frontmatter",
			})
		}

		// 2. 检查双向链接有效性
		for _, link := range kf.Frontmatter.Links {
			totalLinks++
			linkLower := strings.ToLower(link)
			slug := toSimpleSlug(link)
			if !allPaths[linkLower] && !allPaths[slug] {
				brokenLinks++
				issues = append(issues, auditIssue{
					severity: "warning",
					path:     kf.Path,
					message:  fmt.Sprintf("broken link: [[%s]]", link),
				})
			}
		}

		// 3. 检查是否为孤儿页面（没有被任何其他页面链接到）
		if kf.Frontmatter.Type != model.PageTypeProfile && kf.Frontmatter.Type != model.PageTypeInbox {
			titleLower := strings.ToLower(kf.Frontmatter.Title)
			isLinked := false
			for _, other := range allFiles {
				if other.Path == kf.Path {
					continue
				}
				for _, link := range other.Frontmatter.Links {
					if strings.EqualFold(link, titleLower) {
						isLinked = true
						break
					}
				}
				if isLinked {
					break
				}
			}
			if !isLinked && len(allFiles) > 1 {
				orphanCount++
				issues = append(issues, auditIssue{
					severity: "info",
					path:     kf.Path,
					message:  "orphan page: not linked from any other page",
				})
			}
		}
	}

	// 计算健康评分
	score := 100
	for _, issue := range issues {
		switch issue.severity {
		case "error":
			score -= 10
		case "warning":
			score -= 5
		case "info":
			score--
		}
	}
	if score < 0 {
		score = 0
	}

	// 输出报告
	fmt.Printf("📊 Knowledge Base Health: %d/100\n\n", score)

	fmt.Println("📁 Statistics:")
	fmt.Printf("   Total files:   %d\n", len(allFiles))
	for pt, count := range filesByType {
		fmt.Printf("   %-14s %d\n", string(pt)+":", count)
	}
	fmt.Printf("   Total links:   %d\n", totalLinks)
	fmt.Printf("   Broken links:  %d\n", brokenLinks)
	fmt.Printf("   Orphan pages:  %d\n", orphanCount)
	fmt.Println()

	if len(issues) == 0 {
		fmt.Println("✅ No issues found. Knowledge base is healthy!")
	} else {
		fmt.Printf("⚠️  Found %d issue(s):\n\n", len(issues))
		for _, issue := range issues {
			icon := "ℹ️"
			switch issue.severity {
			case "error":
				icon = "❌"
			case "warning":
				icon = "⚠️"
			}
			fmt.Printf("  %s [%s] %s\n     %s\n\n", icon, issue.severity, issue.path, issue.message)
		}
	}

	return nil
}

// installCommand 返回 install 子命令
func installCommand() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "安装 Skill 到 AI 助手的配置目录",
		Description: `将 Synapse Skill 文件安装到目标 AI 助手的用户级配置目录，
让 AI 助手在对话中自动帮你采集、整理、反哺知识。

每个 AI 助手在 skills/<name>/adapter.yaml 中定义差异配置，
新增助手只需添加目录 + adapter.yaml，不改任何 Go 代码。

默认安装到用户主目录下（全局生效），所有项目共享。
使用 --target 可指定安装到特定项目目录（仅该项目生效）。

示例：
  synapse install codebuddy              # 安装到用户主目录（全局）
  synapse install claude                 # 安装到用户主目录（全局）
  synapse install cursor                 # 安装到用户主目录（全局）
  synapse install codebuddy --target .   # 安装到当前项目目录
  synapse install --list                 # 查看所有支持的 AI 助手`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "目标目录（默认用户主目录，全局生效）",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "列出所有支持的 AI 助手",
			},
		},
		Action: installAction,
	}
}

// installAction 执行 install 命令
func installAction(c *cli.Context) error {
	if c.Bool("list") {
		adapters, err := listAdapters()
		if err != nil {
			// 如果 adapter 配置不可用，显示默认列表（这是预期行为，不是错误）
			fmt.Println("📋 Supported AI assistants:")
			fmt.Println()
			fmt.Println("  codebuddy  — CodeBuddy")
			fmt.Println("  claude     — Claude Code")
			fmt.Println("  cursor     — Cursor")
			fmt.Println("\n使用方法：synapse install <assistant>")
			return nil //nolint:nilerr // intentionally return nil as fallback list is shown
		}

		fmt.Println("📋 Supported AI assistants:")
		fmt.Println()
		for name, adapter := range adapters {
			fmt.Printf("  %-12s— 安装到 ~/%s\n", name, adapter.DestPath)
		}
		fmt.Println("\n使用方法：synapse install <assistant>")
		fmt.Println("         synapse install <assistant> --target .  # 安装到当前项目目录")
		return nil
	}

	assistant := c.Args().First()
	if assistant == "" {
		return fmt.Errorf("请指定 AI 助手名称；运行 synapse install --list 查看支持列表")
	}

	// 默认安装到用户主目录（全局生效）
	targetDir := c.String("target")
	if targetDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get user home dir: %w", err)
		}
		targetDir = home
	}

	// 标准化助手名称：支持别名
	adapterKey := normalizeAssistantName(strings.ToLower(assistant))

	return installSkill(adapterKey, targetDir)
}

// normalizeAssistantName 将用户输入的助手名称映射到 adapter key
// 支持别名，如 "claude-code" → "claude"
func normalizeAssistantName(name string) string {
	aliases := map[string]string{
		"claude-code": "claude",
	}
	if mapped, ok := aliases[name]; ok {
		return mapped
	}
	return name
}

// --- Skill 安装核心逻辑 ---

// skillMetadata 定义所有 AI 助手共享的 Skill 公共元数据
// 从 skills/common/metadata.yaml 加载，避免每个平台重复维护
type skillMetadata struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Author      string   `yaml:"author"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Tips        []string `yaml:"tips"`
}

// skillAdapter defines the per-platform adapter configuration.
// Each assistant has its own skills/<name>/adapter.yaml.
type skillAdapter struct {
	Source      string   `yaml:"source"`
	DestPath    string   `yaml:"dest_path"`
	LegacyPaths []string `yaml:"legacy_paths"`
	Frontmatter string   `yaml:"frontmatter"`
	Tips        []string `yaml:"tips"`
}

// skillTemplateData 是传入 Go text/template 渲染的数据
type skillTemplateData struct {
	Frontmatter string
	Source      string
}

// installSkill 通用安装函数：读取通用模板 + 公共 metadata + adapter 配置 → 渲染 → 写入目标路径
func installSkill(assistantKey, targetDir string) error {
	// 1. 加载公共 metadata
	meta, err := loadMetadata()
	if err != nil {
		return fmt.Errorf("load common metadata: %w", err)
	}

	// 2. 加载 adapter 配置
	adapter, err := loadAdapter(assistantKey)
	if err != nil {
		return fmt.Errorf("load adapter config for %q: %w", assistantKey, err)
	}

	// 3. 将公共 metadata 注入 adapter 的 frontmatter 模板
	frontmatter, err := renderFrontmatter(adapter, meta)
	if err != nil {
		return fmt.Errorf("render frontmatter for %q: %w", assistantKey, err)
	}

	// 4. 加载通用模板
	tmplContent, err := findCommonTemplate("SKILL.md.tmpl")
	if err != nil {
		return fmt.Errorf("find common template: %w", err)
	}

	// 5. 渲染模板
	rendered, err := renderTemplate(tmplContent, skillTemplateData{
		Frontmatter: strings.TrimSpace(frontmatter),
		Source:      adapter.Source,
	})
	if err != nil {
		return fmt.Errorf("render template for %q: %w", assistantKey, err)
	}

	// 6. 写入目标文件
	destPath := filepath.Join(targetDir, adapter.DestPath)
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("create dir %s: %w", destDir, err)
	}

	if err := os.WriteFile(destPath, rendered, 0644); err != nil {
		return fmt.Errorf("write skill file: %w", err)
	}

	// 7. 清理旧版本文件
	for _, legacy := range adapter.LegacyPaths {
		oldPath := filepath.Join(targetDir, legacy)
		if _, statErr := os.Stat(oldPath); statErr == nil {
			_ = os.Remove(oldPath)
			fmt.Printf("🧹 已清理旧版 Skill 文件: %s\n", oldPath)
		}
	}

	// 8. 输出安装提示
	fmt.Printf("✅ Skill 已安装到: %s\n", destPath)
	tips := mergedTips(adapter, meta)
	if len(tips) > 0 {
		fmt.Println("\n📖 使用方法：")
		for i, tip := range tips {
			if i < 2 {
				fmt.Printf("   %s\n", tip)
			}
		}
		if len(tips) > 2 {
			fmt.Println("\n💡 提示：")
			for _, tip := range tips[2:] {
				fmt.Printf("   - %s\n", tip)
			}
		}
	}

	return nil
}

// loadAdapter 加载指定助手的 adapter 配置
// 从 skills/<assistantKey>/adapter.yaml 读取
func loadAdapter(assistantKey string) (skillAdapter, error) {
	data, err := findSkillFile(assistantKey, "adapter.yaml")
	if err != nil {
		// adapter.yaml 找不到时，列出可用的助手
		available, listErr := listAdapterNames()
		if listErr != nil {
			return skillAdapter{}, fmt.Errorf("adapter %q not found: %w", assistantKey, err)
		}
		return skillAdapter{}, fmt.Errorf("adapter %q not found; available: %v", assistantKey, available)
	}

	var adapter skillAdapter
	if err := yaml.Unmarshal(data, &adapter); err != nil {
		return skillAdapter{}, fmt.Errorf("parse adapter.yaml for %q: %w", assistantKey, err)
	}

	return adapter, nil
}

// loadMetadata 加载公共 Skill 元数据
// 从 skills/common/metadata.yaml 读取
func loadMetadata() (skillMetadata, error) {
	data, err := findCommonTemplate("metadata.yaml")
	if err != nil {
		return skillMetadata{}, fmt.Errorf("find common metadata.yaml: %w", err)
	}

	var meta skillMetadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return skillMetadata{}, fmt.Errorf("parse metadata.yaml: %w", err)
	}

	return meta, nil
}

// renderFrontmatter injects common metadata into the adapter's frontmatter template.
// adapter.Frontmatter can use Go text/template variables:
//
//	{{.Name}}, {{.Version}}, {{.Author}}, {{.Description}}, {{.Tags}}
func renderFrontmatter(adapter skillAdapter, meta skillMetadata) (string, error) {
	frontmatterTmpl := adapter.Frontmatter

	// 格式化 tags 为 YAML 数组
	tagsStr := "["
	for i, tag := range meta.Tags {
		if i > 0 {
			tagsStr += ", "
		}
		tagsStr += tag
	}
	tagsStr += "]"

	// 构建模板数据
	type frontmatterData struct {
		Name        string
		Version     string
		Author      string
		Description string
		Tags        string
	}

	data := frontmatterData{
		Name:        meta.Name,
		Version:     meta.Version,
		Author:      meta.Author,
		Description: meta.Description,
		Tags:        tagsStr,
	}

	tmpl, err := template.New("frontmatter").Parse(frontmatterTmpl)
	if err != nil {
		return "", fmt.Errorf("parse frontmatter template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute frontmatter template: %w", err)
	}

	return buf.String(), nil
}

// mergedTips combines adapter-specific tips with common metadata tips.
// Adapter tips (e.g. install-path notes) come first, then common tips from metadata.
func mergedTips(adapter skillAdapter, meta skillMetadata) []string {
	var tips []string
	tips = append(tips, adapter.Tips...)
	tips = append(tips, meta.Tips...)
	return tips
}

// renderTemplate 使用 Go text/template 渲染模板内容
func renderTemplate(tmplContent []byte, data skillTemplateData) ([]byte, error) {
	tmpl, err := template.New("skill").Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// findCommonTemplate 查找通用模板文件（从 skills/common/ 目录）
// 搜索顺序：可执行文件目录 → 当前工作目录 → 用户主目录
func findCommonTemplate(filename string) ([]byte, error) {
	var searchPaths []string

	// 1. 可执行文件所在目录的 skills/common/
	if exe, err := os.Executable(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(filepath.Dir(exe), "skills", "common", filename))
	}

	// 2. 当前工作目录的 skills/common/
	if cwd, err := os.Getwd(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(cwd, "skills", "common", filename))
	}

	// 3. 用户主目录的 .synapse/skills/common/
	if home, err := os.UserHomeDir(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(home, ".synapse", "skills", "common", filename))
	}

	for _, p := range searchPaths {
		if data, err := os.ReadFile(p); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("common template %q not found; searched: %v", filename, searchPaths)
}

// listAdapters 列出所有可用的 adapter（用于 install --list）
// 自动扫描 skills/ 目录下含有 adapter.yaml 的子目录
func listAdapters() (map[string]skillAdapter, error) {
	names, err := listAdapterNames()
	if err != nil {
		return nil, err
	}

	result := make(map[string]skillAdapter, len(names))
	for _, name := range names {
		adapter, err := loadAdapter(name)
		if err != nil {
			continue
		}
		result[name] = adapter
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no adapters found")
	}

	return result, nil
}

// listAdapterNames 扫描 skills/ 目录，返回所有含有 adapter.yaml 的子目录名
func listAdapterNames() ([]string, error) {
	skillsDirs := findSkillsDirs()
	if len(skillsDirs) == 0 {
		return nil, fmt.Errorf("skills directory not found")
	}

	seen := make(map[string]bool)
	var names []string

	for _, skillsDir := range skillsDirs {
		entries, err := os.ReadDir(skillsDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() == "common" {
				continue
			}
			name := entry.Name()
			if seen[name] {
				continue
			}
			// 检查该目录下是否有 adapter.yaml
			adapterPath := filepath.Join(skillsDir, name, "adapter.yaml")
			if _, err := os.Stat(adapterPath); err == nil {
				seen[name] = true
				names = append(names, name)
			}
		}
	}

	return names, nil
}

// findSkillsDirs 返回所有可能的 skills/ 根目录（搜索顺序一致）
func findSkillsDirs() []string {
	var dirs []string

	// 1. 可执行文件所在目录的 skills/
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Join(filepath.Dir(exe), "skills")
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			dirs = append(dirs, dir)
		}
	}

	// 2. 当前工作目录的 skills/
	if cwd, err := os.Getwd(); err == nil {
		dir := filepath.Join(cwd, "skills")
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			dirs = append(dirs, dir)
		}
	}

	// 3. 用户主目录的 .synapse/skills/
	if home, err := os.UserHomeDir(); err == nil {
		dir := filepath.Join(home, ".synapse", "skills")
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			dirs = append(dirs, dir)
		}
	}

	return dirs
}

// findSkillFile 查找特定助手的文件（从 skills/<assistant>/ 目录）
// 搜索顺序：可执行文件目录 → 当前工作目录 → 用户主目录
func findSkillFile(assistant, filename string) ([]byte, error) {
	var searchPaths []string

	// 1. 可执行文件所在目录的 skills/<assistant>/
	if exe, err := os.Executable(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(filepath.Dir(exe), "skills", assistant, filename))
	}

	// 2. 当前工作目录的 skills/<assistant>/
	if cwd, err := os.Getwd(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(cwd, "skills", assistant, filename))
	}

	// 3. 用户主目录的 .synapse/skills/<assistant>/
	if home, err := os.UserHomeDir(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(home, ".synapse", "skills", assistant, filename))
	}

	for _, p := range searchPaths {
		if data, err := os.ReadFile(p); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("file %s/%s not found; searched: %v", assistant, filename, searchPaths)
}

// --- 辅助类型和函数 ---

// searchHit 搜索命中结果
type searchHit struct {
	file  model.KnowledgeFile
	info  model.FileInfo
	score float64
}

// auditIssue 审计问题
type auditIssue struct {
	severity string
	path     string
	message  string
}

// hasStdinInput 检查 stdin 是否有管道输入
func hasStdinInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// matchScore 计算知识文件与查询的匹配分数
func matchScore(kf model.KnowledgeFile, queryLower string) float64 {
	score := 0.0

	// 标题匹配（权重最高）
	if strings.Contains(strings.ToLower(kf.Frontmatter.Title), queryLower) {
		score += 3.0
	}

	// 标签匹配
	for _, tag := range kf.Frontmatter.Tags {
		if strings.Contains(strings.ToLower(tag), queryLower) {
			score += 2.0
			break
		}
	}

	// 正文匹配
	if strings.Contains(strings.ToLower(kf.Body), queryLower) {
		score += 1.0
	}

	return score
}

// typeToDir 将页面类型映射到目录
func typeToDir(t string) []string {
	switch t {
	case "topic":
		return []string{"topics"}
	case "entity":
		return []string{"entities"}
	case "concept":
		return []string{"concepts"}
	case "inbox":
		return []string{"inbox"}
	case "journal":
		return []string{"journal"}
	case "profile":
		return []string{"profile"}
	default:
		return []string{"topics", "entities", "concepts", "inbox", "journal", "profile"}
	}
}

// pageTypeIcon 返回页面类型的 emoji 图标
func pageTypeIcon(t string) string {
	switch t {
	case "topic":
		return "📚"
	case "entity":
		return "🏷️"
	case "concept":
		return "💡"
	case "inbox":
		return "📥"
	case "journal":
		return "📅"
	case "profile":
		return "👤"
	default:
		return "📄"
	}
}

// toSimpleSlug 简单的标题转 slug
func toSimpleSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
