package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAgentCommissionAllowsRepeatedPaymentsButIsIdempotentPerTrade(t *testing.T) {
	truncateTables(t)

	agent := &User{
		Username:             "agent-repeat",
		AffCode:              "agent-repeat-code",
		Password:             "password",
		Status:               common.UserStatusEnabled,
		AgentEnabled:         true,
		AgentUseDefaultRates: false,
		AgentFirstTopupRate:  10,
		AgentRepeatTopupRate: 5,
	}
	require.NoError(t, DB.Create(agent).Error)
	invitee := &User{
		Username:     "agent-customer",
		AffCode:      "agent-customer-code",
		Password:     "password",
		Status:       common.UserStatusEnabled,
		InviterId:    agent.Id,
		ReferralMode: ReferralModeAgentDistribution,
	}
	require.NoError(t, DB.Create(invitee).Error)

	var first, repeat, duplicate *AgentCommissionOutcome
	require.NoError(t, DB.Transaction(func(tx *gorm.DB) error {
		var err error
		first, err = HandleAgentCommissionForRedemptionTx(tx, invitee.Id, "redeem-first", 10000)
		return err
	}))
	require.NoError(t, DB.Transaction(func(tx *gorm.DB) error {
		var err error
		repeat, err = HandleAgentCommissionForRedemptionTx(tx, invitee.Id, "redeem-repeat", 10000)
		return err
	}))
	require.NoError(t, DB.Transaction(func(tx *gorm.DB) error {
		var err error
		duplicate, err = HandleAgentCommissionForRedemptionTx(tx, invitee.Id, "redeem-repeat", 10000)
		return err
	}))

	require.NotNil(t, first)
	assert.True(t, first.IsFirstTopupOrder)
	assert.Equal(t, int64(1000), first.CommissionAmount)
	require.NotNil(t, repeat)
	assert.False(t, repeat.IsFirstTopupOrder)
	assert.Equal(t, int64(500), repeat.CommissionAmount)
	assert.Nil(t, duplicate)

	var recordCount int64
	require.NoError(t, DB.Model(&AgentCommissionRecord{}).Where("agent_id = ?", agent.Id).Count(&recordCount).Error)
	assert.Equal(t, int64(2), recordCount)
	require.NoError(t, DB.First(agent, agent.Id).Error)
	assert.Equal(t, int64(1500), agent.AgentCommissionBalance)
	assert.Equal(t, int64(1500), agent.AgentCommissionTotal)
}

func TestAgentDistributionSchemaKeepsOldSQLNames(t *testing.T) {
	require.NoError(t, DB.AutoMigrate(&User{}))
	require.NoError(t, DB.AutoMigrate(&AgentCommissionRecord{}))
	migrator := DB.Migrator()
	assert.True(t, migrator.HasTable(&AgentCommissionRecord{}))
	assert.True(t, migrator.HasIndex(&AgentCommissionRecord{}, "idx_agent_commission_source"))
	for _, column := range []string{
		"referral_mode",
		"agent_enabled",
		"agent_use_default_rates",
		"agent_first_topup_rate",
		"agent_repeat_topup_rate",
		"agent_commission_balance",
		"agent_commission_total",
		"agent_commission_withdrawn",
		"first_payment_at",
		"first_payment_type",
		"first_payment_trade_no",
	} {
		assert.True(t, migrator.HasColumn(&User{}, column), column)
	}
}

func TestDetectsLegacySQLiteAgentDecimalSchema(t *testing.T) {
	require.NoError(t, DB.Exec("CREATE TABLE legacy_agent_rates (id integer, agent_first_topup_rate decimal(10,4))").Error)
	t.Cleanup(func() {
		require.NoError(t, DB.Exec("DROP TABLE legacy_agent_rates").Error)
	})

	skip, err := shouldSkipLegacyAgentSQLiteAutoMigrate("legacy_agent_rates")
	require.NoError(t, err)
	assert.True(t, skip)
}
