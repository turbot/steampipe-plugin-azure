---
title: "Steampipe Table: azure_app_service_environment - Query Azure App Service Environments using SQL"
description: "Allows users to query Azure App Service Environments, providing comprehensive details about the app service environments in your Azure account."
folder: "App Service"
---

# Table: azure_app_service_environment - Query Azure App Service Environments using SQL

Azure App Service Environment is a fully isolated and dedicated environment for securely running App Service apps at high scale. This service is designed for application workloads that require high scale, isolation, and secure network access. It provides a fully isolated and dedicated environment for running applications of almost any scale.

## Table Usage Guide

The `azure_app_service_environment` table provides insights into App Service Environments within Azure. As a DevOps engineer, you can explore specific details about these environments, including the number of workers, the status of the environment, and the virtual network integration. This table is useful for understanding the scale, isolation, and security of your app services, and to identify any potential issues or areas for optimization.

## Examples

### List of app service environments which are not healthy
Uncover the details of Azure App Service environments that are currently not in a healthy state. This can be useful for identifying potential issues that may be affecting the performance or availability of your applications.

```sql+postgres
select
  name,
  is_healthy_environment
from
  azure_app_service_environment
where
  not is_healthy_environment;
```

```sql+sqlite
select
  name,
  is_healthy_environment
from
  azure_app_service_environment
where
  not is_healthy_environment;
```

### Virtual network info of each app service environment
Explore the virtual network configurations of each app service environment to gain insights into the internal load balancing mode and understand the network segregation. This is useful in assessing the security and isolation measures within your Azure App Service Environment.

```sql+postgres
select
  name,
  vnet_name,
  vnet_subnet_name,
  vnet_resource_group_name,
  internal_load_balancing_mode
from
  azure_app_service_environment;
```

```sql+sqlite
select
  name,
  vnet_name,
  vnet_subnet_name,
  vnet_resource_group_name,
  internal_load_balancing_mode
from
  azure_app_service_environment;
```

### List cluster settings details
Explore the configuration details of your Azure App Service Environment to understand the specifics of your cluster settings. This can aid in managing your resources more effectively and troubleshooting potential issues.

```sql+postgres
select
  name,
  id,
  settings ->> 'name' as settings_name,
  settings ->> 'value' as settings_value
from
  azure_app_service_environment,
  jsonb_array_elements(cluster_settings) as settings;
```

```sql+sqlite
select
  name,
  id,
  json_extract(settings.value, '$.name') as settings_name,
  json_extract(settings.value, '$.value') as settings_value
from
  azure_app_service_environment,
  json_each(cluster_settings) as settings;
```