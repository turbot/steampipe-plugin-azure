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
  default     = "f2ecd657-84a7-4445-9b2a-280aa402eb9f"
  description = "Azure environment used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  version         = "=1.36.0"
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azuread_client_config" "current" {}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "West Europe"
}

resource "azurerm_public_ip" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name

  frontend_ip_configuration {
    name                 = var.resource_name
    public_ip_address_id = azurerm_public_ip.named_test_resource.id
  }

  tags = {
    name = var.resource_name
  }
}

output "public_ip_id" {
  value = azurerm_public_ip.named_test_resource.id
}

output "region" {
  value = azurerm_public_ip.named_test_resource.location
}

output "resource_aka" {
  depends_on = [azurerm_lb.named_test_resource]
  value      = "azure://${azurerm_lb.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_lb.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_lb.named_test_resource.id
}
