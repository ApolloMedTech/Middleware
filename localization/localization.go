package localization

//TODO change the package to localManager

import (
	"encoding/json"
	"github.com/ApolloMedTech/Middleware/config"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"time"
)

var bundle *i18n.Bundle
var trie *Trie // Declare trie as a global variable

type Message struct {
	ID    string `json:"id"`
	Other string `json:"other"`
}

func InitLocalization(config config.LocalizationConfig) {
	// Initialize the Bundle with the default language.
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	trie = newTrie() // Initialize the trie

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
			file, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			defer file.Close()

			var messages []Message
			if err := json.NewDecoder(file).Decode(&messages); err != nil {
				return err
			}

			for _, message := range messages {
				// Insert the message ID into the trie
				trie.insert(message.ID)

				// Add the message to the bundle
				bundle.AddMessages(language.English, &i18n.Message{
					ID:    message.ID,
					Other: message.Other,
				})
			}
		}
	}
	return nil
}

/*
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
}*/

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

// LocalizePrefixStrings Modify the  function to use SearchPrefix to find all keys that start with the given prefix, and then localize those strings.
func LocalizePrefixStrings(c *gin.Context, partialID string) map[string]string {
	localizer := GetLocalizer(c)
	localizedStrings := make(map[string]string)
	ids := trie.searchPrefix(partialID)
	for _, id := range ids {
		localizedStrings[id] = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: id})
	}
	return localizedStrings
}
