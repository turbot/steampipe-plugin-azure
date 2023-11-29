---
title: "Steampipe Table: azure_security_center_jit_network_access_policy - Query Azure Security Center Just-In-Time Network Access Policies using SQL"
description: "Allows users to query Just-In-Time Network Access Policies in Azure Security Center."
---

# Table: azure_security_center_jit_network_access_policy - Query Azure Security Center Just-In-Time Network Access Policies using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers and provides advanced threat protection across your hybrid workloads in the cloud. Just-In-Time Network Access Policies in Azure Security Center help you control access to your Azure Virtual Machines by providing a secure way to connect to a VM, reducing exposure to attacks while providing easy access to connect to VMs when needed.

## Table Usage Guide

The 'azure_security_center_jit_network_access_policy' table provides insights into Just-In-Time Network Access Policies within Azure Security Center. As a security engineer, explore policy-specific details through this table, including policy configurations, virtual machine details, and associated metadata. Utilize it to uncover information about policies, such as those with specific IP configurations, the access protocols allowed, and the verification of request status. The schema presents a range of attributes of the Just-In-Time Network Access Policy for your analysis, like the policy ID, provisioning state, location, and associated tags.

## Examples

### List virtual machines with JIT access enabled
Explore which virtual machines have Just-In-Time access enabled. This is particularly beneficial for enhancing security measures by only permitting access when needed.

```sql
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