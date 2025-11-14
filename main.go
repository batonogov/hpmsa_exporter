package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	prefix = "msa_"
)

// XML structures for parsing MSA API responses
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type Object struct {
	Name       string     `xml:"name,attr"`
	Properties []Property `xml:"PROPERTY"`
	Objects    []Object   `xml:"OBJECT"`
}

type Response struct {
	Objects []Object `xml:"OBJECT"`
}

// MetricSource defines how to collect a metric
type MetricSource struct {
	Path              string
	ObjectSelector    string
	PropertySelector  string
	PropertiesAsLabel map[string]string
	Labels            map[string]interface{}
}

// MetricDefinition defines a metric to collect
type MetricDefinition struct {
	Description string
	Sources     []MetricSource
}

// MetricStore manages Prometheus metrics
type MetricStore struct {
	mu      sync.Mutex
	metrics map[string]*prometheus.GaugeVec
}

// NewMetricStore creates a new MetricStore
func NewMetricStore() *MetricStore {
	return &MetricStore{
		metrics: make(map[string]*prometheus.GaugeVec),
	}
}

// GetOrCreate gets or creates a Prometheus gauge metric
func (ms *MetricStore) GetOrCreate(name, description string, labelNames []string) *prometheus.GaugeVec {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	key := name
	if metric, exists := ms.metrics[key]; exists {
		return metric
	}

	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: description,
		},
		labelNames,
	)
	prometheus.MustRegister(metric)
	ms.metrics[key] = metric

	return metric
}

// MSAClient represents a client for the MSA API
type MSAClient struct {
	host       string
	sessionKey string
	login      string
	httpClient *http.Client
	timeout    time.Duration
}

// NewMSAClient creates a new MSA API client
func NewMSAClient(host, login, password string, timeout time.Duration) (*MSAClient, error) {
	// Create HTTP client with insecure TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	client := &MSAClient{
		host:       host,
		login:      login,
		httpClient: httpClient,
		timeout:    timeout,
	}

	// Authenticate
	creds := fmt.Sprintf("%s_%s", login, password)
	hash := sha256.Sum256([]byte(creds))
	hashStr := fmt.Sprintf("%x", hash)

	url := fmt.Sprintf("https://%s/api/login/%s", host, hashStr)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth response: %w", err)
	}

	var response Response
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	// Extract session key
	for _, obj := range response.Objects {
		for _, prop := range obj.Properties {
			if prop.Name == "response" {
				client.sessionKey = prop.Value
				break
			}
		}
	}

	if client.sessionKey == "" {
		return nil, fmt.Errorf("session key not found in response")
	}

	return client, nil
}

