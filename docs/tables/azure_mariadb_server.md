---
title: "Steampipe Table: azure_mariadb_server - Query Azure MariaDB Servers using SQL"
description: "Allows users to query Azure MariaDB Servers."
---

# Table: azure_mariadb_server - Query Azure MariaDB Servers using SQL

Azure Database for MariaDB is a fully managed relational database service provided by Microsoft Azure. It's based on the open-source MariaDB Server engine and allows developers to leverage the capabilities of MariaDB for their applications. The service offers built-in high availability, automatic backups, and scaling of resources in minutes without application downtime.

## Table Usage Guide

The 'azure_mariadb_server' table provides insights into MariaDB servers within Azure Database for MariaDB. As a DevOps engineer, explore server-specific details through this table, including server configurations, performance tiers, and associated metadata. Utilize it to uncover information about servers, such as their performance characteristics, the storage capacity, and the server version. The schema presents a range of attributes of the MariaDB server for your analysis, like the server name, creation date, SKU name, and associated tags.

## Examples

### Basic info
Explore which MariaDB servers in your Azure environment are visible to users. This can help you manage your resources and understand the distribution of your servers across different regions and resource groups.

```sql
select
  name,
  version,
  sku_name,
  user_visible_state,
  region,
  resource_group
from
  azure_mariadb_server;
```

### List servers with Geo-redundant backup disabled
Identify the servers that have their geo-redundant backup feature disabled. This can be useful to ensure all servers are adequately protected and to pinpoint any potential areas of risk.

```sql
select
  name,
  version,
  region,
  geo_redundant_backup_enabled
from
  azure_mariadb_server
where
  geo_redundant_backup_enabled = 'Disabled';
```

### List servers with SSL enabled
Determine the areas in which servers have SSL enabled to enhance security measures within your Azure MariaDB server environment.

```sql
select
  name,
  version,
  region,
  ssl_enforcement
from
  azure_mariadb_server
where
  ssl_enforcement = 'Enabled';
```

### List servers with backup retention days greater than 90 days
Identify Azure MariaDB servers that have a backup retention period of over 90 days. This could be useful in assessing long-term data storage and recovery strategies.

```sql
select
  name,
  version,
  region,
  backup_retention_days
from
  azure_mariadb_server
where
  backup_retention_days > 90;
```