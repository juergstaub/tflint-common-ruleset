package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformResourceNameRule checks whether ...
type TerraformResourceNameRule struct {
	tflint.DefaultRule
}

// TerraformResourceNameRule returns a new rule
func NewTerraformResourceNameRule() *TerraformResourceNameRule {
	return &TerraformResourceNameRule{}
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
				fmt.Sprintf(`Direct assignment of resources names is not allowed`),
				attribute.Expr.Range(),
			)
			continue
		}

		vars := attribute.Expr.Variables()[0]
		expected := [...]string{"data", "namep_azure_name"}
		actual := []string{}
		for _, a := range vars {
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

		success := true
		if len(actual) >= len(expected) {
			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					success = false
					break
				}
			}
		} else {
			success = false
		}

		if !success {
			runner.EmitIssue(
				r,
				fmt.Sprintf(`Resource names must be assigned with namep`),
				attribute.Expr.Range(),
			)
		}

	}
	return nil

}
