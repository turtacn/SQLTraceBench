package plugin

// DatabasePlugin is the interface that we're exposing as a plugin.
// It matches the interface in plugins/interfaces.go but serves as a bridge for hashicorp/go-plugin.
type DatabasePlugin interface {
	GetName() string
	TranslateQuery(sql string) (string, error)
	ConvertSchema(schema string) (string, error)
}
