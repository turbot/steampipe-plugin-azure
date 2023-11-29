---
title: "Steampipe Table: azure_cosmosdb_account - Query Azure Cosmos DB Accounts using SQL"
description: "Allows users to query Azure Cosmos DB Accounts to obtain key information such as account name, resource group, location, and more. The table provides a comprehensive view of these resources, including the account's offer type, IP rules, and virtual network rules."
---

# Table: azure_cosmosdb_account - Query Azure Cosmos DB Accounts using SQL

Azure Cosmos DB is a fully managed NoSQL database service for modern app development with guaranteed single-digit millisecond response times and 99.999-percent availability backed by SLAs, automatic and instant scalability, and open source APIs for MongoDB and Cassandra. It offers multi-mastering feature by automatically indexing all data and allowing massively parallel operations. Azure Cosmos DB provides native support for NoSQL and OSS APIs, including MongoDB, Cassandra, Gremlin, et al.

## Table Usage Guide

The 'azure_cosmosdb_account' table provides insights into Azure Cosmos DB Accounts. As a DevOps engineer, explore account-specific details through this table, including the account's offer type, IP rules, and virtual network rules. Utilize it to uncover information about accounts, such as their locations, enabled capabilities, and associated tags. The schema presents a range of attributes of the Azure Cosmos DB Account for your analysis, like the account name, resource group, read and write locations, and more.

## Examples

### List of database accounts where automatic failover is not enabled
Discover the segments that have automatic failover disabled in their database accounts, which can be critical in maintaining seamless service during unexpected outages. This could be useful in identifying potential vulnerabilities in your database setup.

```sql
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

### List of database accounts which allows traffic from all networks, including the public Internet.
Determine the areas in which database accounts are potentially vulnerable by identifying those that allow traffic from all networks, including the public internet. This can help in enhancing security by restricting access to specific networks.

```sql
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
Discover the segments that have not enabled multiple write locations within their Azure CosmosDB accounts. This can be useful in identifying potential areas of risk or inefficiency, as enabling multiple write locations can increase data redundancy and availability.

```sql
select
  name,
  region,
  enable_multiple_write_locations
from
  azure_cosmosdb_account
where
  not enable_multiple_write_locations;
```

### Failover policy info for the database accounts
Gain insights into the priority and location details of failover policies for your Azure CosmosDB accounts. This helps in strategizing disaster recovery and business continuity plans.

```sql
select
  name,
  fp ->> 'failoverPriority' as failover_priority,
  fp ->> 'locationName' as location_name
from
  azure_cosmosdb_account
  cross join jsonb_array_elements(failover_policies) as fp;
```

### Consistency policy info for each account
Analyze the consistency policy details of each account to understand the maximum interval, staleness prefix, account offer type, and default consistency level. This aids in optimizing data consistency and performance in Azure Cosmos DB accounts.

```sql
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
Determine the areas in which Azure CosmosDB accounts have periodic backups enabled to assess their backup policies. This is useful for understanding the frequency of backups and the retention period, ensuring data safety and compliance with data retention policies.

```sql
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

### Get private endpoint connection details for each account
Explore the status and details of private endpoint connections for each account to understand the connection type, actions required, and current state. This is useful for managing and troubleshooting your private network connections in Azure Cosmos DB.

```sql
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

### Get details of accounts restored from backup
Identify instances where Azure CosmosDB accounts have been restored from backup. This is useful to track restoration activities and ensure data integrity.

```sql
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