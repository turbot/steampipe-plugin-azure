---
title: "Steampipe Table: azure_app_service_function_app - Query Azure App Service Function Apps using SQL"
description: "Allows users to query Azure App Service Function Apps, specifically providing access to configuration details, app settings, and connection strings."
---

# Table: azure_app_service_function_app - Query Azure App Service Function Apps using SQL

Azure App Service Function Apps is a service within Microsoft Azure that allows developers to host and run functions in the cloud without having to manage any infrastructure. It offers an event-driven, compute-on-demand experience that extends the existing Azure App Service platform. With Azure Function Apps, developers can quickly create serverless applications that scale and integrate with other services.

## Table Usage Guide

The `azure_app_service_function_app` table provides insights into Function Apps within Azure App Service. As a developer or DevOps engineer, explore Function App-specific details through this table, including configuration settings, app settings, and connection strings. Utilize it to uncover information about Function Apps, such as their runtime versions, hosting details, and the state of always-on functionality.

## Examples

### List of app functions which accepts HTTP traffic
Determine the areas in which Azure app services function apps are configured to accept HTTP traffic, which can be useful for identifying potential security risks associated with unencrypted data transmission.

```sql+postgres
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_function_app
where
  not https_only;
```

```sql+sqlite
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_function_app
where
  https_only = 0;
```

### List of all unreserved app function
Discover the segments that comprise all unreserved function apps in Azure. This query is useful in managing resources and ensuring optimal performance by identifying potential areas for reallocation.

```sql+postgres
select
  name,
  reserved,
  resource_group
from
  azure_app_service_function_app
where
  not reserved;
```

```sql+sqlite
select
  name,
  reserved,
  resource_group
from
  azure_app_service_function_app
where
  reserved = 0;
```

### Outbound IP addresses and possible outbound IP addresses info of each function app
Gain insights into the outbound IP addresses associated with each function app, as well as potential outbound IP addresses. This information can be useful for managing network security and understanding your app's communication pathways.

```sql+postgres
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_function_app;
```

```sql+sqlite
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_function_app;
```

### List of app functions where client certificate mode is disabled.
Identify instances where the client certificate mode is disabled in your Azure app functions. This can help enhance security by pinpointing areas where client authentication is not enforced.

```sql+postgres
select
  name,
  client_cert_enabled,
  kind,
  region
from
  azure_app_service_function_app
where
  not client_cert_enabled;
```

```sql+sqlite
select
  name,
  client_cert_enabled,
  kind,
  region
from
  azure_app_service_function_app
where
  client_cert_enabled = 0;
```