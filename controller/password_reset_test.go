package controller

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSendPasswordResetEmailReportsDeliveryResult(t *testing.T) {
	cases := []struct {
		name                string
		users               int
		deleted             bool
		databaseUnavailable bool
		rejectSMTP          bool
		wantSuccess         bool
		wantMessage         string
		wantConnection      int32
	}{
		{
			name:        "email has no account",
			wantMessage: "No active account is linked to this email address",
		},
		{
			name:        "email belongs to a deleted account",
			users:       1,
			deleted:     true,
			wantMessage: "No active account is linked to this email address",
		},
		{
			name:        "email belongs to multiple accounts",
			users:       2,
			wantMessage: "This email address is linked to multiple accounts. Please contact the administrator.",
		},
		{
			name:                "database lookup fails",
			databaseUnavailable: true,
			wantMessage:         "sql: database is closed",
		},
		{
			name:           "SMTP rejects recipient",
			users:          1,
			rejectSMTP:     true,
			wantMessage:    "550 mailbox unavailable",
			wantConnection: 1,
		},
		{
			name:           "SMTP accepts reset email",
			users:          1,
			wantSuccess:    true,
			wantConnection: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupPasswordResetTest(t)
			for index := 0; index < tc.users; index++ {
				user := model.User{
					Username: fmt.Sprintf("reset-user-%d", index),
					Password: "existing-password-hash",
					Email:    "RESET@example.com",
					Status:   common.UserStatusEnabled,
					AffCode:  fmt.Sprintf("reset-%d", index),
				}
				require.NoError(t, db.Create(&user).Error)
				if tc.deleted {
					require.NoError(t, db.Delete(&user).Error)
				}
			}
			if tc.databaseUnavailable {
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			}
			server := newPasswordResetSMTPServer(t, tc.rejectSMTP)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			query := url.Values{"email": {" RESET@example.com "}}
			ctx.Request = httptest.NewRequest(http.MethodGet, "/api/reset_password?"+query.Encode(), nil)
			ctx.Request.Header.Set("Accept-Language", "en")

			SendPasswordResetEmail(ctx)

			var response struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}
			require.Equal(t, http.StatusOK, recorder.Code)
			require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, tc.wantSuccess, response.Success)
			assert.Equal(t, tc.wantMessage, response.Message)
			assert.Equal(t, tc.wantConnection, server.connections.Load())
			if tc.wantSuccess {
				require.Len(t, server.messages, 1)
				message := <-server.messages
				assert.Contains(t, message, "To: reset@example.com")
				matches := regexp.MustCompile(`href='([^']+)'`).FindStringSubmatch(message)
				require.Len(t, matches, 2, "reset email must contain a link")
				link, err := url.Parse(matches[1])
				require.NoError(t, err)
				assert.Equal(t, "https", link.Scheme)
				assert.Equal(t, "reset.example.com", link.Host)
				assert.Equal(t, "/user/reset", link.Path)
				assert.Equal(t, "reset@example.com", link.Query().Get("email"))
				token := link.Query().Get("token")
				require.NotEmpty(t, token)
				assert.True(t, common.VerifyCodeWithKey("reset@example.com", token, common.PasswordResetPurpose))
			} else {
				assert.Empty(t, server.messages)
			}
			if tc.databaseUnavailable {
				return
			}
			var users []model.User
			require.NoError(t, db.Unscoped().Find(&users).Error)
			for _, user := range users {
				assert.Equal(t, "existing-password-hash", user.Password, "requesting a reset must not change the password")
			}
		})
	}
}

