package endpoints

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ErrNoIdText(localizer *i18n.Localizer) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrNoIdText",
			Other: "Входные данные пусты",
		},
	})
}
