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
)

var bundle *i18n.Bundle

func InitLocalization(config config.Config) {
	// Initialize the Bundle with the default language.
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load message files.
	err := LoadLocaleFiles(config.Localization.LocalesPath)
	if err != nil {
		logrus.Error("Error loading locale files: ", err)
		return
	}

}

func LoadLocaleFiles(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		logrus.Debug("Error reading locale files: ", err)
		return err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			logrus.Debug("Loading locale file: ", f.Name()
			fullPath := filepath.Join(path, f.Name())
			bundle.MustLoadMessageFile(fullPath)
		}
	}
	return nil
}

// LocalizationMiddleware detects the user's language and initializes a localizer.
func LocalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Detect user's language from the Accept-Language header.
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "en" // default language
		}

		// Create a localizer for the detected language.
		localizer := i18n.NewLocalizer(bundle, lang)

		// Set localizer in Gin's context for use in handlers.
		c.Set("localizer", localizer)

		c.Next()
	}
}

// GetLocalizer retrieves the localizer from Gin's context.
func GetLocalizer(c *gin.Context) *i18n.Localizer {
	localizer, _ := c.Get("localizer")
	return localizer.(*i18n.Localizer)
}

// LocalizeStrings localizes a slice of message IDs and returns a map.
func LocalizeStrings(localizer *i18n.Localizer, messageIDs []string) map[string]string {
	localizedStrings := make(map[string]string)
	for _, id := range messageIDs {
		localizedStrings[id] = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: id})
	}
	return localizedStrings
}

func LocalizeWithArgs(localizer *i18n.Localizer, messageID string, args ...string) string {
	templateData := make(map[string]string)
	for i, arg := range args {
		key := fmt.Sprintf("Arg%d", i+1)
		templateData[key] = arg
	}

	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}
