# AutoSimplex

AutoSimplex is a Go-based web API for solving linear programming problems using the Simplex algorithm. It uses the Gin HTTP framework and is designed to process matrix-based optimization problems.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Bootstrap and Setup
All work happens in the `/backend` directory. Always change to this directory first:
```bash
cd /home/runner/work/autosimplex/autosimplex/backend
```

### Required Dependencies
- Go 1.25.0 (exact version required - see go.mod)
- Internet connection for downloading Go modules

### Build and Test Workflow
Execute commands in this exact order. **NEVER CANCEL** any of these operations - they must complete fully:

1. **Install Dependencies** (takes ~5 seconds):
   ```bash
   go mod tidy
   ```

2. **Build Project** (takes ~18 seconds, NEVER CANCEL - set timeout to 60+ seconds):
   ```bash
   go build -v ./...
   ```

3. **Run Tests** (takes ~4 seconds, NEVER CANCEL - set timeout to 30+ seconds):
   ```bash
   go test -v ./...
   ```

4. **Code Quality Checks** (takes <2 seconds):
   ```bash
   gofmt -l .
   go vet ./...
   ```

### Running the Application
- **Start Server**: `go run .` (NOT `go run main.go` - requires all files)
- **Default Port**: 8080
- **Access URL**: `http://localhost:8080`

The server will show Gin framework debug messages indicating it's ready:
```
[GIN-debug] POST   /process                  --> main.main.process.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

## API Testing and Validation

### Manual Validation Requirements
**ALWAYS** test the API manually after making any changes. Execute these exact scenarios:

#### Test Scenario 1: Valid Matrix Request
```bash
curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '{"matrix": [[1.1, 2.2], [3.3, 4.4]]}'
```
**Expected Response**: `{"received_matrix":[[1.1,2.2],[3.3,4.4]]}`

#### Test Scenario 2: Invalid Data Type
```bash
curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '{"matrix": [[1, 2], [3, "bad"]]}'
```
**Expected Response**: `{"error":"json: cannot unmarshal string into Go struct field MatrixRequest.matrix of type float64"}`

#### Test Scenario 3: Empty Request
```bash
curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '{}'
```
**Expected Response**: `{"received_matrix":null}`

### Validation Checklist
Before completing any changes, verify:
- [ ] All builds complete without errors (18+ seconds expected)
- [ ] All tests pass (4+ seconds expected)  
- [ ] Code formatting is clean (`gofmt -l .` returns nothing)
- [ ] No vet warnings (`go vet ./...` returns nothing)
- [ ] Server starts successfully on port 8080
- [ ] All three API test scenarios return expected responses

## Project Structure

### Key Files and Locations
```
/backend/
├── main.go                    # Server entry point with Gin router
├── handler.go                 # HTTP request handlers  
├── handler_test.go           # Handler unit tests
├── go.mod                    # Go module dependencies (Go 1.25.0)
├── go.sum                    # Dependency checksums
├── internal/models/
│   └── request.go            # Data structures for API requests
└── simplex_draft.md          # Draft implementation notes
```

### Important Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/stretchr/testify` - Testing assertions  
- `gonum.org/v1/gonum` - Mathematical operations library

## Common Tasks

### Adding New Endpoints
1. Define handler function in `handler.go`
2. Add route in `main.go`  
3. Create tests in `handler_test.go`
4. Add data structures to `internal/models/` if needed
5. **ALWAYS** validate with manual API testing

### Debugging Issues
- Check Go version: `go version` (must be 1.25.0)
- Clean build cache: `go clean -cache`  
- Rebuild dependencies: `go mod tidy`
- Check server logs for Gin debug messages

### CI/CD Pipeline
The GitHub workflow (`.github/workflows/go.yml`) runs:
1. `go build -v ./...` (working directory: `./backend`)
2. `go test -v ./...` (working directory: `./backend`)

**CRITICAL**: Always run the complete build and test workflow before committing changes. The CI requires Go 1.25.0 and will fail if tests don't pass.

## Time Expectations and Warnings

| Operation | Expected Time | Timeout Setting | Notes |
|-----------|---------------|-----------------|-------|
| `go mod tidy` | ~5 seconds | 60 seconds | Downloads dependencies |
| `go build -v ./...` | ~18 seconds | 60+ seconds | **NEVER CANCEL** |  
| `go test -v ./...` | ~4 seconds | 30+ seconds | **NEVER CANCEL** |
| `go run .` | Immediate | N/A | Server starts quickly |
| Code quality checks | <2 seconds | 30 seconds | `gofmt`, `go vet` |

**WARNING**: Build times may vary based on system load and network conditions. Always allow commands to complete fully rather than canceling early.

## Frequently Referenced Information

### Repository Root Contents
```
.
├── .git/
├── .github/
│   ├── copilot-instructions.md  # This file
│   └── workflows/
│       └── go.yml              # CI/CD pipeline
├── .gitignore
├── LICENSE  
├── README.md                   # Basic setup instructions (Spanish)
├── backend/                    # Main Go application
└── docs/
    └── request_example.json    # Example API request format
```

### Current API Structure
- **Endpoint**: `POST /process`
- **Request Format**: `{"matrix": [[float64]]}`
- **Response**: `{"received_matrix": [[float64]]}` or `{"error": "message"}`
- **Future**: Additional Simplex-specific endpoints planned (see `simplex_draft.md`)

### Error Resolution
- **Build failures**: Check Go version, run `go clean -cache`, retry `go mod tidy`
- **Test failures**: Review test output, ensure server not running on port 8080
- **Import errors**: Verify all files are in `/backend` directory, not repository root
- **Server won't start**: Check if port 8080 is available, ensure `go run .` not `go run main.go`