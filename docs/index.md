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
### Scope
An azure connection is scoped to a single azure subscription, with a single set of credentials.

The `az` CLI and Azure APIS are inherently global - while resources are created in a region, the commands to manage them are not limited to a single region. The Azure Steampipe plugin will query all Azure regions for the subscription.


### Configuration

Installing the latest azure plugin will create a default connection named `azure`. This connection will dynamically determine the scope and credentials using the same mechanism as the `az` CLI.  In effect, this means that by default Steampipe will execute with the same credentials against the same active subscription as the `az` cli tool. Changing/switching your active settings or subscription will also change the behavior of the azure_default connection.

(Of course this also  implies that the `azure` cli needs to be configured with the proper credentials before the steampipe azure plugin can be used).

