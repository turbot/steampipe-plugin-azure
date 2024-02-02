---
title: "Steampipe Table: azure_spring_cloud_app - Query Azure Spring Cloud Apps using SQL"
description: "Allows users to query Azure Spring Cloud Apps, specifically the details of the application, providing insights into the configuration and state of the Spring Cloud applications."
---

# Table: azure_spring_cloud_app - Query Azure Spring Cloud Apps using SQL

Azure Spring Cloud is a fully managed service for Spring Boot apps that lets you focus on building and running the apps that run your business without the hassle of managing infrastructure. It provides a platform for deploying and managing Spring Boot and Spring Cloud applications in the cloud. The service is jointly built, operated, and supported by Pivotal Software and Microsoft to provide a native platform designed to be easily run and managed on Azure.

## Table Usage Guide

The `azure_spring_cloud_app` table provides insights into Azure Spring Cloud Apps within Microsoft Azure. As a DevOps engineer, explore service-specific details through this table, including configurations, provisioning state, and associated metadata. Utilize it to uncover information about apps, such as service, publoc expose, app URL, HTTPS support, and the verification of app configurations.

## Examples

### Basic info

Explore the comprehensive information about Azure Spring Cloud apps, which is beneficial for administrators and developers managing these applications.

```sql+postgres
select
  id,
  name,
  service_name,
  type,
  provisioning_state,
  https_only,
  active_deployment_name,
  created_time,
  public,
  location
from
  azure_spring_cloud_app;
```

```sql+sqlite
select
  id,
  name,
  service_name,
  type,
  provisioning_state,
  https_only,
  active_deployment_name,
  created_time,
  public,
  location
from
  azure_spring_cloud_app;
```

### List failed applications
 For troubleshooting and maintenance, as it helps administrators quickly identify which apps have failed to provision and may require attention or intervention.

```sql+postgres
select
  id,
  name,
  service_name,
  provisioning_state,
  created_time,
  location
from
  azure_spring_cloud_app
where
  provisioning_state = 'Failed';
```

```sql+sqlite
select
  id,
  name,
  service_name,
  provisioning_state,
  created_time,
  location
from
  azure_spring_cloud_app
where
  provisioning_state = 'Failed';
```

### List applications that are newly added in the last 10 days
The query is valuable for administrators and DevOps teams to oversee recent activities, identify potential issues, and maintain effective control over their Azure Spring Cloud applications.

```sql+postgres
select
  name,
  id,
  service_name,
  public,
  provisioning_state,
  created_time,
  url
from
  azure_spring_cloud_app
where
  created_time >= now() - interval '10' day;
```

```sql+sqlite
select
  name,
  id,
  service_name,
  public,
  provisioning_state,
  created_time,
  url
from
  azure_spring_cloud_app
where
  created_time >= datetime('now', '-10 days');
```

### Get disk details of applications
This query is useful in a cloud computing context, specifically for managing Azure Spring Cloud applications. It extracts detailed information about both temporary and persistent disks associated with these applications, such as their size, used space, and mount paths. Understanding disk usage helps in optimizing resource allocation and performance tuning. Knowing the size and utilization of disks aids in planning for future storage needs.

```sql+postgres
select
  name,
  service_name,
  fqdn,
  temporary_disk ->> 'SizeInGB' as temporary_disk_size_in_gb,
  temporary_disk ->> 'MountPath' as temporary_disk_mount_path,
  persistent_disk ->> 'SizeInGB' as persistent_disk_size_in_gb,
  persistent_disk ->> 'UsedInGB' as persistent_disk_used_in_gb,
  persistent_disk ->> 'MountPath' as persistent_disk_mount_path
from
  azure_spring_cloud_app;
```

```sql+sqlite
select
  name,
  service_name,
  fqdn,
  json_extract(temporary_disk, '$.SizeInGB') as temporary_disk_size_in_gb,
  json_extract(temporary_disk, '$.MountPath') as temporary_disk_mount_path,
  json_extract(persistent_disk, '$.SizeInGB') as persistent_disk_size_in_gb,
  json_extract(persistent_disk, '$.UsedInGB') as persistent_disk_used_in_gb,
  json_extract(persistent_disk, '$.MountPath') as persistent_disk_mount_path
from
  azure_spring_cloud_app;
```

### Get service details of applications
Analyze information between Azure Spring Cloud applications and their associated services. Understanding the relationship between applications and their underlying services. Inclusion of diagnostic settings provides insights for monitoring and troubleshooting.

```sql+postgres
select
  a.name,
  a.id as app_id,
  a.service_name,
  s.service_id,
  s.sku_name as service_sku_name,
  s.sku_tier as service_sku_tier,
  s.version as service_version,
  s.diagnostic_settings
from
  azure_spring_cloud_app as a,
  azure_spring_cloud_service as s
where
  a.service_name = s.name;
```

```sql+sqlite
select
  a.name,
  a.id as app_id,
  a.service_name,
  s.service_id,
  s.sku_name as service_sku_name,
  s.sku_tier as service_sku_tier,
  s.version as service_version,
  s.diagnostic_settings
from
  azure_spring_cloud_app as a
  join azure_spring_cloud_service as s on a.service_name = s.name;
```