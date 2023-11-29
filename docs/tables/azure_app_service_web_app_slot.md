---
title: "Steampipe Table: azure_app_service_web_app_slot - Query Azure App Service Web App Slots using SQL"
description: "Allows users to query Azure App Service Web App Slots"
---

# Table: azure_app_service_web_app_slot - Query Azure App Service Web App Slots using SQL

Azure App Service is a fully managed platform for building, deploying, and scaling web apps. You can host web apps, mobile app back ends, RESTful APIs, or automated business processes. Web App Slots are live apps with their own hostnames that are used to deploy different versions of an app and then swap them to production with zero downtime.

## Table Usage Guide

The 'azure_app_service_web_app_slot' table provides insights into Web App Slots within Azure App Service. As a DevOps engineer, explore slot-specific details through this table, including configuration settings, app service plans, and associated metadata. Utilize it to uncover information about slots, such as those in stopped state, the configuration settings of each slot, and the verification of app service plans. The schema presents a range of attributes of the Web App Slot for your analysis, like the slot name, kind, fully qualified domain name, and associated tags.

## Examples

### Basic info
Explore which web application slots in Azure App Service are currently active and when they were last modified. This can be useful to manage and monitor your application deployment slots.

```sql
select
  name,
  app_name,
  id,
  kind,
  state,
  type,
  last_modified_time_utc,
  repository_site_name,
  enabled
from
  azure_app_service_web_app_slot;
```

### List slots where the apps are enabled
Explore which slots have apps enabled to gain insights into active app usage and distribution. This can be beneficial for managing resources and optimizing app performance.

```sql
select
  name,
  app_name,
  state,
  type,
  reserved,
  server_farm_id,
  target_swap_slot
  enabled
from
  azure_app_service_web_app_slot
where
  enabled;
```

### List slots that accept HTTP traffic (i.e only HTTPS is disabled)
Explore which Azure App Service slots are configured to accept HTTP traffic, allowing you to identify potential security vulnerabilities where HTTPS is not enforced. This could be useful in a security audit to ensure all web applications are using secure protocols.

```sql
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_web_app_slot
where
  not https_only;
```

### Host names of each slot
Explore which web application slots are hosted in different regions and resource groups. This can aid in managing and organizing your Azure App Service resources effectively.

```sql
select
  name,
  host_names,
  kind,
  region,
  resource_group
from
  azure_app_service_web_app_slot;
```

### List enabled host names
Explore which web application slots in your Azure App Service are currently enabled. This can be useful for managing your resources and understanding the active components within your cloud environment.

```sql
select
  name,
  id,
  type,
  kind,
  enabled_host_names
from
  azure_app_service_web_app_slot;
```

### Get slot swap status of each slot
Assess the status of slot swaps within your application, identifying the source and destination of each swap. This allows you to track changes and manage your app's performance effectively.

```sql
select
  name,
  type,
  slot_swap_status ->> 'SlotSwapStatus' as slot_swap_status,
  slot_swap_status ->> 'SourceSlotName' as source_slot_name,
  slot_swap_status ->> 'DestinationSlotName' as destination_slot_name
from
  azure_app_service_web_app_slot;
```

### Get site config details of each slot
Assess the configuration details of each web application slot to gain insights into the number of workers, enabled features, and software versions installed. This can help in managing resources and ensuring optimal performance.

```sql
select
  name,
  id,
  site_config ->> 'NumberOfWorkers' as number_of_workers,
  site_config ->> 'DefaultDocuments' as DefaultDocuments,
  site_config ->> 'NetFrameworkVersion' as NetFrameworkVersion,
  site_config ->> 'PhpVersion' as PhpVersion,
  site_config ->> 'PythonVersion' as PythonVersion,
  site_config ->> 'NodeVersion' as NodeVersion,
  site_config ->> 'PowerShellVersion' as PowerShellVersion,
  site_config ->> 'LinuxFxVersion' as LinuxFxVersion,
  site_config ->> 'WindowsFxVersion' as WindowsFxVersion,
  site_config ->> 'RequestTracingEnabled' as RequestTracingEnabled,
  site_config ->> 'RequestTracingExpirationTime' as RequestTracingExpirationTime,
  site_config ->> 'RemoteDebuggingEnabled' as RemoteDebuggingEnabled,
  site_config ->> 'RemoteDebuggingVersion' as RemoteDebuggingVersion,
  site_config ->> 'HTTPLoggingEnabled' as HTTPLoggingEnabled,
  site_config ->> 'DetailedErrorLoggingEnabled' as DetailedErrorLoggingEnabled,
  site_config ->> 'PublishingUsername' as PublishingUsername,
  site_config ->> 'AppSettings' as AppSettings,
  site_config ->> 'ConnectionStrings' as ConnectionStrings,
  site_config ->> 'MachineKey' as MachineKey,
  site_config ->> 'HandlerMappings' as HandlerMappings,
  site_config ->> 'DocumentRoot' as DocumentRoot
from
  azure_app_service_web_app_slot;
```