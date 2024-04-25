
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
  location            = azurerm_resource_group.named_test_resource.location
  resource_group_name = azurerm_resource_group.named_test_resource.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "named_test_resource" {
  name                 = var.resource_name
  resource_group_name  = azurerm_resource_group.named_test_resource.name
  virtual_network_name = azurerm_virtual_network.named_test_resource.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_app_service_environment" "named_test_resource" {
  name                         = var.resource_name
  resource_group_name          = azurerm_resource_group.named_test_resource.name
  subnet_id                    = azurerm_subnet.named_test_resource.id
  pricing_tier                 = "I1"
  front_end_scale_factor       = 10
  internal_load_balancing_mode = "Web, Publishing"
  allowed_user_ip_cidrs        = ["11.22.33.44/32", "55.66.77.0/24"]

  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = "azure://${azurerm_app_service_environment.named_test_resource.id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_app_service_environment.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = azurerm_app_service_environment.named_test_resource.id
}

output "subscription_id" {
  value = var.azure_subscription
}
