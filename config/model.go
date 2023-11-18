package config

// Config armazena as configurações do aplicativo.
type Config struct {
	Assets        AssetsConfig       `yaml:"assets"`
	StaticConfig  StaticConfig       `yaml:"static_config"`
	TemplatesPath string             `yaml:"templates_path"`
	Database      DatabaseConfig     `yaml:"database"`
	ServerConfig  ServerConfig       `yaml:"server"`
	LogConfig     LogConfig          `yaml:"log"`
	Localization  LocalizationConfig `yaml:"localization"`
}

// AssetsConfig armazena as configurações de ativos.
type AssetsConfig struct {
	Img string `yaml:"images"`
	CSS string `yaml:"css"`
	JS  string `yaml:"js"`
}

type LocalizationConfig struct {
	LocalesPath string `yaml:"locales_path"`
}

// StaticConfig armazena as configurações estáticas.
type StaticConfig struct {
	Path   string `yaml:"path"`
	Prefix string `yaml:"prefix"`
}

// DatabaseConfig armazena as configurações do banco de dados.
type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type LogConfig struct {
	LogLevel    string `yaml:"logLevel"`
	LogPath     string `yaml:"logPath"`
	MaxSizeMB   int    `yaml:"maxSizeMB"`
	MaxBackups  int    `yaml:"maxBackups"`
	MaxAgeDays  int    `yaml:"maxAgeDays"`
	LogToStdout bool   `yaml:"logToStdout"`
}
