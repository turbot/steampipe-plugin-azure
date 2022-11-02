connection "azure" {
  plugin = "azure"

  # The Azure cloud environment to use, defaults to AZUREPUBLICCLOUD
  # Valid environments are AZUREPUBLICCLOUD, AZURECHINACLOUD, AZUREGERMANCLOUD, AZUREUSGOVERNMENTCLOUD
  # If using Azure CLI for authentication, make sure to also set the default environment: https://docs.microsoft.com/en-us/cli/azure/manage-clouds-azure-cli
  # environment = "AZUREPUBLICCLOUD"

  # You can connect to Azure using one of options below:

  # Use client secret authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-2-create-a-new-application-secret)
  # tenant_id       = "00000000-0000-0000-0000-000000000000"
  # subscription_id = "00000000-0000-0000-0000-000000000000"
  # client_id       = "00000000-0000-0000-0000-000000000000"
  # client_secret   = "~dummy@3password"

  # Use client certificate authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-1-upload-a-certificate)
  # tenant_id            = "00000000-0000-0000-0000-000000000000"
  # subscription_id      = "00000000-0000-0000-0000-000000000000"
  # client_id            = "00000000-0000-0000-0000-000000000000"
  # certificate_path     = "~/home/azure_cert.pem"
  # certificate_password = "notreal~pwd"

  # Use resource owner password authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth-ropc)
  # tenant_id       = "00000000-0000-0000-0000-000000000000"
  # subscription_id = "00000000-0000-0000-0000-000000000000"
  # client_id       = "00000000-0000-0000-0000-000000000000"
  # username        = "my-username"
  # password        = "plaintext password"

  # Use a managed identity (https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview)
  # This method is useful with Azure virtual machines
  # tenant_id       = "00000000-0000-0000-0000-000000000000"
  # subscription_id = "00000000-0000-0000-0000-000000000000"
  # client_id       = "00000000-0000-0000-0000-000000000000"

  # If no credentials are specified, the plugin will use Azure CLI authentication

<<<<<<< HEAD
  # The maximum number of attempts (including the initial call) Steampipe will
  # Defaults to 9 and must be greater than or equal to 1.
  #max_error_retry_attempts = 9

  # The minimum retry delay in milliseconds after which retries will be performed.
  # This delay is also used as a base value when calculating the exponential backoff retry times.
  # Defaults to 25ms and must be greater than or equal to 1ms.
  #min_error_retry_delay = 25
=======
  # List of additional Azure error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  #ignore_error_codes = ["NoAuthenticationInformation", "InvalidAuthenticationInfo", "AccountIsDisabled", "UnauthorizedOperation", "UnrecognizedClientException", "AuthorizationError", "AuthenticationFailed", "InsufficientAccountPermissions"]
>>>>>>> 7a2a941e410c5bff5fd49f76e2fa40bd4576f4e0
}
