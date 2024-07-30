package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	errEmptyBody       = errors.New("body must not be empty")
	errBadlyFormedJSON = errors.New("body contains badly-formed JSON")
	errMultipleValues  = errors.New("body must only contain a single JSON value")
)

type so struct {
	offset int64
}

func (s *so) Error() string {
	return fmt.Sprintf("body contains badly-formed JSON (at character %d)", s.offset)
}

type ute struct {
	field  string
	offset int64
}

func (u *ute) Error() string {
	if u.field != "" {
		return fmt.Sprintf("body contains incorrect JSON type for field %q", u.field)
	}

	return fmt.Sprintf("body contains incorrect JSON type (at character %d)", u.offset)
}

type mb struct {
	limit int64
}

func (m *mb) Error() string {
	return fmt.Sprintf("body must not be larger than %d bytes", m.limit)
}

type u struct {
	field string
}

func (u *u) Error() string {
	return fmt.Sprintf("body contains unknown key %s", u.field)
}

func Read(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return &so{offset: syntaxError.Offset}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errBadlyFormedJSON
		case errors.As(err, &unmarshalTypeError):
			return &ute{field: unmarshalTypeError.Field, offset: unmarshalTypeError.Offset}
		case errors.Is(err, io.EOF):
			return errEmptyBody
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return &u{field: fieldName}
		case errors.As(err, &maxBytesError):
			return &mb{limit: maxBytesError.Limit}
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errMultipleValues
	}
	return nil
}
