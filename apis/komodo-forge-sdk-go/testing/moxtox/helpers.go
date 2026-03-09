package moxtox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// MoxtoxConfig represents the parsed YAML config with lookup maps for different modes.
type MoxtoxConfig struct {
	EnableMoxtox        bool     `yaml:"enableMoxtox"`
	Debug               bool     `yaml:"debug"`
	DefaultDelay        int      `yaml:"defaultDelay"`
	PerformanceMode     string  `yaml:"performanceMode"`  // "default", "quick", "dynamic"
	AllowedEnvironments []string `yaml:"allowedEnvironments"`
	IgnoredRoutes       []string `yaml:"ignoredRoutes"`
	Mappings            map[string]interface{} `yaml:"mappings"`  // Raw YAML for parsing
	LookupMap           map[string]map[string][]Scenario `yaml:"-"`  // For default/dynamic: path -> method -> []Scenario
	HashLookupMap       map[string]map[string]map[string]Scenario `yaml:"-"`  // For quick: path -> method -> hash -> Scenario
}

// Mapping represents a path's configuration.
type Mapping struct {
	Methods  map[string]Method `yaml:"methods,omitempty"`
}

// Method represents HTTP method configuration.
type Method struct {
	Scenarios []Scenario `yaml:"scenarios,omitempty"`
}

// Scenario represents a conditional mock response.
type Scenario struct {
	Name       string                 `yaml:"name"`
	File       string                 `yaml:"file"`
	Dynamic    bool                   `yaml:"dynamic,omitempty"`
	Template   string                 `yaml:"template,omitempty"`
	Delay      int                    `yaml:"delay,omitempty"`
	Log        bool                   `yaml:"log,omitempty"`
	Priority   int                    `yaml:"priority,omitempty"`
	Conditions map[string]interface{} `yaml:"conditions,omitempty"`
}

// MockResponse represents the structure of user-defined mock files.
type MockResponse struct {
	Status  int                 `json:"status"`
	Body    string              `json:"body"`
	Headers map[string]string   `json:"headers,omitempty"`
}

// countTotalScenarios counts the total number of scenarios across all mappings for dynamic mode.
func (cnfg *MoxtoxConfig) countTotalScenarios() int {
	count := 0
	for _, rawMapping := range cnfg.Mappings {
		mappingData, ok := rawMapping.(map[interface{}]interface{})
		if !ok { continue }
		mapping := parseMapping(mappingData)
		for _, methodData := range mapping.Methods {
			count += len(methodData.Scenarios)
		}
	}
	return count
}

// parseMapping converts raw YAML map to Mapping struct.
func parseMapping(data map[interface{}]interface{}) Mapping {
	mapping := Mapping{Methods: make(map[string]Method)}

	for key, value := range data {
		keyStr := key.(string)

		if keyStr == "methods" {
			// Handle methods if present: parse as`` map[string]Method
			if methodsData, ok := value.(map[interface{}]interface{}); ok {
				for methodKey, methodValue := range methodsData {
					methodStr := methodKey.(string)
					if methodMap, ok := methodValue.(map[interface{}]interface{}); ok {
						method := parseMethod(methodMap)
						mapping.Methods[methodStr] = method
					}
				}
			}
		} else {
			// Assume sub-paths like /token
			if subData, ok := value.(map[interface{}]interface{}); ok {
				subMapping := parseMapping(subData)
				for method, methodData := range subMapping.Methods {
					mapping.Methods[method] = methodData
				}
			}
		}
	}
	return mapping
}

// parseMethod converts raw method map to Method struct.
func parseMethod(data map[interface{}]interface{}) Method {
	method := Method{Scenarios: []Scenario{}}

	for key, value := range data {
		keyStr := key.(string)

		if keyStr == "scenarios" {
			if scenariosData, ok := value.([]interface{}); ok {
				for _, scenarioData := range scenariosData {
					if scenarioMap, ok := scenarioData.(map[interface{}]interface{}); ok {
						scenario := parseScenario(scenarioMap)
						method.Scenarios = append(method.Scenarios, scenario)
					}
				}
			}
		}
	}
	return method
}

