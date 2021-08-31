variable "resource_name" {
  type        = string
  default     = "turbot-test-azure-lb-backend-address-pool-20210813"
  description = "Name of the resource used throughout the test."
}

variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "3510ae4d-530b-497d-8f30-53c0616fc6c1"
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
}

resource "azurerm_lb" "named_test_resource" {
  name                   = var.resource_name
  location               = azurerm_resource_group.named_test_resource.location
  resource_group_name    = azurerm_resource_group.named_test_resource.name

  frontend_ip_configuration {
    name                 = var.resource_name
    public_ip_address_id = azurerm_public_ip.named_test_resource.id
  }
}

resource "azurerm_lb_backend_address_pool" "named_test_resource" {
  resource_group_name = azurerm_resource_group.named_test_resource.name
  loadbalancer_id = azurerm_lb.named_test_resource.id
  name            = var.resource_name
}

output "resource_aka" {
  depends_on = [azurerm_lb_backend_address_pool.named_test_resource]
  value      = "azure://${azurerm_lb_backend_address_pool.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_lb_backend_address_pool.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_lb_backend_address_pool.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
