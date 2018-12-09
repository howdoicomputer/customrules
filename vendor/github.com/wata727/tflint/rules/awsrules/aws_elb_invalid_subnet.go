package awsrules

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/tflint"
)

// AwsELBInvalidSubnetRule checks whether subnets actually exists
type AwsELBInvalidSubnetRule struct {
	resourceType  string
	attributeName string
	subnets       map[string]bool
	dataPrepared  bool
}

// NewAwsELBInvalidSubnetRule returns new rule with default attributes
func NewAwsELBInvalidSubnetRule() *AwsELBInvalidSubnetRule {
	return &AwsELBInvalidSubnetRule{
		resourceType:  "aws_elb",
		attributeName: "subnets",
		subnets:       map[string]bool{},
		dataPrepared:  false,
	}
}

// Name returns the rule name
func (r *AwsELBInvalidSubnetRule) Name() string {
	return "aws_elb_invalid_subnet"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsELBInvalidSubnetRule) Enabled() bool {
	return true
}

// Type returns the rule severity
func (r *AwsELBInvalidSubnetRule) Type() string {
	return issue.ERROR
}

// Link returns the rule reference link
func (r *AwsELBInvalidSubnetRule) Link() string {
	return ""
}

// Check checks whether `subnets` are included in the list retrieved by `DescribeSubnets`
func (r *AwsELBInvalidSubnetRule) Check(runner *tflint.Runner) error {
	log.Printf("[INFO] Check `%s` rule for `%s` runner", r.Name(), runner.TFConfigPath())

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		if !r.dataPrepared {
			log.Print("[DEBUG] Fetch subnets")
			resp, err := runner.AwsClient.EC2.DescribeSubnets(&ec2.DescribeSubnetsInput{})
			if err != nil {
				err := &tflint.Error{
					Code:    tflint.ExternalAPIError,
					Level:   tflint.ErrorLevel,
					Message: "An error occurred while describing subnets",
					Cause:   err,
				}
				log.Printf("[ERROR] %s", err)
				return err
			}
			for _, subnet := range resp.Subnets {
				r.subnets[*subnet.SubnetId] = true
			}
			r.dataPrepared = true
		}

		return runner.EachStringSliceExprs(attribute.Expr, func(subnet string, expr hcl.Expression) {
			if !r.subnets[subnet] {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is invalid subnet ID.", subnet),
					expr.Range(),
				)
			}
		})
	})
}
