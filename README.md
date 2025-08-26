![image](https://hub.steampipe.io/images/plugins/turbot/azure-social-graphic.png)

# Azure Plugin for Steampipe

Use SQL to query infrastructure including servers, networks, facilities and more from Azure.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/azure)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/azure/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-azure/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install azure
```

Run a query:

```sql
select name, disk_state from azure_compute_disk where disk_state = 'Unattached'
```

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs/steampipe_sqlite/overview) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/overview) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-azure.git
cd steampipe-plugin-azure
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/azure.spc
```

Try it!

```
steampipe query
> .inspect azure
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Storage Data Plane Authentication (Blobs & Queues)

The `azure_storage_blob` and `azure_storage_queue` tables now default to Azure AD (OAuth) authentication using your configured identity (environment variables, managed identity, CLI, Azure CLI login, etc.). The prior implicit Shared Key path (automatic key listing) has been removed to align with hardened environments that disable Shared Key access.

In almost all cases you should rely on Azure AD RBAC (e.g. assign the principal the `Storage Blob Data Reader` or `Storage Queue Data Reader` role). The controls below are advanced / legacy overrides only—avoid using them unless you have a specific need.

Advanced (optional) connection overrides (`azure.spc`):

```
connection "azure" {
	plugin  = "azure"
	# auth_mode can be: aad | shared_key | sas
	auth_mode = "aad"

	# For shared_key mode either provide the key or allow key listing
	# storage_account_key       = "<account key>"
	# allow_storage_key_listing = true

	# For sas mode provide a service SAS (with leading ? optional)
	# storage_sas_token = "?sv=..."
}
```

Notes:
* Default (no settings): Azure AD (`auth_mode` omitted) for both blobs and queues.
* `auth_mode = aad`: Explicit Azure AD (default path).
* `auth_mode = shared_key`: Uses account Shared Key. You must supply `storage_account_key` OR set `allow_storage_key_listing = true`. Fails fast if the account has disabled shared key access (`allowSharedKeyAccess = false`).
* `auth_mode = sas`: Uses a service/account SAS token (`storage_sas_token`). Scope/permissions limited to the SAS grants.
* `auth_mode = auto`: Deprecated; treated as `aad` with a warning.
* Friendly errors are returned for common issues (disabled shared key, permission mismatch / missing role assignments).

Track 2 SDK adoption:
* Blobs: `github.com/Azure/azure-sdk-for-go/sdk/storage/azblob`
* Queues: `github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue`

Future work may extend this pattern to additional storage data-plane surfaces (files, tables) as Track 2 coverage matures.

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Azure Plugin](https://github.com/turbot/steampipe-plugin-azure/labels/help%20wanted)
