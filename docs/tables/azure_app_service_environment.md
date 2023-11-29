---
title: "Steampipe Table: azure_app_service_environment - Query Azure App Service Environments using SQL"
description: "Allows users to query Azure App Service Environments"
---

# Table: azure_app_service_environment - Query Azure App Service Environments using SQL

Azure App Service Environment is a fully isolated and dedicated environment for securely running App Service apps at high scale. It is designed for applications that require secure, scalable and isolated environments for their execution. It provides network isolation and improved scalability capabilities, making it ideal for applications that require high levels of security and isolation, or that run at a large scale.

## Table Usage Guide

The 'azure_app_service_environment' table provides insights into App Service Environments within Azure. As a DevOps engineer, explore environment-specific details through this table, including the environment's capacity, status, and associated metadata. Utilize it to uncover information about the environments, such as their virtual network integration, worker pool specifications, and inbound and outbound IP addresses. The schema presents a range of attributes of the App Service Environment for your analysis, like the environment's ID, location, resource group, and tags.

## Examples

### List of app service environments which are not healthy
Identify the Azure app service environments that are not functioning properly. This is useful for promptly addressing issues and maintaining optimal application performance.

```sql
select
  name,
  is_healthy_environment
from
  azure_app_service_environment
where
  not is_healthy_environment;
```

### Virtual network info of each app service environment
Gain insights into the virtual network configuration of each app service environment to understand the internal load balancing mode and ensure optimal resource allocation.

```sql
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
Explore the configuration details of your Azure app service environment to gain insights into the specific cluster settings. This can help you understand the current setup and make informed decisions on potential modifications.

```sql
select
  name,
  id,
  settings ->> 'name' as settings_name,
  settings ->> 'value' as settings_value
from
  azure_app_service_environment,
  jsonb_array_elements(cluster_settings) as settings;
```