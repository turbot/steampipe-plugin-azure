---
title: "Steampipe Table: azure_cosmosdb_sql_database - Query Azure Cosmos DB SQL Databases using SQL"
description: "Allows users to query Azure Cosmos DB SQL Databases"
---

# Table: azure_cosmosdb_sql_database - Query Azure Cosmos DB SQL Databases using SQL

Azure Cosmos DB is a globally distributed, multi-model database service for managing data at large scale. It provides elastic scalability, high availability, and low latency required for modern applications. SQL API, one of the APIs provided by Azure Cosmos DB, allows you to work with data using SQL queries.

## Table Usage Guide

The 'azure_cosmosdb_sql_database' table provides insights into SQL Databases within Azure Cosmos DB. As a database administrator, explore database-specific details through this table, including the provisioned throughput, partition key path, and associated metadata. Utilize it to uncover information about databases, such as those with high throughput, the partitioning scheme, and the indexing policy. The schema presents a range of attributes of the SQL Database for your analysis, like the database ID, resource group, and associated tags.

## Examples

### Basic info
Explore the configuration of your Azure CosmosDB SQL databases to gain insights into their associated account names, user databases, regions, and resource groups. This can help you manage your resources more effectively and understand where potential issues may arise.

```sql
select
  name,
  account_name,
  database_users,
  region,
  resource_group
from
  azure_cosmosdb_sql_database;
```


### Database count per cosmosdb accounts
Identify the number of databases within each Azure Cosmos DB account. This information can be useful for managing resources and understanding the distribution of databases across different accounts.

```sql
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_sql_database
group by
  account_name;
```


### List of sql databases without application tag key
Identify Azure Cosmos DB SQL databases that are missing an 'application' tag. This can be useful in scenarios where you want to ensure all databases are properly tagged for better management and organization.

```sql
select
  name,
  tags
from
  azure_cosmosdb_sql_database
where
  not tags :: JSONB ? 'application';
```