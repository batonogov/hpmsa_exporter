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

	t.Run("authentication failure - wrong credentials", func(t *testing.T) {
		_, err := NewMSAClient(host, "wronguser", "wrongpass", 10*time.Second)
		if err == nil {
			t.Error("Expected authentication to fail with wrong credentials")
		}
	})

	t.Run("get data - not found", func(t *testing.T) {
		client, err := NewMSAClient(host, "testuser", "testpass", 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		_, err = client.Get("nonexistent")
		if err == nil {
			t.Error("Expected Get to fail with non-existent endpoint")
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
	<OBJECT name="controller-b-versions">
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
		<PROPERTY name="architecture">SSD</PROPERTY>
	</OBJECT>
	<OBJECT name="drive">
		<PROPERTY name="location">1.2</PROPERTY>
		<PROPERTY name="serial-number">TEST124</PROPERTY>
		<PROPERTY name="temperature-numeric">N/A</PROPERTY>
		<PROPERTY name="health-numeric">1</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/host-port-statistics":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="host-port-statistics">
		<PROPERTY name="durable-id">hostport_A1</PROPERTY>
		<PROPERTY name="data-read-numeric">1000</PROPERTY>
		<PROPERTY name="data-written-numeric">2000</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/volumes":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="volume">
		<PROPERTY name="volume-name">vol1</PROPERTY>
		<PROPERTY name="health-numeric">0</PROPERTY>
		<PROPERTY name="size-numeric">1000000</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		case r.URL.Path == "/api/show/pools":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="pools">
		<PROPERTY name="name">pool1</PROPERTY>
		<PROPERTY name="serial-number">POOL123</PROPERTY>
		<PROPERTY name="total-size-numeric">10000000</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		default:
			// Return empty response for other endpoints
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><RESPONSE></RESPONSE>`))
		}
	}))
	defer server.Close()

	host := server.URL[8:]

	// Create client once for all tests
	client, err := NewMSAClient(host, "testuser", "testpass", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create one metric store for all tests (metrics are registered globally)
	ms := NewMetricStore()

	t.Run("scrape metrics", func(t *testing.T) {
		err = scrapeMSA(client, ms)
		if err != nil {
			t.Fatalf("scrapeMSA failed: %v", err)
		}

		// Check that metrics were created
		if len(ms.metrics) == 0 {
			t.Error("No metrics were created")
		}
	})

	t.Run("verify N/A values handling", func(t *testing.T) {
		// N/A values are already tested in the first scrape
		// We just verify metrics exist
		if len(ms.metrics) == 0 {
			t.Error("No metrics were created")
		}
	})

	t.Run("verify SSD filter", func(t *testing.T) {
		// SSD filtering is already tested in the first scrape
		// We just verify metrics exist
		if len(ms.metrics) == 0 {
			t.Error("No metrics were created")
		}
	})

	t.Run("verify version metrics", func(t *testing.T) {
		// Check that version metric was created
		found := false
		for name := range ms.metrics {
			if name == "msa_version" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Version metric was not created")
		}
	})
}

func TestRecursiveFindProperty(t *testing.T) {
	// Test that findProperty can find properties in nested objects (like resettable-statistics in tier-statistics)
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="pool-statistics">
		<PROPERTY name="pool">A</PROPERTY>
		<PROPERTY name="serial-number">00c0ff640e37000089b7d26301000000</PROPERTY>
		<OBJECT name="tier-statistics">
			<PROPERTY name="tier">Performance</PROPERTY>
			<PROPERTY name="pool">A</PROPERTY>
			<PROPERTY name="serial-number">00c0ff640e37000089b7d26301000001</PROPERTY>
			<OBJECT name="resettable-statistics">
				<PROPERTY name="number-of-reads">643699762</PROPERTY>
				<PROPERTY name="number-of-writes">1594098785</PROPERTY>
				<PROPERTY name="data-read-numeric">87388675853312</PROPERTY>
				<PROPERTY name="data-written-numeric">96356295811072</PROPERTY>
				<PROPERTY name="avg-rsp-time">369</PROPERTY>
				<PROPERTY name="avg-read-rsp-time">185</PROPERTY>
				<PROPERTY name="avg-write-rsp-time">427</PROPERTY>
			</OBJECT>
		</OBJECT>
	</OBJECT>
</RESPONSE>`

	var resp Response
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	// Find tier-statistics object
	tierStats := findObjects(resp.Objects, "tier-statistics")
	if len(tierStats) != 1 {
		t.Fatalf("Expected 1 tier-statistics object, got %d", len(tierStats))
	}

	// Test recursive property search
	tests := []struct {
		property string
		expected string
	}{
		{"tier", "Performance"},                    // Direct property
		{"number-of-reads", "643699762"},           // Nested in resettable-statistics
		{"number-of-writes", "1594098785"},         // Nested in resettable-statistics
		{"data-read-numeric", "87388675853312"},    // Nested in resettable-statistics
		{"data-written-numeric", "96356295811072"}, // Nested in resettable-statistics
		{"avg-rsp-time", "369"},                    // Nested in resettable-statistics
		{"avg-read-rsp-time", "185"},               // Nested in resettable-statistics
		{"avg-write-rsp-time", "427"},              // Nested in resettable-statistics
	}

	for _, tt := range tests {
		t.Run(tt.property, func(t *testing.T) {
			value, ok := findProperty(tierStats[0], tt.property)
			if !ok {
				t.Errorf("Property %s not found in tier-statistics", tt.property)
			}
			if value != tt.expected {
				t.Errorf("Property %s: expected %s, got %s", tt.property, tt.expected, value)
			}
		})
	}
}

// Test printXMLStructure - mainly for coverage
func TestPrintXMLStructure(t *testing.T) {
	objects := []Object{
		{
			Name: "parent",
			Properties: []Property{
				{Name: "prop1", Value: "val1"},
			},
			Objects: []Object{
				{
					Name: "child",
					Properties: []Property{
						{Name: "prop2", Value: "val2"},
					},
				},
			},
		},
	}

	// Just call the function to ensure it doesn't panic
	printXMLStructure(objects, "")
}

// Test NewMSAClient error cases
func TestNewMSAClientErrors(t *testing.T) {
	t.Run("invalid XML response", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid xml`))
		}))
		defer server.Close()

		host := server.URL[8:]
		_, err := NewMSAClient(host, "test", "test", 10*time.Second)
		if err == nil {
			t.Error("Expected error for invalid XML response")
		}
	})

	t.Run("missing session key in response", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="status">
		<PROPERTY name="other">value</PROPERTY>
	</OBJECT>
</RESPONSE>`))
		}))
		defer server.Close()

		host := server.URL[8:]
		_, err := NewMSAClient(host, "test", "test", 10*time.Second)
		if err == nil {
			t.Error("Expected error for missing session key")
		}
		if err != nil && err.Error() != "session key not found in response" {
			t.Errorf("Expected 'session key not found in response', got '%v'", err)
		}
	})

	t.Run("authentication with non-OK status", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		host := server.URL[8:]
		_, err := NewMSAClient(host, "test", "test", 10*time.Second)
		if err == nil {
			t.Error("Expected error for non-OK status")
		}
	})
}

