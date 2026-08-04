package setting

import "github.com/QuantumNous/new-api/types"

var systemPromptExemptGroups = types.NewRWMap[string, bool]()

func IsGroupExemptFromSystemPrompt(groupName string) bool {
	exempt, _ := systemPromptExemptGroups.Get(groupName)
	return exempt
}

func SystemPromptExemptGroups2JSONString() string {
	return systemPromptExemptGroups.MarshalJSONString()
}

func UpdateSystemPromptExemptGroupsByJSONString(jsonStr string) error {
	return types.LoadFromJsonString(systemPromptExemptGroups, jsonStr)
}
