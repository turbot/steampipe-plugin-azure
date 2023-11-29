---
title: "Steampipe Table: azure_compute_ssh_key - Query Azure Compute SSH Keys using SQL"
description: "Allows users to query Azure Compute SSH Keys"
---

# Table: azure_compute_ssh_key - Query Azure Compute SSH Keys using SQL

Azure Compute is a service within Microsoft Azure that provides on-demand processing power and infrastructure for applications. It allows you to create and manage virtual machines, containers, and batch jobs, as well as supports remote application access via SSH keys. Azure Compute SSH Keys are used for secure, encrypted connections to your Azure resources.

## Table Usage Guide

The 'azure_compute_ssh_key' table provides insights into SSH Keys within Azure Compute. As a DevOps engineer, explore SSH Key-specific details through this table, including the associated virtual machine, key type, and key data. Utilize it to uncover information about SSH Keys, such as those associated with specific virtual machines, the type of SSH Key being used, and the actual key data for verification purposes. The schema presents a range of attributes of the SSH Key for your analysis, like the virtual machine id, key type, and key data.

## Examples

### Retrieve SSH public key by name
Assess the elements within your Azure Compute resources to identify a specific SSH public key associated with a given name. This can help in verifying access permissions or troubleshooting connectivity issues.

```sql
select
  name,
  public_key
from
  azure_compute_ssh_key
where
  name = 'key-name.';
```

### List compute virtual machines using SSH public key
Explore which virtual machines are using a particular SSH public key. This is useful for managing and securing your virtual machine access by keeping track of the SSH keys in use.

```sql
select
  m.name as machine_name,
  k.name as ssh_key_name
from
  azure_compute_virtual_machine as m,
  jsonb_array_elements(linux_configuration_ssh_public_keys) as s
  left join azure_compute_ssh_key as k on k.public_key = s ->> 'keyData';
```