variable "resource_name" {
  type        = string
  default     = "steampipe-test"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "d46d7416-f95f-4771-bbb5-529d4c766123"
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

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_automation_account" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku_name            = "Basic"

  tags = {
    name = var.resource_name
  }
}

resource "azurerm_automation_variable_string" "named_test_resource" {
  name                    = var.resource_name
  resource_group_name     = azurerm_resource_group.named_test_resource.name
  automation_account_name = azurerm_automation_account.named_test_resource.name
  value                   = "Hello, Terraform Basic Test."
}

output "resource_aka" {
  value = "azure://${azurerm_automation_variable_string.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_automation_variable_string.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_automation_variable_string.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
