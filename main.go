package main

import (
	"log"
	"os"
	// "path"
	// "runtime"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/logutils"
	colorable "github.com/mattn/go-colorable"
	"github.com/wata727/tflint/cmd"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/plugin/shared"
	"github.com/wata727/tflint/rules"
)

// A collection of Rules
type Issues struct {
	logger hclog.Logger
}

var additionalRules = []rules.Rule{
	NewAwsSgCantBeNamedJimRule(),
}

func (r *Issues) Process(files []string) []*issue.Issue {
	cli := cmd.NewCLI(colorable.NewColorable(os.Stdout), colorable.NewColorable(os.Stderr))
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(strings.ToUpper(os.Getenv("TFLINT_LOG"))),
		Writer:   os.Stderr,
	}

	log.SetOutput(filter)
	log.SetFlags(log.Ltime | log.Lshortfile)

	cli.SanityCheck(files)
	cli.ProcessRules(additionalRules...)

	return cli.Issues
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Error,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	ruleCollection := &Issues{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"customRulesPlugin": &rulesplugin.RuleCollectionPlugin{Impl: ruleCollection},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
