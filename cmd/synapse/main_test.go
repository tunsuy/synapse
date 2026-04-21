package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

// TestEnsureGlobalConfigOnStart 测试全局配置初始化函数
func TestEnsureGlobalConfigOnStart(t *testing.T) {
	// 保存原始的 HOME 环境变量
	originalHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", originalHome)
	}()

	// 创建临时目录作为 HOME
	tempDir := t.TempDir()
	_ = os.Setenv("HOME", tempDir)

	// 测试首次调用应该创建配置文件
	t.Run("creates config on first run", func(t *testing.T) {
		// 捕获输出
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		ensureGlobalConfigOnStart()

		_ = w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		output := buf.String()

		// 验证配置文件已创建
		configPath := filepath.Join(tempDir, ".synapse", "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Errorf("expected config file to be created at %s", configPath)
		}

		// 验证输出包含创建提示
		if !strings.Contains(output, "Created global config template") {
			t.Errorf("expected output to contain creation message, got: %s", output)
		}
	})

	// 测试再次调用不应重新创建
	t.Run("does not recreate existing config", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		ensureGlobalConfigOnStart()

		_ = w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		output := buf.String()

		// 第二次调用不应该有创建提示
		if strings.Contains(output, "Created global config template") {
			t.Errorf("should not recreate config, got: %s", output)
		}
	})
}

// TestCollectCommand 测试 collect 命令的创建
func TestCollectCommand(t *testing.T) {
	cmd := collectCommand()

	if cmd.Name != "collect" {
		t.Errorf("expected command name 'collect', got %s", cmd.Name)
	}

	// 验证必要的 flags 存在
	expectedFlags := []string{"config", "content", "title", "topics", "entities", "concepts", "key-points", "source", "session-id"}
	flagNames := make(map[string]bool)
	for _, f := range cmd.Flags {
		for _, name := range f.Names() {
			flagNames[name] = true
		}
	}

	for _, expected := range expectedFlags {
		if !flagNames[expected] {
			t.Errorf("expected flag '%s' not found", expected)
		}
	}
}

// TestSearchCommand 测试 search 命令的创建
func TestSearchCommand(t *testing.T) {
	cmd := searchCommand()

	if cmd.Name != "search" {
		t.Errorf("expected command name 'search', got %s", cmd.Name)
	}

	// 验证必要的 flags 存在
	expectedFlags := []string{"config", "type", "limit"}
	flagNames := make(map[string]bool)
	for _, f := range cmd.Flags {
		for _, name := range f.Names() {
			flagNames[name] = true
		}
	}

	for _, expected := range expectedFlags {
		if !flagNames[expected] {
			t.Errorf("expected flag '%s' not found", expected)
		}
	}
}

// TestAuditCommand 测试 audit 命令的创建
func TestAuditCommand(t *testing.T) {
	cmd := auditCommand()

	if cmd.Name != "audit" {
		t.Errorf("expected command name 'audit', got %s", cmd.Name)
	}
}

// TestCollectActionNoContent 测试 collect 命令在没有内容时返回错误
func TestCollectActionNoContent(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{collectCommand()},
	}

	// 运行无内容的 collect 命令
	err := app.Run([]string{"synapse", "collect"})
	if err == nil {
		t.Error("expected error when no content provided")
	}

	if !strings.Contains(err.Error(), "no content provided") {
		t.Errorf("expected 'no content provided' error, got: %v", err)
	}
}

// TestConfigPathResolution 测试配置路径解析逻辑
func TestConfigPathResolution(t *testing.T) {
	t.Run("uses provided config path", func(t *testing.T) {
		// 创建一个模拟的 cli.Context
		app := &cli.App{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "config",
					Value: "/custom/path/config.yaml",
				},
			},
			Action: func(c *cli.Context) error {
				cfgPath := c.String("config")
				if cfgPath != "/custom/path/config.yaml" {
					t.Errorf("expected custom config path, got %s", cfgPath)
				}
				return nil
			},
		}

		_ = app.Run([]string{"app", "--config", "/custom/path/config.yaml"})
	})

	t.Run("uses default when config not provided", func(t *testing.T) {
		app := &cli.App{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "config",
				},
			},
			Action: func(c *cli.Context) error {
				cfgPath := c.String("config")
				if cfgPath != "" {
					t.Errorf("expected empty config path when not provided, got %s", cfgPath)
				}
				return nil
			},
		}

		_ = app.Run([]string{"app"})
	})
}

// TestInstallCommand 测试 install 命令的创建
func TestInstallCommand(t *testing.T) {
	cmd := installCommand()

	if cmd.Name != "install" {
		t.Errorf("expected command name 'install', got %s", cmd.Name)
	}

	// 验证必要的 flags 存在
	expectedFlags := []string{"target", "list"}
	flagNames := make(map[string]bool)
	for _, f := range cmd.Flags {
		for _, name := range f.Names() {
			flagNames[name] = true
		}
	}

	for _, expected := range expectedFlags {
		if !flagNames[expected] {
			t.Errorf("expected flag '%s' not found", expected)
		}
	}

	// 验证 Action 已设置
	if cmd.Action == nil {
		t.Error("expected Action to be set")
	}
}

// TestPluginCommand 测试 plugin 命令的创建
func TestPluginCommand(t *testing.T) {
	cmd := pluginCommand()

	if cmd.Name != "plugin" {
		t.Errorf("expected command name 'plugin', got %s", cmd.Name)
	}
}

// TestInitCommand 测试 init 命令的创建
func TestInitCommand(t *testing.T) {
	cmd := initCommand()

	if cmd.Name != "init" {
		t.Errorf("expected command name 'init', got %s", cmd.Name)
	}
}

// TestCheckCommand 测试 check 命令的创建
func TestCheckCommand(t *testing.T) {
	cmd := checkCommand()

	if cmd.Name != "check" {
		t.Errorf("expected command name 'check', got %s", cmd.Name)
	}
}
