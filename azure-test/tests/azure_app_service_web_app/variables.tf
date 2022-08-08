
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

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_storage_account" "named_test_resource" {
  name                     = var.resource_name
  resource_group_name      = azurerm_resource_group.named_test_resource.name
  location                 = azurerm_resource_group.named_test_resource.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  sku {
    tier = "Free"
    size = "F1"
  }
}

resource "azurerm_app_service" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  app_service_plan_id = azurerm_app_service_plan.named_test_resource.id
  identity {
    type = "SystemAssigned"
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_app_service.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_app_service.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_app_service.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "tenant_id" {
  value = azurerm_app_service.named_test_resource.identity[0].tenant_id
}

output "principal_id" {
  value = azurerm_app_service.named_test_resource.identity[0].principal_id
}
