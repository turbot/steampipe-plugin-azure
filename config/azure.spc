connection "azure" {
  plugin = "azure"

  # "Defaults to "AZUREPUBLICCLOUD". Can be one of "AZUREPUBLICCLOUD", "AZURECHINACLOUD", "AZUREGERMANCLOUD" and "AZUREUSGOVERNMENTCLOUD"
  # environment = "AZUREPUBLICCLOUD"

  # You may connect to azure using more than one option
  # 1. For client secret authentication, specify TenantID, ClientID and ClientSecret.
  # required options:
  # tenant_id             = "00000000-0000-0000-0000-000000000000"
  # subscription_id       = "00000000-0000-0000-0000-000000000000"
  # client_id             = "00000000-0000-0000-0000-000000000000"
  # client_secret         = "~dummy@3password"


  # 2. client certificate authentication, specify TenantID, ClientID and ClientCertData / ClientCertPath.
  # required options:
  # tenant_id             = "00000000-0000-0000-0000-000000000000"
  # subscription_id       = "00000000-0000-0000-0000-000000000000"
  # client_id             = "00000000-0000-0000-0000-000000000000"
  # certificate_path      = "~/home/azure_cert.pem"
  # certificate_password  = "notreal~pwd"
  #

  # 3. resource owner password
  # required options:
  # tenant_id       = "00000000-0000-0000-0000-000000000000"
  # subscription_id = "00000000-0000-0000-0000-000000000000"
  # client_id       = "00000000-0000-0000-0000-000000000000"
  # username        = "my-username"
  # password        = "plaintext password"

  # 4. Azure CLI authentication (if enabled) is attempted last
}
