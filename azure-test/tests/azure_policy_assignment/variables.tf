
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

data "azurerm_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "azure:///subscriptions/${data.azurerm_client_config.current.subscription_id}"
  }
}

resource "azurerm_policy_definition" "named_test_resource" {
  name         = "SecurityCenterBuiltIn"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "SecurityCenterBuiltIn"
  policy_rule  = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE
  parameters   = <<PARAMETERS
    {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West US"
}

resource "azurerm_policy_assignment" "named_test_resource" {
  name                 = "SecurityCenterBuiltIn"
  scope                = azurerm_resource_group.named_test_resource.id
  policy_definition_id = azurerm_policy_definition.named_test_resource.id
  parameters           = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West US" ]
  }
}
PARAMETERS

}

output "resource_aka" {
  value = "azure:///subscriptions/${var.azure_subscription}/resourceGroups/${var.resource_name}/providers/Microsoft.Authorization/policyAssignments/SecurityCenterBuiltIn"
}

output "resource_aka_lower" {
  value = "azure:///subscriptions/${var.azure_subscription}/resourcegroups/${var.resource_name}/providers/microsoft.authorization/policyassignments/securitycenterbuiltin"
}

output "resource_id" {
  value = "/subscriptions/${var.azure_subscription}/resourceGroups/${var.resource_name}/providers/Microsoft.Authorization/policyAssignments/SecurityCenterBuiltIn"
}

output "resource_name" {
  value = "SecurityCenterBuiltIn"
}

output "subscription_id" {
  value = var.azure_subscription
}
