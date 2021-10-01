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
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
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

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_sql_server" "named_test_resource" {
  name                         = var.resource_name
  resource_group_name          = azurerm_resource_group.named_test_resource.name
  location                     = azurerm_resource_group.named_test_resource.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_elasticpool" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  location            = azurerm_resource_group.named_test_resource.location
  server_name         = azurerm_sql_server.named_test_resource.name
  max_size_gb         = 4.8828125

  sku {
    name     = "BasicPool"
    tier     = "Basic"
    capacity = 50
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 5
  }
}

output "resource_aka" {
  value = "azure://${azurerm_mssql_elasticpool.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_mssql_elasticpool.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_mssql_elasticpool.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
