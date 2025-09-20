package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProcess_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	body, _ := json.Marshal(map[string]interface{}{
		"objective": map[string]interface{}{
			"n":            2,
			"coefficients": []float64{1, 2},
		},
		"constraints": map[string]interface{}{
			"rows": 2,
			"cols": 2,
			"vars": []float64{1, 2, 3, 4},
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "optimal_value")
	assert.Contains(t, resp, "solution")
}

func TestProcess_InvalidConstraints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	// filas y columnas no coinciden con cantidad de vars
	body, _ := json.Marshal(map[string]interface{}{
		"objective": map[string]interface{}{
			"n":            2,
			"coefficients": []float64{1, 2},
		},
		"constraints": map[string]interface{}{
			"rows": 2,
			"cols": 2,
			"vars": []float64{1, 2, 3}, // debería ser 4
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "cantidad de variables")
}

func TestProcess_InvalidObjective(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	body, _ := json.Marshal(map[string]interface{}{
		"objective": map[string]interface{}{
			"n":            2,
			"coefficients": []float64{1}, // longitud incorrecta
		},
		"constraints": map[string]interface{}{
			"rows": 1,
			"cols": 1,
			"vars": []float64{1},
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "coeficientes")
}

func TestProcess_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", process())

	body := []byte(`{"objective": {"n": 2, "coefficients": [1, "bad"]}, "constraints": {"rows": 1, "cols": 1, "vars": [1]}}`)

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
