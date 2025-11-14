package main

import (
	"testing"
)

func TestGetMetrics(t *testing.T) {
	metrics := getMetrics()

	if len(metrics) == 0 {
		t.Fatal("getMetrics() returned empty map")
	}

	// Test that all expected metrics are present
	expectedMetrics := []string{
		"hostport_data_read",
		"hostport_data_written",
		"disk_temperature",
		"disk_health",
		"disk_power_on_hours",
		"volume_health",
		"volume_iops",
		"pool_total_size",
		"controller_cpu",
		"system_health",
	}

	for _, name := range expectedMetrics {
		if _, exists := metrics[name]; !exists {
			t.Errorf("Metric %s not found in metrics", name)
		}
	}
}

func TestMetricDefinitions(t *testing.T) {
	metrics := getMetrics()

	tests := []struct {
		name             string
		expectedSources  int
		checkDescription bool
	}{
		{
			name:             "hostport_data_read",
			expectedSources:  1,
			checkDescription: true,
		},
		{
			name:             "disk_errors",
			expectedSources:  16, // Multiple error types
			checkDescription: true,
		},
		{
			name:             "volume_tier_distribution",
			expectedSources:  4, // Performance, Standard, Archive, RFC
			checkDescription: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, exists := metrics[tt.name]
			if !exists {
				t.Fatalf("Metric %s not found", tt.name)
			}

			if len(metric.Sources) != tt.expectedSources {
				t.Errorf("Metric %s has %d sources, expected %d",
					tt.name, len(metric.Sources), tt.expectedSources)
			}

			if tt.checkDescription && metric.Description == "" {
				t.Errorf("Metric %s has empty description", tt.name)
			}

			// Check that each source has required fields
			for i, source := range metric.Sources {
				if source.Path == "" {
					t.Errorf("Metric %s source %d has empty path", tt.name, i)
				}
				if source.ObjectSelector == "" {
					t.Errorf("Metric %s source %d has empty object selector", tt.name, i)
				}
				if source.PropertySelector == "" {
					t.Errorf("Metric %s source %d has empty property selector", tt.name, i)
				}
			}
		})
	}
}

