package atlas

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (a *Atlas) Validator(data url.Values) *Validation {
	return &Validation{
		Data:   data,
		Errors: make(map[string]string),
	}
}

func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validation) Has(field string, r *http.Request) bool {
	value := r.Form.Get(field)
	if value == "" {
		return false
	}

	return true
}

func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "is required")
		}
	}
}

func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "invalid email address")
	}
}

func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "this field must be an intager")
	}
}

func (v *Validation) IsDateIso(field, value string) {
	_, err := time.Parse("2024-02-12", value)
	if err != nil {
		v.AddError(field, "this feild must be a date in the form of YYYY-MM-DD")
	}
}

func (v *Validation) NoWhiteSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "spaces are not permitted for this field")
	}
}
