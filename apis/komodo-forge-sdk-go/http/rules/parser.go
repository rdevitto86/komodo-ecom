package rules

import (
	"fmt"
	logger "komodo-forge-sdk-go/logging/runtime"
	"os"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	ruleMap       map[string]map[string]EvalRule
	patternRoutes []routePattern // patternRoutes is a list of compiled route patterns for templates (/:id or /{id})
	loadOnce      sync.Once
	configLoaded  bool
)

type routePattern struct {
	template  string
	re        *regexp.Regexp
	methods   map[string]EvalRule
	paramKeys []string
}

// Loads validation rules from a file path or embedded data
func LoadConfig(path ...string) bool {
	loadOnce.Do(func() { 
		var data []byte
		var err error
		var source string

		// Try explicit path or env var
		if len(path) > 0 && path[0] != "" {
			data, err = os.ReadFile(path[0])
			source = path[0]
		} else if envPath := os.Getenv("EVAL_RULES_PATH"); envPath != "" {
			data, err = os.ReadFile(envPath)
			source = envPath
		}

		if err != nil || data == nil {
			logger.Error("failed to load validation rules from" + source, err)
			configLoaded = false
			return
		}

		rt, patterns, parseErr := parseConfigFromData(data)
		if parseErr != nil {
			logger.Error("failed to parse validation rules", parseErr)
			configLoaded = false
			return
		}

		ruleMap = rt
		patternRoutes = patterns
		configLoaded = true

		logger.Info("successfully loaded validation rules from: " + source)
	})
	return configLoaded
}

// Loads validation rules from byte data (for embedding in client services)
func LoadConfigWithData(data []byte) {
	loadOnce.Do(func() {
		rt, patterns, err := parseConfigFromData(data)
		if err != nil {
			logger.Error("failed to parse validation rules from embedded config", err)
			configLoaded = false
			return
		}

		ruleMap = rt
		patternRoutes = patterns
		configLoaded = true

		logger.Info("successfully loaded validation rules from embedded config")
	})
}
	

func IsConfigLoaded() bool { return configLoaded && ruleMap != nil }

func GetRule(pKey string, method string) *EvalRule {
	if pKey == "" || method == "" || ruleMap == nil {
		return nil
	}

	np := normalizePath(pKey)

	// Direct match
	if rules, ok := ruleMap[np]; ok {
		if rule, exists := rules[method]; exists {
			return &rule
		}
	}

	for _, rp := range patternRoutes {
		if rp.re.MatchString(np) {
			if rule, exists := rp.methods[method]; exists {
				return &rule
			}
		}
	}
	return nil
}

func GetRules() RuleConfig { return ruleMap }

func validateAndNormalizeConfig(cfg RuleConfig) error {
	for _, methods := range cfg {
		for method, rule := range methods {
			if rule.Headers == nil {
				rule.Headers = make(Headers)
			}
			if rule.QueryParams == nil {
				rule.QueryParams = make(QueryParams)
			}
			if rule.PathParams == nil {
				rule.PathParams = make(PathParams)
			}
			if rule.Body == nil {
				rule.Body = make(Body)
			}

			methods[method] = rule
		}
	}
	return nil
}

