package onnx

import (
	"errors"
	"fmt"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
	"github.com/OKmeray/chess-engine/internal/usecase"
	ort "github.com/yalue/onnxruntime_go"
)

// Evaluator implements usecase.Evaluator using the ONNX Runtime C-API.
type Evaluator struct {
	session    *ort.DynamicAdvancedSession
	inputName  string
	policyName string
	valueName  string
}

const (
	inputDim1 = 64
	inputDim2 = 20
	policyDim = 4864
	valueDim  = 1
)

var (
	// ErrNotInitialized indicates the ONNX Runtime environment was not set up before creating an Evaluator.
	ErrNotInitialized = errors.New("ONNX Runtime environment is not initialized")
	// ErrNoInputs indicates the loaded ONNX model does not have any input nodes.
	ErrNoInputs = errors.New("model has no inputs")
)

// EvaluatorOption configures the underlying ONNX runtime session.
type EvaluatorOption func(*ort.SessionOptions) error

// WithDirectML configures the evaluator to use DirectML hardware acceleration.
func WithDirectML(deviceID int) EvaluatorOption {
	return func(opts *ort.SessionOptions) error {
		// TODO: Implement DirectML provider options
		return nil
	}
}

// WithCUDA configures the evaluator to use NVIDIA CUDA hardware acceleration.
func WithCUDA(deviceID int) EvaluatorOption {
	return func(opts *ort.SessionOptions) error {
		// TODO: Implement CUDA provider options
		return nil
	}
}

// WithCoreML configures the evaluator to use Apple's CoreML hardware acceleration.
func WithCoreML() EvaluatorOption {
	return func(opts *ort.SessionOptions) error {
		// TODO: Implement CoreML provider options
		return nil
	}
}

// NewEvaluator creates a new ONNX session for the given model.
// The caller must ensure the ONNX runtime environment is initialized via
// ort.SetSharedLibraryPath and ort.InitializeEnvironment prior to this call.
// Additional hardware acceleration can be passed via EvaluatorOption.
func NewEvaluator(modelPath string, opts ...EvaluatorOption) (*Evaluator, error) {
	if !ort.IsInitialized() {
		return nil, ErrNotInitialized
	}

	options, err := ort.NewSessionOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to create session options: %w", err)
	}
	defer func() {
		_ = options.Destroy()
	}()

	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply execution provider: %w", err)
		}
	}

	inputs, _, err := ort.GetInputOutputInfo(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get model info: %w", err)
	}
	if len(inputs) == 0 {
		return nil, ErrNoInputs
	}
	inputName := inputs[0].Name

	session, err := ort.NewDynamicAdvancedSession(
		modelPath,
		[]string{inputName},
		[]string{"policy", "value"},
		options,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic session: %w", err)
	}

	return &Evaluator{
		session:    session,
		inputName:  inputName,
		policyName: "policy",
		valueName:  "value",
	}, nil
}

// Close destroys the ONNX session and frees associated C memory.
func (e *Evaluator) Close() {
	if e.session != nil {
		defer func() {
			_ = e.session.Destroy()
		}()
	}
}

// EvaluateBatch fulfills the usecase.Evaluator interface.
func (e *Evaluator) EvaluateBatch(positions []engine.Position) (policies [][]float32, valueData []float32, err error) {
	batchSize := len(positions)
	if batchSize == 0 {
		return nil, nil, nil
	}

	inputSize := batchSize * inputDim1 * inputDim2
	inputData := make([]float32, inputSize)

	for i := 0; i < batchSize; i++ {
		offset := i * inputDim1 * inputDim2
		usecase.ExtractTransformerFeatures(&positions[i], inputData[offset:offset+(inputDim1*inputDim2)])
	}

	inputShape := []int64{int64(batchSize), inputDim1, inputDim2}
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create input tensor: %w", err)
	}
	defer func() {
		_ = inputTensor.Destroy()
	}()

	policySize := batchSize * policyDim
	policyData := make([]float32, policySize)
	policyShape := []int64{int64(batchSize), policyDim}
	policyTensor, err := ort.NewTensor(policyShape, policyData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create policy tensor: %w", err)
	}
	defer func() {
		_ = policyTensor.Destroy()
	}()

	valueSize := batchSize * valueDim
	valueData = make([]float32, valueSize)
	valueShape := []int64{int64(batchSize), valueDim}
	valueTensor, err := ort.NewTensor(valueShape, valueData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create value tensor: %w", err)
	}
	defer func() {
		_ = valueTensor.Destroy()
	}()

	err = e.session.Run(
		[]ort.Value{inputTensor},
		[]ort.Value{policyTensor, valueTensor},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("session run failed: %w", err)
	}

	// Reshape output slices
	policies = make([][]float32, batchSize)
	for i := 0; i < batchSize; i++ {
		policies[i] = policyData[i*policyDim : (i+1)*policyDim]
	}

	return policies, valueData, nil
}
