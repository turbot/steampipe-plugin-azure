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
  description = "Azure subscription used for the test."
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
}

resource "azurerm_cdn_frontdoor_custom_domain" "named_test_resource" {
  name                     = var.resource_name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.named_test_resource.id
  host_name                = "${var.resource_name}.example.com"

  tls {
    certificate_type = "ManagedCertificate"
  }
}

locals {
  # Azure normalises casing in IDs — replicate the same transforms as the profile test
  resource_id = replace(
    replace(
      azurerm_cdn_frontdoor_custom_domain.named_test_resource.id,
      "resourceGroups",
      "resourcegroups"
    ),
    "customDomains",
    "customdomains"
  )
}

output "resource_id" {
  depends_on = [azurerm_cdn_frontdoor_custom_domain.named_test_resource]
  value      = local.resource_id
}

output "resource_aka" {
  value = "azure://${local.resource_id}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_cdn_frontdoor_custom_domain.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "profile_name" {
  value = azurerm_cdn_frontdoor_profile.named_test_resource.name
}

output "subscription_id" {
  value = var.azure_subscription
}
