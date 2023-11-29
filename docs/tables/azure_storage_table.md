---
title: "Steampipe Table: azure_storage_table - Query Azure Storage Tables using SQL"
description: "Allows users to query Azure Storage Tables."
---

# Table: azure_storage_table - Query Azure Storage Tables using SQL

Azure Storage Tables are a service that stores structured NoSQL data in the cloud, providing a key/attribute store with a schema-less design. Because Table storage is schema-less, it's easy to adapt your data as the needs of your application evolve. Azure Table storage is now part of Azure Cosmos DB.

## Table Usage Guide

The 'azure_storage_table' table provides insights into Azure Storage Tables within Azure Storage Account service. As a DevOps engineer, explore table-specific details through this table, including the table name, resource group, and associated metadata. Utilize it to uncover information about tables, such as those with specific table names, the resource groups associated with the tables, and the region of storage. The schema presents a range of attributes of the Azure Storage Table for your analysis, like the table name, resource group, and region.

## Examples

### Basic info
Explore which Azure storage tables are currently in use across different regions and subscriptions. This can help manage resources more effectively by identifying where storage is allocated.

```sql
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