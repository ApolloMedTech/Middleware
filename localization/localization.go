package localization

import (
	"encoding/json"
	"fmt"
	"github.com/ApolloMedTech/Middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"sync"
)

var (
	locales LocaleData
	mutex   sync.RWMutex
)
var bundle *i18n.Bundle

func InitLocalization() {
	cfg := config.GetConfig().Localization
	locales = make(LocaleData)
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	LoadLocaleFiles(cfg.LocalesPath)
}

func LoadLocaleFiles(path string) {
	mutex.Lock()
	defer mutex.Unlock()

	files, err := os.ReadDir(path)
	if err != nil {
		logrus.Error("Error reading locale files: ", err)
		return
	}

	for _, f := range files {
		logrus.Debugf("Loading locale file: %v", f.Name())
		if filepath.Ext(f.Name()) == ".json" {
			fullPath := filepath.Join(path, f.Name())
			loadLocaleFile(fullPath)
		}
	}
}

func loadLocaleFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		logrus.Error("Error opening locale file: ", err)
		return
	}
	defer file.Close()

	var data Section
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		logrus.Error("Error decoding locale file: ", err)
		return
	}

	locale := filepath.Base(path)
	locale = locale[:len(locale)-len(filepath.Ext(locale))] // Remove extension
	locales[locale] = data
	logrus.Debugf("Loaded locale file: %v", locales)
}

func LocalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, acceptLang)

		// This will find the best match among loaded languages
		lang, _, _ := localizer.LocalizeWithTag(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "language_tag"}})

		c.Set("localizer", lang)
		c.Next()
	}
}

func GetLocalizer(c *gin.Context) string {
	localizer, _ := c.Get("localizer")
	return localizer.(string)
}

func LocalizeSection(lang, section string) (map[string]string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	// Check if the language exists
	sections, ok := locales[lang]
	if !ok {
		return nil, fmt.Errorf("language not found")
	}

	// Check if the section exists
	messages, ok := sections[section]
	if !ok {
		return nil, fmt.Errorf("section not found")
	}

	return messages, nil
}

func ReloadLocalization(config config.Config) {
	LoadLocaleFiles(config.Localization.LocalesPath)
}
