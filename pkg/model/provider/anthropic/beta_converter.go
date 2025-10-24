package anthropic

import (
	"encoding/json"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/docker/cagent/pkg/chat"
	"github.com/docker/cagent/pkg/tools"
)

// convertBetaMessages converts chat messages to Anthropic Beta API format
// Following Anthropic's extended thinking documentation with interleaved thinking enabled:
// - Thinking blocks can appear anywhere in the conversation (not required to be first)
// - Always include the complete, unmodified thinking block from previous assistant turns
// - interleaved parameter is kept for API compatibility but always true
//
// Important: Anthropic API requires that all tool_result blocks corresponding to tool_use
// blocks from the same assistant message MUST be grouped into a single user message.
func convertBetaMessages(messages []chat.Message) []anthropic.BetaMessageParam {
	var betaMessages []anthropic.BetaMessageParam

	for i := 0; i < len(messages); i++ {
		msg := &messages[i]
		if msg.Role == chat.MessageRoleSystem {
			// System messages handled separately
			continue
		}
		if msg.Role == chat.MessageRoleUser {
			// Handle user messages (including images and tool results)
			if len(msg.MultiContent) > 0 {
				contentBlocks := make([]anthropic.BetaContentBlockParamUnion, 0, len(msg.MultiContent))
				for _, part := range msg.MultiContent {
					if part.Type == chat.MessagePartTypeText {
						if txt := strings.TrimSpace(part.Text); txt != "" {
							contentBlocks = append(contentBlocks, anthropic.BetaContentBlockParamUnion{
								OfText: &anthropic.BetaTextBlockParam{Text: txt},
							})
						}
					} else if part.Type == chat.MessagePartTypeImageURL && part.ImageURL != nil {
						if strings.HasPrefix(part.ImageURL.URL, "data:") {
							parts := strings.SplitN(part.ImageURL.URL, ",", 2)
							if len(parts) == 2 {
								mediaTypePart := parts[0]
								base64Data := parts[1]
								var mediaType string
								switch {
								case strings.Contains(mediaTypePart, "image/jpeg"):
									mediaType = "image/jpeg"
								case strings.Contains(mediaTypePart, "image/png"):
									mediaType = "image/png"
								case strings.Contains(mediaTypePart, "image/gif"):
									mediaType = "image/gif"
								case strings.Contains(mediaTypePart, "image/webp"):
									mediaType = "image/webp"
								default:
									mediaType = "image/jpeg"
								}
								imageBlockJSON := map[string]any{
									"type": "image",
									"source": map[string]any{
										"type":       "base64",
										"media_type": mediaType,
										"data":       base64Data,
									},
								}
								jsonBytes, err := json.Marshal(imageBlockJSON)
								if err == nil {
									var imageBlock anthropic.BetaContentBlockParamUnion
									if json.Unmarshal(jsonBytes, &imageBlock) == nil {
										contentBlocks = append(contentBlocks, imageBlock)
									}
								}
							}
						}
					}
				}
				if len(contentBlocks) > 0 {
					betaMessages = append(betaMessages, anthropic.BetaMessageParam{
						Role:    anthropic.BetaMessageParamRoleUser,
						Content: contentBlocks,
					})
				}
			} else if txt := strings.TrimSpace(msg.Content); txt != "" {
				betaMessages = append(betaMessages, anthropic.BetaMessageParam{
					Role: anthropic.BetaMessageParamRoleUser,
					Content: []anthropic.BetaContentBlockParamUnion{
						{OfText: &anthropic.BetaTextBlockParam{Text: txt}},
					},
				})
			}
			continue
		}
		if msg.Role == chat.MessageRoleAssistant {
			contentBlocks := make([]anthropic.BetaContentBlockParamUnion, 0)

			// With interleaved thinking, we can include thinking blocks anywhere
			// If we have thinking content, include it first (conventional order)
			if msg.ReasoningContent != "" && msg.ThinkingSignature != "" {
				contentBlocks = append(contentBlocks,
					anthropic.NewBetaThinkingBlock(msg.ThinkingSignature, msg.ReasoningContent))
			} else if msg.ThinkingSignature != "" {
				// Include redacted thinking placeholder using the original signature
				contentBlocks = append(contentBlocks,
					anthropic.NewBetaRedactedThinkingBlock(msg.ThinkingSignature))
			}

			// Add text content if present
			if txt := strings.TrimSpace(msg.Content); txt != "" {
				contentBlocks = append(contentBlocks, anthropic.BetaContentBlockParamUnion{
					OfText: &anthropic.BetaTextBlockParam{Text: txt},
				})
			}

			// Add tool calls
			if len(msg.ToolCalls) > 0 {
				for _, toolCall := range msg.ToolCalls {
					var inpts map[string]any
					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &inpts); err != nil {
						inpts = map[string]any{}
					}
					contentBlocks = append(contentBlocks, anthropic.BetaContentBlockParamUnion{
						OfToolUse: &anthropic.BetaToolUseBlockParam{
							ID:    toolCall.ID,
							Input: inpts,
							Name:  toolCall.Function.Name,
						},
					})
				}
			}

			if len(contentBlocks) > 0 {
				betaMessages = append(betaMessages, anthropic.BetaMessageParam{
					Role:    anthropic.BetaMessageParamRoleAssistant,
					Content: contentBlocks,
				})
			}
			continue
		}
		if msg.Role == chat.MessageRoleTool {
			// Collect consecutive tool messages and merge them into a single user message
			// This is required by Anthropic API: all tool_result blocks for tool_use blocks
			// from the same assistant message must be in the same user message
			toolResultBlocks := []anthropic.BetaContentBlockParamUnion{
				{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						ToolUseID: msg.ToolCallID,
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{OfText: &anthropic.BetaTextBlockParam{Text: strings.TrimSpace(msg.Content)}},
						},
					},
				},
			}

			// Look ahead for consecutive tool messages and merge them
			j := i + 1
			for j < len(messages) && messages[j].Role == chat.MessageRoleTool {
				toolResultBlocks = append(toolResultBlocks, anthropic.BetaContentBlockParamUnion{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						ToolUseID: messages[j].ToolCallID,
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{OfText: &anthropic.BetaTextBlockParam{Text: strings.TrimSpace(messages[j].Content)}},
						},
					},
				})
				j++
			}

			// Add the merged user message with all tool results
			betaMessages = append(betaMessages, anthropic.BetaMessageParam{
				Role:    anthropic.BetaMessageParamRoleUser,
				Content: toolResultBlocks,
			})

			// Skip the messages we've already processed
			i = j - 1
			continue
		}
	}
	return betaMessages
}

