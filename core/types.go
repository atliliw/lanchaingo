// Package core provides the foundational types and interfaces for langchaingo.
//
// All components in the framework build upon the abstractions defined here:
// Runnable (unified execution interface), tool abstractions,
// language model interfaces, output parsers, and caching.
package core

import "time"

// ============================================================================
// Configuration
// ============================================================================

// RunnableConfig provides execution configuration for Runnable.Invoke/Batch/Stream.
// Maps to Rust langchainrust::core::runnables::RunnableConfig.
type RunnableConfig struct {
	Tags        []string
	RunName     string
	Metadata    map[string]any
	CallbackMgr any // reserved for CallbackManager
}

// NewRunnableConfig creates a RunnableConfig with defaults.
func NewRunnableConfig() *RunnableConfig {
	return &RunnableConfig{
		Tags:     []string{},
		Metadata: make(map[string]any),
	}
}

// WithTag adds a tag and returns the config for chaining.
func (c *RunnableConfig) WithTag(tag string) *RunnableConfig {
	c.Tags = append(c.Tags, tag)
	return c
}

// WithRunName sets the run name and returns the config for chaining.
func (c *RunnableConfig) WithRunName(name string) *RunnableConfig {
	c.RunName = name
	return c
}

// ============================================================================
// Language Model Types
// ============================================================================

// TokenUsage tracks LLM token consumption statistics.
// Maps to Rust langchainrust::core::language_models::chat::TokenUsage.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LLMResult is the standard response from a chat model invocation.
// Maps to Rust langchainrust::core::language_models::chat::LLMResult.
type LLMResult struct {
	Content    string      `json:"content"`
	Model      string      `json:"model"`
	TokenUsage *TokenUsage `json:"token_usage,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
}

// ============================================================================
// Tool Calling Types
// ============================================================================

// ToolCall represents a tool invocation request from the LLM.
// Maps to Rust langchainrust::core::tools::ToolCall.
type ToolCall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON-encoded arguments
}

// ToolDefinition defines a tool's schema for LLM function calling.
// Maps to Rust langchainrust::core::tools::ToolDefinition.
type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// FunctionDefinition defines a function for OpenAI-compatible function calling.
// Maps to Rust langchainrust::core::tools::FunctionDefinition.
type FunctionDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// FunctionCall represents an LLM function call response in OpenAI format.
// Maps to Rust langchainrust::core::tools::FunctionCall.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ============================================================================
// Agent Types
// ============================================================================

// AgentAction represents an action an agent should take (call a tool).
// Maps to Rust langchainrust::agents::AgentAction.
type AgentAction struct {
	Tool      string `json:"tool"`
	ToolInput string `json:"tool_input"`
	Log       string `json:"log"`
}

// AgentFinish represents the final output of an agent.
// Maps to Rust langchainrust::agents::AgentFinish.
type AgentFinish struct {
	Output string `json:"output"`
	Log    string `json:"log"`
}

// AgentStep is one execution step in the agent loop.
// Maps to Rust langchainrust::agents::AgentStep.
type AgentStep struct {
	Action AgentAction `json:"action"`
	Result string      `json:"result"`
}

// ============================================================================
// Document / Retrieval Types
// ============================================================================

// Document represents a document in the retrieval system.
// Maps to Rust langchainrust::vector_stores::Document.
type Document struct {
	PageContent string         `json:"page_content"`
	Metadata    map[string]any `json:"metadata"`
}

// NewDocument creates a Document with initialized Metadata map.
func NewDocument(content string) Document {
	return Document{
		PageContent: content,
		Metadata:    make(map[string]any),
	}
}

// SearchResult is a scored document from a retriever.
// Maps to Rust langchainrust::vector_stores::SearchResult.
type SearchResult struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"`
}

// ============================================================================
// Cache / Expiry
// ============================================================================

// CacheEntry wraps a value with an optional TTL expiry.
type CacheEntry struct {
	Value     any
	ExpiresAt *time.Time
}

// IsExpired returns true if the entry has exceeded its TTL.
func (e *CacheEntry) IsExpired() bool {
	if e.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*e.ExpiresAt)
}
