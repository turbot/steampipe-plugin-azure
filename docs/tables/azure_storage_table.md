---
title: "Steampipe Table: azure_storage_table - Query Azure Storage Tables using SQL"
description: "Allows users to query Azure Storage Tables, providing details about each table's properties, including metadata, resource group, and subscription."
folder: "Storage"
---

# Table: azure_storage_table - Query Azure Storage Tables using SQL

Azure Storage Table is a service in Microsoft Azure that stores structured NoSQL data in the cloud, providing a key/attribute store with a schemaless design. Azure Table storage is now part of Azure Cosmos DB. Because Table storage is schemaless, it's easy to adapt your data as the needs of your application evolve.

## Table Usage Guide

The `azure_storage_table` table provides insights into Azure Storage Tables within Microsoft Azure. As a DevOps engineer, explore table-specific details through this table, including metadata, resource group, and subscription. Utilize it to uncover information about tables, such as their properties, the resources they belong to, and the subscriptions they're part of.

## Examples

### Basic info
This query allows you to gain insights into your Azure storage tables, including their names, IDs, associated storage accounts, resource groups, regions, and subscription IDs. This can be particularly useful when assessing the organization and distribution of your storage resources across different Azure subscriptions and regions.

```sql+postgres
select
  name,
  id,
  storage_account_name,
  resource_group,
  region,
  subscription_id
from
  azure_storage_table;
```

```sql+sqlite
select
  name,
  id,
  storage_account_name,
  resource_group,
  region,
  subscription_id
from
  azure_storage_table;
```