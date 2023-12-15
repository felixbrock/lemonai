package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/felixbrock/lemonai/internal/app"
	"github.com/felixbrock/lemonai/internal/component"
	"github.com/felixbrock/lemonai/internal/persistence"
)

/*
- Check for SQL injection
- Check for XSS
- Address prompt injection
- Block IP addresses: Too many requests
- Implement feedback logic
- Implement parent child tracing
- Remove prompts from code base (.gitignore file)
- update env base and add assistant ids
- end Goroutine after 1 minute
- Implement feedback handling
- Add hints to textfields

- Host (Jack)
- Add Posthog
*/

func config() (*app.Config, error) {

	env, err := os.ReadFile("env.json")
	if err != nil {
		return nil, err
	}

	var config app.Config
	if err := json.Unmarshal(env, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	config, err := config()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	componentBuilder := app.ComponentBuilder{
		Index:   component.Index,
		App:     component.App,
		Draft:   component.DraftModeEditor,
		Edit:    component.EditModeEditor,
		Review:  component.ReviewModeEditor,
		Loading: component.Loading,
		Error:   component.Error,
	}

	dbHeader := []string{
		fmt.Sprintf("apikey:%s", config.DBApiKey),
		fmt.Sprintf("Authorization:Bearer %s", config.DBApiKey)}
	dbUrlBase := "https://cllevlrokigwvbbnbfiu.supabase.co/rest/v1"

	optRepo := persistence.OptimizationRepo{BaseHeaders: dbHeader, BaseUrl: fmt.Sprintf("%s/optimization", dbUrlBase)}
	suggRepo := persistence.SuggestionRepo{BaseHeaders: dbHeader, BaseUrl: fmt.Sprintf("%s/suggestion", dbUrlBase)}
	runRepo := persistence.RunRepo{BaseHeaders: dbHeader, BaseUrl: fmt.Sprintf("%s/run", dbUrlBase)}
	oaiRepo := persistence.OpenAIRepo{BaseHeaders: []string{
		"Content-Type:application/json",
		"OpenAI-Beta:assistants=v1",
		fmt.Sprintf("Authorization: Bearer %s", config.OAIApiKey)}}

	repo := app.Repo{
		OpRepo:   optRepo,
		RunRepo:  runRepo,
		SuggRepo: suggRepo,
		OAIRepo:  oaiRepo,
	}

	a := app.App{
		Repo:             repo,
		ComponentBuilder: componentBuilder,
		Config:           *config,
	}

	a.Start()
}
