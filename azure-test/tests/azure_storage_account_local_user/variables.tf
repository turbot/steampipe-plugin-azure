variable "resource_name" {
  type        = string
  default     = "turbot-test-20240318"
  description = "Name of the resource used throughout the test."
}

variable "azure_location" {
  type        = string
  default     = "eastus"
  description = "Azure region used for the test."
} 