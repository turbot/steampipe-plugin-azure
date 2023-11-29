---
title: "Steampipe Table: azure_firewall_policy - Query Azure Firewall Policies using SQL"
description: "Allows users to query Azure Firewall Policies"
---

# Table: azure_firewall_policy - Query Azure Firewall Policies using SQL

Azure Firewall Policy is a configuration schema for Azure Firewall that can be used across multiple instances. It provides threat intelligence, service tags, application rules, and network rules as top level properties. Firewall policies can be managed independently from firewall instances, allowing for centralized management of your firewall security rules.

## Table Usage Guide

The 'azure_firewall_policy' table provides insights into Firewall Policies within Azure Firewall. As a security engineer, explore policy-specific details through this table, including threat intelligence, service tags, application rules, and network rules. Utilize it to uncover information about policies, such as those associated with specific firewall instances, the rules they enforce, and their overall configuration. The schema presents a range of attributes of the Firewall Policy for your analysis, like the policy ID, name, type, subscription ID, and associated tags.

## Examples

### Basic info
Explore which firewall policies are currently active within your Azure environment. This can help you assess your security measures and identify any areas that may need additional coverage or modifications.

```sql
select
  name,
  id,
  type,
  provisioning_state,
  sku_tier,
  base_policy,
  child_policies,
  region
from
  azure_firewall_policy;
```

### List policies that are in failed state
Identify the firewall policies that are currently in a failed state. This can assist in troubleshooting and maintaining the overall health of your Azure firewall policies.

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
Determine the details of each firewall policy in Azure, including the number of public IP addresses each firewall has. This is useful for understanding the scope and scale of your firewall protection.

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
Explore the DNS settings of each policy to understand whether a proxy is enabled or required for network rules. This can be useful for analyzing and managing network security configurations.

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
Explore the firewall policies that have specific IP addresses whitelisted, aiding in the understanding of threat intelligence and enhancing security measures.

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
Explore the whitelist domain names of firewall policies to understand potential safe sources of traffic in your Azure environment. This can help you maintain a secure network by identifying trusted entities.

```sql
select
  name,
  id,
  f as whitelist_fqdn
from
  azure_firewall_policy,
  jsonb_array_elements_text(threat_intel_whitelist_fqdns) as f;
```