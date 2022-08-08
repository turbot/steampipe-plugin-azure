
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

resource "azurerm_security_center_auto_provisioning" "named_test_resource" {
  auto_provision = "On"
}

output "resource_aka" {
  value = "azure://${azurerm_security_center_auto_provisioning.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_security_center_auto_provisioning.named_test_resource.id)}"
}

output "resource_id" {
  value = azurerm_security_center_auto_provisioning.named_test_resource.id
}

output "resource_name" {
  value = element(split("/", azurerm_security_center_auto_provisioning.named_test_resource.id), 6)
}

output "subscription_id" {
  value = var.azure_subscription
}
