package persistence

import (
	"encoding/json"
	"fmt"

	"github.com/felixbrock/lemonai/internal/domain"
)

type OptimizationRepo struct {
	BaseHeaders []string
	BaseUrl     string
}

func (r OptimizationRepo) Insert(optimization domain.Optimization) error {
	body, err := json.Marshal(optimization)

	if err != nil {
		return err
	}

	_, err = request[domain.Optimization](reqConfig{
		Method:  "POST",
		Url:     r.BaseUrl,
		Body:    body,
		Headers: r.BaseHeaders},
		201)

	if err != nil {
		return err
	}

	return nil
}

func (r OptimizationRepo) Update(id string, state string, optimizedPrompt []byte) error {
	body, err := json.Marshal(fmt.Sprintf(`{"state": "%s", "optimized_prompt": "%s"}`, state, optimizedPrompt))

	if err != nil {
		return err
	}

	_, err = request[domain.Optimization](reqConfig{
		Method:    "PATCH",
		Url:       r.BaseUrl,
		UrlParams: []string{fmt.Sprintf("id=eq.%s", id)},
		Body:      body,
		Headers:   r.BaseHeaders},
		201)

	if err != nil {
		return err
	}

	return nil
}

func (r OptimizationRepo) Read(id string) (*domain.Optimization, error) {
	record, err := request[domain.Optimization](reqConfig{
		Method:    "GET",
		Url:       r.BaseUrl,
		UrlParams: []string{fmt.Sprintf("id=eq.%s", id)},
		Body:      nil,
		Headers:   r.BaseHeaders},
		200)

	if err != nil {
		return nil, err
	}

	return record, nil
}
