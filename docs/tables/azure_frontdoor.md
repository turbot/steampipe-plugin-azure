---
title: "Steampipe Table: azure_frontdoor - Query Azure Front Door Services using SQL"
description: "Allows users to query Azure Front Door Services, providing insights into the routing and acceleration of web traffic."
---

# Table: azure_frontdoor - Query Azure Front Door Services using SQL

Azure Front Door Service is a scalable and secure entry point that uses the global network infrastructure of Azure. It enables you to define, manage, and monitor the global routing for your web traffic by optimizing for best performance and instant global failover for high availability. With Azure Front Door Service, you can transform your global (multi-region) consumer and enterprise applications into robust, high-performance personalized modern applications.

## Table Usage Guide

The `azure_frontdoor` table provides insights into Azure Front Door Services within Azure. As a network engineer or a system administrator, explore details about your front door services such as resource group, subscription ID, and health probe settings through this table. Utilize it to uncover information about your front door services, such as the load balancing settings, routing rules, and backend pools, enabling you to optimize your web traffic routing and acceleration.

## Examples

### Basic info
Explore which Azure Front Door resources are currently being provisioned, providing insights into the state of your network infrastructure and helping to identify any potential issues or bottlenecks.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  cname
from
  azure_frontdoor;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  cname
from
  azure_frontdoor;
```

### List backend pools settings details
Determine the enforcement of certificate name checks and the send/receive timeout settings within your Azure Front Door service. This can help ensure secure connections and manage response times for your web applications.

```sql+postgres
select
  name,
  id,
  backend_pools_settings -> 'enforceCertificateNameCheck' as settings_enforce_certificate_name_check,
  backend_pools_settings -> 'sendRecvTimeoutSeconds' as settings_send_recv_timeout_seconds
from
  azure_frontdoor;
```

```sql+sqlite
select
  name,
  id,
  json_extract(backend_pools_settings, '$.enforceCertificateNameCheck') as settings_enforce_certificate_name_check,
  json_extract(backend_pools_settings, '$.sendRecvTimeoutSeconds') as settings_send_recv_timeout_seconds
from
  azure_frontdoor;
```

### List routing rules details
This query is useful for gaining insights into the specifics of routing rules within the Azure Front Door service. It allows users to analyze factors such as enabled states, resource states, accepted protocols, and route configurations, aiding in the optimization of network traffic routing.

```sql+postgres
select
  name,
  id,
  rule ->> 'id' as rule_id,
  rule ->> 'name' as rule_name,
  rule -> 'properties' ->> 'enabledState' as rule_property_enabled_state,
  rule -> 'properties' ->> 'resourceState' as rule_property_resource_state,
  jsonb_pretty(rule -> 'properties' -> 'acceptedProtocols') as rule_property_accepted_protocols,
  jsonb_pretty(rule -> 'properties' -> 'frontendEndpoints') as rule_property_frontend_endpoints,
  jsonb_pretty(rule -> 'properties' -> 'patternsToMatch') as rule_property_patterns_to_match,
  jsonb_pretty(rule -> 'properties' -> 'routeConfiguration') as rule_property_route_configuration
from
  azure_frontdoor,
  jsonb_array_elements(routing_rules) as rule;
```

```sql+sqlite
select
  name,
  id,
  json_extract(rule.value, '$.id') as rule_id,
  json_extract(rule.value, '$.name') as rule_name,
  json_extract(rule.value, '$.properties.enabledState') as rule_property_enabled_state,
  json_extract(rule.value, '$.properties.resourceState') as rule_property_resource_state,
  rule.value || json_extract(rule.value, '$.properties.acceptedProtocols') as rule_property_accepted_protocols,
  rule.value || json_extract(rule.value, '$.properties.frontendEndpoints') as rule_property_frontend_endpoints,
  rule.value || json_extract(rule.value, '$.properties.patternsToMatch') as rule_property_patterns_to_match,
  rule.value || json_extract(rule.value, '$.properties.routeConfiguration') as rule_property_route_configuration
from
  azure_frontdoor,
  json_each(routing_rules) as rule;
```

### List load balancing settings details
Analyze the settings to understand the details of load balancing configurations in Azure Front Door. This is useful for assessing performance tuning and troubleshooting issues related to load distribution.

```sql+postgres
select
  name,
  id,
  setting ->> 'id' as setting_id,
  setting ->> 'name' as setting_name,
  setting -> 'properties' -> 'additionalLatencyMilliseconds' as setting_property_additional_latency_milliseconds,
  setting -> 'properties' -> 'successfulSamplesRequired' as setting_property_successful_samples_required,
  setting -> 'properties' -> 'sampleSize' as setting_property_sample_size,
  setting -> 'properties' ->> 'resourceState' as setting_property_resource_state
