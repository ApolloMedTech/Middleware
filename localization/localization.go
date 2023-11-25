package localization

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Global variable to hold localization data
var (
	localizationData LocalizationData
	bundle           *i18n.Bundle // i18n bundle as a global variable
)

// Initialize loads localization data from JSON files and initializes the i18n bundle
func Initialize(localesPath string) {
	// Load the localization data into the global variable
	localizationData = LoadLocalizationData(localesPath)

	// Correctly initialize the global bundle variable
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load all JSON files from the locales folder for i18n
	files, err := os.ReadDir(localesPath)
	if err != nil {
		log.Fatalf("Unable to read locales directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(localesPath, file.Name())
			bundle.MustLoadMessageFile(filePath)
		}
	}
}

// LocalizationMiddleware sets the localizer in the Gin context
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

// LoadLocalizationData loads localization data from JSON files in the given directory
func LoadLocalizationData(localesPath string) LocalizationData {
	allData := make(LocalizationData)

	files, err := ioutil.ReadDir(localesPath)
	if err != nil {
		log.Fatalf("Unable to read locales directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(localesPath, file.Name())
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Unable to read file: %v", err)
			}

			var jsonData map[string]SectionData
			if err := json.Unmarshal(fileData, &jsonData); err != nil {
				log.Fatalf("Error parsing JSON: %v", err)
			}

			for section, data := range jsonData {
				if _, exists := allData[section]; !exists {
					allData[section] = data
				} else {
					for key, value := range data {
						allData[section][key] = value
					}
				}
			}
		}
	}

	return allData
}

// GetSection returns all keys and values for a specific section
func GetSection(section string) SectionData {
	return localizationData[section]
}

// LocalizeSection returns a map of localized key-value pairs for a given section
func LocalizeSection(localizer *i18n.Localizer, section string) (map[string]string, error) {
	keys := GetSectionKeys(section)
	localizedSection := make(map[string]string)

	for _, key := range keys {
		localizedValue, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
		if err != nil {
			return nil, err
		}
		localizedSection[key] = localizedValue
	}

	return localizedSection, nil
}

// In your localization package

// GetSectionKeys returns all keys for a specific section
func GetSectionKeys(section string) []string {
	var keys []string
	for key := range localizationData[section] {
		keys = append(keys, section+"."+key)
	}
	return keys
}

// LocalizeSectionWithContext localizes a section based on the Gin context
func LocalizeSectionWithContext(c *gin.Context, section string) (map[string]string, error) {
	localizerObj, exists := c.Get("localizer")
	if !exists {
		return nil, errors.New("localizer not found in context")
	}
	localizer := localizerObj.(*i18n.Localizer)

	return LocalizeSection(localizer, section)
}
