---
title: "Steampipe Table: azure_hpc_cache - Query Azure HPC Cache using SQL"
description: "Allows users to query Azure HPC Cache, specifically the operational details and status of the cache. This can provide insights into cache utilization, performance, and potential issues."
folder: "HPC Cache"
---

# Table: azure_hpc_cache - Query Azure HPC Cache using SQL

Azure HPC Cache is a service within Microsoft Azure that accelerates access to data in Azure Blob Storage for high-performance computing (HPC) applications. It provides a caching layer that allows HPC applications to access data as if it were local, improving performance and reducing latency. Azure HPC Cache is beneficial for workloads that require high-speed access to large datasets, such as genomics, financial risk modeling, and simulation.

## Table Usage Guide

The `azure_hpc_cache` table provides insights into Azure HPC Cache within Azure Storage. As a Data Engineer, explore cache-specific details through this table, including operational details, status, and performance metrics. Utilize it to uncover information about cache utilization, identify potential performance bottlenecks, and monitor cache status for potential issues.

## Examples

### Basic info
Explore which Azure HPC Cache instances are currently deployed in your environment. This is beneficial in understanding the overall usage and configuration of your Azure HPC Cache resources.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  sku_name
from
  azure_hpc_cache;
```

```sql+sqlite
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
This query is used to analyze the network settings for Azure's high-performance cache service. It can help users understand and manage the network configurations for their cache services, including DNS search domains, MTU settings, NTP servers, DNS servers, and utility addresses.

```sql+postgres
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

```sql+sqlite
select
  id,
  name,
  json_extract(network_settings, '$.DNSSearchDomain') as dns_search_domain,
  json_extract(network_settings, '$.Mtu') as mtu,
  json_extract(network_settings, '$.NtpServer') as ntp_server,
  network_settings as dns_servers,
  network_settings as utility_addresses
from
  azure_hpc_cache;
```

### List encryption settings details
Explore the encryption details of your Azure HPC cache to understand its security settings. This can be useful for assessing the security status of your data and ensuring it meets your organization's requirements.

```sql+postgres
select
  id,
  name,
  encryption_settings -> 'keyEncryptionKey' ->> 'keyUrl'  as key_url,
  encryption_settings -> 'keyEncryptionKey' -> 'sourceVault' ->> 'id'  as source_vault_id,
  network_settings -> 'rotationToLatestKeyVersionEnabled' as rotation_to_latest_key_version_enabled
from
  azure_hpc_cache;
```

```sql+sqlite
select
  id,
  name,
  json_extract(json_extract(encryption_settings, '$.keyEncryptionKey'), '$.keyUrl') as key_url,
  json_extract(json_extract(json_extract(encryption_settings, '$.keyEncryptionKey'), '$.sourceVault'), '$.id') as source_vault_id,
  json_extract(network_settings, '$.rotationToLatestKeyVersionEnabled') as rotation_to_latest_key_version_enabled
from
  azure_hpc_cache;
```