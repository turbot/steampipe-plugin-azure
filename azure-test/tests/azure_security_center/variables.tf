
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

variable "provision_name" {
  type        = string
  default     = "default"
  description = "Name of the resource used throughout the test."
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
  path = "${path.cwd}/autoProvision.json"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "az security auto-provisioning-setting show -n ${var.provision_name} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "auto_provision" {
  depends_on = [null_resource.named_test_resource]
  value      = lookup(jsondecode(data.local_file.input.content), "autoProvision", "autoProvision")
}

output "resource_id" {
  value      = lookup(jsondecode(data.local_file.input.content), "id", "id")
}