package plugins

type DatabasePlugin interface {
	GetName() string
	GetVersion() string
	ValidateConnection() error
	SchemaConverter() SchemaConverter // return nil if not implemented
	QueryTranslator() QueryTranslator // return nil if not implemented
	BenchmarkRunner() BenchmarkRunner // return nil if not implemented
}

type SchemaConverter interface {
	ConvertSchema(sql string, from, to string) (string, error)
}

type QueryTranslator interface {
	TranslateQuery(query, from, to string) (string, error)
}

type BenchmarkRunner interface {
	RunBench(queries []string, qps float64, duration int) error
}

type Registry struct {
	plugins map[string]DatabasePlugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]DatabasePlugin)}
}

func (r *Registry) Register(name string, p DatabasePlugin) {
	r.plugins[name] = p
}

//Personal.AI order the ending
