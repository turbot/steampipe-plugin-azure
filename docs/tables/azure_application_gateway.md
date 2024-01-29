---
title: "Steampipe Table: azure_application_gateway - Query Azure Application Gateways using SQL"
description: "Allows users to query Azure Application Gateways, providing detailed information about the configuration and status of each gateway."
---

# Table: azure_application_gateway - Query Azure Application Gateways using SQL

Azure Application Gateway is a web traffic load balancer that enables you to manage traffic to your web applications. It operates at the application layer (Layer 7) of the Open Systems Interconnection (OSI) network stack, and can route traffic based on various attributes of an HTTP request. The gateway also provides SSL offloading, which removes the SSL-based encryption from incoming traffic.

## Table Usage Guide

The `azure_application_gateway` table provides insights into Application Gateways within Azure. As a Network Administrator, explore gateway-specific details through this table, including backend pools, listeners, rules, and associated metadata. Utilize it to uncover information about gateways, such as their health status, configured rules, and the verification of SSL certificates.

## Examples

### Basic info
Explore the general information of your Azure application gateways to gain insights into their types, provisioning states, and regions. This allows you to effectively manage and monitor your gateways, ensuring optimal performance and resource allocation.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  region
from
  azure_application_gateway;
```

```sql+sqlite
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
Determine the configurations of your web application firewall for an application gateway. This query aids in understanding the firewall's operational settings, such as enabled status, file upload limits, and rule set details, which are crucial for maintaining optimal security and performance.

```sql+postgres
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

```sql+sqlite
select
  id,
  name,
  web_application_firewall_configuration as disabled_rule_groups,
  json_extract(web_application_firewall_configuration, '$.enabled') as enabled,
  web_application_firewall_configuration as exclusions,
  json_extract(web_application_firewall_configuration, '$.fileUploadLimitInMb') as file_upload_limit_in_mb,
  json_extract(web_application_firewall_configuration, '$.firewallMode') as firewall_mode,
  json_extract(web_application_firewall_configuration, '$.maxRequestBodySizeInKb') as max_request_body_size_in_kb,
  json_extract(web_application_firewall_configuration, '$.requestBodyCheck') as request_body_check,
  json_extract(web_application_firewall_configuration, '$.ruleSetType') as rule_set_type,
  json_extract(web_application_firewall_configuration, '$.ruleSetVersion') as rule_set_version
from
  azure_application_gateway;
```

### List http listeners for application gateway
Explore the configuration of HTTP listeners in an application gateway to understand the protocol requirements and server name indication settings. This can be particularly useful in identifying potential security weak points and optimizing network performance.

```sql+postgres
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

```sql+sqlite
select
  g.id,
  name,
  json_extract(listeners.value, '$.id') as listener_id,
  json_extract(listeners.value, '$.name') as listener_name,
  json_extract(listeners.value, '$.properties.frontendPort') as listener_frontend_port,
  json_extract(listeners.value, '$.properties.hostNames') as listener_host_names,
  json_extract(listeners.value, '$.properties.protocol') as listener_protocol,
  json_extract(listeners.value, '$.properties.requireServerNameIndication') as listener_require_server_name_indication
from
  azure_application_gateway as g,
  json_each(http_listeners) as listeners;
```

### List backend http settings collection for application gateway
Analyze the settings to understand the configuration of your application gateway's backend HTTP settings. This could be useful for assessing aspects like affinity based on cookies, host name selection from backend address, port, protocol, and request timeout.

```sql+postgres
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

```sql+sqlite
select
  g.id,
  name,
  json_extract(settings.value, '$.id') as settings_id,
  json_extract(settings.value, '$.name') as settings_name,
  json_extract(settings.value, '$.properties.cookieBasedAffinity') as settings_cookie_based_affinity,
  json_extract(settings.value, '$.properties.pickHostNameFromBackendAddress') as settings_pick_host_name_from_backend_address,
  json_extract(settings.value, '$.properties.port') as settings_port,
  json_extract(settings.value, '$.properties.protocol') as settings_protocol,
  json_extract(settings.value, '$.properties.requestTimeout') as settings_request_timeout
from
  azure_application_gateway as g,
  json_each(backend_http_settings_collection) as settings;
```

### List frontend IP configurations for application gateway
This query is useful for gaining insights into the IP configurations of your application gateway in Azure. It allows you to understand both the public and private allocation methods, which is critical for managing network access and security.

```sql+postgres
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

```sql+sqlite
select
  g.id,
  name,
  json_extract(config.value, '$.id') as config_id,
  json_extract(config.value, '$.name') as config_name,
  json_extract(config.value, '$.properties.publicIPAddress') as config_public_ip_address,
  json_extract(config.value, '$.properties.privateIPAllocationMethod') as config_private_ip_allocation_method
from
  azure_application_gateway as g,
  json_each(frontend_ip_configurations) as config;
```