
variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c"
  description = "Azure environment used for the test."
}

provider "azuread" {
  # Cannot be passed as a variable
  environment     = var.azure_environment
  tenant_id = var.azure_subscription
}

data "azuread_client_config" "current" {}

output "subscription_id" {
  value = var.azure_subscription
}
