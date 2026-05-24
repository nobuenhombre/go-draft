package vars

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr error
	}{
		{
			name:    "empty string returns empty map",
			input:   "",
			want:    map[string]string{},
			wantErr: nil,
		},
		{
			name:    "whitespace-only string returns empty map",
			input:   "   ",
			want:    map[string]string{},
			wantErr: nil,
		},
		{
			name:    "single pair without spaces",
			input:   "PROJECT_NAME:hello-world",
			want:    map[string]string{"PROJECT_NAME": "hello-world"},
			wantErr: nil,
		},
		{
			name:    "multiple pairs without spaces",
			input:   "key1:val1,key2:val2,key3:val3",
			want:    map[string]string{"key1": "val1", "key2": "val2", "key3": "val3"},
			wantErr: nil,
		},
		{
			name:    "values with colons",
			input:   "key:val:ue",
			want:    map[string]string{"key": "val:ue"},
			wantErr: nil,
		},
		{
			name:    "trims spaces around keys and values",
			input:   "  key1  :  val1  ,  key2:val2  ",
			want:    map[string]string{"key1": "val1", "key2": "val2"},
			wantErr: nil,
		},
		{
			name:    "missing colon returns error",
			input:   "key1val1",
			want:    nil,
			wantErr: ErrorInvalidKeyValuePair,
		},
		{
			name:    "empty key returns error",
			input:   ":val1,key2:val2",
			want:    nil,
			wantErr: ErrorEmptyKey,
		},
		{
			name:    "empty key in first pair only",
			input:   ":val1",
			want:    nil,
			wantErr: ErrorEmptyKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			got, err := s.Parse(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.wantErr),
					"expected error %v, got %v", tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewReturnsService(t *testing.T) {
	s := New()
	require.NotNil(t, s)

	// Verify it implements the Service interface
	_, ok := s.(Service)
	require.True(t, ok, "New() should return a Service")
}

func TestProviderConcreteType(t *testing.T) {
	s := New()
	_, ok := s.(*Provider)
	require.True(t, ok, "New() should return *Provider")
}
