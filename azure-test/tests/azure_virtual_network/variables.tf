
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

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  address_space       = ["10.1.2.0/24"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  tags = {
    name = var.resource_name
  }

  subnet {
    name           = "subnet1"
    address_prefix = "10.1.2.0/25"
  }
}

output "resource_aka" {
  value = "azure://${azurerm_virtual_network.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_virtual_network.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_virtual_network.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
