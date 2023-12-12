---
title: "Steampipe Table: azure_cosmosdb_account - Query Azure Cosmos DB Accounts using SQL"
description: "Allows users to query Azure Cosmos DB Accounts, providing detailed information about each account's configurations, properties, and associated resources."
---

# Table: azure_cosmosdb_account - Query Azure Cosmos DB Accounts using SQL

Azure Cosmos DB is a globally distributed, multi-model database service designed for scalable and high-performance modern applications. It provides native support for NoSQL and OSS APIs, including MongoDB, Cassandra, Gremlin, et al. With turnkey global distribution and transparent multi-master replication, it offers single-digit millisecond latency, and 99.999% availability.

## Table Usage Guide

The `azure_cosmosdb_account` table provides insights into Azure Cosmos DB Accounts within Azure's database services. As a database administrator or developer, explore account-specific details through this table, including configurations, properties, and associated resources. Utilize it to uncover information about accounts, such as their replication policies, failover policies, and the verification of virtual network rules.

## Examples

### List of database accounts where automatic failover is not enabled
Explore which database accounts in Azure CosmosDB do not have automatic failover enabled. This is useful to identify potential risks and ensure high availability and disaster recovery in your database setup.

```sql+postgres
select
  name,
  region,
  enable_automatic_failover,
  resource_group
from
  azure_cosmosdb_account
where
  not enable_automatic_failover;
```

```sql+sqlite
select
  name,
  region,
  enable_automatic_failover,
  resource_group
from
  azure_cosmosdb_account
where
  enable_automatic_failover = 0;
```

### List of database accounts which allows traffic from all networks, including the public Internet.
Explore which database accounts are potentially exposed to security risks by allowing traffic from all networks, including the public internet. This can be useful to identify potential vulnerabilities and improve security measures.

```sql+postgres
select
  name,
  region,
  virtual_network_rules
from
  azure_cosmosdb_account
where
  virtual_network_rules = '[]';
```

```sql+sqlite
select
  name,
  region,
  virtual_network_rules
from
  azure_cosmosdb_account
where
  virtual_network_rules = '[]';
```

### List of database accounts where multiple write location is not enabled
Identify the Azure Cosmos DB accounts that do not have multiple write locations enabled. This can be useful in pinpointing potential areas of risk or inefficiency, as enabling multiple write locations can enhance data redundancy and availability.

```sql+postgres
select
  name,
  region,
  enable_multiple_write_locations
from
  azure_cosmosdb_account
where
  not enable_multiple_write_locations;
```

```sql+sqlite
select
  name,
  region,
  enable_multiple_write_locations
from
  azure_cosmosdb_account
where
  enable_multiple_write_locations = 0;
```

### Failover policy info for the database accounts
Determine the areas in which your Azure CosmosDB accounts have their failover policies set. This helps in understanding the priority and location of failover events, thereby assisting in ensuring high availability and disaster recovery strategies.

```sql+postgres
select
  name,
  fp ->> 'failoverPriority' as failover_priority,
  fp ->> 'locationName' as location_name
from
  azure_cosmosdb_account
  cross join jsonb_array_elements(failover_policies) as fp;
```

```sql+sqlite
select
  name,
  json_extract(fp.value, '$.failoverPriority') as failover_priority,
  json_extract(fp.value, '$.locationName') as location_name
from
  azure_cosmosdb_account,
  json_each(failover_policies) as fp;
```

### Consistency policy info for each account
Discover the segments that detail the consistency policy for each account, useful for understanding the database account offer type and the default consistency level. This aids in managing data consistency and staleness across your Azure Cosmos DB accounts.

```sql+postgres
select
  name,
  consistency_policy_max_interval,
  consistency_policy_max_staleness_prefix,
  database_account_offer_type,
  default_consistency_level
from
  azure_cosmosdb_account;
```

```sql+sqlite
select
  name,
  consistency_policy_max_interval,
  consistency_policy_max_staleness_prefix,
  database_account_offer_type,
  default_consistency_level
from
  azure_cosmosdb_account;
```

### Get backup policy for accounts having periodic backups enabled
Determine the backup policies of accounts with periodic backups enabled. This is useful for understanding the frequency and retention of backups, as well as the redundancy of storage, ensuring data safety and availability.

