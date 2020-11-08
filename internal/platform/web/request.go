package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

// Holds the settings and caches for validating request struct field values.
var validate = validator.New()

//translator is a cache for locale and translation information.
var translator *ut.UniversalTranslator

func init() {
	// Locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale.
	// Provide one ore more locales to additional language support.
	translator := ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Decodes request body into dest.
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at the
		// Accept-Language header if you intend to support multiple languages.
		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}
