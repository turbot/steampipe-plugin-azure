
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