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

resource "azuread_group" "named_test_resource" {
  name        = var.resource_name
  description = "For testing tbt cli"
}

output "resource_aka" {
  depends_on = [azuread_group.named_test_resource]
  value      = "azure:///group/${azuread_group.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "object_id" {
  value = azuread_group.named_test_resource.id
}
