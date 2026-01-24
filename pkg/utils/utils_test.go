package utils_test

import (
	"testing"
	"time"

	"github.com/pewpowder/url-shortener/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseOptionalTimeReturnError(t *testing.T) {
	invalidTimeValue := "2026-01-02T15:04:05"
	timeValue := "2026-01-02T15:04:05Z"
	parsedTime := time.Date(2026, 01, 02, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name        string
		value       *string
		want        *time.Time
		wantErr     bool
		errContains string
	}{
		{
			name:    "check_nil_value",
			value:   nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:        "check_invalid_time",
			value:       &invalidTimeValue,
			want:        nil,
			wantErr:     true,
			errContains: "invalid time format",
		},
		{
			name:    "parse_success",
			value:   &timeValue,
			want:    &parsedTime,
			wantErr: false,
		},
	}
	
	

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseOptionalTime(tt.value)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errContains)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			if assert.NotNil(t, got) {
				assert.True(t, got.Equal(*tt.want))
			}
		})
	}
}