// Parses YAML data and returns the rule map plus compiled templates.
func parseConfigFromData(data []byte) (map[string]map[string]EvalRule, []routePattern, error) {
	var root struct {
		Rules RuleConfig `yaml:"rules"`
	}

	if err := yaml.Unmarshal(data, &root); err != nil {
		logger.Error("failed to parse validation rules from embedded config", err)
		return nil, nil, fmt.Errorf("invalid yaml: %w", err)
	}

	cfg := root.Rules
	
	// Validate and normalize the configuration
	if err := validateAndNormalizeConfig(cfg); err != nil {
		logger.Error("validation rules configuration is invalid", err)
		return nil, nil, err
	}
	
	patterns := make([]routePattern, 0)

	// Build patterns for templates with dynamic segments
	for tpl, methods := range cfg {
		if strings.Contains(tpl, ":") || strings.Contains(tpl, "{") || strings.Contains(tpl, "*") {
			reStr, keys := templateToRegex(tpl)
			re, err := regexp.Compile("^" + reStr + "$")
			if err != nil {
				logger.Error("invalid route pattern " + tpl, err)
				return nil, nil, fmt.Errorf("invalid route pattern %s: %w", tpl, err)
			}

			patterns = append(patterns, routePattern{
				template:  tpl,
				re:        re,
				methods:   methods,
				paramKeys: keys,
			})
		}
	}

	// Sort patterns by specificity
	specificityScore := func(tpl string) int {
		parts := strings.Split(strings.TrimPrefix(tpl, "/"), "/")
		literal, wild := 0, 0

		for _, p := range parts {
			if p == "*" || strings.HasPrefix(p, ":") || (strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}")) {
				wild++
			} else if p != "" {
				literal++
			}
		}
		return literal * 10 - wild
	} // end

	for i := 0; i < len(patterns); i++ {
		for j := i + 1; j < len(patterns); j++ {
			if specificityScore(patterns[j].template) > specificityScore(patterns[i].template) {
				patterns[i], patterns[j] = patterns[j], patterns[i]
			}
		}
	}

	return cfg, patterns, nil
}

// Converts a route template into a regex string
func templateToRegex(tpl string) (string, []string) {
	// ensure we only operate on path part (strip query if present)
	if idx := strings.Index(tpl, "?"); idx != -1 {
		tpl = tpl[:idx]
	}
	// trim trailing slash handling
	tpl = strings.TrimSuffix(tpl, "/")

	parts := strings.Split(strings.TrimPrefix(tpl, "/"), "/")
	regexParts := make([]string, 0, len(parts))
	keys := make([]string, 0)

	for _, p := range parts {
		if p == "*" {
			regexParts = append(regexParts, ".*")
			continue
		}
		// :param or {param}
		if strings.HasPrefix(p, ":") {
			key := strings.TrimPrefix(p, ":")
			keys = append(keys, key)
			regexParts = append(regexParts, `(?P<`+key+`>[^/]+)`)
			continue
		}
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(p, "{"), "}")
			keys = append(keys, key)
			regexParts = append(regexParts, `(?P<`+key+`>[^/]+)`)
			continue
		}
		// literal segment - escape regexp metacharacters
		regexParts = append(regexParts, regexp.QuoteMeta(p))
	}
	return "/" + strings.Join(regexParts, "/"), keys
}

// Strips version prefixes and ensures leading slash
func normalizePath(p string) string {
	if p == "" {
		return p
	}
	// drop query
	if idx := strings.Index(p, "?"); idx != -1 {
		p = p[:idx]
	}
	p = strings.TrimSpace(p)
	if p == "" {
		return "/"
	}

	// remove trailing slash (but keep root)
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}

	// strip version prefix like /v1 or /v1.2
	trimmed := strings.TrimPrefix(p, "/")
	segs := strings.Split(trimmed, "/")

	if len(segs) > 0 && len(segs[0]) > 0 && segs[0][0] == 'v' {
		// basic check: next characters are digits or digit+dot
		if regexp.MustCompile(`^v[0-9]`).MatchString(segs[0]) {
			// remove first segment
			segs = segs[1:]
			p = "/" + strings.Join(segs, "/")
		}
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

// Finds the first pattern that matches and extracts params
func matchRouteAndExtractParams(path string) (*routePattern, map[string]string) {
	np := normalizePath(path)

	for _, rp := range patternRoutes {
		if !rp.re.MatchString(np) {
			continue
		}
		matches := rp.re.FindStringSubmatch(np)
		names := rp.re.SubexpNames()
		params := make(map[string]string)

		for i, name := range names {
			if i != 0 && name != "" && i < len(matches) {
				params[name] = matches[i]
			}
		}
		return &rp, params
	}
	return nil, nil
}
