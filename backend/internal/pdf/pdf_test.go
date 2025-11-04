package pdf

import (
    "bytes"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"

    "autosimplex/internal/handler"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// TestProcessPDFIntegration carga el ejemplo de request y solicita el PDF al handler
func TestProcessPDFIntegration(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/process", handler.Process())

    // cargar archivo de ejemplo desde docs
    repoRoot := ".." // estamos en backend/internal/pdf, subir dos niveles
    reqPath := filepath.Join(repoRoot, "..", "docs", "request_example.json")
    // corregir ruta relativa desde el paquete de tests
    // intentar rutas comunes
    possible := []string{
        filepath.Join("..", "..", "docs", "request_example.json"),
        filepath.Join("..", "docs", "request_example.json"),
        filepath.Join("docs", "request_example.json"),
    }

    var data []byte
    var err error
    for _, p := range possible {
        if _, statErr := os.Stat(p); statErr == nil {
            data, err = os.ReadFile(p)
            break
        }
    }
    if data == nil || err != nil {
        t.Fatalf("no se pudo leer docs/request_example.json: %v", err)
    }

    req, _ := http.NewRequest(http.MethodPost, "/process?format=pdf", bytes.NewBuffer(data))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    // verificar headers
    assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
    assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")

    // el cuerpo debe tener datos (PDF mínimo)
    body := w.Body.Bytes()
    assert.True(t, len(body) > 0)
    // PDF típico comienza con %PDF
    assert.True(t, bytes.HasPrefix(body, []byte{"%", "P", "D", "F"}))

    // adicional: intentar generar PDF directamente con GenerateSimplexPDF usando un reader
    // leer un ejemplo más pequeño
    // ...existing code...
    var buf bytes.Buffer
    // usar valores de ejemplo: valor óptimo 1.23, solución [1,2,3]
    err = GenerateSimplexPDF(1.23, []float64{1, 2, 3}, &buf)
    assert.NoError(t, err)
    assert.True(t, buf.Len() > 0)
    assert.True(t, bytes.HasPrefix(buf.Bytes(), []byte{"%", "P", "D", "F"}))
}

// helper para leer desde io.Reader a bytes
func readAll(r io.Reader) ([]byte, error) {
    var b bytes.Buffer
    _, err := io.Copy(&b, r)
    return b.Bytes(), err
}
