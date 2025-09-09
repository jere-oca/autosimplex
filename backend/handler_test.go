package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProcess_ValidMatrix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	matrix := [][]float64{
		{1.1, 2.2},
		{3.3, 4.4},
	}
	body, _ := json.Marshal(gin.H{"matrix": matrix})

	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string][][]float64
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, matrix, resp["received_matrix"])
}

func TestProcess_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	body := []byte(`{"matrix": [ [1, 2], [3, "bad"] ]}`) // "bad" is not a float

	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, ok := resp["error"]
	assert.True(t, ok)
}

func TestProcess_MissingMatrix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	body := []byte(`{}`)

	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string][][]float64
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, [][]float64(nil), resp["received_matrix"])
}
