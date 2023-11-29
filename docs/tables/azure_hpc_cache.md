---
title: "Steampipe Table: azure_hpc_cache - Query Azure Storage Cache using SQL"
description: "Allows users to query Azure Storage Caches"
---

# Table: azure_hpc_cache - Query Azure Storage Cache using SQL

Azure HPC Cache is a service that provides low-latency, high-throughput access to data located in Azure Blob storage. It creates a caching layer between compute clusters and storage to help you run more jobs, more iterations, and get results faster. It is designed to support high-performance computing (HPC) scenarios where data is read from and written to Azure Blob storage.

## Table Usage Guide

The 'azure_hpc_cache' table provides insights into the Azure HPC Cache within Azure Storage. As a DevOps engineer, explore cache-specific details through this table, including cache size, health, provisioning state, and associated metadata. Utilize it to uncover information about caches, such as their network settings, subnet ID, and usage model. The schema presents a range of attributes of the Azure HPC Cache for your analysis, like the cache ID, creation time, health, provisioning state, and associated tags.

## Examples

### Basic info
Explore which High Performance Computing (HPC) caches are currently active in your Azure environment and understand their types and provisioning states. This can help in assessing their performance and managing resources efficiently.

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
Explore the network settings of your Azure HPC Cache to gain insights into configurations such as DNS search domain, MTU, NTP server, DNS servers, and utility addresses. This can help you understand and manage your network's performance, security, and reliability.

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
Explore the encryption settings of your Azure HPC Cache to understand the configuration of your key encryption and network settings. This can be useful for maintaining security standards and ensuring proper data protection.

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