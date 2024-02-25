package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/micahasowata/jason"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string
		body string
		code int
	}{
		{
			name: "valid",
			body: `{"name": "tbd", "email":"me@tbd.com", "password":"://me@tbd.com://"}`,
			code: http.StatusOK,
		},
		{
			name: "bad request body",
			body: `{"name": "}`,
			code: http.StatusBadRequest,
		},
		{
			name: "invalid data",
			body: `{"name": "tbd","email": "metbd.com"}`,
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "duplicate data",
			body: `{"name": "tbd", "email":"me@tbd.com", "password":"://me@tbd.com://"}`,
			code: http.StatusInternalServerError,
		},
	}

	_, ts := setUpTest()
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := ts.Client().Post(ts.URL+"/v1/users/create", jason.ContentTypeJSON, strings.NewReader(tt.body))
			require.Nil(t, err)

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
