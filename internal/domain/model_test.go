package domain

import (
	"reflect"
	"testing"
)

func TestNNModelDetails_Value(t *testing.T) {
	details := NNModelDetails{
		Name:         "test-model",
		Architecture: "CNN",
		CNN: &CNNDetails{
			Filters:   64,
			ResBlocks: 6,
		},
	}

	val, err := details.Value()
	if err != nil {
		t.Fatalf("Value(): unexpected error: %v", err)
	}

	b, ok := val.([]byte)
	if !ok {
		t.Fatalf("Value() type = %T, want []byte", val)
	}

	expectedJSON := `{"name":"test-model","architecture":"CNN","cnn":{"filters":64,"res_blocks":6,"se_blocks":false}}`
	if string(b) != expectedJSON {
		t.Errorf("Value() = %q, want %q", string(b), expectedJSON)
	}
}

func TestNNModelDetails_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    NNModelDetails
		wantErr bool
	}{
		{
			name:  "Valid byte slice",
			input: []byte(`{"name":"test-model","architecture":"CNN","cnn":{"filters":64,"res_blocks":6,"se_blocks":false}}`),
			want: NNModelDetails{
				Name:         "test-model",
				Architecture: "CNN",
				CNN: &CNNDetails{
					Filters:   64,
					ResBlocks: 6,
				},
			},
			wantErr: false,
		},
		{
			name:  "Valid string",
			input: `{"name":"test-model","architecture":"Transformer"}`,
			want: NNModelDetails{
				Name:         "test-model",
				Architecture: "Transformer",
			},
			wantErr: false,
		},
		{
			name:    "Invalid type",
			input:   123,
			want:    NNModelDetails{},
			wantErr: true,
		},
		{
			name:    "Empty byte slice",
			input:   []byte{},
			want:    NNModelDetails{},
			wantErr: false,
		},
		{
			name:    "Invalid json",
			input:   []byte(`{"name": "broken`),
			want:    NNModelDetails{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NNModelDetails
			err := got.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan(%#v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan(%#v):\ngot:  %#v\nwant: %#v", tt.input, got, tt.want)
			}
		})
	}
}
