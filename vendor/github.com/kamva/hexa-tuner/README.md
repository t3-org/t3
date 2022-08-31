#### hexa-tuner tune Hexa services.

#### Install
```
go get github.com/kamva/hexa-tuner
```

#### Used config variables:
```
// Translator config variables
translate.files (optional) : []string translation files.
translate.fallback.langs (optional,default:en): []string fallback langues

e.g environtment variable list in viper driver of cofig:
TRANSLATE.FILES="fa en"
TRANSLATE.FALLBACK.LANGS="en fa"
```

#### Todo:
- [ ] Write Tests
- [ ] Add badges to readme.
- [ ] CI 