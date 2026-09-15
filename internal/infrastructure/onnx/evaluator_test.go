package onnx

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
	ort "github.com/yalue/onnxruntime_go"
)

// initTestEnvironment initializes the ONNX Runtime C-API for testing.
// It skips the test if running in short mode or if the DLL cannot be located.
func initTestEnvironment(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if ort.IsInitialized() {
		return
	}

	rootDir := filepath.Join("..", "..", "..")
	dllPath, _ := filepath.Abs(filepath.Join(rootDir, "onnxruntime.dll"))

	if _, err := os.Stat(dllPath); os.IsNotExist(err) {
		t.Skipf("Skipping ONNX tests: onnxruntime.dll not found at %s", dllPath)
	}

	ort.SetSharedLibraryPath(dllPath)
	if err := ort.InitializeEnvironment(); err != nil {
		t.Skipf("Skipping ONNX tests: failed to initialize environment: %v", err)
	}
}

// getTestModelPath returns the absolute path to the transformer NN test model.
// It skips the test if the test data file is missing.
func getTestModelPath(t *testing.T) string {
	t.Helper()
	modelPath, _ := filepath.Abs(filepath.Join("testdata", "test_tr.onnx"))

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skipf("Skipping ONNX tests: test model not found at %s", modelPath)
	}
	return modelPath
}

func TestNewEvaluator_Uninitialized(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Destroy the global environment if it's already initialized
	if ort.IsInitialized() {
		_ = ort.DestroyEnvironment()
	}

	_, err := NewEvaluator("dummy.onnx")
	if !errors.Is(err, ErrNotInitialized) {
		t.Errorf("NewEvaluator(%q) error = %v, wantErr %v", "dummy.onnx", err, ErrNotInitialized)
	}
}

func TestNewEvaluator(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	evaluator, err := NewEvaluator(modelPath)
	if err != nil {
		t.Fatalf("failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	if evaluator.session == nil {
		t.Errorf("NewEvaluator(%q) session = nil, want non-nil", modelPath)
	}
}

func TestNewEvaluator_BadPath(t *testing.T) {
	initTestEnvironment(t)
	badPath := "this_file_does_not_exist.onnx"
	_, err := NewEvaluator(badPath)
	if err == nil {
		t.Errorf("NewEvaluator(%q) error = nil, wantErr %v", badPath, true)
	}
}

func TestNewEvaluator_CorruptFile(t *testing.T) {
	initTestEnvironment(t)
	corruptModelPath := filepath.Join(t.TempDir(), "bad.onnx")
	_ = os.WriteFile(corruptModelPath, []byte("this is not a real onnx model"), 0644)

	_, err := NewEvaluator(corruptModelPath)
	if err == nil {
		t.Errorf("NewEvaluator(%q) error = nil, wantErr %v", corruptModelPath, true)
	}
}

func TestNewEvaluator_OptionError(t *testing.T) {
	initTestEnvironment(t)

	errMockHardware := errors.New("hardware unsupported")

	// Create a mock option that intentionally fails with our sentinel error
	badOption := func(opts *ort.SessionOptions) error {
		return errMockHardware
	}

	_, err := NewEvaluator("dummy.onnx", badOption)
	if !errors.Is(err, errMockHardware) {
		t.Errorf("NewEvaluator(badOption) error = %v, wantErr %v", err, errMockHardware)
	}
}

func TestEvaluator_EvaluateBatch_Empty(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	evaluator, err := NewEvaluator(modelPath)
	if err != nil {
		t.Fatalf("failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	policies, values, err := evaluator.EvaluateBatch([]engine.Position{})
	if err != nil {
		t.Fatalf("EvaluateBatch(empty): unexpected error: %v", err)
	}
	if policies != nil || values != nil {
		t.Errorf("EvaluateBatch(empty) = %v, %v; want nil, nil", policies, values)
	}
}

func TestEvaluator_EvaluateBatch_OutputLengths(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	evaluator, err := NewEvaluator(modelPath)
	if err != nil {
		t.Fatalf("failed to create evaluator: %v", err)
	}
	defer evaluator.Close()

	batchSize := 2
	positions := make([]engine.Position, batchSize)
	for i := 0; i < batchSize; i++ {
		positions[i] = *engine.NewPosition()
	}

	policies, values, err := evaluator.EvaluateBatch(positions)
	if err != nil {
		t.Fatalf("EvaluateBatch(%v): unexpected error: %v", batchSize, err)
	}

	if n := len(policies); n != batchSize {
		t.Errorf("Len: got %d; want %d", n, batchSize)
	}
	if n := len(values); n != batchSize {
		t.Errorf("Len: got %d; want %d", n, batchSize)
	}

	if n := len(policies[0]); n != policyDim {
		t.Errorf("Len: got %d; want %d", n, policyDim)
	}
}

func TestEvaluator_EvaluateBatch_SessionFailure(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	evaluator, err := NewEvaluator(modelPath)
	if err != nil {
		t.Fatalf("failed to create evaluator: %v", err)
	}

	// Deliberately close the session to force an ONNX C-API error during evaluation
	evaluator.Close()

	batchSize := 1
	positions := make([]engine.Position, batchSize)
	positions[0] = *engine.NewPosition()

	_, _, err = evaluator.EvaluateBatch(positions)
	if err == nil {
		t.Errorf("EvaluateBatch() error = nil, wantErr %v", true)
	}
}
