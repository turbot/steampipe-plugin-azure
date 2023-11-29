---
title: "Steampipe Table: azure_mssql_elasticpool - Query Azure SQL Database Elastic Pools using SQL"
description: "Allows users to query Azure SQL Database Elastic Pools"
---

# Table: azure_mssql_elasticpool - Query Azure SQL Database Elastic Pools using SQL

An Azure SQL Database Elastic Pool is a shared resource model that enables higher resource utilization efficiency. It allows for the management of multiple databases that have varying and unpredictable usage demands. The databases within an elastic pool are on a single Azure SQL Database server and share a set number of resources at a set price.

## Table Usage Guide

The 'azure_mssql_elasticpool' table provides insights into Elastic Pools within Azure SQL Database. As a database administrator, explore details specific to each Elastic Pool through this table, including the number of databases, storage limit, and associated metadata. Utilize it to uncover information about each Elastic Pool, such as the maximum and minimum data storage capacity, the number of databases it contains, and its resource usage statistics. The schema presents a range of attributes of the Elastic Pool for your analysis, like the pool's ID, name, type, region, and associated tags.

## Examples

### Basic info
Gain insights into the status and redundancy of your Microsoft SQL Server elastic pools in Azure. This can help you manage resources and ensure your databases are resilient and available.

```sql
select
  name,
  id,
  state,
  dtu,
  zone_redundant
from
  azure_mssql_elasticpool;
```

### List zone redundant elastic pools
Identify the state and capacity of your elastic pools in Azure SQL Database that are configured for zone redundancy. This can help ensure high availability and disaster recovery for your databases.

```sql
select
  name,
  id,
  state,
  dtu,
  zone_redundant
from
  azure_mssql_elasticpool
where
  zone_redundant;
```