---
title: "Steampipe Table: azure_mariadb_server - Query Azure MariaDB Servers using SQL"
description: "Allows users to query Azure MariaDB Servers, offering insights into the configuration and status of these managed database service instances."
folder: "MariaDB"
---

# Table: azure_mariadb_server - Query Azure MariaDB Servers using SQL

Azure MariaDB Server is a fully managed relational database service provided by Microsoft Azure. It is a scalable and flexible service that allows users to deploy highly available MariaDB databases in the cloud. Azure MariaDB Server provides automatic backups, patching, monitoring, and scaling of resources to ensure optimal performance and reliability.

## Table Usage Guide

The `azure_mariadb_server` table offers insights into the configuration and status of Azure MariaDB Server instances. As a database administrator, you can utilize this table to monitor and manage your MariaDB instances, including their performance, security settings, and backup configurations. This table is also beneficial for auditing purposes, allowing you to track changes and maintain compliance with organizational policies and standards.

## Examples

### Basic info
Explore the overall status and location of your Azure MariaDB servers. This query is useful in gaining insights into the versions, pricing tiers, visibility, and geographical distribution of your databases, aiding in resource management and cost optimization.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have disabled geo-redundant backup, a feature essential for data protection and disaster recovery, on their servers. This assists in identifying potential vulnerabilities in the system and aids in enhancing data security.

```sql+postgres
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

```sql+sqlite
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
Identify instances where your Azure MariaDB servers have SSL enabled. This is useful for ensuring that your data transmissions are secure and encrypted.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which servers on Azure MariaDB are configured to retain backups for more than 90 days. This can be useful for identifying servers with potentially excessive storage use or for compliance purposes.

```sql+postgres
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

```sql+sqlite
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