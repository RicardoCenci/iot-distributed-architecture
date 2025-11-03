package logger

type Config struct {
	Level  string
	Source SourceConfig
}

type SourceConfig struct {
	Enabled  bool
	Relative bool
	AsJSON   bool
}
