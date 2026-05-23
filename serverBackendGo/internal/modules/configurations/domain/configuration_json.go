package domain

import (
	"encoding/json"
	"strings"
)

// configurationScalarKeys are persisted as SQL columns or nested arrays, not in settingsjson alone.
var configurationScalarKeys = map[string]struct{}{
	"id": {}, "name": {}, "description": {}, "type": {}, "deviceCount": {},
	"password": {}, "backgroundColor": {}, "textColor": {}, "backgroundImageUrl": {},
	"qrCodeKey": {}, "baseUrl": {}, "mainAppId": {}, "contentAppId": {},
	"defaultFilePath": {}, "permissive": {},
	"applications": {}, "files": {}, "applicationSettings": {}, "policyLocks": {},
}

// SetPolicyFromJSON stores MDM policy keys from settingsjson for merge on API responses.
func (c *Configuration) SetPolicyFromJSON(raw []byte) {
	if len(raw) == 0 {
		return
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if _, skip := configurationScalarKeys[k]; skip {
			continue
		}
		var val any
		if err := json.Unmarshal(v, &val); err == nil {
			out[k] = val
		}
	}
	if raw, ok := m[policyLocksKey]; ok {
		var locks map[string]bool
		if err := json.Unmarshal(raw, &locks); err == nil {
			c.PolicyLocks = NormalizePolicyLocks(locks)
		}
		delete(m, policyLocksKey)
	}
	if len(out) > 0 {
		c.Policy = out
	}
}

// BuildSettingsJSON returns settingsjson bytes: policy map merged over scalar fields from cfg.
func (c Configuration) BuildSettingsJSON() ([]byte, error) {
	m := make(map[string]any)
	if locks := NormalizePolicyLocks(c.PolicyLocks); locks != nil {
		m[policyLocksKey] = locks
	}
	if c.Policy != nil {
		for k, v := range c.Policy {
			m[k] = v
		}
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var scalars map[string]any
	if err := json.Unmarshal(b, &scalars); err != nil {
		return json.Marshal(m)
	}
	for k, v := range scalars {
		if _, skip := configurationScalarKeys[k]; skip {
			continue
		}
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && strings.TrimSpace(s) == "" {
			continue
		}
		m[k] = v
	}
	return json.Marshal(m)
}

// ConfigurationResponseMap flattens cfg and Policy for Headwind-compatible JSON.
func ConfigurationResponseMap(cfg *Configuration) map[string]any {
	if cfg == nil {
		return map[string]any{}
	}
	b, _ := json.Marshal(cfg)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if cfg.Policy != nil {
		for k, v := range cfg.Policy {
			m[k] = v
		}
	}
	if cfg.PolicyLocks != nil {
		m["policyLocks"] = cfg.PolicyLocks
	}
	delete(m, "policy")
	return m
}

// ParseConfigurationBody decodes PUT/POST body preserving MDM keys in Policy.
func ParseConfigurationBody(raw []byte) (Configuration, error) {
	var cfg Configuration
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return cfg, nil
	}
	policy := make(map[string]any)
	for k, v := range m {
		if _, skip := configurationScalarKeys[k]; skip {
			continue
		}
		var val any
		if err := json.Unmarshal(v, &val); err == nil {
			policy[k] = val
		}
	}
	if raw, ok := m[policyLocksKey]; ok {
		var locks map[string]bool
		if err := json.Unmarshal(raw, &locks); err == nil {
			cfg.PolicyLocks = NormalizePolicyLocks(locks)
		}
		delete(policy, policyLocksKey)
	}
	if len(policy) > 0 {
		cfg.Policy = policy
	}
	return cfg, nil
}
