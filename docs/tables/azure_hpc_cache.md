# Table: azure_hpc_cache

Azure HPC Cache speeds access to your data for high-performance computing (HPC) tasks. By caching files in Azure, Azure HPC Cache brings the scalability of cloud computing to your existing workflow. This service can be used even for workflows where your data is stored across WAN links, such as in your local datacenter network-attached storage (NAS) environment.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  sku_name
from
  azure_hpc_cache;
```

### List network settings details

```sql
select
  id,
  name,
  network_settings ->> 'DNSSearchDomain' as dns_search_domain,
  network_settings -> 'Mtu' as mtu,
  network_settings ->> 'NtpServer' as ntp_server,
  jsonb_pretty(network_settings -> 'DNSServers') as dns_servers,
  jsonb_pretty(network_settings -> 'UtilityAddresses') as utility_addresses
from
  azure_hpc_cache;
```

### List encryption settings details

```sql
select
  id,
  name,
  encryption_settings -> 'keyEncryptionKey' ->> 'keyUrl'  as key_url,
  encryption_settings -> 'keyEncryptionKey' -> 'sourceVault' ->> 'id'  as source_vault_id,
  network_settings -> 'rotationToLatestKeyVersionEnabled' as rotation_to_latest_key_version_enabled
from
  azure_hpc_cache;
```
