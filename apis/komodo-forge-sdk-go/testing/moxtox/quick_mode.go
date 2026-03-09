package moxtox

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// buildHashLookupMap builds HashLookupMap for quick mode.
func (cnfg *MoxtoxConfig) buildHashLookupMap() {
	cnfg.HashLookupMap = make(map[string]map[string]map[string]Scenario)

	for path, rawMapping := range cnfg.Mappings {
		mappingData, ok := rawMapping.(map[interface{}]interface{})
		if !ok { continue }

		mapping := parseMapping(mappingData)
		cnfg.HashLookupMap[path] = make(map[string]map[string]Scenario)

		for method, methodData := range mapping.Methods {
			cnfg.HashLookupMap[path][method] = make(map[string]Scenario)

			for _, scenario := range methodData.Scenarios {
				hash := hashConditions(scenario.Conditions)
				cnfg.HashLookupMap[path][method][hash] = scenario
			}
		}
	}
}

// hashConditions creates a unique hash for conditions map.
func hashConditions(conditions map[string]interface{}) string {
	if len(conditions) == 0 { return "default" }

	// Sort keys for consistency
	keys := make([]string, 0, len(conditions))
	for k := range conditions { keys = append(keys, k) }

	sort.Strings(keys)
	var sorted []string
	for _, k := range keys {
		sorted = append(sorted, fmt.Sprintf("%s:%v", k, conditions[k]))
	}

	data := strings.Join(sorted, "|")
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// matchesRequestHash checks if the request matches any scenario using hash-based lookup.
// Returns the matching scenario if found.
func matchesRequestHash(req *http.Request) (Scenario, bool) {
	path := req.URL.Path
	method := req.Method

	// Extract conditions from request and hash them
	requestConditions := extractRequestConditions(req)
	hash := hashConditions(requestConditions)

	if methodMap, ok := config.HashLookupMap[path]; ok {
		if scenarioMap, ok := methodMap[method]; ok {
			if scenario, ok := scenarioMap[hash]; ok {
				if config.Debug || scenario.Log {
					fmt.Printf("[::Moxtox::] matched scenario '%s' for %s %s (quick mode)\n", scenario.Name, method, path)
				}
				return scenario, true
			}
		}
		// Check "*" method
		if scenarioMap, ok := methodMap["*"]; ok {
			if scenario, ok := scenarioMap[hash]; ok {
				if config.Debug || scenario.Log {
					fmt.Printf("[::Moxtox::] matched scenario '%s' for %s %s (quick mode)\n", scenario.Name, "*", path)
				}
				return scenario, true
			}
		}
	}
	if config.Debug {
		fmt.Printf("[::Moxtox::] no match for %s %s (quick mode)\n", method, path)
	}
	return Scenario{}, false
}

// extractRequestConditions extracts conditions from the request for hashing.
func extractRequestConditions(req *http.Request) map[string]interface{} {
	conditions := make(map[string]interface{})

	// Body (simplified)
	if req.Body != nil && strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		var body map[string]interface{}
		json.NewDecoder(req.Body).Decode(&body)
		conditions["body"] = body
	}

	// Query
	query := make(map[string]interface{})
	for k, v := range req.URL.Query() {
		if len(v) > 0 {
			query[k] = v[0]
		}
	}
	if len(query) > 0 {
		conditions["query"] = query
	}

	// Headers
	headers := make(map[string]interface{})
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[strings.ToLower(k)] = v[0]
		}
	}
	if len(headers) > 0 {
		conditions["headers"] = headers
	}

	// Path
	conditions["path"] = map[string]interface{}{"path": req.URL.Path}

	return conditions
}
