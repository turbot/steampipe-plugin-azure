variable "resource_name" {
  type        = string
  default     = "turbot-test-azure-lb-rule-20200805"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "cdffd708-7da0-4cea-abeb-0a4c334d7f64"
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
  location = "West US"
}

resource "azurerm_public_ip" "named_test_resource" {
  name                = var.resource_name
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "named_test_resource" {
  name                   = var.resource_name
  location               = azurerm_resource_group.named_test_resource.location
  resource_group_name    = azurerm_resource_group.named_test_resource.name
  sku                    = "Standard"

  frontend_ip_configuration {
    name                 = var.resource_name
    public_ip_address_id = azurerm_public_ip.named_test_resource.id
  }
}

resource "azurerm_lb_backend_address_pool" "named_test_resource" {
  loadbalancer_id     = azurerm_lb.named_test_resource.id
  name                = var.resource_name
}

resource "azurerm_lb_outbound_rule" "named_test_resource" {
  loadbalancer_id         = azurerm_lb.named_test_resource.id
  name                    = var.resource_name
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.named_test_resource.id

  frontend_ip_configuration {
    name = var.resource_name
  }
}

output "resource_aka" {
  depends_on = [azurerm_lb_rule.named_test_resource]
  value      = "azure://${azurerm_lb_outbound_rule.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_lb_outbound_rule.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_lb_outbound_rule.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
