# Table: azure_firewall_policy

Azure Firewall Policy is the recommended method to configure your Azure Firewall. It's a global resource that can be used across multiple Azure Firewall instances in Secured Virtual Hubs and Hub Virtual Networks. Policies work across regions and subscriptions.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state,
  sku_tier,
  base_policy,,
  child_policies,
  location,
from
  azure_firewall_policy;
```

### List policies that are in failed state

```sql
select
  name,
  id,
  dns_settings,
  firewalls
from
  azure_firewall_policy
where
  provisioning_state = 'Failed';
```

### Get firewall details of each policy

```sql
select
  p.name as firewall_policy_name,
  p.id as firewall_policy_id,
  f.id as firewall_id,
  f.hub_private_ip_address,
  f.hub_public_ip_address_count
from
  azure_firewall_policy as p,
  jsonb_array_elements(p.firewalls) as firewall,
  azure_firewall as f
where
  f.id = firewall ->> 'ID';
```

### Get DNS setting details of each policy

```sql
select
  name,
  id,
  dns_settings ->> 'Servers' as servers,
  dns_settings ->> 'EnableProxy' as enable_proxy,
  dns_settings ->> 'RequireProxyForNetworkRules' as require_proxy_for_network_rules
from
  azure_firewall_policy;
```

### List threat intel whitelist IP addresses of firewall policies

```sql
select
  name,
  id,
  i as whitelist_ip_address
from
  azure_firewall_policy,
  jsonb_array_elements_text(threat_intel_whitelist_ip_addresses) as i;
```

### List threat intel whitelist FQDNs of firewall policies

```sql
select
  name,
  id,
  f as whitelist_fqdn
from
  azure_firewall_policy,
  jsonb_array_elements_text(threat_intel_whitelist_fqdns) as f;
```