---
title: "Steampipe Table: azure_app_service_plan - Query Azure App Service Plans using SQL"
description: "Allows users to query Azure App Service Plans."
---

# Table: azure_app_service_plan - Query Azure App Service Plans using SQL

Azure App Service Plan is a component of Azure App Service, the platform that runs and manages web applications. This service plan allocates the resources your web app will use. It determines the number of VM instances that will be used to run the app and it affects the cost.

## Table Usage Guide

The 'azure_app_service_plan' table provides insights into App Service Plans within Azure App Service. As a DevOps engineer, explore service plan-specific details through this table, including the number of workers, kind of operating system, and associated metadata. Utilize it to uncover information about service plans, such as the maximum number of workers, the reserved status, and the targeted worker size. The schema presents a range of attributes of the App Service Plan for your analysis, like the resource group, kind, status, and associated tags.

## Examples

### App service plan SKU info
Gain insights into the various specifications of your Azure App Service Plan, such as the SKU family, name, size, tier, and capacity. This is useful in understanding the resources allocated to your application, which can help in optimizing performance and cost.

```sql
select
  name,
  sku_family,
  sku_name,
  sku_size,
  sku_tier,
  sku_capacity
from
  azure_app_service_plan;
```


### List of Hyper-V container app service plan
Explore which Azure app service plans are using Hyper-V containers and understand their distribution across different regions. This can be useful for assessing the distribution and usage of Hyper-V containers in your Azure environment.

```sql
select
  name,
  hyper_v,
  kind,
  region
from
  azure_app_service_plan
where
  hyper_v;
```


### List of App service plan that owns spot instances
Explore which Azure App Service plans own spot instances to better manage your resources and costs in different regions and resource groups. This is particularly useful for identifying potential cost savings and optimizing resource allocation.

```sql
select
  name,
  is_spot,
  kind,
  region,
  resource_group
from
  azure_app_service_plan
where
  is_spot;
```