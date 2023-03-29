# Table: azure_cosmosdb_account

Azure Cosmos DB is a fully managed NoSQL database service for modern app development.

## Examples

### List of database accounts where automatic failover is not enabled

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