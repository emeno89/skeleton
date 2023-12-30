package graph

import "github.com/nicksnyder/go-i18n/v2/i18n"

func ExampleTitle(localizer *i18n.Localizer, id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ExampleTitle",
			Other: "Пользователь {{.Id}}",
		},
		TemplateData: map[string]any{
			"Id": id,
		},
	})
}
