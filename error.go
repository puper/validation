// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"encoding/json"
	"fmt"

	"github.com/puper/orderedmap"
)

type (
	// Errors represents the validation errors that are indexed by struct field names, map or slice keys.
	Errors struct {
		orderedmap.OrderedMap
	}

	// InternalError represents an error that should NOT be treated as a validation error.
	InternalError interface {
		error
		InternalError() error
	}

	internalError struct {
		error
	}
)

func NewErrors() *Errors {
	return &Errors{
		OrderedMap: *orderedmap.New(),
	}
}

// NewInternalError wraps a given error into an InternalError.
func NewInternalError(err error) InternalError {
	return &internalError{error: err}
}

// InternalError returns the actual error that it wraps around.
func (e *internalError) InternalError() error {
	return e.error
}

// Error returns the error string of Errors.
func (es *Errors) Error() string {
	if len(es.Keys()) == 0 {
		return ""
	}
	s := ""
	for i, key := range es.Keys() {
		if i > 0 {
			s += "; "
		}
		tmp, _ := es.Get(key)
		if errs, ok := tmp.(*Errors); ok {
			s += fmt.Sprintf("%v: (%v)", key, errs)
		} else {
			s += fmt.Sprintf("%v: %v", key, tmp.(error).Error())
		}
	}
	return s + "."
}

// Filter removes all nils from Errors and returns back the updated Errors as an error.
// If the length of Errors becomes 0, it will return nil.
func (es *Errors) Filter() error {
	for _, key := range es.Keys() {
		value, _ := es.Get(key)
		if value == nil {
			es.Delete(key)
		}
	}
	if len(es.Keys()) == 0 {
		return nil
	}
	return es
}

func (es *Errors) MarshalJSON() ([]byte, error) {
	errs := orderedmap.New()
	for _, key := range es.Keys() {
		err, _ := es.Get(key)
		if ms, ok := err.(json.Marshaler); ok {
			errs.Set(key, ms)
		} else {
			errs.Set(key, err.(error).Error())
		}
	}
	return json.Marshal(errs)
}