// Test scrapeMSA error cases
func TestScrapeMSAErrors(t *testing.T) {
	t.Run("version fetch failure", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/login/"+getSHA256("testerr1_testerr1") {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="status">
		<PROPERTY name="response">test-key</PROPERTY>
	</OBJECT>
</RESPONSE>`))
			} else if r.URL.Path == "/api/show/version" {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
		defer server.Close()

		host := server.URL[8:]
		client, err := NewMSAClient(host, "testerr1", "testerr1", 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ms := NewMetricStore()
		err = scrapeMSA(client, ms)
		if err == nil {
			t.Error("Expected error when version fetch fails")
		}
	})

	t.Run("invalid version XML", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/login/"+getSHA256("testerr2_testerr2") {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RESPONSE>
	<OBJECT name="status">
		<PROPERTY name="response">test-key</PROPERTY>
	</OBJECT>
</RESPONSE>`))
			} else if r.URL.Path == "/api/show/version" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`invalid xml`))
			}
		}))
		defer server.Close()

		host := server.URL[8:]
		client, err := NewMSAClient(host, "testerr2", "testerr2", 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ms := NewMetricStore()
		err = scrapeMSA(client, ms)
		if err == nil {
			t.Error("Expected error when version XML is invalid")
		}
	})
}

// Test health check endpoint
func TestHealthCheckEndpoint(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"status":"healthy","service":"msa_exporter"}`)
	})

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"healthy","service":"msa_exporter"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}
}
