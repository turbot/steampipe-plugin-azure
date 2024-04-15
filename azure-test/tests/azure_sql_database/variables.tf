
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

variable "azure_location" {
  type        = string
  default     = "East US"
  description = "Azure location where the resource will be created."
}

provider "azurerm" {
  # Cannot be passed as a variable
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

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = var.azure_location
}

resource "azurerm_sql_server" "named_test_resource" {
  name                         = var.resource_name
  resource_group_name          = azurerm_resource_group.named_test_resource.name
  location                     = azurerm_resource_group.named_test_resource.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  server_name         = azurerm_sql_server.named_test_resource.name
  tags = {
    foo = "bar"
  }
}

output "resource_aka" {
  value = "azure://${azurerm_sql_database.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_sql_database.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_sql_database.named_test_resource.id
}

output "resource_location" {
  value = var.azure_location
}

output "resource_region" {
  value = lower(var.azure_location)
}

output "resource_creation_date" {
  value = azurerm_sql_database.named_test_resource.creation_date
}

output "resource_default_secondary_location" {
  value = azurerm_sql_database.named_test_resource.default_secondary_location
}

output "subscription_id" {
  value = var.azure_subscription
}
