---
title: "Steampipe Table: azure_mssql_elasticpool - Query Azure SQL Database Elastic Pools using SQL"
description: "Allows users to query Azure SQL Database Elastic Pools, providing insights into their configuration, performance, and usage statistics."
folder: "SQL Server"
---

# Table: azure_mssql_elasticpool - Query Azure SQL Database Elastic Pools using SQL

Azure SQL Database Elastic Pools are a simple, cost-effective solution for managing and scaling multiple databases that have varying and unpredictable usage demands. They provide a resource model that allows databases to use resources as needed, within certain limits, while also providing a level of isolation from other databases. Azure SQL Database Elastic Pools are particularly useful for SaaS providers who need to manage and scale multiple databases with varying and unpredictable usage.

## Table Usage Guide

The `azure_mssql_elasticpool` table provides insights into Azure SQL Database Elastic Pools within Azure. As a database administrator or DevOps engineer, explore details about each elastic pool, including its configuration, performance metrics, and usage statistics. Utilize it to understand the resource usage and performance of your elastic pools, and to identify potential areas for optimization or scaling.

## Examples

### Basic info
Explore which Microsoft SQL Server elastic pools in your Azure environment are zone redundant and their current state to manage resource allocation effectively. This query is useful for assessing the distribution of Database Transaction Units (DTUs) across your environment.

```sql+postgres
select
  name,
  id,
  state,
  dtu,
  zone_redundant
from
  azure_mssql_elasticpool;
```

```sql+sqlite
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
Explore which elastic pools in Azure SQL are zone redundant. This query is useful for understanding the distribution and resilience of your database resources across different zones.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  state,
  dtu,
  zone_redundant
from
  azure_mssql_elasticpool
where
  zone_redundant = 1;
```