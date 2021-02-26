---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/azure.svg"
brand_color: "#0089D6"
display_name: "Azure"
name: "azure"
description: "Steampipe plugin for 30+ Azure services and resource types."
---

# Azure

The Azure plugin is used to interact with the many resources supported by Microsoft Azure.

### Installation

To download and install the latest azure plugin:

```bash
$ ./steampipe plugin install azure
Installing plugin azure...
$
```

Installing the latest azure plugin will create a default connection named `azure`. This connection will dynamically determine the scope and credentials using the same mechanism as the `az` CLI. In effect, this means that by default Steampipe will execute with the same credentials against the same active subscription as the `az` cli tool. Changing/switching your active settings or subscription will also change the behavior of the default `azure` connection.

(Of course this also implies that the `azure` cli needs to be configured with the proper credentials before the steampipe azure plugin can be used).

## Connection Configuration

Connection configurations are defined using HCL in one or more Steampipe config files. Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.

### Scope

An azure connection is scoped to a single Azure Subscription, with a single set of credentials.  If no subscription id is specified, the current active subscription per the `az` cli will be used.

The `az` CLI and Azure APIS are inherently global - while resources are created in a region, the commands to manage them are not limited to a single region. The Azure Steampipe plugin will query all Azure regions for the subscription.

### Credentials

The Azure plugin support multiple formats / authentication mechanisms, and they are tried in the standard order:

1. **Client Credentials** if set; otherwise
2. **Client Certificate** if set; otherwise
3. **Resource Owner Password** if set; otherwise
4. **Azure Managed Service Identity** if set; otherwise
5. If no credentials are supplied, then the `az` cli credentials are used

## Connection Arguments

If connection arguments are provided, they will always take precedence over [Azure SDK environment variables](https://github.com/Azure/azure-sdk-for-go#more-authentication-details), and they are tried in the standard order:

1. **Client Credentials**: Azure AD Application ID and Secret.

   - `tenant_id`: Specifies the Tenant to which to authenticate.
   - `client_id`: Specifies the app client ID to use.
   - `client_secret`: Specifies the app secret to use.
   - `subscription_id`: Specifies the subscription to connect to.

   ```hcl
   connection "azure_via_sp_secret" {
     plugin            = "azure"
     subscription_id   = "00000000-0000-0000-0000-000000000000"
     client_id         = "00000000-0000-0000-0000-000000000000"
     client_secret     = "my plaintext password"
     tenant_id         = "00000000-0000-0000-0000-000000000000"
   }
   ```

2. **Client Certificate**: Azure AD Application ID and X.509 Certificate.

   - `tenant_id`: Specifies the Tenant to which to authenticate.
   - `client_id`: Specifies the app client ID to use.
   - `client_certificate_path`: Specifies the certificate Path to use.
   - `client_certificate_password`: Specifies the certificate password to use.
   - `subscription_id`: Specifies the subscription to connect to/

   ```hcl
   connection "azure_via_sp_cert" {
     plugin                      = "azure"
     subscription_id             = "00000000-0000-0000-0000-000000000000"
     client_id                   = "00000000-0000-0000-0000-000000000000"
     client_certificate_path     = "parth/to/file.pfx"
     client_certificate_password = "my plaintext password"
     tenant_id                   = "00000000-0000-0000-0000-000000000000"
   }
   ```

3. **Resource Owner Password**: Azure AD User and Password. This grant type is _not recommended_, use device login instead if you need interactive login.

   - `tenant_id`: Specifies the Tenant to which to authenticate.
   - `client_id`: Specifies the app client ID to use.
   - `username`: Specifies the username to use.
   - `password`: Specifies the password to use.
   - `subscription_id`: Specifies the subscription to connect to.

   ```hcl
   connection "password_not_recommended" {
     plugin                      = "azure"
     subscription_id             = "00000000-0000-0000-0000-000000000000"
     client_id                   = "00000000-0000-0000-0000-000000000000"
     tenant_id                   = "00000000-0000-0000-0000-000000000000"
     username                    = "my-username"
     password                    = "plaintext password"
   }
   ```

4. **Azure Managed Service Identity**: Delegate credential management to the platform. Requires that code is running in Azure, e.g. on a VM. All configuration is handled by Azure. See [Azure Managed Service Identity](https://docs.microsoft.com/azure/active-directory/msi-overview) for more details.

   - `use_msi`: pass `true` to use msi credentials
   - `subscription_id`: Specifies the subscription to connect to/
   - `tenant_id`: Specifies the Tenant to which to authenticate.

   ```hcl
   connection "azure_msi" {
     plugin            = "azure"
     use_msi           = true
     subscription_id   = "00000000-0000-0000-0000-000000000000"
     tenant_id         = "11111111-1111-1111-1111-111111111111"
   }
   ```

5. If no credentials are specified and the SDK environment variables are not set, the plugin will use the active credentials from the `az` cli (as used in the [Azure Default Connection](#azure-default-connection)). You can run `az login` to set up these credentials.

   ```hcl
   connection "azure" {
     plugin = "azure"
   }
   ```

Note that you can optionally specify the `subscription_id` as well, to use the CLI login but set a specific subscription:

  ```hcl
  connection "azure_sub_000000" {
    plugin = "azure"
    subscription_id = "00000000-0000-0000-0000-000000000000"
  }
  ```
