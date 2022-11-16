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
  default     = "d46d7416-f95f-4771-bbb5-529d4c76659c1"
  description = "Azure environment used for the test."
}

provider "azuread" {
  # Cannot be passed as a variable
  version     = "=1.2.2"
  environment = var.azure_environment
}

data "azuread_client_config" "current" {}

resource "azuread_user" "named_test_resource" {
  user_principal_name = "${var.resource_name}@turbotad.onmicrosoft.com"
  display_name        = var.resource_name
  given_name          = var.resource_name
  mail_nickname       = "jdoe"
  password            = "SecretP@sswd99!"
  account_enabled     = false
}

output "resource_aka" {
  depends_on = [azuread_user.named_test_resource]
  value      = "azure:///user/${azuread_user.named_test_resource.id}"
}

output "resource_name" {
  value = azuread_user.named_test_resource.display_name
}

output "object_id" {
  value = azuread_user.named_test_resource.id
}
