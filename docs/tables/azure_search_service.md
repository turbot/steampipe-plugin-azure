---
title: "Steampipe Table: azure_search_service - Query Azure Search Services using SQL"
description: "Allows users to query Azure Search Services."
---

# Table: azure_search_service - Query Azure Search Services using SQL

Azure Search Service is a cloud-based search-as-a-service solution that delegates server and infrastructure management to Microsoft, leaving you with a ready-to-use service that you can populate with your data and then use to add search to your web or mobile application. Azure Search Service supports a wide variety of features to provide a rich search experience, including full-text search, filters and facets, typeaheads, hit highlighting, and suggestions. It also provides capabilities for tuning the relevance of search results and offers a simple query syntax for a wide range of query types.

## Table Usage Guide

The 'azure_search_service' table provides insights into Search Services within Azure. As a DevOps engineer, explore service-specific details through this table, including the service name, resource group, subscription ID, and associated metadata. Utilize it to uncover information about search services, such as the service tier, the number of replicas and partitions, and the verification of public network access. The schema presents a range of attributes of the Search Service for your analysis, like the service name, location, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which Azure Search services are currently active and assess their configuration, including the number of replicas. This is useful for managing resources and understanding the scale of your Azure Search services.

```sql
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
Discover the segments that have publicly accessible search services enabled within the Azure platform. This is useful for assessing potential security risks and ensuring appropriate access controls are in place.

```sql
select
  name,
  id,
  public_network_access
from
  azure_search_service
where
  public_network_access = 'Enabled';
```