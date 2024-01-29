---
title: "Steampipe Table: azure_storage_table_service - Query Azure Storage Table Services using SQL"
description: "Allows users to query Azure Storage Table Services, specifically the details of a table service within a storage account, providing insights into its properties and settings."
---

# Table: azure_storage_table_service - Query Azure Storage Table Services using SQL

Azure Storage Table Service is a NoSQL datastore providing a key/attribute store with a schemaless design. This service allows you to store large amounts of structured data, providing a flexible schema for data. Azure Storage Table Services are ideal for storing structured, non-relational data.

## Table Usage Guide

The `azure_storage_table_service` table provides insights into Azure Storage Table Services within Microsoft Azure. As a Data Engineer or Developer, you can explore service-specific details through this table, including properties, settings, and associated metadata. Utilize it to uncover information about table services, such as their properties, the storage account they belong to, and the configuration settings applied to them.

## Examples

### Basic info
Gain insights into the association between storage account names and their corresponding regions and resource groups. This information can be useful for managing resources and understanding the distribution of storage accounts across different regions and groups.

```sql+postgres
select
  name,
  storage_account_name,
  region,
  resource_group
from
  azure_storage_table_service;
```

```sql+sqlite
select
  name,
  storage_account_name,
  region,
  resource_group
from
  azure_storage_table_service;
```

### CORS rules info of each storage table service
Explore the Cross-Origin Resource Sharing (CORS) rules of your Azure Storage Table services. This query helps you understand the CORS configurations in place, including allowed headers, methods, origins, exposed headers, and the maximum age in seconds, providing insights into how your resources interact with requests from different origins.

```sql+postgres
select
  name,
  storage_account_name,
  cors -> 'allowedHeaders' as allowed_headers,
  cors -> 'allowedMethods' as allowed_methods,
  cors -> 'allowedOrigins' as allowed_origins,
  cors -> 'exposedHeaders' as exposed_eaders,
  cors -> 'maxAgeInSeconds' as max_age_in_seconds
from
  azure_storage_table_service,
  jsonb_array_elements(cors_rules) as cors;
```

```sql+sqlite
select
  name,
  storage_account_name,
  json_extract(cors.value, '$.allowedHeaders') as allowed_headers,
  json_extract(cors.value, '$.allowedMethods') as allowed_methods,
  json_extract(cors.value, '$.allowedOrigins') as allowed_origins,
  json_extract(cors.value, '$.exposedHeaders') as exposed_headers,
  json_extract(cors.value, '$.maxAgeInSeconds') as max_age_in_seconds
from
  azure_storage_table_service,
  json_each(cors_rules) as cors;
```