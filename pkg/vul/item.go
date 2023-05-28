package vul

import (
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	printer  = message.NewPrinter(language.English)
	replacer = strings.NewReplacer("[", "", "]", "")
)

type Item struct{}

func (i *Item) FormatTime(s string, v time.Time) string {
	return clean(v.Format(s))
}

func (i *Item) Printf(m message.Reference, v ...interface{}) string {
	return clean(printer.Sprintf(m, v))
}

func (i *Item) Print(v ...interface{}) string {
	return clean(printer.Sprintln(v))
}

// HACK: remove [ and ] from string
func clean(s string) string {
	return replacer.Replace(s)
}
