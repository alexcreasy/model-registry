package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"net/url"
)

const registerModelPath = "/registered_models"

type RegisteredModelInterface interface {
	GetAllRegisteredModels(client integrations.HTTPClientInterface) (*openapi.RegisteredModelList, error)
	CreateRegisteredModel(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.RegisteredModel, error)
	GetRegisteredModel(client integrations.HTTPClientInterface, id string) (*openapi.RegisteredModel, error)
}

type RegisteredModel struct {
	RegisteredModelInterface
}

func (m RegisteredModel) GetAllRegisteredModels(client integrations.HTTPClientInterface) (*openapi.RegisteredModelList, error) {

	responseData, err := client.GET(registerModelPath)
	if err != nil {
		return nil, fmt.Errorf("error fetching registered models: %w", err)
	}

	var modelList openapi.RegisteredModelList
	if err := json.Unmarshal(responseData, &modelList); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &modelList, nil
}

func (m RegisteredModel) CreateRegisteredModel(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.RegisteredModel, error) {
	responseData, err := client.POST(registerModelPath, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("error posting registered model: %w", err)
	}

	var model openapi.RegisteredModel
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}

func (m RegisteredModel) GetRegisteredModel(client integrations.HTTPClientInterface, id string) (*openapi.RegisteredModel, error) {
	path, err := url.JoinPath(registerModelPath, id)
	if err != nil {
		return nil, err
	}
	responseData, err := client.GET(path)

	if err != nil {
		return nil, fmt.Errorf("error fetching registered model: %w", err)
	}

	var model openapi.RegisteredModel
	if err := json.Unmarshal(responseData, &model); err != nil {
		return nil, fmt.Errorf("error decoding response data: %w", err)
	}

	return &model, nil
}
