package localization

// SectionData represents the data of a specific section in your JSON files
type SectionData map[string]string

// LocalizationData holds all sections with their respective translations
type LocalizationData map[string]SectionData
