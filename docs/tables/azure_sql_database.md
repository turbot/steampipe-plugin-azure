---
title: "Steampipe Table: azure_sql_database - Query Azure SQL Databases using SQL"
description: "Allows users to query Azure SQL Databases, specifically providing details on database status, server information, and configuration settings."
---

# Table: azure_sql_database - Query Azure SQL Databases using SQL

Azure SQL Database is a fully managed platform as a service (PaaS) Database Engine that handles most of the database management functions such as upgrading, patching, backups, and monitoring without user involvement. It is always running on the latest stable version of the SQL Server database engine and patched OS with 99.99% availability. Azure SQL Database is based on the latest stable version of the Microsoft SQL Server database engine.

## Table Usage Guide

The `azure_sql_database` table provides insights into SQL databases within Microsoft Azure. As a Database Administrator, explore database-specific details through this table, including status, server information, and configuration settings. Utilize it to uncover information about databases, such as their current status, the server they are hosted on, and specific configuration settings.

## Examples

### Basic info
Explore the basic details of your Azure SQL databases such as name, id, server name, location, and edition. This query can be utilized to better understand your SQL database configuration and assess any potential changes or updates that may be necessary.

```sql+postgres
select
  name,
  id,
  server_name,
  location,
  edition
from
  azure_sql_database;
```

```sql+sqlite
select
  name,
  id,
  server_name,
  location,
  edition
from
  azure_sql_database;
```

### List databases that are not online
Discover the segments that consist of databases that are not currently online. This is particularly useful for identifying potential issues and ensuring the smooth functioning of your system.

```sql+postgres
select
  name,
  id,
  server_name,
  location,
  edition,
  status
from
  azure_sql_database
where
  status != 'Online';
```

```sql+sqlite
select
  name,
  id,
  server_name,
  location,
  edition,
  status
from
  azure_sql_database
where
  status != 'Online';
```

### List databases that are not encrypted
Identify instances where Azure SQL databases are unencrypted. This is crucial for assessing potential security vulnerabilities in your database infrastructure.

```sql+postgres
select
  name,
  id,
  server_name,
  location,
  edition,
  transparent_data_encryption ->> 'status' as encryption_status
from
  azure_sql_database
where
  transparent_data_encryption ->> 'status' != 'Enabled';
```

```sql+sqlite
select
  name,
  id,
  server_name,
  location,
  edition,
  json_extract(transparent_data_encryption, '$.status') as encryption_status
from
  azure_sql_database
where
  json_extract(transparent_data_encryption, '$.status') != 'Enabled';
```