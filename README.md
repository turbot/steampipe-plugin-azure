![image](https://hub.steampipe.io/images/plugins/turbot/azure-social-graphic.png)

# Azure Plugin for Steampipe

Use SQL to query infrastructure including servers, networks, facilities and more from Azure.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/azure)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/azure/tables)
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
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

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-azure/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Azure Plugin](https://github.com/turbot/steampipe-plugin-azure/labels/help%20wanted)
