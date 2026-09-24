package controller

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAdminBindUserInviter(t *testing.T) {
	for _, tc := range []struct {
		name       string
		actorRole  int
		targetRole int
		bound      bool
		body       string
		wantOK     bool
		wantError  string
	}{
		{name: "admin can bind common user", actorRole: common.RoleAdminUser, targetRole: common.RoleCommonUser, body: `{"aff_code":"code"}`, wantOK: true},
		{name: "root can bind admin", actorRole: common.RoleRootUser, targetRole: common.RoleAdminUser, body: `{"aff_code":"code"}`, wantOK: true},
		{name: "admin cannot bind peer", actorRole: common.RoleAdminUser, targetRole: common.RoleAdminUser, body: `{"aff_code":"code"}`, wantError: i18n.MsgUserNoPermissionHigherLevel},
		{name: "admin cannot bind root", actorRole: common.RoleAdminUser, targetRole: common.RoleRootUser, body: `{"aff_code":"code"}`, wantError: i18n.MsgUserNoPermissionHigherLevel},
		{name: "existing inviter is immutable", actorRole: common.RoleRootUser, targetRole: common.RoleCommonUser, bound: true, body: `{"aff_code":"code"}`, wantError: i18n.MsgUserInviterAlreadyBound},
		{name: "invalid payload", actorRole: common.RoleRootUser, targetRole: common.RoleCommonUser, body: `{"aff_code":1}`, wantError: i18n.MsgInvalidParams},
		{name: "missing invitation code", actorRole: common.RoleRootUser, targetRole: common.RoleCommonUser, body: `{}`, wantError: i18n.MsgUserInvalidInviteCode},
	} {
		t.Run(tc.name, func(t *testing.T) {
			originalDB, originalLogDB := model.DB, model.LOG_DB
			originalRedis, originalTranslate := common.RedisEnabled, common.TranslateMessage
			originalMainType, originalLogType := common.MainDatabaseType(), common.LogDatabaseType()
			t.Cleanup(func() {
				model.DB, model.LOG_DB = originalDB, originalLogDB
				common.RedisEnabled, common.TranslateMessage = originalRedis, originalTranslate
				common.SetDatabaseTypes(originalMainType, originalLogType)
			})
			db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "inviter.db")), &gorm.Config{})
			require.NoError(t, err)
			sqlDB, err := db.DB()
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, sqlDB.Close()) })
			model.DB, model.LOG_DB = db, db
			common.RedisEnabled = false
			common.SetDatabaseTypes(common.DatabaseTypeSQLite, common.DatabaseTypeSQLite)
			require.NoError(t, db.AutoMigrate(&model.User{}, &model.Log{}))
			require.NoError(t, i18n.Init())
			common.TranslateMessage = i18n.T
			actor := model.User{Username: "admin", AffCode: "admin", Role: tc.actorRole}
			inviter := model.User{Username: "inviter", AffCode: "code", Status: common.UserStatusEnabled}
			target := model.User{Username: "invitee", AffCode: "own", Role: tc.targetRole}
			require.NoError(t, db.Create(&actor).Error)
			require.NoError(t, db.Create(&inviter).Error)
			if tc.bound {
				target.InviterId = actor.Id
			}
			require.NoError(t, db.Create(&target).Error)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Params = gin.Params{{Key: "id", Value: strconv.Itoa(target.Id)}}
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/user/"+strconv.Itoa(target.Id)+"/inviter", strings.NewReader(tc.body))
			ctx.Request.Header.Set("Accept-Language", "en")
			ctx.Set("id", actor.Id)
			ctx.Set("role", actor.Role)
			ctx.Set("username", actor.Username)

			AdminBindUserInviter(ctx)

			var response struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
				Data    struct {
					InviterId int `json:"inviter_id"`
				} `json:"data"`
			}
			require.Equal(t, http.StatusOK, recorder.Code)
			require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, tc.wantOK, response.Success)
			var stored model.User
			require.NoError(t, db.First(&stored, target.Id).Error)
			var logs []model.Log
			require.NoError(t, db.Find(&logs).Error)
			if !tc.wantOK {
				assert.Equal(t, i18n.T(ctx, tc.wantError), response.Message)
				assert.Equal(t, target.InviterId, stored.InviterId)
				assert.Empty(t, logs)
				return
			}
			assert.Equal(t, inviter.Id, response.Data.InviterId)
			assert.Equal(t, inviter.Id, stored.InviterId)
			require.Len(t, logs, 1)
			assert.Equal(t, actor.Id, logs[0].UserId)
			var audit struct {
				Op struct {
					Action string `json:"action"`
					Params struct {
						TargetUserId int `json:"target_user_id"`
						InviterId    int `json:"inviter_id"`
					} `json:"params"`
				} `json:"op"`
			}
			require.NoError(t, common.UnmarshalJsonStr(logs[0].Other, &audit))
			assert.Equal(t, "user.inviter_bind", audit.Op.Action)
			assert.Equal(t, target.Id, audit.Op.Params.TargetUserId)
			assert.Equal(t, inviter.Id, audit.Op.Params.InviterId)
		})
	}
}
