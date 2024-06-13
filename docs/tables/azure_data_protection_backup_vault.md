---
title: "Steampipe Table: azure_data_protection_backup_vault - Query Azure Data Protection Backup Vaults using SQL"
description: "Allows users to query Data Protection Backup Vaults in Azure, providing detailed information about each backup vault, including its location, identity, provisioning state, and storage settings."
---

# Table: azure_data_protection_backup_vault - Query Azure Data Protection Backup Vaults using SQL

A Backup Vault in Azure Data Protection is a storage entity in Azure used to manage and store backups. Backup Vaults help in organizing your backups and configuring backup policies for your Azure resources.

## Table Usage Guide

The `azure_data_protection_backup_vault` table provides insights into Backup Vaults within Azure. As an Infrastructure Engineer, explore detailed information about each backup vault through this table, including its identity, provisioning state, and storage settings. Use this table to manage and optimize your backup strategies, ensuring secure and efficient data protection for your Azure resources.

## Examples

### Basic backup vault information
Retrieve basic information about your Azure Backup Vaults, including their names, locations, and provisioning states.

```sql+postgres
select
  name,
  location,
  provisioning_state
from
  azure_data_protection_backup_vault;
```

```sql+sqlite
select
  name,
  location,
  provisioning_state
from
  azure_data_protection_backup_vault;
```

### Backup vault storage settings
Explore the storage settings of your Azure Backup Vaults. This can help you understand the storage configuration and optimize your backup storage strategies.

```sql+postgres
select
  name,
  s -> 'datastoreType' as datastore_type,
  s -> 'type' as storage_type
from
  azure_data_protection_backup_vault,
  jsonb_array_elements(storage_settings) as s;
```

```sql+sqlite
select
  name,
  json_extract(s.value, '$.datastoreType') as datastore_type,
  json_extract(s.value, '$.type') as storage_type
from
  azure_data_protection_backup_vault,
  json_each(storage_settings) as s;
```

### Managed identity details
List the managed identity details for each Backup Vault, which can be useful for managing access and security configurations.

```sql+postgres
select
  name,
  identity ->> 'type' as identity_type,
  identity -> 'principalId' as principal_id
from
  azure_data_protection_backup_vault;
```

```sql+sqlite
select
  name,
  json_extract(identity, '$.type') as identity_type,
  json_extract(identity, '$.principalId') as principal_id
from
  azure_data_protection_backup_vault;
```

### Backup vault tags
Retrieve tags associated with each Backup Vault to help with resource organization and management.

```sql+postgres
select
  name,
  jsonb_each_text(tags) as tag
from
  azure_data_protection_backup_vault;
```

```sql+sqlite
select
  name,
  json_each(tags) as tag
from
  azure_data_protection_backup_vault;
```

### Backup vaults with failed provisioning state
Retrieve information about backup vaults that are in a failed provisioning state. This can help in identifying and troubleshooting issues with backup vault deployment.

```sql+postgres
select
  name,
  s ->> 'datastoreType' as datastore_type,
  s ->> 'type' as storage_type
from
  azure_data_protection_backup_vault,
  jsonb_array_elements(storage_settings) as s
where
  provisioning_state = 'Failed';
```

```sql+sqlite
select
  name,
  json_extract(s.value, '$.datastoreType') as datastore_type,
  json_extract(s.value, '$.type') as storage_type
from
  azure_data_protection_backup_vault,
  json_each(storage_settings) as s
where
  provisioning_state = 'Failed';
```
