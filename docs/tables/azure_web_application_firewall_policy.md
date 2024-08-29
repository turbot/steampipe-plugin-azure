---
title: "Steampipe Table: azure_web_application_firewall_policy - Query Azure Web Application Firewall Policies using SQL"
description: "Allows users to query Azure Web Application Firewall (WAF) policies, providing detailed information on policy configurations, rules, and associated resources."
---

# Table: azure_web_application_firewall_policy - Query Azure Web Application Firewall Policies using SQL

Azure Web Application Firewall (WAF) is a service that helps protect your web applications by filtering and monitoring HTTP traffic between a web application and the Internet. The `azure_web_application_firewall_policy` table in Steampipe allows you to query information about WAF policies in Azure, including their custom and managed rules, associated application gateways, HTTP listeners, and more.

## Table Usage Guide

The `azure_web_application_firewall_policy` table enables cloud administrators and security engineers to gather detailed insights into their WAF policies. You can query various aspects of the policies, such as their custom rules, managed rules, provisioning state, and associated application gateways. This table is particularly useful for monitoring the security of web applications, managing WAF policy configurations, and ensuring that your web applications are protected against threats.

## Examples

### Basic info
Retrieve basic information about your Azure WAF policies, including their name, resource group, and region.

```sql+postgres
select
  name,
  id,
  resource_group,
  region,
  provisioning_state
from
  azure_web_application_firewall_policy;
```

```sql+sqlite
select
  name,
  id,
  resource_group,
  region,
  provisioning_state
from
  azure_web_application_firewall_policy;
```

### List policies with custom rules
Fetch WAF policies that include custom rules, which can be useful for identifying policies with specific, user-defined security measures.

```sql+postgres
select
  name,
  custom_rules
from
  azure_web_application_firewall_policy
where
  custom_rules is not null;
```

```sql+sqlite
select
  name,
  custom_rules
from
  azure_web_application_firewall_policy
where
  json_extract(custom_rules, '$[0]') is not null;
```

### List policies by provisioning state
Identify WAF policies based on their provisioning state, such as those that are currently updating or have failed.

```sql+postgres
select
  name,
  provisioning_state,
  region
from
  azure_web_application_firewall_policy
where
  provisioning_state = 'ProvisioningStateUpdating';
```

```sql+sqlite
select
  name,
  provisioning_state,
  region
from
  azure_web_application_firewall_policy
where
  provisioning_state = 'ProvisioningStateUpdating';
```

### Get managed rules for each policy
Retrieve the managed rules associated with each WAF policy to understand the built-in protections that are applied to your web applications.

```sql+postgres
select
  name,
  managed_rules -> 'Exclusions' as exclusions,
  managed_rules -> 'ManagedRuleSets' as managed_rule_sets
from
  azure_web_application_firewall_policy;
```

```sql+sqlite
select
  name,
  json_extract(managed_rules, '$.Exclusions') as exclusions,
  json_extract(managed_rules, '$.ManagedRuleSets') as managed_rule_sets
from
  azure_web_application_firewall_policy;
```

### List policies with associated application gateways
Identify WAF policies that are linked to specific application gateways, which can help in managing and securing web traffic.

```sql+postgres
select
  name,
  application_gateways
from
  azure_web_application_firewall_policy
where
  application_gateways is not null;
```

```sql+sqlite
select
  name,
  application_gateways
from
  azure_web_application_firewall_policy
where
  json_extract(application_gateways, '$[0]') is not null;
```

### List policies with HTTP listeners
Fetch WAF policies that are associated with specific HTTP listeners, which can be important for understanding how traffic is being monitored and filtered.

```sql+postgres
select
  name,
  http_listeners
from
  azure_web_application_firewall_policy
where
  http_listeners is not null;
```

```sql+sqlite
select
  name,
  http_listeners
from
  azure_web_application_firewall_policy
where
  json_extract(http_listeners, '$[0]') is not null;
```

### Get application gateway details that is associated with the firewall policy
Get application gateway associated with WAF policies and configurations

```sql+postgres
select
  a.name as application_name,
  a.provisioning_state as application_provisioning_state,
  a.enable_fips,
  a.autoscale_configuration,
  p.name as policy_name,
  p.policy_settings
from
  azure_application_gateway as a
  join azure_web_application_firewall_policy as p on (a.firewall_policy ->> 'id') = p.id;
```

```sql+sqlite
select
  a.name as application_name,
  a.provisioning_state as application_provisioning_state,
  a.enable_fips,
  a.autoscale_configuration,
  p.name as policy_name,
  p.policy_settings
from
  azure_application_gateway as a
  join azure_web_application_firewall_policy as p on json_extract(a.firewall_policy, '$.id') = p.id;
```