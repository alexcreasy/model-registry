package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet,
		"/api/v1/model_registry/model-registry/registered_models/1", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelRes RegisteredModelEnvelope
	err = json.Unmarshal(body, &registeredModelRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockModel := mocks.GetRegisteredModelMocks()[0]

	var expected = RegisteredModelEnvelope{
		Data: &mockModel,
	}

	//TODO assert the full structure, I couldn't get unmarshalling to work for the full customProperties values
	// this issue is in the test only
	assert.Equal(t, expected.Data.Name, registeredModelRes.Data.Name)
}

func TestGetAllRegisteredModelsHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet,
		"/api/v1/model_registry/model-registry/registered_models", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetAllRegisteredModelsHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelsListRes RegisteredModelListEnvelope
	err = json.Unmarshal(body, &registeredModelsListRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	modelList := mocks.GetRegisteredModelListMock()

	var expected = RegisteredModelListEnvelope{
		Data: &modelList,
	}

	assert.Equal(t, expected.Data.Size, registeredModelsListRes.Data.Size)
	assert.Equal(t, expected.Data.PageSize, registeredModelsListRes.Data.PageSize)
	assert.Equal(t, expected.Data.NextPageToken, registeredModelsListRes.Data.NextPageToken)
	assert.Equal(t, len(expected.Data.Items), len(registeredModelsListRes.Data.Items))
}

func TestCreateRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	newModel := openapi.NewRegisteredModel("Model One")
	newEnvelope := RegisteredModelEnvelope{Data: newModel}

	newModelJSON, err := json.Marshal(newEnvelope)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(newModelJSON)

	req, err := http.NewRequest(http.MethodPost,
		"/api/v1/model_registry/model-registry/registered_models", reqBody)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.CreateRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual RegisteredModelEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var expected = mocks.GetRegisteredModelMocks()[0]

	assert.Equal(t, expected.Name, actual.Data.Name)
	assert.NotEmpty(t, rs.Header.Get("location"))
}

func TestUpdateRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	newModel := openapi.NewRegisteredModel("Model One")
	newEnvelope := RegisteredModelEnvelope{Data: newModel}

	newEnvelopeJSON, err := json.Marshal(newEnvelope)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(newEnvelopeJSON)

	req, err := http.NewRequest(http.MethodPatch,
		"/api/v1/model_registry/model-registry/registered_models/1", reqBody)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.UpdateRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual RegisteredModelEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedModel := mocks.GetRegisteredModelMocks()[0]
	expected := RegisteredModelEnvelope{Data: &expectedModel}

	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}

func TestGetAllModelVersionsForRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1/versions", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetAllModelVersionsForRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual ModelVersionListEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	versionList := mocks.GetModelVersionListMock()

	expected := ModelVersionListEnvelope{
		Data: &versionList,
	}

	assert.Equal(t, expected.Data.Size, actual.Data.Size)
	assert.Equal(t, expected.Data.PageSize, actual.Data.PageSize)
	assert.Equal(t, expected.Data.NextPageToken, actual.Data.NextPageToken)
	assert.Equal(t, len(expected.Data.Items), len(actual.Data.Items))
}

func TestCreateModelVersionForRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	newVersion := openapi.NewModelVersion("Model One", "1")
	reqEnvelope := ModelVersionEnvelope{Data: newVersion}

	reqJSON, err := json.Marshal(reqEnvelope)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(reqJSON)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/model_registry/model-registry/registered_models/1/versions", reqBody)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	ps := httprouter.Params{
		httprouter.Param{
			Key:   ModelRegistryId,
			Value: "model-registry",
		},
		httprouter.Param{
			Key:   RegisteredModelId,
			Value: "1",
		},
	}

	testApp.CreateModelVersionForRegisteredModelHandler(rr, req, ps)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual ModelVersionEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)

	expectedVersion := mocks.GetModelVersionMocks()[0]

	expected := ModelVersionEnvelope{Data: &expectedVersion}

	assert.Equal(t, expected.Data.Name, actual.Data.Name)
	assert.Equal(t, rs.Header.Get("Location"), "/api/v1/model_registry/model-registry/model_versions/1")
}
