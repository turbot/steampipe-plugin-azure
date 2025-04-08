---
title: "Steampipe Table: azure_security_center_jit_network_access_policy - Query Azure Security Center Just-In-Time Network Access Policies using SQL"
description: "Allows users to query Just-In-Time Network Access Policies in Azure Security Center, providing insights into policy configurations and associated resources."
folder: "Policy"
---

# Table: azure_security_center_jit_network_access_policy - Query Azure Security Center Just-In-Time Network Access Policies using SQL

Azure Security Center Just-In-Time Network Access Policies are resources within Microsoft Azure that provide controlled access to Azure VMs. They reduce exposure to attacks by enabling access to VMs only when needed and from specific, approved IP addresses. Azure JIT Network Access Policies help maintain a secure environment by minimizing the potential attack surface.

## Table Usage Guide

The `azure_security_center_jit_network_access_policy` table provides insights into Just-In-Time Network Access Policies within Azure Security Center. As a security analyst, you can explore policy-specific details through this table, including policy configurations, associated resources, and access controls. Utilize it to uncover information about policies, such as their status, provisioned locations, and the resources they are associated with.

## Examples

### List virtual machines with JIT access enabled
The query is useful for identifying virtual machines that have Just-In-Time (JIT) access enabled, a feature that can help enhance security by limiting open ports. This can be particularly helpful in managing security risks and ensuring that only necessary access points are open.

```sql+postgres
select
  vm.name,
  vm.id,
  jsonb_pretty(vms -> 'ports') as ports
from
  azure_security_center_jit_network_access_policy,
  jsonb_array_elements(virtual_machines) as vms,
  azure_compute_virtual_machine as vm
where
  lower(vms ->> 'id') = lower(vm.id);
```

```sql+sqlite
select
  vm.name,
  vm.id,
  vms.value as ports
from
  azure_security_center_jit_network_access_policy,
  json_each(virtual_machines) as vms,
  azure_compute_virtual_machine as vm
where
  lower(json_extract(vms.value, '$.id')) = lower(vm.id);
```