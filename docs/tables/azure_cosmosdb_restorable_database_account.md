---
title: "Steampipe Table: azure_cosmosdb_restorable_database_account - Query Azure Cosmos DB Restorable Database Accounts using SQL"
description: "Allows users to query Azure Cosmos DB Restorable Database Accounts, providing insights into the restorable Azure Cosmos DB accounts within a specified time range."
---

# Table: azure_cosmosdb_restorable_database_account - Query Azure Cosmos DB Restorable Database Accounts using SQL

Azure Cosmos DB is a fully managed NoSQL database service for modern app development. It offers turnkey global distribution, guarantees single-digit millisecond latencies at the 99th percentile, and guarantees high availability with multi-homing capabilities and low latencies anywhere in the world. A Restorable Database Account in Azure Cosmos DB is a resource that can be used to restore the account to a previous state.

## Table Usage Guide

The `azure_cosmosdb_restorable_database_account` table provides insights into restorable Azure Cosmos DB accounts within a specified time range. As a database administrator, explore account-specific details through this table, including the creation time, deletion time, and restorable time range. Utilize it to uncover information about accounts, such as those that are recently deleted, the time range within which the account can be restored, and the verification of restore locations.

## Examples

### Basic Info
Explore which Azure Cosmos DB accounts can be restored, pinpointing their specific locations and the time they were created. This is useful for assessing the elements within your resource group and planning for disaster recovery scenarios.

```sql+postgres
select
  name,
  region,
  account_name,
  creation_time,
  resource_group
from
  azure_cosmosdb_restorable_database_account;
```

```sql+sqlite
select
  name,
  region,
  account_name,
  creation_time,
  resource_group
from
  azure_cosmosdb_restorable_database_account;
```

### Get the regions that the database accounts can be restored from
Explore which regions your database accounts can be restored from, providing useful insights for disaster recovery planning and risk management. This allows you to identify potential fallback locations in case of regional outages or disruptions.

```sql+postgres
select
  name,
  region,
  restorable_locations ->> 'LocationName' as restorable_location,
  restorable_locations ->> 'CreationTime' as regional_database_account_creation_time,
  restorable_locations ->> 'RegionalDatabaseAccountInstanceID' as restorable_location_database_instance_id
from
  azure_cosmosdb_restorable_database_account;
```

```sql+sqlite
select
  name,
  region,
  json_extract(restorable_locations, '$.LocationName') as restorable_location,
  json_extract(restorable_locations, '$.CreationTime') as regional_database_account_creation_time,
  json_extract(restorable_locations, '$.RegionalDatabaseAccountInstanceID') as restorable_location_database_instance_id
from
  azure_cosmosdb_restorable_database_account;
```

### Get the accounts having point-in-time recovery enabled
Discover the Azure CosmosDB accounts that have point-in-time recovery enabled. This is useful for identifying accounts that may require additional backup strategies or have higher potential for data recovery in the event of data loss.

```sql+postgres
select
  ra.account_name,
  ra.name as restorable_database_account_name,
  creation_time,
  ra.id as restorable_database_account_id
from
  azure_cosmosdb_restorable_database_account ra,
  azure_cosmosdb_account a
where
  ra.account_name = a.name
  and ra.subscription_id = a.subscription_id;
```

```sql+sqlite
select
  ra.account_name,
  ra.name as restorable_database_account_name,
  creation_time,
  ra.id as restorable_database_account_id
from
  azure_cosmosdb_restorable_database_account ra
join
  azure_cosmosdb_account a
on
  ra.account_name = a.name
  and ra.subscription_id = a.subscription_id;
```

### Get the restorable account count per api type
Determine the number of restorable accounts for each API type to manage and optimize your Azure Cosmos DB resources. This can be useful for understanding your capacity and planning for potential disaster recovery scenarios.

```sql+postgres
select
  api_type,
  count(ra.*) as accounts
from
  azure_cosmosdb_restorable_database_account ra
group by
  api_type;
```

```sql+sqlite
select
  api_type,
  count(*) as accounts
from
  azure_cosmosdb_restorable_database_account ra
group by
  api_type;
```