// Get performs a GET request to the MSA API
func (c *MSAClient) Get(path string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/show/%s", c.host, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("sessionKey", c.sessionKey)
	req.AddCookie(&http.Cookie{Name: "wbisessionkey", Value: c.sessionKey})
	req.AddCookie(&http.Cookie{Name: "wbiusername", Value: c.login})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Helper function to find objects by name
func findObjects(objects []Object, name string) []Object {
	var result []Object
	for _, obj := range objects {
		if obj.Name == name {
			result = append(result, obj)
		}
		// Recursively search nested objects
		result = append(result, findObjects(obj.Objects, name)...)
	}
	return result
}

// Helper function to find property by name
func findProperty(obj Object, name string) (string, bool) {
	for _, prop := range obj.Properties {
		if prop.Name == name {
			return prop.Value, true
		}
	}
	return "", false
}

// Helper function to extract labels
func extractLabels(obj Object, mapping map[string]string) map[string]string {
	labels := make(map[string]string)
	for _, prop := range obj.Properties {
		if labelName, ok := mapping[prop.Name]; ok {
			labels[labelName] = prop.Value
		}
	}
	return labels
}

// scrapeMSA collects metrics from MSA storage
func scrapeMSA(client *MSAClient, metricStore *MetricStore) error {
	pathCache := make(map[string][]byte)

	// Collect firmware version
	versionData, err := client.Get("version")
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	var versionResp Response
	if err := xml.Unmarshal(versionData, &versionResp); err != nil {
		return fmt.Errorf("failed to parse version: %w", err)
	}

	// Process firmware versions
	for _, controller := range []string{"controller-a-versions", "controller-b-versions"} {
		for _, obj := range findObjects(versionResp.Objects, controller) {
			labels := map[string]string{"controller": controller}
			for _, version := range []string{"bundle-version", "bundle-base-version", "sc-fw", "mc-fw", "pld-rev"} {
				if val, ok := findProperty(obj, version); ok {
					labelKey := ""
					switch version {
					case "bundle-version":
						labelKey = "bundle_version"
					case "bundle-base-version":
						labelKey = "bundle_base_version"
					case "sc-fw":
						labelKey = "sc_fw"
					case "mc-fw":
						labelKey = "mc_fw"
					case "pld-rev":
						labelKey = "pld_rev"
					}
					labels[labelKey] = val
				}
			}
			labelNames := make([]string, 0, len(labels))
			for k := range labels {
				labelNames = append(labelNames, k)
			}
			metric := metricStore.GetOrCreate(prefix+"version", "Firmware Versions", labelNames)
			metric.With(labels).Set(1)
		}
	}

	// Process all metrics
	for name, metricDef := range getMetrics() {
		metricName := prefix + name
		for _, source := range metricDef.Sources {
			// Get or cache the path data
			if _, ok := pathCache[source.Path]; !ok {
				data, err := client.Get(source.Path)
				if err != nil {
					log.Printf("Failed to get %s: %v", source.Path, err)
					continue
				}
				pathCache[source.Path] = data
			}

			var resp Response
			if err := xml.Unmarshal(pathCache[source.Path], &resp); err != nil {
				log.Printf("Failed to parse %s: %v", source.Path, err)
				continue
			}

			// Find objects matching the selector
			objects := findObjects(resp.Objects, source.ObjectSelector)

			// Special handling for complex selectors (like SSD-only disks)
			if source.ObjectSelector == "drive" {
				// Check for SSD architecture filter
				filtered := []Object{}
				for _, obj := range objects {
					if arch, ok := findProperty(obj, "architecture"); ok && arch == "SSD" {
						filtered = append(filtered, obj)
					} else if arch == "" {
						// No filter, include all
						filtered = append(filtered, obj)
					}
				}
				if len(filtered) > 0 {
					objects = filtered
				}
			}

			for _, obj := range objects {
				// Extract labels
				labels := extractLabels(obj, source.PropertiesAsLabel)

				// Add static labels from source
				for k, v := range source.Labels {
					labels[k] = fmt.Sprint(v)
				}

				// Find the value
				value, ok := findProperty(obj, source.PropertySelector)
				if !ok {
					continue
				}

				// Handle N/A values
				if value == "N/A" {
					value = "NaN"
				}

				// Parse value
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					if value == "NaN" {
						floatValue = math.NaN()
					} else {
						log.Printf("Failed to parse value %s: %v", value, err)
						continue
					}
				}

				// Get label names
				labelNames := make([]string, 0, len(labels))
				for k := range labels {
					labelNames = append(labelNames, k)
				}

				// Set the metric
				metric := metricStore.GetOrCreate(metricName, metricDef.Description, labelNames)
				metric.With(labels).Set(floatValue)
			}
		}
	}

	return nil
}

func main() {
	// Parse command line arguments
	hostname := flag.String("hostname", "", "MSA storage hostname")
	login := flag.String("login", "", "MSA storage login")
	password := flag.String("password", "", "MSA storage password")
	port := flag.Int("port", 8000, "Exporter port")
	interval := flag.Int("interval", 60, "Scrape interval in seconds")
	timeout := flag.Int("timeout", 60, "Scrape timeout in seconds")

	flag.Parse()

	// Handle positional arguments for backward compatibility
	args := flag.Args()
	if len(args) >= 3 {
		*hostname = args[0]
		*login = args[1]
		*password = args[2]
	}

	if *hostname == "" || *login == "" || *password == "" {
		log.Fatal("hostname, login, and password are required")
	}

	fmt.Printf("Starting MSA exporter on port %d\n", *port)
	fmt.Printf("Connecting to %s as %s\n", *hostname, *login)
	fmt.Printf("Scraping every %d seconds with timeout %d seconds\n", *interval, *timeout)

	// Create metric store
	metricStore := NewMetricStore()

	// Start Prometheus HTTP server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		addr := fmt.Sprintf(":%d", *port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Main scraping loop
	timeoutDuration := time.Duration(*timeout) * time.Second
	intervalDuration := time.Duration(*interval) * time.Second

	for {
		client, err := NewMSAClient(*hostname, *login, *password, timeoutDuration)
		if err != nil {
			log.Printf("Failed to create client: %v", err)
		} else {
			if err := scrapeMSA(client, metricStore); err != nil {
				log.Printf("Failed to scrape: %v", err)
			}
		}
		time.Sleep(intervalDuration)
	}
}
