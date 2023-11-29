package localization

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
	trie = newTrie() // Initialize the trie
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

type TrieNode struct {
	children map[rune]*TrieNode
	endOfKey bool
}

type Trie struct {
	root *TrieNode
}

func newTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

func (t *Trie) insert(key string) {
	node := t.root
	for _, char := range key {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.endOfKey = true
}

func (t *Trie) searchPrefix(prefix string) []string {
	node := t.root
	for _, char := range prefix {
		if _, ok := node.children[char]; !ok {
			return nil
		}
		node = node.children[char]
	}
	var results []string
	t.collectAllKeys(node, prefix, &results)
	return results
}

func (t *Trie) collectAllKeys(node *TrieNode, prefix string, results *[]string) {
	if node.endOfKey {
		*results = append(*results, prefix)
	}
	for char, childNode := range node.children {
		t.collectAllKeys(childNode, prefix+string(char), results)
	}
}
