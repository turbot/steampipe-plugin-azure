---
title: "Steampipe Table: azure_search_service - Query Azure Search Services using SQL"
description: "Allows users to query Azure Search Services, specifically the details regarding each search service in an Azure subscription. This provides insights into the properties, settings, and status of these services."
folder: "Search Service"
---

# Table: azure_search_service - Query Azure Search Services using SQL

Azure Search Service is a fully managed cloud search service provided by Microsoft Azure. It offers scalable and secure search capabilities across all your data. With Azure Search Service, you can quickly add sophisticated search capabilities to your applications, making it easier for users to find the information they are looking for.

## Table Usage Guide

The `azure_search_service` table provides insights into the Search Services within Microsoft Azure. As a developer or system administrator, you can explore service-specific details through this table, including properties, settings, and status. Utilize it to uncover information about each search service, such as its provisioning state, SKU, and network rules, to manage and optimize your application's search capabilities effectively.

## Examples

### Basic info
Explore the status and configuration of your Azure Search Services to assess resource allocation and utilization. This can help in identifying areas for optimization and managing your resources efficiently.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  status,
  sku_name,
  replica_count
from
  azure_search_service;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  status,
  sku_name,
  replica_count
from
  azure_search_service;
```

### List publicly accessible search services
Determine the areas in which publicly accessible search services are enabled. This is useful in identifying potential security risks and ensuring appropriate access controls are in place.

```sql+postgres
select
  name,
  id,
  public_network_access
from
  azure_search_service
where
  public_network_access = 'Enabled';
```

```sql+sqlite
select
  name,
  id,
  public_network_access
from
  azure_search_service
where
  public_network_access = 'Enabled';
```