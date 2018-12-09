package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/tflint"
)

type AwsSgCantBeNamedJimRule struct {
	ResourceType  string
	AttributeName string
}

func NewAwsSgCantBeNamedJimRule() *AwsSgCantBeNamedJimRule {
	rule := &AwsSgCantBeNamedJimRule{
		ResourceType:  "aws_security_group",
		AttributeName: "name",
	}

	return rule
}

func (r *AwsSgCantBeNamedJimRule) Name() string {
	return "aws_security_group_cant_be_named_jim"
}

func (r *AwsSgCantBeNamedJimRule) Enabled() bool {
	return true
}

func (r *AwsSgCantBeNamedJimRule) Link() string {
	return "foobar.com"
}

func (r *AwsSgCantBeNamedJimRule) Type() string {
	return issue.ERROR
}

func (r *AwsSgCantBeNamedJimRule) Check(runner *tflint.Runner) error {
	log.Printf("[INFO] Check `%s` rule for `%s` runner", r.Name(), runner.TFConfigPath())

	return runner.WalkResourceAttributes(r.ResourceType, r.AttributeName, func(attribute *hcl.Attribute) error {
		var name string

		err := runner.EvaluateExpr(attribute.Expr, &name)

		return runner.EnsureNoError(err, func() error {
			if name == "jim" {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is named jim.", r.ResourceType),
					attribute.Expr.Range(),
				)
			}

			return nil
		})
	})
}
