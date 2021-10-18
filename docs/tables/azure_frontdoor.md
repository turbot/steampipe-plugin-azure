# Table: azure_frontdoor

Azure Front Door is a global, scalable entry-point that uses the Microsoft global edge network to create fast, secure, and widely scalable web applications. With Front Door, you can transform your global consumer and enterprise applications into robust, high-performing personalized modern applications with contents that reach a global audience through Azure.

## Examples

### Basic info

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
