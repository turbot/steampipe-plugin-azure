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
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = var.resource_name
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}
resource "azurerm_private_dns_zone" "named_test_resource" {
  name                = "test.${var.resource_name}.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "named_test_resource" {
  name                  = "${var.resource_name}.com"
  private_dns_zone_name = azurerm_private_dns_zone.named_test_resource.name
  virtual_network_id    = azurerm_virtual_network.named_test_resource.id
  resource_group_name   = azurerm_resource_group.named_test_resource.name
  depends_on            = [azurerm_subnet.named_test_resource]
}

resource "azurerm_postgresql_flexible_server" "named_test_resource" {
  name                          = var.resource_name
  resource_group_name           = azurerm_resource_group.named_test_resource.name
  location                      = azurerm_resource_group.named_test_resource.location
  version                       = "12"
  delegated_subnet_id           = azurerm_subnet.named_test_resource.id
  private_dns_zone_id           = azurerm_private_dns_zone.named_test_resource.id
  administrator_login           = "psqladmin"
  administrator_password        = "H@Sh1CoR3!"
  zone                          = "1"

  storage_mb   = 32768

  sku_name   = "GP_Standard_D4s_v3"
  depends_on = [azurerm_private_dns_zone_virtual_network_link.named_test_resource]
  tags = {
    name = var.resource_name
  }

}

output "resource_aka" {
  value = "azure://${azurerm_postgresql_flexible_server.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_postgresql_flexible_server.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_postgresql_flexible_server.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}