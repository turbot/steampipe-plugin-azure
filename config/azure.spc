connection "azure" {
  plugin = "azure"

  # "Defaults to "AZUREPUBLICCLOUD". Valid environments are "AZUREPUBLICCLOUD", "AZURECHINACLOUD", "AZUREGERMANCLOUD" and "AZUREUSGOVERNMENTCLOUD"
  # environment = "AZUREPUBLICCLOUD"

  # You can connect to Azure using one of options below:

  # Use client secret authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-2-create-a-new-application-secret)
  # tenant_id             = "00000000-0000-0000-0000-000000000000"
  # subscription_id       = "00000000-0000-0000-0000-000000000000"
  # client_id             = "00000000-0000-0000-0000-000000000000"
  # client_secret         = "~dummy@3password"


  # Use client certificate authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#option-1-upload-a-certificate)
  # tenant_id             = "00000000-0000-0000-0000-000000000000"
  # subscription_id       = "00000000-0000-0000-0000-000000000000"
  # client_id             = "00000000-0000-0000-0000-000000000000"
  # certificate_path      = "~/home/azure_cert.pem"
  # certificate_password  = "notreal~pwd"
  #

  # Use resource owner password authentication (https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth-ropc)
  # tenant_id       = "00000000-0000-0000-0000-000000000000"
  # subscription_id = "00000000-0000-0000-0000-000000000000"
  # client_id       = "00000000-0000-0000-0000-000000000000"
  # username        = "my-username"
  # password        = "plaintext password"

  # If no credentials are specified, the plugin will use Azure CLI authentication
}
