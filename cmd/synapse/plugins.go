package main

// 通过 blank import 激活所有内置扩展点实现
// 每个 package 的 init() 会自动注册到全局 Registry
import (
	// Source 实现
	_ "github.com/tunsuy/synapse/internal/source/skill"

	// Processor 实现
	_ "github.com/tunsuy/synapse/internal/processor/skill"

	// Store 实现
	_ "github.com/tunsuy/synapse/internal/store/github"
	_ "github.com/tunsuy/synapse/internal/store/local"
)
