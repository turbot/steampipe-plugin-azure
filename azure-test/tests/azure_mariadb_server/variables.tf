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
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  version = "=2.41.0"
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
  location = "East US"
}

resource "azurerm_mariadb_server" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  administrator_login          = "steampipeadmin"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_1"
  storage_mb = 5120
  version    = "10.2"

  auto_grow_enabled             = false
  backup_retention_days         = 7
  geo_redundant_backup_enabled  = false
  public_network_access_enabled = true
  ssl_enforcement_enabled       = true
}

output "resource_aka" {
  value = "azure://${azurerm_mariadb_server.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_mariadb_server.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_mariadb_server.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "location" {
  value = azurerm_resource_group.named_test_resource.location
}

output "fqdn" {
  value = azurerm_mariadb_server.named_test_resource.fqdn
}
