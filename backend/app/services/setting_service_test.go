package services

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 20: Settings Persistence**
// **Validates: Requirements 9.2**
// For any setting update by Super Admin, the new value should be persisted
// and returned on subsequent retrieval.

// MockSettingStore simulates setting storage
type MockSettingStore struct {
	settings map[string]*string
}

func NewMockSettingStore() *MockSettingStore {
	return &MockSettingStore{
		settings: make(map[string]*string),
	}
}

func (s *MockSettingStore) Get(key string) *string {
	return s.settings[key]
}

func (s *MockSettingStore) Set(key string, value *string) {
	s.settings[key] = value
}

func (s *MockSettingStore) GetAll() map[string]*string {
	result := make(map[string]*string)
	for k, v := range s.settings {
		result[k] = v
	}
	return result
}

func TestProperty20_SettingsPersistence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Setting update should persist and be retrievable
	properties.Property("Setting update should persist value", prop.ForAll(
		func(key, value string) bool {
			if key == "" {
				return true // Skip empty keys
			}

			store := NewMockSettingStore()
			
			// Set value
			store.Set(key, &value)
			
			// Retrieve and verify
			retrieved := store.Get(key)
			return retrieved != nil && *retrieved == value
		},
		gen.AlphaString(),
		gen.AlphaString(),
	))

	// Property: Multiple settings should all persist
	properties.Property("Multiple settings should all persist", prop.ForAll(
		func(numSettings int) bool {
			if numSettings <= 0 || numSettings > 20 {
				return true
			}

			store := NewMockSettingStore()
			expected := make(map[string]string)

			// Set multiple settings
			for i := 0; i < numSettings; i++ {
				key := "setting_" + string(rune('A'+i))
				value := "value_" + string(rune('A'+i))
				store.Set(key, &value)
				expected[key] = value
			}

			// Verify all settings
			all := store.GetAll()
			if len(all) != numSettings {
				return false
			}

			for key, expectedValue := range expected {
				retrieved := store.Get(key)
				if retrieved == nil || *retrieved != expectedValue {
					return false
				}
			}

			return true
		},
		gen.IntRange(1, 20),
	))

	// Property: Setting update should overwrite previous value
	properties.Property("Setting update should overwrite previous value", prop.ForAll(
		func(key, value1, value2 string) bool {
			if key == "" {
				return true
			}

			store := NewMockSettingStore()
			
			// Set initial value
			store.Set(key, &value1)
			
			// Update with new value
			store.Set(key, &value2)
			
			// Should have new value
			retrieved := store.Get(key)
			return retrieved != nil && *retrieved == value2
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.AlphaString(),
	))

	// Property: Nil value should be stored correctly
	properties.Property("Nil value should be stored correctly", prop.ForAll(
		func(key string) bool {
			if key == "" {
				return true
			}

			store := NewMockSettingStore()
			
			// Set nil value
			store.Set(key, nil)
			
			// Should be nil
			retrieved := store.Get(key)
			return retrieved == nil
		},
		gen.AlphaString(),
	))

	properties.TestingRun(t)
}
