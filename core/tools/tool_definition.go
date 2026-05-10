package tools

import "github.com/atliliw/lanchaingo/core"

// ToolDefinition defines a tool's schema for LLM function calling.
// Re-exported from core for convenience.
type ToolDefinition = core.ToolDefinition

// ToolCall represents a tool invocation from the LLM.
// Re-exported from core for convenience.
type ToolCall = core.ToolCall

// FunctionDefinition defines a function in OpenAI format.
// Re-exported from core for convenience.
type FunctionDefinition = core.FunctionDefinition

// FunctionCall represents a function call in OpenAI format.
// Re-exported from core for convenience.
type FunctionCall = core.FunctionCall