func setupPasswordResetTest(t *testing.T) *gorm.DB {
	t.Helper()
	originalDB := model.DB
	originalMode := gin.Mode()
	originalTranslate := common.TranslateMessage
	originalServerAddress := system_setting.ServerAddress
	originalSystemName := common.SystemName
	originalSMTPServer, originalSMTPPort := common.SMTPServer, common.SMTPPort
	originalSMTPSSL, originalSMTPStartTLS := common.SMTPSSLEnabled, common.SMTPStartTLSEnabled
	originalSMTPInsecure, originalSMTPForceLogin := common.SMTPInsecureSkipVerify, common.SMTPForceAuthLogin
	originalSMTPAccount, originalSMTPFrom, originalSMTPToken := common.SMTPAccount, common.SMTPFrom, common.SMTPToken
	t.Cleanup(func() {
		model.DB = originalDB
		gin.SetMode(originalMode)
		common.TranslateMessage = originalTranslate
		system_setting.ServerAddress = originalServerAddress
		common.SystemName = originalSystemName
		common.SMTPServer, common.SMTPPort = originalSMTPServer, originalSMTPPort
		common.SMTPSSLEnabled, common.SMTPStartTLSEnabled = originalSMTPSSL, originalSMTPStartTLS
		common.SMTPInsecureSkipVerify, common.SMTPForceAuthLogin = originalSMTPInsecure, originalSMTPForceLogin
		common.SMTPAccount, common.SMTPFrom, common.SMTPToken = originalSMTPAccount, originalSMTPFrom, originalSMTPToken
		common.DeleteKey("reset@example.com", common.PasswordResetPurpose)
	})

	require.NoError(t, i18n.Init())
	common.TranslateMessage = i18n.T
	gin.SetMode(gin.TestMode)
	system_setting.ServerAddress = "https://reset.example.com"
	common.SystemName = "Password Reset Test"
	common.SMTPSSLEnabled, common.SMTPStartTLSEnabled = false, false
	common.SMTPInsecureSkipVerify, common.SMTPForceAuthLogin = false, false
	common.SMTPAccount, common.SMTPFrom, common.SMTPToken = "", "sender@example.com", ""
	common.DeleteKey("reset@example.com", common.PasswordResetPurpose)

	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "password-reset.db")), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, sqlDB.Close()) })
	require.NoError(t, db.AutoMigrate(&model.User{}))
	model.DB = db
	return db
}

type passwordResetSMTPServer struct {
	listener    net.Listener
	connections atomic.Int32
	messages    chan string
	reject      bool
}

func newPasswordResetSMTPServer(t *testing.T, reject bool) *passwordResetSMTPServer {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := &passwordResetSMTPServer{listener: listener, messages: make(chan string, 1), reject: reject}
	done := make(chan error, 1)
	go func() { done <- server.serve() }()
	t.Cleanup(func() {
		require.NoError(t, listener.Close())
		if err := <-done; !errors.Is(err, net.ErrClosed) {
			require.NoError(t, err)
		}
	})
	host, port, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)
	common.SMTPServer = host
	common.SMTPPort, err = strconv.Atoi(port)
	require.NoError(t, err)
	return server
}

func (server *passwordResetSMTPServer) serve() error {
	conn, err := server.listener.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	server.connections.Add(1)
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}
	reader := textproto.NewReader(bufio.NewReader(conn))
	writer := textproto.NewWriter(bufio.NewWriter(conn))
	if err := writer.PrintfLine("220 reset.test ESMTP"); err != nil {
		return err
	}
	for {
		command, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		var response string
		switch {
		case strings.HasPrefix(command, "EHLO "), strings.HasPrefix(command, "HELO "):
			response = "250 reset.test"
		case strings.HasPrefix(command, "MAIL FROM:"):
			response = "250 sender accepted"
		case strings.HasPrefix(command, "RCPT TO:"):
			response = "250 recipient accepted"
			if server.reject {
				response = "550 mailbox unavailable"
			}
		case command == "DATA":
			if err := writer.PrintfLine("354 send message"); err != nil {
				return err
			}
			message, err := reader.ReadDotBytes()
			if err != nil {
				return err
			}
			server.messages <- string(message)
			response = "250 message accepted"
		case command == "QUIT":
			return writer.PrintfLine("221 goodbye")
		default:
			return fmt.Errorf("unexpected SMTP command: %s", command)
		}
		if err := writer.PrintfLine("%s", response); err != nil {
			return err
		}
	}
}
