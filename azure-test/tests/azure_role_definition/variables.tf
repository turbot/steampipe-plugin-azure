
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
  default     = "3510ae4d-530b-497d-8f30-53b9616fc6c1"
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

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "named_test_resource" {
  name        = var.resource_name
  scope       = data.azurerm_subscription.primary.id
  description = "This is a custom role created via Terraform"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = ["/subscriptions/${var.azure_subscription}"]
}

output "resource_aka" {
  value = "azure://${azurerm_role_definition.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_role_definition.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_role_definition.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
