package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGetFirstGroupSupportsLegacyMultiGroupValues(t *testing.T) {
	tests := []struct {
		name  string
		group string
		want  string
	}{
		{name: "empty inherits user group", group: "", want: ""},
		{name: "single group", group: "default", want: "default"},
		{name: "single group trims whitespace", group: " vip ", want: "vip"},
		{name: "legacy multi group uses first", group: `["default","vip"]`, want: "default"},
		{name: "legacy first group trims whitespace", group: `[" vip ","default"]`, want: "vip"},
		{name: "legacy empty list inherits user group", group: `[]`, want: ""},
		{name: "malformed JSON preserves original value", group: `["default"`, want: `["default"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &Token{Group: tt.group}
			assert.Equal(t, tt.want, token.GetFirstGroup())
		})
	}
}
