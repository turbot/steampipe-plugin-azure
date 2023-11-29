---
title: "Steampipe Table: azure_data_lake_store - Query Azure Data Lake Store using SQL"
description: "Allows users to query Azure Data Lake Stores"
---

# Table: azure_data_lake_store - Query Azure Data Lake Stores using SQL

Azure Data Lake Store is a scalable and secure data lake that allows you to store and analyze large amounts of data. It is built to handle high volumes of small writes at low latency and is optimized for analytics. Azure Data Lake Store supports standard Hadoop Distributed File System (HDFS) interfaces.

## Table Usage Guide

The 'azure_data_lake_store' table provides insights into Data Lake Stores within Azure. As a Data Engineer, explore store-specific details through this table, including encryption settings, firewall rules, and associated metadata. Utilize it to uncover information about stores, such as those with specific firewall rules, the encryption type used, and the verification of virtual network rules. The schema presents a range of attributes of the Data Lake Store for your analysis, like the store name, creation date, encryption settings, and associated tags.

## Examples

### Basic info
Explore which Azure Data Lake stores are currently provisioned to gain insights into your data storage utilization and management. This can help you identify instances where resources may be underutilized or over-provisioned, aiding in efficient resource allocation.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store;
```

### List data lake stores with encryption disabled
This query helps identify Azure Data Lake stores where encryption is disabled, allowing you to pinpoint potential security vulnerabilities and take necessary measures to enhance data protection. It's a practical tool for maintaining the integrity of your stored data and ensuring compliance with data security standards.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  encryption_state = 'Disabled';
```

### List data lake stores with firewall disabled
Explore which Azure Data Lake stores have their firewall disabled. This is crucial for identifying potential security vulnerabilities within your system.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  firewall_state = 'Disabled';
```