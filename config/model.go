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

type LoginRequest struct {
	Email    string `JSON:"email"`
	Password string `JSON:"password"`
}

type LogConfig struct {
	LogLevel    string `yaml:"logLevel"`
	LogPath     string `yaml:"logPath"`
	MaxSizeMB   int    `yaml:"maxSizeMB"`
	MaxBackups  int    `yaml:"maxBackups"`
	MaxAgeDays  int    `yaml:"maxAgeDays"`
	LogToStdout bool   `yaml:"logToStdout"`
}

// AUTHBOSS INTERFACES
type ApolloUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
	// Password string `json:"password"`
	// needed for AuthBoss
	// PID             string    `json:"pid"`
	// RecoverSelector string    `json:"recoverselector"`
	// RecoverVerifier string    `json:"recoververifier"`
	// RecoverExpiry   time.Time `json:"recoveryexpiry"`
	// authboss.User
	// authboss.AuthableUser
	// authboss.RecoverableUser
}

//  AB - AUTHBOSS USER

// // func (u *ApolloUser) GetPID() string {
// // 	return u.PID
// // }

// // func (u *ApolloUser) PutPID(pid string) {
// // 	u.PID = pid
// // }

// // AB - AUTHABLE USER

// func (u *ApolloUser) GetPassword() (password string) {
// 	return u.Password
// }

// func (u *ApolloUser) PutPassword(password string) {
// 	u.Password = password
// }

// // AB - RECOVERABLE USER

// func (u *ApolloUser) GetEmail() (email string) {
// 	return u.Email
// }

// func (u *ApolloUser) GetRecoverSelector() (selector string) {
// 	return u.RecoverSelector
// }

// func (u *ApolloUser) GetRecoverVerifier() (verifier string) {
// 	return u.RecoverVerifier
// }

// func (u *ApolloUser) GetRecoverExpiry() (expiry time.Time) {
// 	return u.RecoverExpiry
// }

// func (u *ApolloUser) PutEmail(email string) {
// 	u.Email = email
// }

// func (u *ApolloUser) PutRecoverSelector(selector string) {
// 	u.RecoverSelector = selector
// }

// func (u *ApolloUser) PutRecoverVerifier(verifier string) {
// 	u.RecoverVerifier = verifier
// }

// func (u *ApolloUser) PutRecoverExpiry(expiry time.Time) {
// 	u.RecoverExpiry = expiry
// }
