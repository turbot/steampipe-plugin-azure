---
title: "Steampipe Table: azure_role_assignment - Query Azure Role Assignments using SQL"
description: "Allows users to query Azure Role Assignments, specifically the details of role assignments for users, groups, and service principals in Azure Active Directory."
folder: "IAM"
---

# Table: azure_role_assignment - Query Azure Role Assignments using SQL

Azure Role Assignments are a critical component of Azure's access control capabilities. They determine what actions a security principal (like a user, group, or service principal) can perform on a specific Azure resource. Each role assignment is a combination of a security principal, a role definition, and a scope.

## Table Usage Guide

The `azure_role_assignment` table provides insights into role assignments within Azure. As a security administrator, you can explore details of role assignments through this table, including the assigned roles, the associated security principals, and the scope of the assignments. Use it to monitor and manage access control within your Azure environment, ensuring that only the appropriate users, groups, or service principals have access to specific resources.

## Examples

### Basic info
Explore which roles are assigned to different principals in your Azure environment. This can help you manage access control and understand who has permissions to what resources, enhancing your security posture.

```sql+postgres
select
  name,
  id,
  principal_id,
  principal_type
from
  azure_role_assignment;
```

```sql+sqlite
select
  name,
  id,
  principal_id,
  principal_type
from
  azure_role_assignment;
```

### List of role assignments which has permission at root level
Discover the segments that are assigned roles with root level access. This is useful for auditing security and access controls in your Azure environment.

```sql+postgres
select
  name,
  id,
  scope
from
  azure_role_assignment
where
  scope = '/';
```

```sql+sqlite
select
  name,
  id,
  scope
from
  azure_role_assignment
where
  scope = '/';
```

### List of role assignments which has subscription level permission and full access to the subscription
This query is useful for identifying roles that have full access permissions at the subscription level within your Azure environment. It helps in maintaining security and managing access by revealing potential over-permissions.

```sql+postgres
select
  ra.name as roll_assignment_name,
  rd.role_name
from
  azure_role_assignment ra
  join azure_role_definition rd on ra.role_definition_id = rd.id
  cross join jsonb_array_elements(rd.permissions) as perm
where
  ra.scope like '/subscriptions/%'
  and perm -> 'actions' = '["*"]';
```

```sql+sqlite
select
  ra.name as roll_assignment_name,
  rd.role_name
from
  azure_role_assignment ra
  join azure_role_definition rd on ra.role_definition_id = rd.id,
  json_each(rd.permissions) as perm
where
  ra.scope like '/subscriptions/%'
  and json_extract(perm.value, '$.actions') = '["*"]';
```

### Role assignments at the resource group scope
Retrieve all role assignments that grant permissions at a specific resource group level.

```sql+postgres
select
  name,
  scope,
  type,
  principal_id.
  principal_type
from
  azure_role_assignment
where
  scope = '/subscriptions/abcdef12-3456-7890-abcd-ef1234567890/resourceGroups/my-rg'
```

```sql+sqlite
select
  name,
  scope,
  type,
  principal_id.
  principal_type
from
  azure_role_assignment
where
  scope = '/subscriptions/abcdef12-3456-7890-abcd-ef1234567890/resourceGroups/my-rg'
```

### Role assignments at the management group scope
List all role assignments that apply at the management group level, determining access across multiple subscriptions.

```sql+postgres
select
  name,
  scope,
  type,
  principal_id.
  principal_type
from
  azure_role_assignment
where
  scope = '/providers/Microsoft.Management/managementGroups/12345678-90ab-cdef-1234-567890abcdef'
```

```sql+sqlite
select
  name,
  scope,
  type,
  principal_id.
  principal_type
from
  azure_role_assignment
where
  scope = '/providers/Microsoft.Management/managementGroups/12345678-90ab-cdef-1234-567890abcdef'
```

### Role assignments at the storage account scope
Identify role assignments that provide permissions at a storage account level, enabling access control for storage resources.

```sql+postgres
select
  name,
  scope,
  type,
  principal_id.
  principal_type
from
  azure_role_assignment
where
  scope = '/subscriptions/abcdef12-3456-7890-abcd-ef1234567890/resourceGroups/nist-test_group/providers/Microsoft.Storage/storageAccounts/testimmutablecontainer'
```

```sql+sqlite
select
  name,
  scope,
  type,
  principal_id.
  principal_type
from
  azure_role_assignment
where
  scope = '/subscriptions/abcdef12-3456-7890-abcd-ef1234567890/resourceGroups/nist-test_group/providers/Microsoft.Storage/storageAccounts/testimmutablecontainer'
```