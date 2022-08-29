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

provider "azuread" {
  # Cannot be passed as a variable
  version         = "=0.10.0"
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azuread_client_config" "current" {}

resource "azuread_application" "named_test_resource" {
  name                       = var.resource_name
  homepage                   = "http://homepage"
  identifier_uris            = ["http://${var.resource_name}"]
  reply_urls                 = ["http://replyurl"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}

resource "azuread_service_principal" "named_test_resource" {
  application_id               = azuread_application.named_test_resource.application_id
  app_role_assignment_required = false
}

output "resource_aka" {
  depends_on = [azuread_service_principal.named_test_resource]
  value      = "azure:///serviceprincipal/${azuread_service_principal.named_test_resource.id}"
}

output "resource_name" {
  value = azuread_service_principal.named_test_resource.display_name
}

output "object_id" {
  value = azuread_service_principal.named_test_resource.id
}
