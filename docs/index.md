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

[Azure](https://azure.microsoft.com) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis. 

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

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

| Item        | Description                                                                                                                                                                                                                     |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Use the `az login` command to setup your [Azure Default Connection](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli).                                                                                         |
| Permissions | Assign the `Reader` role to your user or service principal in the subscription.                                                                                                                                                                              |
| Radius      | Each connection represents a single Azure subscription.                                                                                                                                                                         |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/azure.spc`).<br />2. Credentials specified in [environment variables](#credentials-from-environment-variables), e.g., `AZURE_SUBSCRIPTION_ID`.<br />3. Credentials from the Azure CLI. |

### Configuration

Installing the latest azure plugin will create a config file (~/.steampipe/config/azure.spc) with a single connection named azure:

```hcl
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

  # List of additional azure error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  #ignore_error_codes = ["NoAuthenticationInformation", "InvalidAuthenticationInfo", "AccountIsDisabled", "UnauthorizedOperation", "UnrecognizedClientException", "AuthorizationError", "AuthenticationFailed", "InsufficientAccountPermissions"]
}
```

## Multi-Subscription Connections

You may create multiple azure connections:

```hcl
connection "azure_all" {
  type        = "aggregator"
  plugin      = "azure"
  connections = ["azure_*"]
}

connection "azure_sub_1" {
  plugin          = "azure"
  subscription_id = "azure_01"
}

connection "azure_sub_2" {
  plugin          = "azure"
  subscription_id = "azure_02"
}

connection "azure_sub_3" {
  plugin          = "azure"
  subscription_id = "azure_03"
}
```

Depending on the mode of authentication, a multi-subscription configuration can also look like:

```hcl
connection "azure_all" {
  type        = "aggregator"
  plugin      = "azure"
  connections = ["azure_*"]
}

connection "azure_sub_1" {
  plugin          = "azure"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  client_secret   = "~dummy@3password"
}

connection "azure_sub_2" {
  plugin          = "azure"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  client_secret   = "~dummy@3password"
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from azure_sub_1.azure_subscription
```

Alternatively, you can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/using-steampipe/managing-connections#setting-the-search-path):

```sql
select * from azure_subscription
```

You can create multi-subscription connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators). Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection:

```hcl
connection "azure_all" {
  plugin      = "azure"
  type        = "aggregator"
  connections = ["azure_sub_1", "azure_sub_2", "azure_sub_3"]
}
```

Querying tables from this connection will return results from the `azure_sub_1`, `azure_sub_2`, and `azure_sub_3` connections:

```sql
select * from azure_all.azure_subscription
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the Azure plugin connections whose names begin with `azure_`:

```hcl
connection "azure_all" {
  type        = "aggregator"
  plugin      = "azure"
  connections = ["azure_*"]
}
```

## Configuring Azure Credentials

The Azure plugin support multiple formats/authentication mechanisms and they are tried in the below order:

1. [Client Secret Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-saml-bearer-assertion#prerequisites) if set; otherwise
2. [Client Certificate Credentials](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-microsoft-identity-platform) if set; otherwise
3. [Resource Owner Password](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth-ropc) if set; otherwise
4. If no credentials are supplied, then the [az cli](https://docs.microsoft.com/en-us/cli/azure/#:~:text=The%20Azure%20command%2Dline%20interface,with%20an%20emphasis%20on%20automation.) credentials are used

If connection arguments are provided, they will always take precedence over [Azure SDK environment variables](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/new-version-quickstart.md#setting-environment-variables), and they are tried in the below order:

### Client Secret Credentials

You may specify the tenant ID, subscription ID, client ID, and client secret to authenticate:

- `tenant_id`: Specify the tenant to authenticate with.
- `subscription_id`: Specify the subscription to query.
- `client_id`: Specify the app client ID to use.
- `client_secret`: Specify the app secret to use.

```hcl
connection "azure_via_sp_secret" {
  plugin            = "azure"
  tenant_id         = "00000000-0000-0000-0000-000000000000"
  subscription_id   = "00000000-0000-0000-0000-000000000000"
  client_id         = "00000000-0000-0000-0000-000000000000"
  client_secret     = "my plaintext password"
}
```

### Client Certificate Credentials

You may specify the tenant ID, subscription ID, client ID, certificate path, and certificate password to authenticate:

- `tenant_id`: Specify the tenant to authenticate with.
- `subscription_id`: Specify the subscription to query.
- `client_id`: Specify the app client ID to use.
- `certificate_path`: Specify the certificate path to use.
- `certificate_password`: Specify the certificate password to use.

```hcl
connection "azure_via_sp_cert" {
  plugin               = "azure"
  tenant_id            = "00000000-0000-0000-0000-000000000000"
  subscription_id      = "00000000-0000-0000-0000-000000000000"
  client_id            = "00000000-0000-0000-0000-000000000000"
  certificate_path     = "path/to/file.pem"
  certificate_password = "my plaintext password"
}
```

### Resource Owner Password

**Note:** This grant type is _not recommended_, use device login instead if you need interactive login.

You may specify the tenant ID, subscription ID, client ID, username, and password to authenticate:

- `tenant_id`: Specify the tenant to authenticate with.
- `subscription_id`: Specify the subscription to query.
- `client_id`: Specify the app client ID to use.
- `username`: Specify the username to use.
- `password`: Specify the password to use.

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

### Azure Managed Identity

Steampipe works with managed identities (formerly known as Managed Service Identity), provided it is running in Azure, e.g., on a VM. All configuration is handled by Azure. See [Azure Managed Identities](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview) for more details.

- `tenant_id`: Specify the tenant to authenticate with.
- `subscription_id`: Specify the subscription to query.
- `client_id`: Specify the app client ID of managed identity to use.

```hcl
connection "azure_msi" {
  plugin          = "azure"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "00000000-0000-0000-0000-000000000000"
}
```

### Azure CLI

If no credentials are specified and the SDK environment variables are not set, the plugin will use the active credentials from the Azure CLI. You can run `az login` to set up these credentials.

- `subscription_id`: Specifies the subscription to connect to. If not set, use the subscription ID set in the Azure CLI (`az account show`)

```hcl
connection "azure" {
  plugin = "azure"
}
```

### Credentials from Environment Variables

The Azure AD plugin will use the standard Azure environment variables to obtain credentials **only if other arguments (`tenant_id`, `client_id`, `client_secret`, `certificate_path`, etc..) are not specified** in the connection:

```sh
export AZURE_ENVIRONMENT="AZUREPUBLICCLOUD" # Defaults to "AZUREPUBLICCLOUD". Valid environments are "AZUREPUBLICCLOUD", "AZURECHINACLOUD", "AZUREGERMANCLOUD" and "AZUREUSGOVERNMENTCLOUD"
export AZURE_TENANT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
export AZURE_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_CLIENT_SECRET="my plaintext secret"
export AZURE_CERTIFICATE_PATH="path/to/file.pem"
export AZURE_CERTIFICATE_PASSWORD="my plaintext password"
```

```hcl
connection "azure" {
  plugin = "azure"
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-azure
- Community: [Slack Channel](https://steampipe.io/community/join)
