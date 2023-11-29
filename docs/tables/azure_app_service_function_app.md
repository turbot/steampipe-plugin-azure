---
title: "Steampipe Table: azure_app_service_function_app - Query Azure App Service Function Apps using SQL"
description: "Allows users to query Azure App Service Function Apps."
---

# Table: azure_app_service_function_app - Query Azure App Service Function Apps using SQL

Azure App Service is a fully managed platform for building, deploying, and scaling web apps. Azure Function Apps, a part of Azure App Service, is a serverless compute service that lets you run event-triggered code without having to provision or manage infrastructure. It enables developers to host and run chunks of code, or "functions," in the cloud, without needing to create a virtual machine or publish a web application.

## Table Usage Guide

The 'azure_app_service_function_app' table provides insights into Function Apps within Azure App Service. As a DevOps engineer, explore Function App-specific details through this table, including App settings, connection strings, default hostname, and associated metadata. Utilize it to uncover information about Function Apps, such as those with specific configurations, the relationships between apps, and the verification of connection strings. The schema presents a range of attributes of the Function App for your analysis, like the app service plan id, creation date, default hostname, and associated tags.

## Examples

### List of app functions which accepts HTTP traffic
Identify Azure app functions that accept HTTP traffic to assess potential security risks and ensure secure communication protocols are in place.

```sql
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


### List of all unreserved app function
Explore which Azure app service function apps are not reserved. This is useful for identifying potential resource allocation inefficiencies and optimizing your cloud infrastructure.

```sql
select
  name,
  reserved,
  resource_group
from
  azure_app_service_function_app
where
  not reserved;
```


### Outbound IP addresses and possible outbound IP addresses info of each function app
Explore the outbound IP addresses associated with each function app to understand potential network communication paths. This is useful in identifying and managing the network traffic routes for your application.

```sql
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_function_app;
```


### List of app functions where client certificate mode is disabled.
Explore which Azure app service functions have the client certificate mode disabled. This can be useful for identifying potential security vulnerabilities in your application services.

```sql
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