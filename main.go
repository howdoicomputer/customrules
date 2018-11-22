package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/wata727/tflint/rules"
)

// Rules that are returned to the main cli process
var defaultRules []rules.Rule
var deepCheckRUles []rules.Rule

// A collection of Rules
type RuleCollection struct {
	rules []rules.Rule
}

func (r *RuleCollection) Present() []rules.Rule {
	return r.rules
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	rule_collection := &RuleCollection{
		rules: []rules.Rule{
			NewAwsSgCantBeNamedJimRule(),
		},
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"customrules": &rules.RulePlugin{Impl: rule_collection},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
