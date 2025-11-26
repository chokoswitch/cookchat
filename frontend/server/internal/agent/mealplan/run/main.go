// Copyright (c) CurioSwitch (choko@curioswitch.org)
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"

	"github.com/curioswitch/cookchat/frontend/server/internal/agent/mealplan"
	"github.com/curioswitch/cookchat/frontend/server/internal/agent/webrecipes"
)

func doMain() error {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return fmt.Errorf("mealplan:run: creating model: %w", err)
	}

	webRecipesAgent, err := webrecipes.NewAgent(model)
	if err != nil {
		return fmt.Errorf("mealplan:run: creating web recipes agent: %w", err)
	}
	mealPlanAgent, err := mealplan.NewAgent(model, webRecipesAgent)
	if err != nil {
		return fmt.Errorf("mealplan:run: creating meal plan agent: %w", err)
	}

	lConfig := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(mealPlanAgent),
	}
	l := full.NewLauncher()
	if err := l.Execute(ctx, lConfig, os.Args[1:]); err != nil {
		return fmt.Errorf("mealplan:run: executing launcher: %w", err)
	}
	return nil
}

func main() {
	if err := doMain(); err != nil {
		panic(err)
	}
}
