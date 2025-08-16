package models

import (
	"testing"
	"time"

	"subconverter-go/pkg/constants"
	"subconverter-go/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestConvertRequest_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		request  *ConvertRequest
		expected bool
	}{
		{
			name: "valid request",
			request: &ConvertRequest{
				URL:    "https://example.com/subscription",
				Target: "clash",
			},
			expected: true,
		},
		{
			name: "invalid - missing URL",
			request: &ConvertRequest{
				Target: "clash",
			},
			expected: false,
		},
		{
			name: "invalid - missing target",
			request: &ConvertRequest{
				URL: "https://example.com/subscription",
			},
			expected: false,
		},
		{
			name:     "nil request",
			request:  nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertRequest_SetDefaults(t *testing.T) {
	request := &ConvertRequest{
		URL:    "https://example.com/subscription",
		Target: "clash",
	}

	request.SetDefaults()

	assert.Equal(t, constants.DefaultRulesetInterval, request.Interval)
	assert.False(t, request.CreatedAt.IsZero())
}

func TestConvertRequest_GetContentType(t *testing.T) {
	tests := []struct {
		target   string
		expected string
	}{
		{"clash", constants.ContentTypeYAML},
		{"surge", constants.ContentTypeConf},
		{"ss", constants.ContentTypeText},
		{"singbox", constants.ContentTypeJSON},
		{"unknown", constants.ContentTypeYAML},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			request := &ConvertRequest{Target: tt.target}
			result := request.GetContentType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response *ConvertResponse
		expected bool
	}{
		{
			name: "successful response",
			response: &ConvertResponse{
				Success:    true,
				StatusCode: 200,
			},
			expected: true,
		},
		{
			name: "failed response",
			response: &ConvertResponse{
				Success:    false,
				StatusCode: 400,
			},
			expected: false,
		},
		{
			name:     "nil response",
			response: nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.IsSuccess()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertResponse_SetError(t *testing.T) {
	response := &ConvertResponse{}
	err := types.NewConvertError(types.ErrorCodeInvalidURL, "Invalid URL provided")
	err.RequestID = "test-123"

	response.SetError(err)

	assert.False(t, response.Success)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, "Invalid URL provided", response.Message)
	assert.Equal(t, "test-123", response.RequestID)
}

func TestRequestStats_SuccessRate(t *testing.T) {
	stats := &RequestStats{
		TotalRequests:   100,
		SuccessRequests: 85,
		ErrorRequests:   15,
	}

	rate := stats.SuccessRate()
	assert.Equal(t, 85.0, rate)

	// 测试零请求
	emptyStats := &RequestStats{}
	emptyRate := emptyStats.SuccessRate()
	assert.Equal(t, 0.0, emptyRate)
}

func TestRequestStats_AddRequest(t *testing.T) {
	stats := &RequestStats{}

	// 添加成功请求
	stats.AddRequest(true, 100*time.Millisecond)

	assert.Equal(t, int64(1), stats.TotalRequests)
	assert.Equal(t, int64(1), stats.SuccessRequests)
	assert.Equal(t, int64(0), stats.ErrorRequests)
	assert.Equal(t, 100*time.Millisecond, stats.AverageTime)
	assert.False(t, stats.LastRequest.IsZero())

	// 添加失败请求
	stats.AddRequest(false, 200*time.Millisecond)

	assert.Equal(t, int64(2), stats.TotalRequests)
	assert.Equal(t, int64(1), stats.SuccessRequests)
	assert.Equal(t, int64(1), stats.ErrorRequests)
	assert.Equal(t, 150*time.Millisecond, stats.AverageTime)
}