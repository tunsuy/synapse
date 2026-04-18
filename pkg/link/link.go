// Package link 提供双向链接解析工具
// 支持 [[wiki-link]] 格式，兼容 Obsidian
package link

import (
	"regexp"
	"strings"
)

// linkPattern 匹配 [[双向链接]] 的正则表达式
// 支持 [[page-id]] 和 [[page-id|显示文本]] 两种格式
var linkPattern = regexp.MustCompile(`\[\[([^\]|]+?)(?:\|([^\]]+?))?\]\]`)

// Link 表示一个双向链接
type Link struct {
	// Target 链接目标（page-id）
	Target string

	// Display 显示文本（可选，为空时使用 Target）
	Display string

	// Position 在原文中的起始位置
	Position int
}

// DisplayText 返回链接的显示文本
func (l Link) DisplayText() string {
	if l.Display != "" {
		return l.Display
	}
	return l.Target
}

// Parse 从 Markdown 文本中解析出所有双向链接
func Parse(content string) []Link {
	matches := linkPattern.FindAllStringSubmatchIndex(content, -1)
	links := make([]Link, 0, len(matches))

	for _, match := range matches {
		target := content[match[2]:match[3]]
		link := Link{
			Target:   strings.TrimSpace(target),
			Position: match[0],
		}

		// 如果有显示文本（match[4] != -1 表示第二个捕获组匹配到了）
		if match[4] != -1 {
			link.Display = strings.TrimSpace(content[match[4]:match[5]])
		}

		links = append(links, link)
	}

	return links
}

// ExtractTargets 从 Markdown 文本中提取所有链接目标（去重）
func ExtractTargets(content string) []string {
	links := Parse(content)
	seen := make(map[string]bool, len(links))
	targets := make([]string, 0, len(links))

	for _, l := range links {
		if !seen[l.Target] {
			seen[l.Target] = true
			targets = append(targets, l.Target)
		}
	}

	return targets
}

// ReplaceLinks 将双向链接替换为指定格式
// formatter 接收 target 和 display，返回替换后的文本
func ReplaceLinks(content string, formatter func(target, display string) string) string {
	return linkPattern.ReplaceAllStringFunc(content, func(match string) string {
		sub := linkPattern.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		target := strings.TrimSpace(sub[1])
		display := ""
		if len(sub) > 2 {
			display = strings.TrimSpace(sub[2])
		}
		return formatter(target, display)
	})
}
