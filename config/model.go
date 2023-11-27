package config

// Config armazena as configurações do aplicativo.
type Config struct {
	Assets       AssetsConfig       `yaml:"assets"`
	StaticConfig StaticConfig       `yaml:"static_config"`
	Templates    TemplatesConfig    `yaml:"templates"`
	Database     DatabaseConfig     `yaml:"database"`
	ServerConfig ServerConfig       `yaml:"server"`
	LogConfig    LogConfig          `yaml:"log"`
	Localization LocalizationConfig `yaml:"localization"`
}

type TemplatesConfig struct {
	Path string `yaml:"path"`
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

type ApolloUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PID      string `json:"pid"` // needed for AuthBoss
}

func (u *ApolloUser) GetPID() string {
	return u.PID
}

func (u *ApolloUser) PutPID(pid string) {
	u.PID = pid
}
