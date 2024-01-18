---
title: "Steampipe Table: azure_cosmosdb_sql_database - Query Azure Cosmos DB SQL Databases using SQL"
description: "Allows users to query Azure Cosmos DB SQL Databases, providing data on the properties, configurations, and usage metrics of each database."
---

# Table: azure_cosmosdb_sql_database - Query Azure Cosmos DB SQL Databases using SQL

Azure Cosmos DB is a globally distributed, multi-model database service designed for scalable and high-performance modern applications. It is a fully managed NoSQL database for modern app development with guaranteed single-digit millisecond response times and 99.999-percent availability backed by SLAs, automatic and instant scalability, and open source APIs for MongoDB and Cassandra. A SQL Database in Azure Cosmos DB is a schema-less JSON database engine with SQL querying capabilities.

## Table Usage Guide

The `azure_cosmosdb_sql_database` table provides detailed insights into SQL Databases within Azure Cosmos DB. As a database administrator or developer, you can explore database-specific details through this table, including throughput settings, indexing policies, and associated metadata. Utilize it to monitor database performance, manage configurations, and ensure optimal resource utilization.

## Examples

### Basic info
Explore which Azure CosmosDB SQL databases are tied to specific accounts and regions. This can be helpful in managing resources and understanding the distribution of databases across different regions and accounts.

```sql+postgres
select
  name,
  account_name,
  database_users,
  region,
  resource_group
from
  azure_cosmosdb_sql_database;
```

```sql+sqlite
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
Determine the number of databases associated with each CosmosDB account to better manage resources and plan for scaling needs.

```sql+postgres
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_sql_database
group by
  account_name;
```

```sql+sqlite
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_sql_database
group by
  account_name;
```

### List of sql databases without application tag key
Identify Azure Cosmos DB SQL databases that have not been tagged with an 'application' key. This can be useful in managing and organizing databases, particularly in larger systems where proper tagging can streamline operations and maintenance.

```sql+postgres
select
  name,
  tags
from
  azure_cosmosdb_sql_database
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  azure_cosmosdb_sql_database
where
  json_extract(tags, '$.application') is null;
```