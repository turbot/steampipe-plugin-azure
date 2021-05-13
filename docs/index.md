---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/azure.svg"
brand_color: "#0089D6"
display_name: "Azure"
name: "azure"
description: "Steampipe plugin for querying resource groups, virtual machines, storage accounts and more from Azure."
og_description: Query Azure with SQL! Open source CLI. No DB required. 
og_image: "/images/plugins/turbot/azure-social-graphic.png"
---

# Azure + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[Azure](https://azure.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis. 

For example:

```sql
select 
  display_name, 
  user_type 
from 
  azure_ad_user 
```

```
+----------------------+-----------+
|     display_name     | user_type |
+----------------------+-----------+
| Dwight Schrute       | Member    |
| Jim Halpert          | Member    |
| Pam Beesly           | Member    |
| Michael Scott        | Member    |
| Stanley Hudson       | Member    |
+----------------------+-----------+
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

| Item | Description |
| - | - |
| Credentials | Use the `az login` command to setup your [Azure Default Connection](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli) |
| Permissions | Grant the `Global Reader` permission to your user. |
| Radius | Each connection represents a single Azure Subscription. |
| Resolution |  1. Client Credentials<br />2. Client Certificate<br />3. Resource Owner Password<br />4. Azure Managed Service Identity<br />5. If no credentials are supplied, then the az cli credentials are used |

### Configuration

Installing the latest azure plugin will create a config file (`~/.steampipe/config/azure.spc`) with a single connection named `azure`:

```hcl
connection "azure_sub_000000" {
  plugin = "azure"
  subscription_id = "00000000-0000-0000-0000-000000000000"
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-azure
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)
