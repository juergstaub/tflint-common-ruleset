package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformResourceNameRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "issue found",
			Content: `
data "namep_azure_name" "rgp-np" {
	type = "azurerm_resource_group"
	name = "np"
}

variable "supername" {
	type = string
}

resource "azurerm_resource_group" "rg1" {
  name = "rg1"
  location = var.location
}

resource "azurerm_resource_group" "rg2" {
	name = var.supername
	location = var.location
}

resource "azurerm_resource_group" "rg3" {
	name = data.namep_azure_name.rgp-np.result
	location = var.location
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceNameRule(),
					Message: "Direct assignment of resources names is not allowed",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 10},
						End:      hcl.Pos{Line: 12, Column: 15},
					},
				},
				{
					Rule:    NewTerraformResourceNameRule(),
					Message: "Resource names must be assigned with namep",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 17, Column: 9},
						End:      hcl.Pos{Line: 17, Column: 22},
					},
				},
			},
		},
	}

	rule := NewTerraformResourceNameRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