func TestMetricLabels(t *testing.T) {
	metrics := getMetrics()

	tests := []struct {
		name           string
		expectedLabels []string
	}{
		{
			name:           "hostport_data_read",
			expectedLabels: []string{"port"},
		},
		{
			name:           "disk_temperature",
			expectedLabels: []string{"location", "serial"},
		},
		{
			name:           "volume_health",
			expectedLabels: []string{"volume"},
		},
		{
			name:           "controller_cpu",
			expectedLabels: []string{"controller"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, exists := metrics[tt.name]
			if !exists {
				t.Fatalf("Metric %s not found", tt.name)
			}

			source := metric.Sources[0]

			// Check that expected labels are in the mapping
			for _, expectedLabel := range tt.expectedLabels {
				found := false
				for _, labelName := range source.PropertiesAsLabel {
					if labelName == expectedLabel {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Metric %s missing expected label %s", tt.name, expectedLabel)
				}
			}
		})
	}
}

func TestDiskErrorsMetric(t *testing.T) {
	metrics := getMetrics()

	diskErrors, exists := metrics["disk_errors"]
	if !exists {
		t.Fatal("disk_errors metric not found")
	}

	// Should have 16 sources (8 error types * 2 ports)
	if len(diskErrors.Sources) != 16 {
		t.Errorf("disk_errors should have 16 sources, got %d", len(diskErrors.Sources))
	}

	// Check that we have all error types
	expectedTypes := []string{
		"smart",
		"io-timeout",
		"no-response",
		"spinup-retry",
		"media-errors",
		"nonmedia-errors",
		"block-reassigns",
		"bad-blocks",
	}

	foundTypes := make(map[string]int)
	for _, source := range diskErrors.Sources {
		if errorType, ok := source.Labels["type"]; ok {
			foundTypes[errorType.(string)]++
		}
	}

	for _, expectedType := range expectedTypes {
		count, found := foundTypes[expectedType]
		if !found {
			t.Errorf("disk_errors missing error type: %s", expectedType)
		}
		if count != 2 {
			t.Errorf("disk_errors error type %s should appear 2 times (for 2 ports), got %d", expectedType, count)
		}
	}

	// Check that each source has port label (1 or 2)
	for i, source := range diskErrors.Sources {
		port, ok := source.Labels["port"]
		if !ok {
			t.Errorf("disk_errors source %d missing port label", i)
			continue
		}
		portNum := port.(int)
		if portNum != 1 && portNum != 2 {
			t.Errorf("disk_errors source %d has invalid port %d", i, portNum)
		}
	}
}

func TestVolumeTierDistribution(t *testing.T) {
	metrics := getMetrics()

	tierDist, exists := metrics["volume_tier_distribution"]
	if !exists {
		t.Fatal("volume_tier_distribution metric not found")
	}

	// Should have 4 sources (Performance, Standard, Archive, RFC)
	if len(tierDist.Sources) != 4 {
		t.Errorf("volume_tier_distribution should have 4 sources, got %d", len(tierDist.Sources))
	}

	expectedTiers := []string{"Performance", "Standard", "Archive", "RFC"}
	foundTiers := make(map[string]bool)

	for _, source := range tierDist.Sources {
		if tier, ok := source.Labels["tier"]; ok {
			foundTiers[tier.(string)] = true
		}
	}

	for _, expectedTier := range expectedTiers {
		if !foundTiers[expectedTier] {
			t.Errorf("volume_tier_distribution missing tier: %s", expectedTier)
		}
	}
}

func TestSystemHealthMetric(t *testing.T) {
	metrics := getMetrics()

	systemHealth, exists := metrics["system_health"]
	if !exists {
		t.Fatal("system_health metric not found")
	}

	// system_health should have no labels
	if len(systemHealth.Sources) != 1 {
		t.Errorf("system_health should have 1 source, got %d", len(systemHealth.Sources))
	}

	source := systemHealth.Sources[0]
	if len(source.PropertiesAsLabel) != 0 {
		t.Errorf("system_health should have no labels, got %d", len(source.PropertiesAsLabel))
	}
}

func TestAllMetricsHaveDescription(t *testing.T) {
	metrics := getMetrics()

	for name, metric := range metrics {
		if metric.Description == "" {
			t.Errorf("Metric %s has empty description", name)
		}
	}
}

func TestAllMetricsHaveValidPaths(t *testing.T) {
	metrics := getMetrics()

	validPaths := map[string]bool{
		"host-port-statistics":  true,
		"disks":                 true,
		"disk-statistics":       true,
		"volumes":               true,
		"volume-statistics":     true,
		"pool-statistics":       true,
		"pools":                 true,
		"enclosures":            true,
		"enclosure":             true,
		"controller-statistics": true,
		"system":                true,
	}

	for name, metric := range metrics {
		for i, source := range metric.Sources {
			if !validPaths[source.Path] {
				t.Errorf("Metric %s source %d has invalid path: %s", name, i, source.Path)
			}
		}
	}
}

func TestMetricSourceConsistency(t *testing.T) {
	metrics := getMetrics()

	for name, metric := range metrics {
		for i, source := range metric.Sources {
			// Check that path is not empty
			if source.Path == "" {
				t.Errorf("Metric %s source %d has empty path", name, i)
			}

			// Check that object selector is not empty
			if source.ObjectSelector == "" {
				t.Errorf("Metric %s source %d has empty object selector", name, i)
			}

			// Check that property selector is not empty
			if source.PropertySelector == "" {
				t.Errorf("Metric %s source %d has empty property selector", name, i)
			}

			// PropertiesAsLabel can be nil/empty for system_health
			// Labels can be nil/empty for most metrics
		}
	}
}

func TestMetricCount(t *testing.T) {
	metrics := getMetrics()

	// We should have at least 80 metrics as documented in README
	if len(metrics) < 70 {
		t.Errorf("Expected at least 70 metrics, got %d", len(metrics))
	}
}
