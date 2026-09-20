package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenLegacyMultiGroupCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		token      Token
		wantGroup  string
		wantBackup string
	}{
		{name: "empty inherits user group", token: Token{}, wantGroup: "", wantBackup: ""},
		{name: "single group", token: Token{Group: "default"}, wantGroup: "default", wantBackup: ""},
		{name: "single group trims whitespace", token: Token{Group: " vip "}, wantGroup: "vip", wantBackup: ""},
		{name: "legacy groups use first and second", token: Token{Group: `["default","vip"]`}, wantGroup: "default", wantBackup: "vip"},
		{name: "legacy groups trim whitespace", token: Token{Group: `[" vip "," backup "]`}, wantGroup: "vip", wantBackup: "backup"},
		{name: "dedicated backup takes precedence", token: Token{Group: `["default","legacy"]`, BackupGroup: "current"}, wantGroup: "default", wantBackup: "current"},
		{name: "legacy empty list inherits user group", token: Token{Group: `[]`}, wantGroup: "", wantBackup: ""},
		{name: "malformed JSON preserves primary value", token: Token{Group: `["default"`}, wantGroup: `["default"`, wantBackup: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantGroup, tt.token.GetFirstGroup())
			assert.Equal(t, tt.wantBackup, tt.token.GetBackupGroup())
		})
	}
}