from
  azure_frontdoor,
  jsonb_array_elements(load_balancing_settings) as setting;
```

```sql+sqlite
select
  name,
  id,
  json_extract(setting.value, '$.id') as setting_id,
  json_extract(setting.value, '$.name') as setting_name,
  json_extract(setting.value, '$.properties.additionalLatencyMilliseconds') as setting_property_additional_latency_milliseconds,
  json_extract(setting.value, '$.properties.successfulSamplesRequired') as setting_property_successful_samples_required,
  json_extract(setting.value, '$.properties.sampleSize') as setting_property_sample_size,
  json_extract(setting.value, '$.properties.resourceState') as setting_property_resource_state
from
  azure_frontdoor,
  json_each(load_balancing_settings) as setting;
```

### List frontend endpoints details
Explore the specifics of your frontend endpoints to understand their configuration and properties. This can be useful for assessing the status, security settings, and session affinity details of your web application's frontend endpoints.

```sql+postgres
select
  name,
  id,
  endpoint ->> 'id' as endpoint_id,
  endpoint ->> 'name' as endpoint_name,
  endpoint -> 'properties' ->> 'hostName' as endpoint_property_host_name,
  endpoint -> 'properties' ->> 'sessionAffinityEnabledState' as endpoint_property_session_affinity_enabled_state,
  endpoint -> 'properties' -> 'sessionAffinityTtlSeconds' as endpoint_property_session_affinity_ttl_seconds,
  endpoint -> 'properties' ->> 'resourceState' as endpoint_property_resource_state,
  jsonb_pretty(endpoint -> 'properties' -> 'webApplicationFirewallPolicyLink') as endpoint_property_web_application_firewall_policy_link
from
  azure_frontdoor,
  jsonb_array_elements(frontend_endpoints) as endpoint;
```

```sql+sqlite
select
  name,
  id,
  json_extract(endpoint.value, '$.id') as endpoint_id,
  json_extract(endpoint.value, '$.name') as endpoint_name,
  json_extract(endpoint.value, '$.properties.hostName') as endpoint_property_host_name,
  json_extract(endpoint.value, '$.properties.sessionAffinityEnabledState') as endpoint_property_session_affinity_enabled_state,
  json_extract(endpoint.value, '$.properties.sessionAffinityTtlSeconds') as endpoint_property_session_affinity_ttl_seconds,
  json_extract(endpoint.value, '$.properties.resourceState') as endpoint_property_resource_state,
  endpoint.value as endpoint_property_web_application_firewall_policy_link
from
  azure_frontdoor,
  json_each(frontend_endpoints) as endpoint;
```

### List health probe settings details
Discover the specifics of health probe settings in your Azure Front Door service. This can help identify potential issues and optimize your network's performance by understanding the intervals, methods, and states of your health probes.

```sql+postgres
select
  name,
  id,
  setting ->> 'id' as setting_id,
  setting ->> 'name' as setting_name,
  setting -> 'properties' -> 'intervalInSeconds' as setting_property_interval_in_seconds,
  setting -> 'properties' ->> 'healthProbeMethod' as setting_property_health_probe_method,
  setting -> 'properties' ->> 'enabledState' as setting_property_enabled_state,
  setting -> 'properties' ->> 'path' as setting_property_path,
  setting -> 'properties' ->> 'protocol' as setting_property_protocol,
  setting -> 'properties' ->> 'resourceState' as setting_property_resource_state
from
  azure_frontdoor,
  jsonb_array_elements(health_probe_settings) as setting;
```

```sql+sqlite
select
  name,
  id,
  json_extract(setting.value, '$.id') as setting_id,
  json_extract(setting.value, '$.name') as setting_name,
  json_extract(setting.value, '$.properties.intervalInSeconds') as setting_property_interval_in_seconds,
  json_extract(setting.value, '$.properties.healthProbeMethod') as setting_property_health_probe_method,
  json_extract(setting.value, '$.properties.enabledState') as setting_property_enabled_state,
  json_extract(setting.value, '$.properties.path') as setting_property_path,
  json_extract(setting.value, '$.properties.protocol') as setting_property_protocol,
  json_extract(setting.value, '$.properties.resourceState') as setting_property_resource_state
from
  azure_frontdoor,
  json_each(health_probe_settings) as setting;
```