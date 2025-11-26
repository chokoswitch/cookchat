// Copyright (c) CurioSwitch (choko@curioswitch.org)
// SPDX-License-Identifier: BUSL-1.1

package webrecipes

import (
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const prompt = `
You are a web searcher that finds cooking recipes tailored to particular requirements. 

# Requirements

- The source of a recipe must be in the same language as the user's.
- When returning recipe content such as the description and steps, generate the content yourself referring to the recipe. 
  Do not copy the text from the recipe as-is.
`

func NewAgent(model model.LLM) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Model:       model,
		Name:        "web_recipes_agent",
		Description: "Finds cooking recipes from the web based on a user's preferences.",
		Instruction: prompt,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		AfterAgentCallbacks: []agent.AfterAgentCallback{
			func(cc agent.CallbackContext) (*genai.Content, error) {
				cc.UserContent()
				fmt.Println(cc)
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("webrecipes: creating agent: %w", err)
	}
	return agent, nil
}
