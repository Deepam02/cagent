package session

import (
	"testing"

	"github.com/rumpl/cagent/pkg/agent"
	"github.com/rumpl/cagent/pkg/chat"
	"github.com/rumpl/cagent/pkg/tools"
	"github.com/stretchr/testify/assert"
)

func TestTrimMessages(t *testing.T) {
	// Create a slice of messages that exceeds maxMessages
	messages := make([]chat.Message, maxMessages+10)

	// Fill with some basic messages
	for i := range messages {
		messages[i] = chat.Message{
			Role:    "user",
			Content: "message",
		}
	}

	// Test basic trimming
	result := trimMessages(messages)
	assert.Equal(t, maxMessages, len(result), "should trim to maxMessages")
}

func TestTrimMessagesWithToolCalls(t *testing.T) {
	messages := []chat.Message{
		{
			Role:    "user",
			Content: "first message",
		},
		{
			Role:    "assistant",
			Content: "response with tool",
			ToolCalls: []tools.ToolCall{
				{
					ID: "tool1",
				},
			},
		},
		{
			Role:       "tool",
			Content:    "tool result",
			ToolCallID: "tool1",
		},
		{
			Role:    "user",
			Content: "second message",
		},
		{
			Role:    "assistant",
			Content: "another response",
			ToolCalls: []tools.ToolCall{
				{
					ID: "tool2",
				},
			},
		},
		{
			Role:       "tool",
			Content:    "tool result 2",
			ToolCallID: "tool2",
		},
	}

	// Set maxMessages to 3 to force trimming
	oldMax := maxMessages
	maxMessages = 3
	defer func() { maxMessages = oldMax }()

	result := trimMessages(messages)

	// Should keep last 3 messages, but ensure tool call consistency
	assert.Equal(t, 3, len(result), "should have exactly 3 messages")

	// Verify we don't have any orphaned tool results
	toolCalls := make(map[string]bool)
	for _, msg := range result {
		if msg.Role == "assistant" {
			for _, tool := range msg.ToolCalls {
				toolCalls[tool.ID] = true
			}
		}
		if msg.Role == "tool" {
			assert.True(t, toolCalls[msg.ToolCallID], "tool result should have corresponding assistant message")
		}
	}
}

func TestGetMessages(t *testing.T) {
	// Create a test agent
	testAgent := &agent.Agent{}

	// Create a session with many messages
	s := New(map[string]*agent.Agent{"test": testAgent})

	// Add more than maxMessages to the session
	for i := 0; i < maxMessages+10; i++ {
		s.Messages = append(s.Messages, AgentMessage{
			Agent: testAgent,
			Message: chat.Message{
				Role:    "user",
				Content: "test message",
			},
		})
	}

	// Get messages for the agent
	messages := s.GetMessages(testAgent)

	// Verify we get at most maxMessages
	assert.LessOrEqual(t, len(messages), maxMessages, "should not exceed maxMessages")
}

func TestGetMessagesWithToolCalls(t *testing.T) {
	// Create a test agent
	testAgent := &agent.Agent{}

	// Create a session
	s := New(map[string]*agent.Agent{"test": testAgent})

	// Add a sequence of messages with tool calls
	s.Messages = append(s.Messages, AgentMessage{
		Agent: testAgent,
		Message: chat.Message{
			Role:    "user",
			Content: "test message",
		},
	})

	s.Messages = append(s.Messages, AgentMessage{
		Agent: testAgent,
		Message: chat.Message{
			Role:    "assistant",
			Content: "using tool",
			ToolCalls: []tools.ToolCall{
				{
					ID: "test-tool",
				},
			},
		},
	})

	s.Messages = append(s.Messages, AgentMessage{
		Agent: testAgent,
		Message: chat.Message{
			Role:       "tool",
			ToolCallID: "test-tool",
			Content:    "tool result",
		},
	})

	// Set maxMessages to 2 to force trimming
	oldMax := maxMessages
	maxMessages = 2
	defer func() { maxMessages = oldMax }()

	messages := s.GetMessages(testAgent)

	// Verify tool call consistency
	toolCalls := make(map[string]bool)
	for _, msg := range messages {
		if msg.Role == "assistant" {
			for _, tool := range msg.ToolCalls {
				toolCalls[tool.ID] = true
			}
		}
		if msg.Role == "tool" {
			assert.True(t, toolCalls[msg.ToolCallID], "tool result should have corresponding assistant message")
		}
	}
}
