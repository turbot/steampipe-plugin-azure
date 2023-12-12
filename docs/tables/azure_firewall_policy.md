---
title: "Steampipe Table: azure_firewall_policy - Query Azure Firewall Policies using SQL"
description: "Allows users to query Azure Firewall Policies, providing insights into the rules and settings that govern network traffic flow at the application and network level."
---

# Table: azure_firewall_policy - Query Azure Firewall Policies using SQL

Azure Firewall Policy is a resource in Microsoft Azure that allows you to create, enforce, and log application and network connectivity policies across subscriptions and virtual networks. It provides centralized network and application rule collections that can be referenced by multiple Azure Firewalls. Azure Firewall Policy simplifies management and reduces errors with its ability to manage all Azure Firewalls through Azure Policy and Azure Management Groups.

## Table Usage Guide

The `azure_firewall_policy` table provides insights into Firewall Policies within Microsoft Azure. As a Network Administrator, explore policy-specific details through this table, including rules, settings, and associated metadata. Utilize it to uncover information about policies, such as those governing network traffic flow at the application and network level, providing a centralized way to manage and enforce network connectivity policies.

## Examples

### Basic info
Explore the characteristics of your Azure firewall policies such as their provisioning state, tier, base and child policies, and the region they're set up in. This helps in understanding the configuration and status of your firewall policies, assisting in security management and planning.

```sql+postgres
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

```sql+sqlite
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
Identify instances where certain policies have not been provisioned successfully. This is useful for troubleshooting and rectifying issues to ensure all policies are active and functional.

```sql+postgres
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

```sql+sqlite
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
This query is used to explore the firewall details associated with each policy in Azure. It provides valuable insights into the private and public IP addresses associated with each firewall, aiding in network security management and policy review.

```sql+postgres
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

```sql+sqlite
select
  p.name as firewall_policy_name,
  p.id as firewall_policy_id,
  f.id as firewall_id,
  f.hub_private_ip_address,
  f.hub_public_ip_address_count
from
  azure_firewall_policy as p,
  json_each(p.firewalls) as firewall,
  azure_firewall as f
where
  f.id = json_extract(firewall.value, '$.ID');
```

### Get DNS setting details of each policy
This query helps to analyze the DNS settings for each policy in your Azure firewall. It's useful in understanding the server configurations, whether a proxy is enabled, and if a proxy is required for network rules, thus aiding in security and network management.

```sql+postgres
select
  name,
  id,
  dns_settings ->> 'Servers' as servers,
  dns_settings ->> 'EnableProxy' as enable_proxy,
  dns_settings ->> 'RequireProxyForNetworkRules' as require_proxy_for_network_rules
from
  azure_firewall_policy;
```

```sql+sqlite
select
  name,
  id,
  json_extract(dns_settings, '$.Servers') as servers,
  json_extract(dns_settings, '$.EnableProxy') as enable_proxy,
  json_extract(dns_settings, '$.RequireProxyForNetworkRules') as require_proxy_for_network_rules
from
  azure_firewall_policy;
```

### List threat intel whitelist IP addresses of firewall policies
Determine the areas in which firewall policies have whitelisted IP addresses, which is beneficial for understanding potential security vulnerabilities and ensuring your network is protected from known threats.

```sql+postgres
select
  name,
  id,
  i as whitelist_ip_address
from
  azure_firewall_policy,
  jsonb_array_elements_text(threat_intel_whitelist_ip_addresses) as i;
```

```sql+sqlite
select
  name,
  p.id,
  i.value as whitelist_ip_address
from
  azure_firewall_policy as p,
  json_each(threat_intel_whitelist_ip_addresses) as i;
```

### List threat intel whitelist FQDNs of firewall policies
Explore which firewall policies have specific domains whitelisted, providing a way to identify potential security vulnerabilities or unnecessary exceptions in your Azure firewall configuration.

```sql+postgres
select
  name,
  id,
  f as whitelist_fqdn
from
  azure_firewall_policy,
  jsonb_array_elements_text(threat_intel_whitelist_fqdns) as f;
```

```sql+sqlite
select
  name,
  p.id,
  f.value as whitelist_fqdn
from
  azure_firewall_policy as p,
  json_each(threat_intel_whitelist_fqdns) as f;
```