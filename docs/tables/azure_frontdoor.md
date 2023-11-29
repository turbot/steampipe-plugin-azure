---
title: "Steampipe Table: azure_frontdoor - Query Azure Front Door Services using SQL"
description: "Allows users to query Azure Front Door Services."
---

# Table: azure_frontdoor - Query Azure Front Door Services using SQL

Azure Front Door is a scalable and secure entry point for fast delivery of your global web applications. It provides global load balancing and site acceleration service for fast and reliable application delivery at global scale. It offers SSL offload, path-based routing, fast failover, and many more capabilities.

## Table Usage Guide

The 'azure_frontdoor' table provides insights into Front Door Services within Azure. As a DevOps engineer, explore service-specific details through this table, including routing rules, backend pools, frontend endpoints, and associated metadata. Utilize it to uncover information about services, such as those with specific routing rules, the health probes between backend pools, and the verification of frontend endpoints. The schema presents a range of attributes of the Front Door Service for your analysis, like the service ID, creation date, enabled state, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure Front Door service to understand its current state and type. This can help you assess the overall setup and configuration for effective resource management.

```sql
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
Explore the configuration settings of backend pools in an Azure Front Door service. This allows you to assess security measures, like enforcing certificate name checks, and performance parameters, such as send/receive timeout durations.

```sql
select
  name,
  id,
  backend_pools_settings -> 'enforceCertificateNameCheck' as settings_enforce_certificate_name_check,
  backend_pools_settings -> 'sendRecvTimeoutSeconds' as settings_send_recv_timeout_seconds
from
  azure_frontdoor;
```

### List routing rules details
Determine the specific details of routing rules, such as their enabled state, resource state, accepted protocols, and associated endpoints. This can assist in understanding how traffic is being directed and managed within your Azure Front Door service.

```sql
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

### List load balancing settings details
Explore the specifics of load balancing settings to assess their properties and understand their configuration, which is crucial for managing traffic distribution and ensuring efficient resource utilization.

```sql
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

### List frontend endpoints details
Explore the details of frontend endpoints to gain insights into their properties such as host name, session affinity enabled state, and resource state. This can be useful in understanding and managing the configuration of these endpoints, especially in terms of their security settings like the web application firewall policy link.

```sql
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

### List health probe settings details
Explore the configuration of health probe settings to understand how they are set up and functioning. This can help in assessing the performance and reliability of your network connections.

```sql
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