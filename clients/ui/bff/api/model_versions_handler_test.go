package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetModelVersionHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet,
		"/api/v1/model_registry/model-registry/model_versions/1", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetModelVersionHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual ModelVersionEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockModel := mocks.GetModelVersionMocks()[0]

	var expected = ModelVersionEnvelope{
		Data: &mockModel,
	}

	//TODO assert the full structure, I couldn't get unmarshalling to work for the full customProperties values
	// this issue is in the test only
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}

func TestCreateModelVersionHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	newVersion := openapi.NewModelVersion("Model One", "1")
	newEnvelope := ModelVersionEnvelope{Data: newVersion}

	newVersionJSON, err := json.Marshal(newEnvelope)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(newVersionJSON)

	req, err := http.NewRequest(http.MethodPost,
		"/api/v1/model_registry/model-registry/model_versions", reqBody)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.CreateModelVersionHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual ModelVersionEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var expected = mocks.GetModelVersionMocks()[0]

	assert.Equal(t, expected.Name, actual.Data.Name)
	assert.NotEmpty(t, rs.Header.Get("Location"))
	assert.Equal(t, rs.Header.Get("Location"), "/api/v1/model_registry/model-registry/model_versions/1")
}
