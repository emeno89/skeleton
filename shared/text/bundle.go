package text

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	bundle.MustLoadMessageFile("./translation/active.en.yaml")

	return bundle
}
