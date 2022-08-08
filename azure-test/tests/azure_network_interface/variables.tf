
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
  location = "West US"
}

resource "azurerm_virtual_network" "named_test_resource" {
  name                = var.resource_name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = var.resource_name
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  ip_configuration {
    name                          = var.resource_name
    subnet_id                     = azurerm_subnet.named_test_resource.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_network_interface.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_network_interface.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_network_interface.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
