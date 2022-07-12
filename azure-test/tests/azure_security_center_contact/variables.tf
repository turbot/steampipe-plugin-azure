
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
  default    = "d46d7416-f95f-4771-bbb5-529d4c76659c"
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

resource "azurerm_security_center_contact" "named_test_resource" {
  email               = "contact@example.com"
  phone               = "+1-555-555-5555"
  alert_notifications = true
  alerts_to_admins    = true
}

output "resource_aka" {
  value = "azure://${azurerm_security_center_contact.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_security_center_contact.named_test_resource.id)}"
}

output "resource_id" {
  value = azurerm_security_center_contact.named_test_resource.id
}

output "resource_name" {
  value = element(split("/", azurerm_security_center_contact.named_test_resource.id), 6)
}

output "subscription_id" {
  value = var.azure_subscription
}
