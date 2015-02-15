package ar

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

type Validator struct {
	rule *Validation
}

func NewValidator(rule *Validation) Validator {
	return Validator{rule}
}

type CustomValidator func(errors *Errors)

func (v Validator) IsValid(value interface{}) (bool, []error) {
	result := true
	errors := []error{}
	if v.rule.presence != nil {
		if ok, err := v.isPersistent(value); !ok {
			result = false
			errors = append(errors, err)
		}
	}
	if v.rule.format != nil {
		if ok, err := v.isFormatted(value); !ok {
			result = false
			errors = append(errors, err)
		}
	}
	if v.rule.length != nil {
		if ok, err := v.isMinimumLength(value); !ok {
			result = false
			errors = append(errors, err)
		}
		if ok, err := v.isMaximumLength(value); !ok {
			result = false
			errors = append(errors, err)
		}
		if ok, err := v.isLength(value); !ok {
			result = false
			errors = append(errors, err)
		}
		if ok, err := v.inLength(value); !ok {
			result = false
			errors = append(errors, err)
		}
	}
	if v.rule.numericality != nil {
		if ok, err := v.isNumericality(value); !ok {
			result = false
			errors = append(errors, err)
		} else {
			if ok, err := v.greaterThan(value); !ok {
				result = false
				errors = append(errors, err)
			}
			if ok, err := v.greaterThanOrEqualTo(value); !ok {
				result = false
				errors = append(errors, err)
			}
			if ok, err := v.equalTo(value); !ok {
				result = false
				errors = append(errors, err)
			}
			if ok, err := v.lessThan(value); !ok {
				result = false
				errors = append(errors, err)
			}
			if ok, err := v.lessThanOrEqualTo(value); !ok {
				result = false
				errors = append(errors, err)
			}
			if ok, err := v.odd(value); !ok {
				result = false
				errors = append(errors, err)
			}
			if ok, err := v.even(value); !ok {
				result = false
				errors = append(errors, err)
			}
		}
	}
	return result, errors
}

func (v Validator) isPersistent(value interface{}) (bool, error) {
	if IsZero(value) {
		return false, fmt.Errorf("%s", v.rule.presence.message)
	}
	return true, nil
}

func (v Validator) isFormatted(value interface{}) (bool, error) {
	with := v.rule.format.with
	if with.regexp == "" {
		return true, nil
	}
	s, ok := value.(string)
	if !ok {
		return false, fmt.Errorf(with.message)
	}
	match, _ := regexp.MatchString(with.regexp, s)
	if !match {
		return false, fmt.Errorf(with.message)
	}
	return true, nil
}

func (v Validator) isMinimumLength(value interface{}) (bool, error) {
	minimum := v.rule.length.minimum
	if minimum.number == 0 {
		return true, nil
	}
	result := utf8.RuneCountInString(fmt.Sprintf("%s", value)) <= minimum.number
	if !result {
		return false, fmt.Errorf(minimum.message, minimum.number)
	}
	return true, nil
}

func (v Validator) isMaximumLength(value interface{}) (bool, error) {
	maximum := v.rule.length.maximum
	if maximum.number == 0 {
		return true, nil
	}
	result := utf8.RuneCountInString(fmt.Sprintf("%s", value)) >= maximum.number
	if !result {
		return false, fmt.Errorf(maximum.message, maximum.number)
	}
	return true, nil
}

func (v Validator) isLength(value interface{}) (bool, error) {
	is := v.rule.length.is
	if is.number == 0 {
		return true, nil
	}
	result := utf8.RuneCountInString(fmt.Sprintf("%s", value)) == is.number
	if !result {
		return false, fmt.Errorf(is.message, is.number)
	}
	return true, nil
}

func (v Validator) inLength(value interface{}) (bool, error) {
	ok, err := v.isMinimumLength(value)
	if !ok {
		return false, err
	}
	ok, err = v.isMaximumLength(value)
	if !ok {
		return false, err
	}
	return true, nil
}

func (v Validator) isNumericality(value interface{}) (bool, error) {
	numericality := v.rule.numericality
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true, nil
	case float32, float64:
		if numericality.onlyInteger.bool {
			return false, fmt.Errorf(numericality.onlyInteger.message)
		} else {
			return true, nil
		}
	}
	return false, fmt.Errorf(numericality.message)
}

func (v Validator) greaterThan(value interface{}) (bool, error) {
	greaterThan := v.rule.numericality.greaterThan
	i, _ := value.(int)
	if i <= greaterThan.number {
		return false, fmt.Errorf(greaterThan.message, value)
	}
	return true, nil
}

func (v Validator) greaterThanOrEqualTo(value interface{}) (bool, error) {
	greaterThanOrEqualTo := v.rule.numericality.greaterThanOrEqualTo
	i, _ := value.(int)
	if i < greaterThanOrEqualTo.number {
		return false, fmt.Errorf(greaterThanOrEqualTo.message, value)
	}
	return true, nil
}

func (v Validator) equalTo(value interface{}) (bool, error) {
	equalTo := v.rule.numericality.equalTo
	i, _ := value.(int)
	if i != equalTo.number {
		return false, fmt.Errorf(equalTo.message, value)
	}
	return true, nil
}

func (v Validator) lessThan(value interface{}) (bool, error) {
	lessThan := v.rule.numericality.lessThan
	i, _ := value.(int)
	if i >= lessThan.number {
		return false, fmt.Errorf(lessThan.message, value)
	}
	return true, nil
}

func (v Validator) lessThanOrEqualTo(value interface{}) (bool, error) {
	lessThanOrEqualTo := v.rule.numericality.lessThanOrEqualTo
	i, _ := value.(int)
	if i > lessThanOrEqualTo.number {
		return false, fmt.Errorf(lessThanOrEqualTo.message, value)
	}
	return true, nil
}

func (v Validator) odd(value interface{}) (bool, error) {
	odd := v.rule.numericality.odd
	if !odd.bool {
		return true, nil
	}
	i, _ := value.(int)
	if i%2 == 0 {
		return false, fmt.Errorf(odd.message)
	}
	return true, nil
}

func (v Validator) even(value interface{}) (bool, error) {
	if ok, _ := v.odd(value); !ok {
		return false, fmt.Errorf(v.rule.numericality.even.message)
	}
	return true, nil
}
