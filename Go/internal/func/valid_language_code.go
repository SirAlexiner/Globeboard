// Package _func provides self-defined functions for use in the API.
package _func

import "errors"

// IsValidLanguages checks if provided languages are valid based on specified criteria.
// It takes a list of languages, a flag to allow empty input, and a flag to allow multiple languages.
// It returns an error if any criteria are not met.
func IsValidLanguages(languages []string, allowEmpty bool, allowMultiple bool) error {
	err, done := isEmptyDisallowedEmpty(languages, allowEmpty)
	if done {
		return err
	}

	err2, done2 := isMultipleAllowedAboveOne(languages, allowMultiple)
	if done2 {
		return err2
	}

	languageMap := make(map[string]bool)
	for _, lang := range languages {
		if allowEmpty {
			if len(lang) > 2 || len(lang) == 1 {
				err := errors.New("language parameter must be in the format of two-letter language code")
				return err
			}
		} else if len(lang) != 2 {
			err := errors.New("language parameter must be in the format of two-letter language code and not empty")
			return err
		}
		if languageMap[lang] {
			err := errors.New("same language parameter cannot be given twice")
			return err
		}
		languageMap[lang] = true
	}
	return nil
}

// isEmptyDisallowedEmpty checks if empty input is disallowed and the input is empty.
// It returns an error if so, along with a boolean indicating if the check is completed.
func isEmptyDisallowedEmpty(languages []string, allowEmpty bool) (error, bool) {
	if !allowEmpty {
		if len(languages) == 0 {
			err := errors.New("no language parameters specified")
			return err, true
		}
	}
	return nil, false
}

// isMultipleAllowedAboveOne checks if multiple languages are disallowed and more than one language is provided.
// It returns an error if so, along with a boolean indicating if the check is completed.
func isMultipleAllowedAboveOne(languages []string, allowMultiple bool) (error, bool) {
	if !allowMultiple {
		if len(languages) > 1 {
			err := errors.New("only one language parameter can be provided")
			return err, true
		}
	}
	return nil, false
}
