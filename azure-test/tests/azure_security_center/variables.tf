
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
  description = "Azure subscription used for the test."
}

provider "azurerm" {
  # Cannot be passed as a variable
  version         = "=2.43.0"
  features {}
  environment     = var.azure_environment
  subscription_id = var.azure_subscription
}

data "azurerm_client_config" "current" {}

locals {
  provision_path = "${path.cwd}/autoProvision.json"
  contact_path = "${path.cwd}/contact.json"

}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "az security auto-provisioning-setting list > ${local.provision_path}"
  }
  provisioner "local-exec" {
    command = "az security contact list > ${local.contact_path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.provision_path
}

data "local_file" "input_contact" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.contact_path
}

output "auto_provision" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(jsondecode(data.local_file.input.content)[0], "autoProvision", "autoProvision")
}

output "provision_id" {
   depends_on = [null_resource.named_test_resource]
   value      = lookup(jsondecode(data.local_file.input.content)[0], "id", "id")
}

output "contact_id" {
   depends_on = [null_resource.named_test_resource]
   value      = lookup(jsondecode(data.local_file.input_contact.content)[0], "id", "id")
}

output "email" {
   depends_on = [null_resource.named_test_resource]
   value      = lookup(jsondecode(data.local_file.input_contact.content)[0], "email", "email")
}

output "alert_notification" {
   depends_on = [null_resource.named_test_resource]
   value      = lookup(jsondecode(data.local_file.input_contact.content)[0], "alertNotifications", "alertNotifications")
}

output "alert_to_admin" {
   depends_on = [null_resource.named_test_resource]
   value      = lookup(jsondecode(data.local_file.input_contact.content)[0], "alertsToAdmins", "alertsToAdmins")
}

output "name" {
   depends_on = [null_resource.named_test_resource]
   value      = lookup(jsondecode(data.local_file.input_contact.content)[0], "name", "name")
}