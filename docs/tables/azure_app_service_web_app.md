---
title: "Steampipe Table: azure_app_service_web_app - Query Azure App Service Web Apps using SQL"
description: "Allows users to query Azure App Service Web Apps"
---

# Table: azure_app_service_web_app - Query Azure App Service Web Apps using SQL

Azure App Service is a fully managed platform for building, deploying, and scaling web apps. You can host and scale web apps in Azure with minimal to zero code changes. Azure App Service not only adds the power of Microsoft Azure to your application, such as security, load balancing, and automated management, but also provides the ability to build a web app in your favorite language, be it .NET, .NET Core, Java, Ruby, Node.js, PHP, or Python.

## Table Usage Guide

The 'azure_app_service_web_app' table provides insights into web apps within Azure App Service. As a DevOps engineer, explore web app-specific details through this table, including app settings, configuration details, and associated metadata. Utilize it to uncover information about web apps, such as those with specific configurations, the relationships between different apps, and the verification of app settings. The schema presents a range of attributes of the web app for your analysis, like the app name, resource group, kind, location, and associated tags.

## Examples

### Outbound IP addresses and possible outbound IP addresses info of each web app
Analyze the settings to understand the outbound IP addresses currently in use and potential future ones for each web application. This can help in planning and managing network configurations for improved security and performance.

```sql
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_web_app;
```

### List web apps which accepts HTTP traffics (i.e HTTPS only is disabled)
Discover the segments of your web applications that are potentially insecure by identifying which ones are accepting HTTP traffic. This is useful for understanding where your system may be vulnerable to unencrypted data transfer, aiding in enhancing your overall security measures.

```sql
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

### List of web app where client certificate mode is disabled
Explore which web applications on Azure App Service have the client certificate mode disabled. This can be useful in identifying potential security risks, as applications without client certificates may be more vulnerable to unauthorized access.

```sql
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

### Host names of each web app
Analyze the settings to understand the geographical distribution and organization of your Azure web applications. This can help you manage resources more effectively and plan for scalability.

```sql
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
Explore which Azure web apps are enabled with the latest HTTP version. This can be useful in identifying potential updates needed for apps running on older HTTP versions.

```sql
select
  name,
  enabled,
  region
from
  azure_app_service_web_app
where
  (configuration -> 'properties' ->> 'http20Enabled')::boolean;
```

### List web apps that have FTP deployments set to disabled
Determine the areas in which web apps are operating with FTP deployments disabled. This can be beneficial for auditing security measures and ensuring compliance with company policies that disallow FTP deployments.

```sql
select
  name,
  configuration -> 'properties' ->> 'ftpsState' as ftps_state
from
  azure_app_service_web_app
where
  configuration -> 'properties' ->> 'ftpsState' <> 'AllAllowed';
```

### List web apps that have managed service identity disabled
Discover the segments that have the managed service identity feature disabled in your web applications. This is useful in identifying potential security risks as it allows you to pinpoint applications that might not be properly utilizing Azure's built-in identity management features.

```sql
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
Explore the storage details linked to a specific application within a designated resource group to better manage and allocate resources. This is particularly useful for optimizing storage utilization and planning for future capacity needs.

```sql
select
  name,
  nabled,
  region,
  identity
  storage_info_value
from
  azure_app_service_web_app
where
  resource_group = 'demo'
  and name = 'web-app-test-storage-info';
```