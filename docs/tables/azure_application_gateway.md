---
title: "Steampipe Table: azure_application_gateway - Query Azure Network Application Gateways using SQL"
description: "Allows users to query Azure Network Application Gateways"
---

# Table: azure_application_gateway - Query Azure Network Application Gateways using SQL

An Azure Application Gateway is a web traffic load balancer that enables you to manage traffic to your web applications. It operates at the application layer (Layer 7) of the Open Systems Interconnection (OSI) model. This service provides routing capabilities and can make routing decisions based on additional attributes of an HTTP request, for instance, URI path or host headers.

## Table Usage Guide

The 'azure_application_gateway' table provides insights into Application Gateways within Azure Network. As a Network Engineer, explore Application Gateway-specific details through this table, including backend configurations, SSL policy, and associated metadata. Utilize it to uncover information about Application Gateways, such as their SKU, operational state, and the verification of SSL policies. The schema presents a range of attributes of the Application Gateway for your analysis, like the gateway's ID, name, type, region, and associated tags.

## Examples

### Basic info
Explore which application gateways in your Azure environment are currently being provisioned and where they are located. This is beneficial for keeping track of your network resources and their geographical distribution.

```sql
select
  id,
  name,
  type,
  provisioning_state,
  region
from
  azure_application_gateway;
```

### List web application firewall configurations for application gateway 
Analyze the settings to understand the configurations of your web application firewall for application gateways. This can help you assess its current status, identify any disabled rule groups, exclusions, and understand the limitations such as file upload limit and maximum request body size.

```sql
select
  id,
  name,
  jsonb_pretty(web_application_firewall_configuration -> 'disabledRuleGroups') as disabled_rule_groups,
  web_application_firewall_configuration -> 'enabled' as enabled,
  jsonb_pretty(web_application_firewall_configuration -> 'exclusions') as exclusions,
  web_application_firewall_configuration -> 'fileUploadLimitInMb' as file_upload_limit_in_mb,
  web_application_firewall_configuration -> 'firewallMode' as firewall_mode,
  web_application_firewall_configuration -> 'maxRequestBodySizeInKb' as max_request_body_size_in_kb,
  web_application_firewall_configuration -> 'requestBodyCheck' as request_body_check,
  web_application_firewall_configuration -> 'ruleSetType' as rule_set_type,
  web_application_firewall_configuration -> 'ruleSetVersion' as rule_set_version
from
  azure_application_gateway;
```

### List http listeners for application gateway 
Determine the areas in which HTTP listeners for the application gateway are configured. This is useful for understanding the setup and configuration of your application gateway, particularly for troubleshooting or optimizing network traffic management.

```sql
select
  id,
  name,
  listeners -> 'id' as listener_id,
  listeners -> 'name' as listener_name,
  jsonb_pretty(listeners -> 'properties' -> 'frontendPort') as listener_frontend_port,
  jsonb_pretty(listeners -> 'properties' -> 'hostNames') as listener_host_names,
  listeners -> 'properties' -> 'protocol' as listener_protocol,
  listeners -> 'properties' -> 'requireServerNameIndication' as listener_require_server_name_indication
from
  azure_application_gateway,
  jsonb_array_elements(http_listeners) as listeners;
```

### List backend http settings collection for application gateway 
This query aids in gaining insights into the backend HTTP settings for an application gateway. It's particularly useful for understanding settings such as cookie-based affinity, host name selection, port, protocol, and request timeout, which can help optimize the application gateway's performance and security.

```sql
select
  id,
  name,
  settings -> 'id' as settings_id,
  settings -> 'name' as settings_name,
  settings -> 'properties' -> 'cookieBasedAffinity' as settings_cookie_based_affinity,
  settings -> 'properties' -> 'pickHostNameFromBackendAddress' as settings_pick_host_name_from_backend_address,
  settings -> 'properties' -> 'port' as settings_port,
  settings -> 'properties' -> 'protocol' as settings_protocol,
  settings -> 'properties' -> 'requestTimeout' as settings_request_timeout
from
  azure_application_gateway,
  jsonb_array_elements(backend_http_settings_collection) as settings;
```

### List frontend IP configurations for application gateway 
This example helps you identify the different frontend IP configurations for your application gateway in Azure. It's useful for managing and understanding the various IP settings associated with your application gateway, including public and private IP allocations.

```sql
select
  id,
  name,
  config -> 'id' as config_id,
  config -> 'name' as config_name,
  jsonb_pretty(config -> 'properties' -> 'publicIPAddress') as config_public_ip_address,
  config -> 'properties' -> 'privateIPAllocationMethod' as config_private_ip_allocation_method
from
  azure_application_gateway,
  jsonb_array_elements(frontend_ip_configurations) as config;
```