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
  default     = "cdffd708-7da0-4cea-abeb-0a4c334d7f64"
  description = "Azure environment used for the test."
}

variable "azure_tenant" {
  type        = string
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c"
  description = "Azure environment used for the test."
}

provider "azuread" {
  # Cannot be passed as a variable
  # version         = "=0.10.0"
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
  tenant_id       = var.azure_tenant
}

data "azurerm_client_config" "current" {
}

provider "azurerm" {
  features {}
}

output "current_tenant_display_name" {
  value = data.azurerm_client_config.current.tenant_id
}

output "tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}
