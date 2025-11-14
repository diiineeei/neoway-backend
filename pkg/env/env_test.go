package env

import (
	"os"
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	key := "TEST_STRING"
	expectedValue := "test-value"
	defaultValue := "default"

	result := GetString(key, defaultValue)

	if result != defaultValue {
		t.Errorf("Expected %s, got %s", defaultValue, result)
	}

	os.Setenv(key, expectedValue)
	defer os.Unsetenv(key)
	result = GetString(key, defaultValue)

	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue int
		expected     int
	}{
		{"Sem env var", "", 10, 10},
		{"Com valor válido", "42", 10, 42},
		{"Com valor inválido", "invalid", 10, 10},
	}

	key := "TEST_INT"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
				defer os.Unsetenv(key)
			} else {
				os.Unsetenv(key)
			}

			result := GetInt(key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{"Sem env var", "", false, false},
		{"Com true", "true", false, true},
		{"Com 1", "1", false, true},
		{"Com false", "false", true, false},
		{"Com valor inválido", "invalid", true, true},
	}

	key := "TEST_BOOL"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
				defer os.Unsetenv(key)
			} else {
				os.Unsetenv(key)
			}

			result := GetBool(key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{"Sem env var", "", 5 * time.Second, 5 * time.Second},
		{"Com valor válido", "10s", 5 * time.Second, 10 * time.Second},
		{"Com minutos", "2m", 5 * time.Second, 2 * time.Minute},
		{"Com valor inválido", "invalid", 5 * time.Second, 5 * time.Second},
	}

	key := "TEST_DURATION"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
				defer os.Unsetenv(key)
			} else {
				os.Unsetenv(key)
			}

			result := GetDuration(key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue []string
		separator    string
		expected     []string
	}{
		{
			"Sem env var",
			"",
			[]string{"default1", "default2"},
			",",
			[]string{"default1", "default2"},
		},
		{
			"Com vírgula",
			"value1,value2,value3",
			[]string{},
			",",
			[]string{"value1", "value2", "value3"},
		},
		{
			"Com espaços",
			"value1 , value2 , value3",
			[]string{},
			",",
			[]string{"value1", "value2", "value3"},
		},
		{
			"Com ponto e vírgula",
			"value1;value2;value3",
			[]string{},
			";",
			[]string{"value1", "value2", "value3"},
		},
	}

	key := "TEST_SLICE"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
				defer os.Unsetenv(key)
			} else {
				os.Unsetenv(key)
			}

			result := GetStringSlice(key, tt.defaultValue, tt.separator)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("At index %d: expected %s, got %s", i, tt.expected[i], v)
				}
			}
		})
	}
}

func TestMustGetString(t *testing.T) {
	key := "TEST_MUST_GET"
	expectedValue := "required-value"

	// Test com variável definida
	os.Setenv(key, expectedValue)
	defer os.Unsetenv(key)

	result := MustGetString(key)
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}

	// Test sem variável definida (deve entrar em pânico)
	os.Unsetenv(key)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but didn't panic")
		}
	}()

	MustGetString(key) // Deve entrar em pânico
}
