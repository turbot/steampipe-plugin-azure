
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
  version         = "=1.36.0"
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azurerm_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "azure:///subscriptions/${data.azurerm_client_config.current.subscription_id}"
  }
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "named_test_resource" {
  name  = var.resource_name
  scope = data.azurerm_subscription.primary.id

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}

resource "azurerm_role_assignment" "named_test_resource" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.named_test_resource.id
  principal_id       = data.azurerm_client_config.current.object_id
}

output "resource_aka" {
  value = "azure://${azurerm_role_assignment.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_role_assignment.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_role_assignment.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "role_definition_id" {
  value = azurerm_role_definition.named_test_resource.id
}

output "principal_id" {
  value = data.azurerm_client_config.current.object_id
}

output "principal_type" {
  value = azurerm_role_assignment.named_test_resource.principal_type
}
