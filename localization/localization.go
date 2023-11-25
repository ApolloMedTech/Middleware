package localization

import (
	"encoding/json"
	"fmt"
	"github.com/ApolloMedTech/Middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

var (
	locales LocaleData
	mutex   sync.RWMutex
)

func InitLocalization(config config.Config) {
	locales = make(LocaleData)
	LoadLocaleFiles(config.Localization.LocalesPath)
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
}

func LocalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en"
		}

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
