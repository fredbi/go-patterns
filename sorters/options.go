package sorters

import (
	"golang.org/x/text/language"
)

// CompareOption sets options for sorting.
type CompareOption func(*compareOptions)

type compareOptions struct {
	Locale *language.Tag
}

// WithLocale sets the language specific collation for ordering strings.
func WithLocale(locale string) CompareOption {
	return func(opts *compareOptions) {
		opts.Locale = tagFromLocale(locale)
	}
}

// WithLocaleTag sets the language specific collation for ordering strings.
func WithLocaleTag(localeTag language.Tag) CompareOption {
	return func(opts *compareOptions) {
		tag := localeTag
		opts.Locale = &tag
	}
}

func compareOptionsWithDefault(opts []CompareOption) *compareOptions {
	options := &compareOptions{}
	for _, apply := range opts {
		apply(options)
	}

	return options
}

func tagFromLocale(locale string) *language.Tag {
	tag, err := language.Parse(locale)
	if err != nil {
		// defaults to english if we can't parse
		tag = language.AmericanEnglish
	}

	return &tag
}
