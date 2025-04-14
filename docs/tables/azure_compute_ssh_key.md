---
title: "Steampipe Table: azure_compute_ssh_key - Query Azure Compute SSH Keys using SQL"
description: "Allows users to query Azure Compute SSH Keys, providing insights into the SSH keys associated with virtual machines in Azure Compute."
folder: "Compute"
---

# Table: azure_compute_ssh_key - Query Azure Compute SSH Keys using SQL

Azure Compute SSH Key is a resource in Microsoft Azure that allows users to manage SSH keys for virtual machines. These keys are used for secure shell login to VM instances. Azure Compute SSH Key provides a secure way to access VMs without needing to manage passwords.

## Table Usage Guide

The `azure_compute_ssh_key` table enables users to gain insights into the SSH keys associated with their Azure Compute virtual machines. As a system administrator or DevOps engineer, leverage this table to manage and audit SSH keys, ensuring secure and appropriate access to VM instances. This table is beneficial in maintaining security best practices, identifying unused or unnecessary keys, and enforcing compliance with organizational access policies.

## Examples

### Retrieve SSH public key by name
Discover the segments that have specific SSH public keys associated with them in your Azure Compute instances. This helps ensure secure access to your instances by verifying the SSH keys in use.

```sql+postgres
select
  name,
  public_key
from
  azure_compute_ssh_key
where
  name = 'key-name.';
```

```sql+sqlite
select
  name,
  public_key
from
  azure_compute_ssh_key
where
  name = 'key-name.';
```

### List compute virtual machines using SSH public key
The query is used to identify which virtual machines are utilizing a specific SSH public key. This can be useful for security audits, ensuring only authorized keys are in use.

```sql+postgres
select
  m.name as machine_name,
  k.name as ssh_key_name
from
  azure_compute_virtual_machine as m,
  jsonb_array_elements(linux_configuration_ssh_public_keys) as s
  left join azure_compute_ssh_key as k on k.public_key = s ->> 'keyData';
```

```sql+sqlite
select
  m.name as machine_name,
  k.name as ssh_key_name
from
  azure_compute_virtual_machine as m,
  json_each(linux_configuration_ssh_public_keys) as s
  left join azure_compute_ssh_key as k on k.public_key = json_extract(s.value, '$.keyData');
```