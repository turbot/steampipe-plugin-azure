---
title: "Steampipe Table: azure_storage_table_service - Query Azure Storage Table Services using SQL"
description: "Allows users to query Azure Storage Table Services."
---

# Table: azure_storage_table_service - Query Azure Storage Table Services using SQL

The Azure Storage Table service is a NoSQL datastore providing a key-attribute store with a schemaless design. This service allows users to store large amounts of structured data. The service is a non-relational data store that allows for rapid development and fast access to data by scaling as needed.

## Table Usage Guide

The 'azure_storage_table_service' table provides insights into Azure Storage Table Services. As a DevOps engineer, explore specific details about this service through this table, including the storage account name, resource group, and subscription ID. Utilize it to uncover information about the service, such as the CORS (Cross-Origin Resource Sharing) rules, hour metrics, minute metrics, and the retention policy. The schema presents a range of attributes of the Azure Storage Table Service for your analysis, like the storage account ID, CORS rules, hour metrics enabled status, minute metrics enabled status, and the retention policy days.

## Examples

### Basic info
Explore which storage services are being utilized in your Azure environment. This can help in managing resources and optimizing storage allocation across different regions and resource groups.

```sql
select
  name,
  storage_account_name,
  region,
  resource_group
from
  azure_storage_table_service;
```

### CORS rules info of each storage table service
Discover the segments that have specific Cross-Origin Resource Sharing (CORS) rules in each Azure storage table service. This can be useful in understanding the security measures in place for data access and transfer across different origins.

```sql
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