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
  default     = "cdffd708-7da0-4cea-abeb-0a4c334d7f64"
  description = "Azure environment used for the test."
}

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.78.0"
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

resource "azurerm_frontdoor" "named_test_resource" {
  name                                         = var.resource_name
  resource_group_name                          = azurerm_resource_group.named_test_resource.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "${var.resource_name}RoutingRule1"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = ["${var.resource_name}FrontendEndpoint1"]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "${var.resource_name}BackendBing"
    }
  }

  backend_pool_load_balancing {
    name = "${var.resource_name}LoadBalancingSettings1"
  }

  backend_pool_health_probe {
    name = "${var.resource_name}HealthProbeSetting1"
  }

  backend_pool {
    name = "${var.resource_name}BackendBing"
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = "${var.resource_name}LoadBalancingSettings1"
    health_probe_name   = "${var.resource_name}HealthProbeSetting1"
  }

  frontend_endpoint {
    name      = "${var.resource_name}FrontendEndpoint1"
    host_name = "${var.resource_name}.azurefd.net"
  }
}

output "resource_id" {
  depends_on = [azurerm_frontdoor.named_test_resource]
  value      = replace(replace(azurerm_frontdoor.named_test_resource.id, "resourceGroups", "resourcegroups"), "frontDoors", "frontdoors")
}

output "resource_aka" {
  value = "azure://${replace(replace(azurerm_frontdoor.named_test_resource.id, "resourceGroups", "resourcegroups"), "frontDoors", "frontdoors")}"
}

output "resource_aka_lower" {
  value = "azure://${lower(azurerm_frontdoor.named_test_resource.id)}"
}

output "resource_name" {
  value = var.resource_name
}

output "subscription_id" {
  value = var.azure_subscription
}
