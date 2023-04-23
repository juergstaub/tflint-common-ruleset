package rules

import (
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformResourceNameRule checks whether ...
type TerraformResourceNameRule struct {
	tflint.DefaultRule
	ResourceNameExpression regexp.Regexp
}

// TerraformResourceNameRule returns a new rule
func NewTerraformResourceNameRule() *TerraformResourceNameRule {
	return &TerraformResourceNameRule{
		ResourceNameExpression: *regexp.MustCompile(`data\.namep_azure_name\..*\.result`)}
}

// Name returns the rule name
func (r *TerraformResourceNameRule) Name() string {
	return ""
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformResourceNameRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformResourceNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *TerraformResourceNameRule) Link() string {
	return "azurerm_resource_group"
}

// Check checks whether ...
func (r *TerraformResourceNameRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("azurerm_resource_group", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "name"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["name"]
		if !exists {
			continue
		}

		if len(attribute.Expr.Variables()) == 0 {
			runner.EmitIssue(
				r,
				`Direct assignment of resources names is not allowed`,
				attribute.Expr.Range(),
			)
			continue
		}

		actual := []string{}
		for _, a := range attribute.Expr.Variables()[0] {
			at, ok := a.(hcl.TraverseAttr)
			if ok {
				actual = append(actual, at.Name)
			} else {
				at, ok := a.(hcl.TraverseRoot)
				if ok {
					actual = append(actual, at.Name)
				}
			}
		}

		if !r.ResourceNameExpression.MatchString(strings.Join(actual, ".")) {
			runner.EmitIssue(
				r,
				`Resource names must be assigned with namep`,
				attribute.Expr.Range(),
			)
		}

	}
	return nil

}
