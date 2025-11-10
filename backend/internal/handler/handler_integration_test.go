package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"autosimplex/internal/simplex"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestProcess_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", Process())

	body, _ := json.Marshal(map[string]any{
		"objective": map[string]any{
			"n":            2,
			"coefficients": []float64{1, 2},
			"type":         "maximize",
		},
		"constraints": map[string]any{
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
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "optimal_value")
	assert.Contains(t, resp, "solution")
}

func TestProcess_InvalidConstraints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", Process())

	// filas y columnas no coinciden con cantidad de vars
	body, _ := json.Marshal(map[string]any{
		"objective": map[string]any{
			"n":            2,
			"coefficients": []float64{1, 2},
		},
		"constraints": map[string]any{
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
	router.POST("/process", Process())

	body, _ := json.Marshal(map[string]any{
		"objective": map[string]any{
			"n":            2,
			"coefficients": []float64{1}, // longitud incorrecta
		},
		"constraints": map[string]any{
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
	router.POST("/process", Process())

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

func TestProcess_MinimizeHandlerMatchesManualConversion(t *testing.T) {
	// Ensure that handler's minimize flow (negating coefficients + inverting result)
	// matches the manual conversion using the simplex solver.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", Process())

	coefs := []float64{5, 4, 3}
	constraintsVars := []float64{
		2, 3, 1, 5,
		4, 1, 2, 11,
		3, 4, 2, 8,
	}

	// Call handler with minimize request
	body, _ := json.Marshal(map[string]any{
		"objective": map[string]any{
			"n":            3,
			"coefficients": coefs,
			"type":         "minimize",
		},
		"constraints": map[string]any{
			"rows": 3,
			"cols": 4,
			"vars": constraintsVars,
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Manual conversion: Solve with maximize vector = -coefs, then invert sign
	maximizeVec := mat.NewVecDense(3, []float64{-5, -4, -3})
	constraintMatrix := mat.NewDense(3, 4, constraintsVars)
	manualMax, _, _ := simplex.Solve(maximizeVec, constraintMatrix)
	expectedMin := -manualMax

	// Response optimal_value is float64
	val, ok := resp["optimal_value"].(float64)
	assert.True(t, ok)
	assert.Equal(t, expectedMin, val)
}

func TestProcess_MinimizeWithGreaterEqual(t *testing.T) {
	// Test a minimization problem with >= constraints submitted to the handler.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/process", Process())

	coefs := []float64{4, 5}
	constraintsVars := []float64{
		2, 1, 8,
		1, 3, 12,
	}
	signs := []string{">=", ">="}

	body, _ := json.Marshal(map[string]any{
		"objective": map[string]any{
			"n":            2,
			"coefficients": coefs,
			"type":         "minimize",
		},
		"constraints": map[string]any{
			"rows":  2,
			"cols":  3,
			"vars":  constraintsVars,
			"signs": signs,
		},
	})

	req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Manual conversion: negate coefficients to maximize and use SolveWithSigns,
	// then invert the result
	maximizeVec := mat.NewVecDense(2, []float64{-4, -5})
	constraintMatrix := mat.NewDense(2, 3, constraintsVars)
	manualMax, _, _ := simplex.SolveWithSigns(maximizeVec, constraintMatrix, signs)
	expectedMin := -manualMax

	val, ok := resp["optimal_value"].(float64)
	assert.True(t, ok)
	assert.Equal(t, expectedMin, val)
}
