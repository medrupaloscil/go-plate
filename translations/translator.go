package translations

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var Translations = make(map[string]map[string]string)

func LoadTranslations() error {
	languages := []string{"en", "fr"}

	for _, lang := range languages {
		filePath := fmt.Sprintf("./translations/%s.yaml", lang)
		file, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error loading %s: %w", filePath, err)
		}

		var translations map[string]string
		if err := yaml.Unmarshal(file, &translations); err != nil {
			return fmt.Errorf("error parsing %s: %w", filePath, err)
		}

		Translations[lang] = translations
	}

	return nil
}

func T(lang string, key string) string {
	if translations, ok := Translations[lang]; ok {
		if value, exists := translations[key]; exists {
			return value
		}
	}

	if value, exists := Translations["en"][key]; exists {
		return value
	}

	return key
}