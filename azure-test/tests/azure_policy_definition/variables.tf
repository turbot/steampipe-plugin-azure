
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
  default     = "cdffd708-7da0-4cea-abeb-0a4c334d7f64"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  version = "=2.43.0"
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azurerm_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "azure:///subscriptions/${data.azurerm_client_config.current.subscription_id}"
  }
}

resource "azurerm_policy_definition" "named_test_resource" {
  name         = var.resource_name
  policy_type  = "Custom"
  mode         = "All"
  display_name = var.resource_name
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

output "resource_aka" {
  value = "azure://${azurerm_policy_definition.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_policy_definition.named_test_resource.id)}"
}

output "resource_id" {
  value = azurerm_policy_definition.named_test_resource.id
}

output "resource_name" {
  value = azurerm_policy_definition.named_test_resource.name
}

output "subscription_id" {
  value = var.azure_subscription
}
