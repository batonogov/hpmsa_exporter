package main

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test helper functions
func TestFindObjects(t *testing.T) {
	tests := []struct {
		name     string
		objects  []Object
		search   string
		expected int
	}{
		{
			name: "find direct object",
			objects: []Object{
				{Name: "test1"},
				{Name: "test2"},
			},
			search:   "test1",
			expected: 1,
		},
		{
			name: "find nested objects",
			objects: []Object{
				{
					Name: "parent",
					Objects: []Object{
						{Name: "child1"},
						{Name: "child2"},
					},
				},
			},
			search:   "child1",
			expected: 1,
		},
		{
			name: "find multiple objects",
			objects: []Object{
				{Name: "disk"},
				{Name: "disk"},
				{Name: "volume"},
			},
			search:   "disk",
			expected: 2,
		},
		{
			name:     "no match",
			objects:  []Object{{Name: "test"}},
			search:   "nomatch",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findObjects(tt.objects, tt.search)
			if len(result) != tt.expected {
				t.Errorf("findObjects() returned %d objects, expected %d", len(result), tt.expected)
			}
		})
	}
}

func TestFindProperty(t *testing.T) {
	tests := []struct {
		name          string
		obj           Object
		propertyName  string
		expectedValue string
		expectedFound bool
	}{
		{
			name: "property exists",
			obj: Object{
				Properties: []Property{
					{Name: "name", Value: "test-value"},
					{Name: "id", Value: "123"},
				},
			},
			propertyName:  "name",
			expectedValue: "test-value",
			expectedFound: true,
		},
		{
			name: "property not found",
			obj: Object{
				Properties: []Property{
					{Name: "id", Value: "123"},
				},
			},
			propertyName:  "name",
			expectedValue: "",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := findProperty(tt.obj, tt.propertyName)
			if found != tt.expectedFound {
				t.Errorf("findProperty() found = %v, expected %v", found, tt.expectedFound)
			}
			if value != tt.expectedValue {
				t.Errorf("findProperty() value = %v, expected %v", value, tt.expectedValue)
			}
		})
	}
}

func TestExtractLabels(t *testing.T) {
	obj := Object{
		Properties: []Property{
			{Name: "durable-id", Value: "A"},
			{Name: "serial-number", Value: "12345"},
			{Name: "other", Value: "ignored"},
		},
	}

	mapping := map[string]string{
		"durable-id":    "controller",
		"serial-number": "serial",
	}

	labels := extractLabels(obj, mapping)

	if len(labels) != 2 {
		t.Errorf("extractLabels() returned %d labels, expected 2", len(labels))
	}

	if labels["controller"] != "A" {
		t.Errorf("extractLabels() controller = %v, expected 'A'", labels["controller"])
	}

	if labels["serial"] != "12345" {
		t.Errorf("extractLabels() serial = %v, expected '12345'", labels["serial"])
	}

	if _, exists := labels["other"]; exists {
		t.Errorf("extractLabels() should not include 'other' label")
	}
}

