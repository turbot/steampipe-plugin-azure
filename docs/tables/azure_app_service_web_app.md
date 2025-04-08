---
title: "Steampipe Table: azure_app_service_web_app - Query Azure App Service Web Apps using SQL"
description: "Allows users to query Azure App Service Web Apps, providing insights into the configuration, status, and metadata of web applications hosted on the Azure App Service platform."
folder: "App Service"
---

# Table: azure_app_service_web_app - Query Azure App Service Web Apps using SQL

Azure App Service is a fully managed platform for building, deploying, and scaling web applications. It supports a variety of programming languages, tools, and frameworks, including both Microsoft-specific and third-party software and systems. With Azure App Service, you can quickly build, deploy, and scale enterprise-grade web, mobile, and API apps running on any platform.

## Table Usage Guide

The `azure_app_service_web_app` table provides insights into web applications hosted on Azure App Service. As a developer or system administrator, you can use this table to examine the configuration, status, and metadata of these applications. It can be particularly useful for monitoring and managing your web applications, ensuring they are correctly configured, running smoothly, and adhering to your organization's operational and security policies.

## Examples

### Outbound IP addresses and possible outbound IP addresses info of each web app
Explore which web applications in your Azure App Service have specific outbound IP addresses. This is useful for understanding the network behavior of your applications, particularly for security monitoring or compliance purposes.

```sql+postgres
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_web_app;
```

```sql+sqlite
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_web_app;
```

### List web apps which accepts HTTP traffics (i.e HTTPS only is disabled)
Determine the areas in which web applications are accepting HTTP traffic, indicating that the more secure HTTPS-only mode is disabled. This can be useful for identifying potential security risks in your Azure App Service.

```sql+postgres
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_web_app
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
  azure_app_service_web_app
where
  https_only = 0;
```

### List of web app where client certificate mode is disabled
Determine the areas in which web applications are potentially vulnerable due to disabled client certificate mode. This is crucial for enhancing security measures and ensuring data protection.

```sql+postgres
select
  name,
  client_cert_enabled,
  kind,
  region
from
  azure_app_service_web_app
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
  azure_app_service_web_app
where
  client_cert_enabled = 0;
```

### Host names of each web app
Determine the areas in which your web applications are hosted. This aids in understanding their geographical distribution and aids in resource management.

```sql+postgres
select
  name,
  host_names,
  kind,
  region,
  resource_group
from
  azure_app_service_web_app;
```

```sql+sqlite
select
  name,
  host_names,
  kind,
  region,
  resource_group
from
  azure_app_service_web_app;
```

### List web apps with latest HTTP version
Determine the areas in which web applications are running on the latest HTTP version across different regions. This can be useful for ensuring applications are up-to-date and taking advantage of the latest protocol features for performance and security.

```sql+postgres
select
  name,
  enabled,
  region
from
  azure_app_service_web_app
where
  (configuration -> 'properties' ->> 'http20Enabled')::boolean;
```

```sql+sqlite
select
  name,
  enabled,
  region
from
  azure_app_service_web_app
where
  json_extract(configuration, '$.properties.http20Enabled') = 'true';
```

### List web apps that have FTP deployments set to disabled
Determine the areas in which web applications have FTP deployments disabled, allowing for a better understanding of security measures in place and potential areas of vulnerability.

```sql+postgres
select
  name,
  configuration -> 'properties' ->> 'ftpsState' as ftps_state
from
  azure_app_service_web_app
where
  configuration -> 'properties' ->> 'ftpsState' <> 'AllAllowed';
```

```sql+sqlite
select
  name,
  json_extract(json_extract(configuration, '$.properties'), '$.ftpsState') as ftps_state
from
  azure_app_service_web_app
where
  json_extract(json_extract(configuration, '$.properties'), '$.ftpsState') <> 'AllAllowed';
```

### List web apps that have managed service identity disabled
Determine the areas in which web apps are operating without a managed service identity, which is a key security feature. This could be used to identify potential vulnerabilities and improve overall system security.

```sql+postgres
select
  name,
  enabled,
  region,
  identity
from
  azure_app_service_web_app
where
  identity = '{}';
```

```sql+sqlite
select
  name,
  enabled,
  region,
  identity
from
  azure_app_service_web_app
where
  identity = '{}';
```

### Get the storage information associated to a particular app
Explore the storage details linked to a specific application in Azure's App Service. This can help you understand the configuration and enablement status of your storage in a particular region, which can be crucial for optimizing resource allocation and management.

```sql+postgres
select
  name,
  enabled,
  region,
  identity
  storage_info_value
from
  azure_app_service_web_app
where
  resource_group = 'demo'
  and name = 'web-app-test-storage-info';
```

```sql+sqlite
select
  name,
  enabled,
  region,
  identity,
  storage_info_value
from
  azure_app_service_web_app
where
  resource_group = 'demo'
  and name = 'web-app-test-storage-info';
```