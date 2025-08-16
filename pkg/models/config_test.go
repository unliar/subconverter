package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   *ServerConfig
		expected bool
	}{
		{
			name: "valid config",
			config: &ServerConfig{
				ListenAddress: "0.0.0.0",
				ListenPort:    25500,
			},
			expected: true,
		},
		{
			name: "invalid - missing address",
			config: &ServerConfig{
				ListenPort: 25500,
			},
			expected: false,
		},
		{
			name: "invalid - port too low",
			config: &ServerConfig{
				ListenAddress: "0.0.0.0",
				ListenPort:    0,
			},
			expected: false,
		},
		{
			name:     "nil config",
			config:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServerConfig_SetDefaults(t *testing.T) {
	config := &ServerConfig{}

	config.SetDefaults()

	assert.Equal(t, "0.0.0.0", config.ListenAddress)
	assert.Equal(t, 25500, config.ListenPort)
	assert.Equal(t, 10240, config.MaxPendingConns)
	assert.Equal(t, 4, config.MaxConcurThreads)
	assert.False(t, config.CreatedAt.IsZero())
	assert.False(t, config.UpdatedAt.IsZero())
}

func TestConverterConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   *ConverterConfig
		expected bool
	}{
		{
			name:     "valid empty config",
			config:   &ConverterConfig{},
			expected: true,
		},
		{
			name: "valid config with rulesets",
			config: &ConverterConfig{
				CustomRulesets: []RulesetConfig{
					{
						Group: "test",
						URL:   "https://example.com/rules.list",
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   *TemplateConfig
		expected bool
	}{
		{
			name: "valid config",
			config: &TemplateConfig{
				Name: "test-template",
				Path: "/path/to/template.j2",
			},
			expected: true,
		},
		{
			name: "invalid - missing name",
			config: &TemplateConfig{
				Path: "/path/to/template.j2",
			},
			expected: false,
		},
		{
			name: "invalid - missing path",
			config: &TemplateConfig{
				Name: "test-template",
			},
			expected: false,
		},
		{
			name:     "nil config",
			config:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestApplicationConfig_SetDefaults(t *testing.T) {
	config := &ApplicationConfig{}

	config.SetDefaults()

	assert.NotNil(t, config.Server)
	assert.NotNil(t, config.Converter)
	assert.False(t, config.CreatedAt.IsZero())
	assert.False(t, config.UpdatedAt.IsZero())
}