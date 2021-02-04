## v0.2.0 [2021-01-28]

_What's new?_

- New Tables added

  - azure_storage_blob [PR #15](https://github.com/turbot/steampipe-plugin-azure/pull/15)
  - azure_storage_table [PR #12](https://github.com/turbot/steampipe-plugin-azure/pull/12)

- Updates
  - Added columns power_state, private_ips and public_ips to azure_compute_virtual_machine table [PR #17](https://github.com/turbot/steampipe-plugin-azure/pull/17)

_Bug fixes_

- Breaking changes

  - Renamed earlier `azure_storage_blob` table to `azure_storage_blob_service` [PR #7](https://github.com/turbot/steampipe-plugin-azure/pull/7)
  - Renamed earlier `azure_storage_table` table to `azure_storage_table_service` [PR #10](https://github.com/turbot/steampipe-plugin-azure/pull/10)
  - Removed columns managed_disk_storage_account_type and os_disk_size_gb from `azure_compute_virtual_machine` table [PR #17](https://github.com/turbot/steampipe-plugin-azure/pull/17)
