---
title: "Steampipe Table: azure_storage_account_local_user - Query Azure Storage Account Local Users using SQL"
description: "Allows users to query Azure Storage Account Local Users, providing insights into local user configurations, permissions, and security settings for SFTP and NFS access."
folder: "Storage"
---

# Table: azure_storage_account_local_user - Query Azure Storage Account Local Users using SQL

Azure Storage Account Local Users are user accounts that provide secure access to Azure Storage resources using SFTP and NFS protocols. These local users can be configured with SSH keys or passwords, assigned specific permissions and access scopes, and given granular control over storage resource access. Local users are essential for scenarios requiring secure file transfer protocols (SFTP) for Azure Blob Storage or NFS access to Azure Files.

## Table Usage Guide

The `azure_storage_account_local_user` table provides insights into local users within Azure Storage Accounts. As a storage administrator, security engineer, or compliance officer, explore user-specific details through this table, including authentication methods, home directories, permission scopes, and access configurations. Utilize it to manage and monitor local user access, verify security settings, ensure proper SSH key configurations, and maintain compliance with organizational security policies.

## Examples

### Basic info
Explore the basic information about local users in Azure storage accounts to understand their authentication methods and security configurations.

```sql+postgres
select
  name,
  storage_account_name,
  resource_group,
  has_ssh_key,
  has_ssh_password,
  has_shared_key,
  home_directory
from
  azure_storage_account_local_user;
```

```sql+sqlite
select
  name,
  storage_account_name,
  resource_group,
  has_ssh_key,
  has_ssh_password,
  has_shared_key,
  home_directory
from
  azure_storage_account_local_user;
```

### List local users with SSH key access and their authorized keys
Identify local users who have SSH key access configured and examine their authorized keys, which is crucial for security auditing and access management.

```sql+postgres
select
  name,
  storage_account_name,
  resource_group,
  has_ssh_key,
  ssh_authorized_keys,
  home_directory
from
  azure_storage_account_local_user
where
  has_ssh_key = true;
```

```sql+sqlite
select
  name,
  storage_account_name,
  resource_group,
  has_ssh_key,
  ssh_authorized_keys,
  home_directory
from
  azure_storage_account_local_user
where
  has_ssh_key = 1;
```

### List local users with NFSv3 access enabled and their permissions
Determine which local users have NFSv3 access enabled and examine their access permissions, helping to manage file sharing capabilities and access control.

```sql+postgres
select
  name,
  storage_account_name,
  resource_group,
  is_nfsv3_enabled,
  home_directory,
  permission_scopes,
  user_id,
  group_id
from
  azure_storage_account_local_user
where
  is_nfsv3_enabled = true;
```

```sql+sqlite
select
  name,
  storage_account_name,
  resource_group,
  is_nfsv3_enabled,
  home_directory,
  permission_scopes,
  user_id,
  group_id
from
  azure_storage_account_local_user
where
  is_nfsv3_enabled = 1;
```

### Analyze permission scopes for all local users
Examine the detailed permission configurations for all local users to understand their access levels and restrictions across different storage services.

```sql+postgres
select
  name,
  storage_account_name,
  resource_group,
  jsonb_pretty(permission_scopes) as formatted_permission_scopes,
  home_directory
from
  azure_storage_account_local_user;
```

```sql+sqlite
select
  name,
  storage_account_name,
  resource_group,
  permission_scopes,
  home_directory
from
  azure_storage_account_local_user;
```

### Find local users with extended group memberships
Identify local users who have additional group memberships configured, which is useful for understanding access patterns and group-based permissions.

```sql+postgres
select
  name,
  storage_account_name,
  group_id,
  jsonb_pretty(extended_groups) as group_memberships,
  allow_acl_authorization
from
  azure_storage_account_local_user
where
  extended_groups is not null;
```

```sql+sqlite
select
  name,
  storage_account_name,
  group_id,
  extended_groups as group_memberships,
  allow_acl_authorization
from
  azure_storage_account_local_user
where
  extended_groups is not null;
```

### List local users with password-based authentication
Identify local users configured with password authentication, which may require additional security review in compliance-sensitive environments.

```sql+postgres
select
  name,
  storage_account_name,
  resource_group,
  has_ssh_password,
  has_shared_key,
  allow_acl_authorization
from
  azure_storage_account_local_user
where
  has_ssh_password = true;
```

```sql+sqlite
select
  name,
  storage_account_name,
  resource_group,
  has_ssh_password,
  has_shared_key,
  allow_acl_authorization
from
  azure_storage_account_local_user
where
  has_ssh_password = 1;
``` 