```sql+postgres
select
  name,
  region,
  backup_policy -> 'periodicModeProperties' ->> 'backupIntervalInMinutes' as backup_interval_mins,
  backup_policy -> 'periodicModeProperties' ->> 'backupRetentionIntervalInHours' as backup_retention_interval_hrs,
  backup_policy -> 'periodicModeProperties' ->> 'backupStorageRedundancy' as backup_storage_redundancy
from
  azure_cosmosdb_account
where
  backup_policy ->> 'type' = 'Periodic';
```

```sql+sqlite
select
  name,
  region,
  json_extract(backup_policy, '$.periodicModeProperties.backupIntervalInMinutes') as backup_interval_mins,
  json_extract(backup_policy, '$.periodicModeProperties.backupRetentionIntervalInHours') as backup_retention_interval_hrs,
  json_extract(backup_policy, '$.periodicModeProperties.backupStorageRedundancy') as backup_storage_redundancy
from
  azure_cosmosdb_account
where
  json_extract(backup_policy, '$.type') = 'Periodic';
```

### Get private endpoint connection details for each account
Explore the connection details of each private endpoint linked to your account. This can help you assess the status and type of each connection, enabling better management and troubleshooting of your network resources.

```sql+postgres
select
  c ->> 'PrivateEndpointConnectionName' as private_endpoint_connection_name,
  c ->> 'PrivateEndpointConnectionType' as private_endpoint_connection_type,
  c ->> 'PrivateEndpointId' as private_endpoint_id,
  c ->> 'PrivateLinkServiceConnectionStateActionsRequired' as private_link_service_connection_state_actions_required,
  c ->> 'PrivateLinkServiceConnectionStateDescription' as private_link_service_connection_state_description,
  c ->> 'PrivateLinkServiceConnectionStateStatus' as private_link_service_connection_state_status,
  c ->> 'ProvisioningState' as provisioning_state,
  c ->> 'PrivateEndpointConnectionId' as private_endpoint_connection_id
from
  azure_cosmosdb_account,
  jsonb_array_elements(private_endpoint_connections) as c;
```

```sql+sqlite
select
  json_extract(c.value, '$.PrivateEndpointConnectionName') as private_endpoint_connection_name,
  json_extract(c.value, '$.PrivateEndpointConnectionType') as private_endpoint_connection_type,
  json_extract(c.value, '$.PrivateEndpointId') as private_endpoint_id,
  json_extract(c.value, '$.PrivateLinkServiceConnectionStateActionsRequired') as private_link_service_connection_state_actions_required,
  json_extract(c.value, '$.PrivateLinkServiceConnectionStateDescription') as private_link_service_connection_state_description,
  json_extract(c.value, '$.PrivateLinkServiceConnectionStateStatus') as private_link_service_connection_state_status,
  json_extract(c.value, '$.ProvisioningState') as provisioning_state,
  json_extract(c.value, '$.PrivateEndpointConnectionId') as private_endpoint_connection_id
from
  azure_cosmosdb_account,
  json_each(private_endpoint_connections) as c;
```

### Get details of accounts restored from backup
The example demonstrates how to identify the instances where Azure Cosmos DB accounts have been restored from a backup. This can be particularly useful for auditing purposes, to ensure data integrity and to track any unauthorized restorations.

```sql+postgres
select
  name,
  restore_parameters ->> 'restoreMode' as restore_mode,
  restore_parameters ->> 'restoreSource' as restore_source,
  d ->> 'databaseName' as restored_database_name,
  c as restored_collection_name
from
  azure_cosmosdb_account,
  jsonb_array_elements(restore_parameters -> 'databasesToRestore') d,
  jsonb_array_elements_text(d -> 'collectionNames') c;
```

```sql+sqlite
select
  a.name,
  json_extract(a.restore_parameters, '$.restoreMode') as restore_mode,
  json_extract(a.restore_parameters, '$.restoreSource') as restore_source,
  json_extract(d.value, '$.databaseName') as restored_database_name,
  json_extract(c.value, '$') as restored_collection_name
from
  azure_cosmosdb_account a,
  json_each(json_extract(a.restore_parameters, '$.databasesToRestore')) as d,
  json_each(json_extract(d.value, '$.collectionNames')) as c;
```