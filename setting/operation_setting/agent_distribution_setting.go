package operation_setting

import "github.com/QuantumNous/new-api/setting/config"

type AgentDistributionSetting struct {
	DefaultFirstTopupRate  float64 `json:"default_first_topup_rate"`
	DefaultRepeatTopupRate float64 `json:"default_repeat_topup_rate"`
}

var agentDistributionSetting = AgentDistributionSetting{}

func init() {
	config.GlobalConfig.Register("agent_distribution_setting", &agentDistributionSetting)
}

func GetAgentDistributionSetting() *AgentDistributionSetting {
	return &agentDistributionSetting
}