// extractBetaSystemBlocks extracts system messages for Beta API format
func extractBetaSystemBlocks(messages []chat.Message) []anthropic.BetaTextBlockParam {
	regularBlocks := extractSystemBlocks(messages)
	betaBlocks := make([]anthropic.BetaTextBlockParam, len(regularBlocks))
	for i, block := range regularBlocks {
		betaBlocks[i] = anthropic.BetaTextBlockParam{Text: block.Text}
	}
	return betaBlocks
}

// convertBetaTools converts tools to Beta API format
func convertBetaTools(t []tools.Tool) ([]anthropic.BetaToolUnionParam, error) {
	betaTools := make([]anthropic.BetaToolUnionParam, len(t))

	for i, tool := range t {
		inputSchema, err := ConvertParametersToSchema(tool.Parameters)
		if err != nil {
			return nil, err
		}

		// Convert to BetaToolInputSchemaParam
		var betaInputSchema anthropic.BetaToolInputSchemaParam
		if err := tools.ConvertSchema(inputSchema, &betaInputSchema); err != nil {
			return nil, err
		}

		// Create BetaToolParam and wrap it in BetaToolUnionParam
		betaTools[i] = anthropic.BetaToolUnionParam{
			OfTool: &anthropic.BetaToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: betaInputSchema,
			},
		}
	}

	return betaTools, nil
}