// Test XML parsing
func TestXMLParsing(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="drive">
		<PROPERTY name="location">1.1</PROPERTY>
		<PROPERTY name="serial-number">ABC123</PROPERTY>
		<PROPERTY name="temperature-numeric">45</PROPERTY>
	</OBJECT>
	<OBJECT name="drive">
		<PROPERTY name="location">1.2</PROPERTY>
		<PROPERTY name="serial-number">ABC124</PROPERTY>
		<PROPERTY name="temperature-numeric">46</PROPERTY>
	</OBJECT>
</RESPONSE>`

	var response Response
	err := xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	if len(response.Objects) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(response.Objects))
	}

	drives := findObjects(response.Objects, "drive")
	if len(drives) != 2 {
		t.Errorf("Expected 2 drives, got %d", len(drives))
	}

	temp, found := findProperty(drives[0], "temperature-numeric")
	if !found {
		t.Error("temperature-numeric property not found")
	}
	if temp != "45" {
		t.Errorf("Expected temperature 45, got %s", temp)
	}
}

// Test MetricStore
func TestMetricStore(t *testing.T) {
	ms := NewMetricStore()

	t.Run("create new metric", func(t *testing.T) {
		labels := []string{"controller", "serial"}
		metric := ms.GetOrCreate("test_metric", "Test metric description", labels)

		if metric == nil {
			t.Fatal("GetOrCreate returned nil")
		}

		// Try to get the same metric again
		metric2 := ms.GetOrCreate("test_metric", "Test metric description", labels)
		if metric2 != metric {
			t.Error("GetOrCreate should return the same metric instance")
		}
	})

	t.Run("set metric value", func(t *testing.T) {
		labels := []string{"host"}
		metric := ms.GetOrCreate("test_gauge", "Test gauge", labels)

		labelValues := map[string]string{"host": "localhost"}
		metric.With(labelValues).Set(42.5)

		// Metric is set successfully if no panic occurs
	})
}

// Test MSA Client with mock server
func TestMSAClient(t *testing.T) {
	// Create a mock server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/login/"+getSHA256("testuser_testpass"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="status">
		<PROPERTY name="response-type">success</PROPERTY>
		<PROPERTY name="response">test-session-key-123</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/version":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="controller-a-versions">
		<PROPERTY name="bundle-version">1.2.3</PROPERTY>
		<PROPERTY name="bundle-base-version">1.2.0</PROPERTY>
		<PROPERTY name="sc-fw">1.0</PROPERTY>
		<PROPERTY name="mc-fw">2.0</PROPERTY>
		<PROPERTY name="pld-rev">3.0</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/disks":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="drive">
		<PROPERTY name="location">1.1</PROPERTY>
		<PROPERTY name="serial-number">ABC123</PROPERTY>
		<PROPERTY name="temperature-numeric">45</PROPERTY>
		<PROPERTY name="health-numeric">0</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Extract host from server URL (remove https://)
	host := server.URL[8:]

	t.Run("authentication", func(t *testing.T) {
		client, err := NewMSAClient(host, "testuser", "testpass", 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if client.sessionKey == "" {
			t.Error("Session key is empty")
		}

		if client.sessionKey != "test-session-key-123" {
			t.Errorf("Expected session key 'test-session-key-123', got '%s'", client.sessionKey)
		}
	})

	t.Run("get data", func(t *testing.T) {
		client, err := NewMSAClient(host, "testuser", "testpass", 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		data, err := client.Get("version")
		if err != nil {
			t.Fatalf("Failed to get version: %v", err)
		}

		if len(data) == 0 {
			t.Error("Response data is empty")
		}

		var response Response
		err = xml.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if len(response.Objects) == 0 {
			t.Error("No objects in response")
		}
	})
}

// Helper function for SHA256 (for testing)
func getSHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func TestScrapeMSA(t *testing.T) {
	// Create a mock server with complete responses
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/login/"+getSHA256("testuser_testpass"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="status">
		<PROPERTY name="response">test-session-key</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/version":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="controller-a-versions">
		<PROPERTY name="bundle-version">1.0</PROPERTY>
		<PROPERTY name="bundle-base-version">1.0</PROPERTY>
		<PROPERTY name="sc-fw">1.0</PROPERTY>
		<PROPERTY name="mc-fw">1.0</PROPERTY>
		<PROPERTY name="pld-rev">1.0</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/disks":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="drive">
		<PROPERTY name="location">1.1</PROPERTY>
		<PROPERTY name="serial-number">TEST123</PROPERTY>
		<PROPERTY name="temperature-numeric">40</PROPERTY>
		<PROPERTY name="health-numeric">0</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		default:
			// Return empty response for other endpoints
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><RESPONSE></RESPONSE>`))
		}
	}))
	defer server.Close()

	host := server.URL[8:]

	t.Run("scrape metrics", func(t *testing.T) {
		client, err := NewMSAClient(host, "testuser", "testpass", 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ms := NewMetricStore()
		err = scrapeMSA(client, ms)
		if err != nil {
			t.Fatalf("scrapeMSA failed: %v", err)
		}

		// Check that metrics were created
		if len(ms.metrics) == 0 {
			t.Error("No metrics were created")
		}
	})
}
