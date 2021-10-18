# Table: azure_application_gateway

Azure Application Gateway is a web traffic load balancer that enables you to manage traffic to your web applications. Traditional load balancers operate at the transport layer (OSI layer 4 - TCP and UDP) and route traffic based on source IP address and port, to a destination IP address and port. Application Gateway can make routing decisions based on additional attributes of an HTTP request, for example URI path or host headers.

## Examples

### Basic info

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
