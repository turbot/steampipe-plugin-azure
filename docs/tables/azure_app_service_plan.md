---
title: "Steampipe Table: azure_app_service_plan - Query Azure App Service Plans using SQL"
description: "Allows users to query Azure App Service Plans, providing insights into the capacity and scale of the app services running in an Azure subscription."
folder: "App Service"
---

# Table: azure_app_service_plan - Query Azure App Service Plans using SQL

Azure App Service Plan is a service within Microsoft Azure that defines a set of compute resources for a web app to run. These compute resources are analogous to the server farm in conventional web hosting. It specifies the number of VM instances to allocate, the size of each instance, and the pricing tier.

## Table Usage Guide

The `azure_app_service_plan` table provides insights into the App Service Plans within Microsoft Azure. As a Cloud Engineer, explore App Service Plan-specific details through this table, including the number of web apps, capacity, maximum number of workers, and other associated metadata. Utilize it to uncover information about each App Service Plan, such as its current status, tier, and the geographical location of the data center where the plan is running.

## Examples

### App service plan SKU info
Explore the details of your Azure App Service Plan to understand the specifics of your service tier and capacity. This can help you assess if your current plan aligns with your application's requirements and if there is a need for scaling or downgrading.

```sql+postgres
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

```sql+sqlite
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
Explore which Azure App Service Plans are using Hyper-V containers. This can help determine the areas in which these specific types of containers are being utilized, aiding in resource management and optimization.

```sql+postgres
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

```sql+sqlite
select
  name,
  hyper_v,
  kind,
  region
from
  azure_app_service_plan
where
  hyper_v = 1;
```

### List of App service plan that owns spot instances
Explore which Azure App Service Plans are utilizing spot instances. This is useful for managing costs and understanding the distribution of your resources.

```sql+postgres
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

```sql+sqlite
select
  name,
  is_spot,
  kind,
  region,
  resource_group
from
  azure_app_service_plan
where
  is_spot = 1;
```