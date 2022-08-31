package huner

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"strings"
)

type TranslateOpts struct {
	Files         []string
	FallbackLangs []string
}

// return new translator service.
func NewTranslator(pathPrefix string, cfg TranslateOpts) hexa.Translator {
	defaultLang := language.English

	if len(cfg.FallbackLangs) >= 1 {
		defaultLang = language.MustParse(cfg.FallbackLangs[0])
	}

	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	loadLangFiles(bundle, pathPrefix, cfg.Files)

	localizer := i18n.NewLocalizer(bundle, cfg.FallbackLangs...)

	return hexatranslator.NewI18nDriver(bundle, localizer, cfg.FallbackLangs)
}

func loadLangFiles(bundle *i18n.Bundle, prefix string, files []string) {
	for _, file := range files {
		f := fmt.Sprintf("%s/%s.toml", strings.TrimRight(prefix, "/"), strings.Trim(file, "/"))
		bundle.MustLoadMessageFile(f)
	}
}
