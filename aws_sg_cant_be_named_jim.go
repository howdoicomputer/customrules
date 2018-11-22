package main

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/tflint"
	"log"
)

type AwsSgCantBeNamedJimRule struct {
	resourceType  string
	attributeName string
}

func NewAwsSgCantBeNamedJimRule() *AwsSgCantBeNamedJimRule {
	rule := &AwsSgCantBeNamedJimRule{
		resourceType:  "aws_security_group",
		attributeName: "name",
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

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var name string

		err := runner.EvaluateExpr(attribute.Expr, &name)

		return runner.EnsureNoError(err, func() error {
			if name == "jim" {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is named jim.", r.resourceType),
					attribute.Expr.Range(),
				)
			}

			return nil
		})
	})
}
