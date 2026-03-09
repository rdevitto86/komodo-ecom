package rules

const (
	LevelIgnore = "ignore"
	LevelLenient = "lenient"
	LevelStrict  = "strict"
)

type Headers map[string]struct {
	Type     string 	`yaml:"type,omitempty"`
	Required bool   	`yaml:"required,omitempty"` // defaults to false (optional)
	Value    string 	`yaml:"value,omitempty"`    // exact value or pattern (e.g., "Bearer *")
	Pattern  string 	`yaml:"pattern,omitempty"`  // regex pattern
	Enum     []string `yaml:"enum,omitempty"`     // list of allowed values
	MinLen   int    	`yaml:"min_len,omitempty"`
	MaxLen   int    	`yaml:"max_len,omitempty"`
}

type PathParams map[string]struct {
	Type     string 	`yaml:"type,omitempty"` // "string","int","bool","object","array"
	Required bool   	`yaml:"required,omitempty"` // defaults to false (optional)
	Value    any    	`yaml:"value,omitempty"`
	Enum     []string `yaml:"enum,omitempty"` // list of allowed values
	Pattern  string 	`yaml:"pattern,omitempty"` // regex pattern
	MinLen   int    	`yaml:"min_len,omitempty"`
	MaxLen   int    	`yaml:"max_len,omitempty"`
}

type QueryParams map[string]struct {
	Type     string 	`yaml:"type,omitempty"` // "string","int","bool","object","array"
	Required bool   	`yaml:"required,omitempty"` // defaults to false (optional)
	Value    any    	`yaml:"value,omitempty"`
	Enum     []string `yaml:"enum,omitempty"` // list of allowed values
	Pattern  string 	`yaml:"pattern,omitempty"` // regex pattern
	MinLen   int    	`yaml:"min_len,omitempty"`
	MaxLen   int    	`yaml:"max_len,omitempty"`
}

type Body map[string]struct {
	Type     string   	`yaml:"type,omitempty"` // "string","int","bool","object","array"
	Required bool     	`yaml:"required,omitempty"` // defaults to false (optional)
	Value    any     		`yaml:"value,omitempty"`
	Enum     []string 	`yaml:"enum,omitempty"` // list of allowed values
	Pattern  string   	`yaml:"pattern,omitempty"` // regex pattern
}

type EvalRule struct {
	Toggle       		bool         	`yaml:"toggle,omitempty"`
	Level        		string       	`yaml:"level,omitempty"`
	OriginTypes  		[]string     	`yaml:"originTypes,omitempty"`
	Headers      		Headers  			`yaml:"headers,omitempty"`
	PathParams   		PathParams   	`yaml:"params,omitempty"`
	QueryParams  		QueryParams  	`yaml:"query,omitempty"`
	Body         		Body         	`yaml:"body,omitempty"`
	RequiredVersion int						`yaml:"requiredVersion,omitempty"`
}

type RuleConfig map[string]map[string]EvalRule
