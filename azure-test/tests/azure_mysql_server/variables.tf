
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

resource "azurerm_mysql_server" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  administrator_login          = "mradministrator"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_1"
  storage_mb = 5120
  version    = "5.6"

  auto_grow_enabled                 = false
  backup_retention_days             = 7
  geo_redundant_backup_enabled      = false
  infrastructure_encryption_enabled = false
  public_network_access_enabled     = true
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"

  tags = {
    "name" = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_mysql_server.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_mysql_server.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_mysql_server.named_test_resource.id
}

output "location" {
  value = lower(azurerm_resource_group.named_test_resource.location)
}

output "subscription_id" {
  value = var.azure_subscription
}

output "server_fqdn" {
  value = azurerm_mysql_server.named_test_resource.fqdn
}
