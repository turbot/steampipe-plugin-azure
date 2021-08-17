---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/azure.svg"
brand_color: "#0089D6"
display_name: "Azure"
name: "azure"
description: "Steampipe plugin for querying resource groups, virtual machines, storage accounts and more from Azure."
og_description: "Query Azure with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/azure-social-graphic.png"
---

# Azure + Steampipe

[Azure](https://azure.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

[Steampipe](https://steampipe.io) is an open-source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  name,
  access_tier,
  sku_name,
  resource_group
from
  azure_storage_account;
```

```
+-------------------------+-------------+--------------+----------------------------------+
| name                    | access_tier | sku_name     | resource_group                   |
+-------------------------+-------------+--------------+----------------------------------+
| parkerrajmodtesting2021 | Hot         | Standard_LRS | azurebackuprg_westus2_1          |
| testsumitsa             | Cool        | Standard_LRS | test_sumit_rg                    |
| sqlvaskpahgwu6znae      | Hot         | Standard_LRS | lalit_test                       |
| sqlvaoggbf26f2ajye      | Hot         | Standard_LRS | turbot_rg                        |
| csg1003200098033c2d     | Hot         | Standard_LRS | cloud-shell-storage-centralindia |
+-------------------------+-------------+--------------+----------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/azure/tables)**

## Get started

### Install

Download and install the latest Azure plugin:

```bash
steampipe plugin install azure
```

### Credentials

| Item        | Description                                                                                                                                                                                                                   |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Use the `az login` command to setup your [Azure Default Connection](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli)                                                                                        |
| Permissions | Grant the `Global Reader` permission to your user.                                                                                                                                                                            |
| Radius      | Each connection represents a single Azure Subscription.                                                                                                                                                                       |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/azuread.spc`).<br />2. Credentials specified in [environment variables](#credentials-from-environment-variables) e.g. `AZURE_SUBSCRIPTION_ID`. |

### Configuration

Installing the latest azure plugin will create a config file (~/.steampipe/config/azure.spc) with a single connection named azure:

```hcl
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
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-azure
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Configuring Azure Credentials

The Azure plugin support multiple formats/authentication mechanisms and they are tried in the below order:

1. [Client Secret Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-saml-bearer-assertion#prerequisites) if set; otherwise
2. [Client Certificate Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-microsoft-identity-platform) if set; otherwise
3. [Resource Owner Password](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth-ropc) if set; otherwise
4. If no credentials are supplied, then the [az cli](https://docs.microsoft.com/en-us/cli/azure/#:~:text=The%20Azure%20command%2Dline%20interface,with%20an%20emphasis%20on%20automation.) credentials are used

If connection arguments are provided, they will always take precedence over [Azure SDK environment variables](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/new-version-quickstart.md#setting-environment-variables), and they are tried in the below order:

### Client Secret Credentials: Azure AD Application ID and Secret.

- `tenant_id`: Specifies the Tenant to which to authenticate.
- `client_id`: Specifies the app client ID to use.
- `client_secret`: Specifies the app secret to use.
- `subscription_id`: Specifies the subscription to connect to.

```hcl
connection "azure_via_sp_secret" {
  plugin            = "azure"
  tenant_id         = "00000000-0000-0000-0000-000000000000"
  subscription_id   = "00000000-0000-0000-0000-000000000000"
  client_id         = "00000000-0000-0000-0000-000000000000"
  client_secret     = "my plaintext password"
}
```

### Client Certificate Credentials: Azure AD Application ID and X.509 Certificate.

- `tenant_id`: Specifies the Tenant to which to authenticate.
- `client_id`: Specifies the app client ID to use.
- `certificate_path`: Specifies the certificate Path to use.
- `certificate_password`: Specifies the certificate password to use.
- `subscription_id`: Specifies the subscription to connect to/

```hcl
connection "azure_via_sp_cert" {
  plugin               = "azure"
  tenant_id            = "00000000-0000-0000-0000-000000000000"
  subscription_id      = "00000000-0000-0000-0000-000000000000"
  client_id            = "00000000-0000-0000-0000-000000000000"
  certificate_path     = "parth/to/file.pem"
  certificate_password = "my plaintext password"
}
```

### Resource Owner Password: Azure AD User and Password. This grant type is _not recommended_, use device login instead if you need interactive login.

- `tenant_id`: Specifies the Tenant to which to authenticate.
- `client_id`: Specifies the app client ID to use.
- `username`: Specifies the username to use.
- `password`: Specifies the password to use.
- `subscription_id`: Specifies the subscription to connect to.

```hcl
connection "password_not_recommended" {
  plugin          = "azure"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  username        = "my-username"
  password        = "plaintext password"
}
```

### Credentials from Environment Variables

The Azure AD plugin will use the standard Azure environment variables to obtain credentials **only if other arguments (`tenant_id`, `client_id`, `client_secret`, `certificate_path`, etc..) are not specified** in the connection:

```sh
export AZURE_TENANT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
export AZURE_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_CLIENT_SECRET="my plaintext secret"
export AZURE_CERTIFICATE_PATH=path/to/file.pem
export AZURE_CERTIFICATE_PASSWORD="my plaintext password"
export AZURE_ENVIRONMENT="AZUREPUBLICCLOUD" # 	Default to "AZUREPUBLICCLOUD". Can be one of "AZUREPUBLICCLOUD", "AZURECHINACLOUD", "AZUREGERMANCLOUD" and "AZUREUSGOVERNMENTCLOUD
```

```hcl
connection "azure" {
  plugin = "azure"
}
```

### Azure CLI: If no credentials are specified

If no credentials are specified and the SDK environment variables are not set, the plugin will use the active credentials from the `az` CLI. You can run `az login` to set up these credentials.

- `subscription_id`: Specifies the subscription to connect to; Otherwise uses the subscription_id set for azure CLI(`az account show`)

```hcl
connection "azure" {
  plugin = "azure"
}
```
