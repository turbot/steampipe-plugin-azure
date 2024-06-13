variable "resource_name" {
  type        = string
  default     = "turbot-test-20200930-create-update"
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
  description = "Azure environment used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.107.0"
    }
  }
}

provider "azurerm" {
  # Cannot be passed as a variable
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  features {}
}

resource "azurerm_resource_group" "named_test_resource" {
  name     = var.resource_name
  location = "East US"
}

resource "azurerm_cdn_frontdoor_profile" "named_test_resource" {
  name                = var.resource_name
  resource_group_name = azurerm_resource_group.named_test_resource.name
  sku_name            = "Standard_AzureFrontDoor"

  tags = {
    environment = "Production"
  }
}

output "resource_id" {
  depends_on = [azurerm_cdn_frontdoor_profile.named_test_resource]
  value      = replace(replace(azurerm_cdn_frontdoor_profile.named_test_resource.id, "resourceGroups", "resourcegroups"), "Profiles", "profiles")
}

output "resource_aka" {
  value = "azure://${replace(replace(azurerm_cdn_frontdoor_profile.named_test_resource.id, "resourceGroups", "resourcegroups"), "Profiles", "profiles")}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_cdn_frontdoor_profile.named_test_resource.id)}"
}

output "resource_name" {
  value = "${var.resource_name}"
}

output "subscription_id" {
  value = var.azure_subscription
}
