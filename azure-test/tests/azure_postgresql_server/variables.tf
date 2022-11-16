
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
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c1"
  description = "Azure subscription used for the test."
}

provider "azurerm" {
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

resource "azurerm_postgresql_server" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  sku_name   = "B_Gen5_1"
  version    = "9.5"
  storage_mb = 5120

  backup_retention_days             = 7
  geo_redundant_backup_enabled      = false
  auto_grow_enabled                 = false
  infrastructure_encryption_enabled = false

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"

  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_2"

  tags = {
    name = var.resource_name
  }
}

resource "azurerm_postgresql_configuration" "named_test_resource" {
  name                = "log_checkpoints"
  resource_group_name = azurerm_resource_group.named_test_resource.name
  server_name         = azurerm_postgresql_server.named_test_resource.name
  value               = "on"
}

resource "azurerm_postgresql_firewall_rule" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  server_name         = azurerm_postgresql_server.named_test_resource.name
  start_ip_address    = "40.112.8.12"
  end_ip_address      = "40.112.8.12"
}

resource "azurerm_postgresql_active_directory_administrator" "named_test_resource" {
  server_name         = azurerm_postgresql_server.named_test_resource.name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.object_id
}

output "resource_aka" {
  value = "azure://${azurerm_postgresql_server.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_postgresql_server.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_postgresql_server.named_test_resource.id
}

output "location" {
  value = lower(azurerm_resource_group.named_test_resource.location)
}

output "subscription_id" {
  value = var.azure_subscription
}

output "server_fqdn" {
  value = azurerm_postgresql_server.named_test_resource.fqdn
}
