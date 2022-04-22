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

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.3.0"
    }
  }
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_mysql_flexible_server" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  administrator_login    = "mradministrator"
  administrator_password = "H@Sh1CoR3!"

  sku_name = "GP_Standard_D2ds_v4"
  version  = "5.7"

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_mysql_flexible_server.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_mysql_flexible_server.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_mysql_flexible_server.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}

output "server_fqdn" {
  value = azurerm_mysql_flexible_server.named_test_resource.fqdn
}
