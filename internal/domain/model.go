package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ErrModelNotFound is returned when a requested NN model cannot be found in the repository.
var ErrModelNotFound = errors.New("model not found")

// NNModel represents a neural network configuration used for position evaluation.
type NNModel struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Details   NNModelDetails `json:"details"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
}

// NNModelDetails contains the architectural and training metadata of the neural network.
type NNModelDetails struct {
	Name         string `json:"name"`
	Architecture string `json:"architecture"` // "CNN" or "Transformer"

	TestPolicyLoss float64 `json:"test_policy,omitempty"`
	TestValueLoss  float64 `json:"test_value,omitempty"`
	MAE            float64 `json:"mae,omitempty"`
	TestTop1       float64 `json:"test_top1,omitempty"`
	TestTop3       float64 `json:"test_top3,omitempty"`
	TestTop5       float64 `json:"test_top5,omitempty"`

	TrainingMeta *TrainingMeta       `json:"training_meta,omitempty"`
	CNN          *CNNDetails         `json:"cnn,omitempty"`
	Transformer  *TransformerDetails `json:"transformer,omitempty"`
}

// Value implements the driver.Valuer interface for database serialization.
func (d NNModelDetails) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NNModelDetails) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("failed to scan NNModelDetails: expected []byte or string, got %T", value)
		}
		b = []byte(s)
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, d)
}

// TrainingMeta encapsulates hyperparameters utilized during the model's training phase.
type TrainingMeta struct {
	BatchSize       int     `json:"batch_size"`
	Epochs          int     `json:"epochs"`
	MaxLR           float64 `json:"max_lr"`
	MinLR           float64 `json:"min_lr"`
	DatasetSize     int     `json:"dataset_size"`
	ValidationSplit float64 `json:"validation_split"`
	WeightDecay     float64 `json:"weight_decay"`
	LabelSmoothing  float64 `json:"label_smoothing,omitempty"`
}

// CNNDetails defines hyperparameters specific to Convolutional Neural Networks.
type CNNDetails struct {
	Filters   int  `json:"filters"`
	ResBlocks int  `json:"res_blocks"`
	SEBlocks  bool `json:"se_blocks"`
}

// TransformerDetails defines hyperparameters specific to Transformer Neural Networks.
type TransformerDetails struct {
	InputDim        int     `json:"input_dim"`
	DModel          int     `json:"d_model"`
	NHead           int     `json:"nhead"`
	NumLayers       int     `json:"num_layers"`
	FeedforwardMult int     `json:"feedforward_mult"`
	Dropout         float64 `json:"dropout"`
}
