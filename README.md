<p align="center">
    <h1 align="center">Azure Plugin for Steampipe</h1>
</p>

<p align="center">
  <a aria-label="Steampipe logo" href="https://steampipe.io">
    <img src="https://steampipe.io/images/steampipe_logo_wordmark_padding.svg" height="28">
  </a>
  <a aria-label="License" href="LICENSE">
    <img alt="" src="https://img.shields.io/static/v1?label=license&message=MPL-2.0&style=for-the-badge&labelColor=777777&color=F3F1F0">
  </a>
</p>

## Query Azure with SQL

Use SQL to query compute disks, virtual machines and more from your Azure subsription. For example:

```sql
select
  name,
  disk_state
from
  azure_compute_disk
where
  disk_state = 'Unattached';
```

Learn about [Steampipe](https://steampipe.io/).

## Get started

**[Table documentation and examples &rarr;](https://hub.steampipe.io/plugins/turbot/azure)**

Install the plugin:

```shell
steampipe plugin install azure
```

## Get involved

### Community

The Steampipe community can be found on [GitHub Discussions](https://github.com/turbot/steampipe/discussions), where you can ask questions, voice ideas, and share your projects.

Our [Code of Conduct](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) applies to all Steampipe community channels.

### Contributing

Please see [CONTRIBUTING.md](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md).
