package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/tidwall/gjson"
	"os"
)

var authToken = ""

func init() {
	if authToken == "" {
		authToken = os.Getenv("OPENAI_API_KEY")
	}
}

func get_current_weather(location string, unit string) string {
	// 这里可以调用天气API获取天气信息
	// 随机返回一个温度
	return fmt.Sprintf("It's 20 %s in %s", unit, location)
}

func main() {
	fmt.Println(authToken)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "What's the weather like in Tokyo, Beijing",
		},
	}
	// 调用function-calling
	client := openai.NewClient(authToken)
	resp1, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
			Tools: []openai.Tool{
				{
					Type: openai.ToolTypeFunction,
					Function: &openai.FunctionDefinition{
						Name:        "get_current_weather",
						Description: "Get the current weather in a given location",
						Parameters: jsonschema.Definition{
							Type: jsonschema.Object,
							Properties: map[string]jsonschema.Definition{
								"location": {
									Type:        jsonschema.String,
									Description: "The city and state, e.g. San Francisco, CA",
								},
								"unit": {
									Type: jsonschema.String,
									Enum: []string{"celsius", "fahrenheit"},
								},
							},
							Required: []string{"location", "unit"},
						},
					},
				},
			},
			ToolChoice: "auto",
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	// res content: {"id":"chatcmpl-9Xz3YHAfOJ29hN1qq5zrIO3jRlzzs","object":"chat.completion","created":1717886552,"model":"gpt-3.5-turbo-0125","choices":[{"index":0,"message":{"role":"assistant","content":"","tool_calls":[{"id":"call_yDGvI8ua7vvMHKcvgizvUGoc","type":"function","function":{"name":"get_current_weather","arguments":"{\"location\": \"San Francisco\", \"unit\": \"celsius\"}"}},{"id":"call_kCAd0zauU8GR5jZPDHZ2tzHm","type":"function","function":{"name":"get_current_weather","arguments":"{\"location\": \"Tokyo\", \"unit\": \"celsius\"}"}},{"id":"call_kIcqwXieMj3a1zRf3TgN8WSj","type":"function","function":{"name":"get_current_weather","arguments":"{\"location\": \"Paris\", \"unit\": \"celsius\"}"}}]},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":78,"completion_tokens":77,"total_tokens":155},"system_fingerprint":""}
	toolCalls := resp1.Choices[0].Message.ToolCalls
	if len(toolCalls) > 0 {
		// 添加回复
		messages = append(messages, resp1.Choices[0].Message)
		for _, toolCall := range toolCalls {
			if toolCall.Function.Name == "get_current_weather" {
				location := gjson.Get(toolCall.Function.Arguments, "location").String()
				unit := gjson.Get(toolCall.Function.Arguments, "unit").String()
				// 调用function
				funcResponse := get_current_weather(location, unit)
				// 添加function的返回信息
				messages = append(messages, openai.ChatCompletionMessage{
					// role必须是tool
					Role:       openai.ChatMessageRoleTool,
					Content:    funcResponse,
					ToolCallID: toolCall.ID,
					Name:       toolCall.Function.Name,
				})
			}
		}
	}

	// 带着function调用结果再次call llm
	resp2, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Printf(resp2.Choices[0].Message.Content)
}
