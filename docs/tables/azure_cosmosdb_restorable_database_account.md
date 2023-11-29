---
title: "Steampipe Table: azure_cosmosdb_restorable_database_account - Query Azure Cosmos DB Restorable Database Accounts using SQL"
description: "Allows users to query Azure Cosmos DB Restorable Database Accounts"
---

# Table: azure_cosmosdb_restorable_database_account - Query Azure Cosmos DB Restorable Database Accounts using SQL

Azure Cosmos DB is a globally distributed, multi-model database service designed for scalable and high performance modern applications. It is a fully managed NoSQL database service built for fast and predictable performance, high availability, elastic scaling, global distribution, and ease of development. A restorable database account represents a Cosmos DB account that can be restored to any point in time within its retention period.

## Table Usage Guide

The 'azure_cosmosdb_restorable_database_account' table provides insights into restorable database accounts within Azure Cosmos DB. As a DevOps engineer, explore account-specific details through this table, including locations, enabled capabilities, and associated metadata. Utilize it to uncover information about accounts, such as those with specific capabilities, the locations of accounts, and the verification of failover policies. The schema presents a range of attributes of the restorable database account for your analysis, like the account name, creation date, enabled capabilities, and associated tags.

## Examples

### Basic Info
Explore which Azure Cosmos DB accounts are available for restoration, along with their associated details such as region, account name, and creation time. This is particularly useful for assessing recovery options and planning for potential disaster recovery scenarios.

```sql
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
Explore which regions your database accounts can be restored from to ensure business continuity and disaster recovery. This query aids in identifying the locations where your database backups are stored, helping you plan your restoration strategy effectively.

```sql
select
  name,
  region,
  restorable_locations ->> 'LocationName' as restorable_location,
  restorable_locations ->> 'CreationTime' as regional_database_account_creation_time,
  restorable_locations ->> 'RegionalDatabaseAccountInstanceID' as restorable_location_database_instance_id
from
  azure_cosmosdb_restorable_database_account;
```

### Get the accounts having point-in-time recovery enabled
Discover the segments that have point-in-time recovery enabled for Azure CosmosDB accounts. This query can be useful in instances where you need to analyze the safety measures of your data, ensuring that it can be restored to a specific point in time if needed.

```sql
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

### Get the restorable account count per api type
Analyze the settings to understand the distribution of restorable accounts across different API types in Azure CosmosDB. This can be beneficial for assessing the balance of your account types and identifying any potential vulnerabilities or over-reliances.

```sql
select
  api_type,
  count(ra.*) as accounts
from
  azure_cosmosdb_restorable_database_account ra
group by
  api_type;
```