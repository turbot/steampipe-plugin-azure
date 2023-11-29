---
title: "Steampipe Table: azure_sql_database - Query Azure SQL Databases using SQL"
description: "Allows users to query Azure SQL Databases."
---

# Table: azure_sql_database - Query Azure SQL Databases using SQL

Azure SQL Database is a managed cloud database provided as part of Microsoft Azure. A high-performance, reliable, and secure database you can use to build data-driven applications and websites in the programming language of your choice, without needing to manage infrastructure.

## Table Usage Guide

The 'azure_sql_database' table provides insights into SQL databases within Azure. As a DevOps engineer, you can explore database-specific details through this table, including server details, collation, status, and associated metadata. Utilize it to uncover information about databases, such as those with specific collation, the status of the databases, and the verification of server details. The schema presents a range of attributes of the SQL database for your analysis, like the database ID, creation date, server name, and associated tags.

## Examples

### Basic info
Explore the general attributes of your Azure SQL databases, such as their names, IDs, server names, locations, and editions. This is useful for gaining a broad overview of your database configurations and locations.

```sql
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
Discover the databases that are not currently online, enabling you to identify potential issues or areas for maintenance within your Azure SQL server. This can be useful for troubleshooting, ensuring optimal performance, and managing resources.

```sql
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
Explore which databases in your Azure SQL server are not encrypted. This can help in identifying potential security risks and ensuring data protection compliance.

```sql
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