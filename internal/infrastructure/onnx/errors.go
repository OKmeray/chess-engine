package onnx

import "errors"

var (
	// ErrFetchModelMetadata is returned when the manager cannot retrieve the model path from the database.
	ErrFetchModelMetadata = errors.New("failed to fetch model metadata")

	// ErrLoadONNXFile is returned when the ONNX runtime fails to parse or load the model file.
	ErrLoadONNXFile = errors.New("failed to load onnx file")
)