// parseScenario converts raw scenario map to Scenario struct.
func parseScenario(data map[interface{}]interface{}) Scenario {
	scenario := Scenario{Conditions: make(map[string]interface{})}
	for key, value := range data {
		keyStr := key.(string)

		switch keyStr {
			case "name":
				if str, ok := value.(string); ok {
					scenario.Name = str
				}
			case "file":
				if str, ok := value.(string); ok {
					scenario.File = str
				}
			case "dynamic":
				if b, ok := value.(bool); ok {
					scenario.Dynamic = b
				}
			case "template":
				if str, ok := value.(string); ok {
					scenario.Template = str
				}
			case "delay":
				if i, ok := value.(int); ok {
					scenario.Delay = i
				}
			case "log":
				if b, ok := value.(bool); ok {
					scenario.Log = b
				}
			case "priority":
				if i, ok := value.(int); ok {
					scenario.Priority = i
				}
			case "conditions":
				if condMap, ok := value.(map[interface{}]interface{}); ok {
					for condKey, condValue := range condMap {
						scenario.Conditions[condKey.(string)] = condValue
					}
				}
		}
	}
	return scenario
}

// matchesRequest checks if the request matches any scenario.
// Dispatches to slice or hash based on mode.
// Returns the matching scenario if found.
func matchesRequest(r *http.Request) (Scenario, bool) {
	if config.PerformanceMode == "quick" || (config.PerformanceMode == "dynamic" && config.countTotalScenarios() > 10) {
		return matchesRequestHash(r)
	}
	return matchesRequestSlice(r)
}

// matchesConditions checks if the request matches the given conditions.
func matchesConditions(r *http.Request, cond map[string]interface{}) bool {
	for condType, condMap := range cond {
		switch condType {
			case "body":
				if !matchesBody(r, condMap.(map[interface{}]interface{})) {
					return false
				}
			case "query":
				if !matchesQuery(r, condMap.(map[interface{}]interface{})) {
					return false
				}
			case "headers":
				if !matchesHeaders(r, condMap.(map[interface{}]interface{})) {
					return false
				}
			case "path":
				// For path params, simple check (extend for regex if needed)
				if !matchesPath(r.URL.Path, condMap.(map[interface{}]interface{})) {
					return false
				}
		}
	}
	return true
}

// matchesBody checks body conditions.
func matchesBody(r *http.Request, condMap map[interface{}]interface{}) bool {
	if r.Body == nil || !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return false
	}

	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)

	for k, v := range condMap {
		key := k.(string)
		if bodyVal, ok := body[key]; !ok || fmt.Sprintf("%v", bodyVal) != fmt.Sprintf("%v", v) {
			return false
		}
	}
	return true
}

// matchesQuery checks query conditions.
func matchesQuery(r *http.Request, condMap map[interface{}]interface{}) bool {
	query := r.URL.Query()

	for k, v := range condMap {
		key := k.(string)
		if queryVal := query.Get(key); queryVal != fmt.Sprintf("%v", v) {
			return false
		}
	}
	return true
}

// matchesHeaders checks header conditions.
func matchesHeaders(r *http.Request, condMap map[interface{}]interface{}) bool {
	for k, v := range condMap {
		key := k.(string)
		if headerVal := r.Header.Get(key); headerVal != fmt.Sprintf("%v", v) {
			return false
		}
	}
	return true
}

// matchesPath checks path conditions (simple for now).
func matchesPath(path string, condMap map[interface{}]interface{}) bool {
	// Placeholder: assume conditions are key-value for path segments
	for k, v := range condMap {
		key := k.(string)
		// Simple check, e.g., if key is "id", check if path contains /id/v
		if !strings.Contains(path, fmt.Sprintf("/%s/%v", key, v)) {
			return false
		}
	}
	return true
}

// injectMock reads the mock file or generates dynamic response and applies it.
func injectMock(w http.ResponseWriter, r *http.Request, scenario Scenario) error {
	delay := scenario.Delay
	if delay == 0 {
		delay = config.DefaultDelay
	}
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	var responseBody string
	var status int
	headers := make(map[string]string)

	if scenario.Dynamic {
		// Generate dynamic response from template
		responseBody = generateDynamicResponse(r, scenario.Template)
		status = 200  // Default for dynamic
	} else {
		// Read from file
		data, err := os.ReadFile(scenario.File)
		if err != nil {
			return fmt.Errorf("failed to read mock file %s: %v", scenario.File, err)
		}
		var mock MockResponse
		if err := json.Unmarshal(data, &mock); err != nil {
			return fmt.Errorf("failed to parse mock file %s: %v", scenario.File, err)
		}
		responseBody = mock.Body
		status = mock.Status
		headers = mock.Headers
	}

	// Apply headers
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write([]byte(responseBody))
	return nil
}

// generateDynamicResponse replaces placeholders in template with request data.
func generateDynamicResponse(r *http.Request, template string) string {
	// Extract body
	var body map[string]interface{}
	if r.Body != nil && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		json.NewDecoder(r.Body).Decode(&body)
	}

	// Simple replacement: {{.body.key}} -> body[key]
	result := template
	for key, val := range body {
		placeholder := fmt.Sprintf("{{.body.%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", val))
	}
	return result
}

// checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
