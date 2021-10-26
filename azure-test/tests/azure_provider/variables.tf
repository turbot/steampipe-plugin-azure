
variable "azure_environment" {
  type        = string
  default     = "public"
  description = "Azure environment used for the test."
}

variable "azure_subscription" {
  type        = string
  default     = "d7245080-b4ae-4fe5-b6fa-2e71b3dae6c8"
  description = "Azure environment used for the test."
}

provider "azuread" {
  # Cannot be passed as a variable
  version         = "=0.10.0"
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azuread_client_config" "current" {}

output "subscription_id" {
  value = var.azure_subscription
}
