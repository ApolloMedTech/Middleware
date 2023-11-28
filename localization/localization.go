package localization

import (
	"encoding/json"
	"github.com/ApolloMedTech/Middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"time"
)

var bundle *i18n.Bundle

func InitLocalization(config config.LocalizationConfig) {
	// Initialize the Bundle with the default language.
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load message files.
	err := LoadLocaleFiles(config.LocalesPath)
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
			logrus.Debug("Loading locale file: ", f.Name())
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

// LocalizePlural is a helper function that localizes a message ID with pluralization
func LocalizePlural(c *gin.Context, messageID string, count int) string {
	localizedString, err := GetLocalizer(c).Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageID,
		},
		PluralCount: count,
	})

	if err != nil {
		logrus.Error("Error localizing string: ", err)
		return "" // or return a default message
	}

	return localizedString
}

// LocalizeDate formats a date according to the language in the localizer.
func LocalizeDate(c *gin.Context, year int, month int, day int) string {
	localizer := GetLocalizer(c)
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	dateString := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "DateFormat",
		},
		TemplateData: map[string]interface{}{
			"Date": t,
		},
	})
	return dateString
}

// LocalizeCurrency formats a currency value according to the language in the localizer and currency unit.
func LocalizeCurrency(c *gin.Context, value float64, currencyUnit currency.Unit) string {
	localizer := GetLocalizer(c)
	currencyString := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Currency",
			Other: currency.Symbol(currencyUnit.Amount(value)),
		},
	})
	return currencyString
}
