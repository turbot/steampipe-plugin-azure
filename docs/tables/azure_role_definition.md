---
title: "Steampipe Table: azure_role_definition - Query Azure Role Definitions using SQL"
description: "Allows users to query Role Definitions in Azure, specifically the permissions, trust policies, and associated metadata, providing insights into role-specific details."
folder: "IAM"
---

# Table: azure_role_definition - Query Azure Role Definitions using SQL

Azure Role Definition is a resource within Microsoft Azure that represents a collection of permissions. It's used to provide access to Azure resources that the role is assigned to. Role Definitions help you manage access to your Azure resources by providing a way to group together permissions into roles.

## Table Usage Guide

The `azure_role_definition` table provides insights into Role Definitions within Microsoft Azure. As a DevOps engineer, explore role-specific details through this table, including permissions, trust policies, and associated metadata. Utilize it to manage access to your Azure resources, group together permissions into roles, and gain insights into role-specific details.

## Examples

### List the custom roles
Explore which custom roles have been defined in your Azure environment. This is beneficial to understand and manage the unique permissions and restrictions applied within your system.

```sql+postgres
select
  name,
  description,
  role_name,
  role_type,
  title
from
  azure_role_definition
where
  role_type = 'CustomRole';
```

```sql+sqlite
select
  name,
  description,
  role_name,
  role_type,
  title
from
  azure_role_definition
where
  role_type = 'CustomRole';
```

### List of roles whose assignable scope is set to root('/') scope
Explore the roles in your Azure environment that have been given broad permissions, as indicated by their assignable scope being set to root. This can be useful for identifying potential security risks and ensuring that permissions are appropriately restricted.

```sql+postgres
select
  name,
  role_name,
  scope
from
  azure_role_definition,
  jsonb_array_elements_text(assignable_scopes) as scope
where
  scope = '/';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```
### Permissions of all custom roles
Explore which permissions are assigned to all custom roles within your Azure environment. This can help in maintaining security standards and ensuring that roles are not granted excessive permissions.

```sql+postgres
select
  name,
  role_name,
  role_type,
  permission -> 'actions' as action,
  permission -> 'dataActions' as data_action,
  permission -> 'notActions' as no_action,
  permission -> 'notDataActions' as not_data_actions
from
  azure_role_definition
  cross join jsonb_array_elements(permissions) as permission
where
  role_type = 'CustomRole';
```

```sql+sqlite
select
  name,
  role_name,
  role_type,
  json_extract(permission.value, '$.actions') as action,
  json_extract(permission.value, '$.dataActions') as data_action,
  json_extract(permission.value, '$.notActions') as no_action,
  json_extract(permission.value, '$.notDataActions') as not_data_actions
from
  azure_role_definition,
  json_each(permissions) as permission
where
  role_type = 'CustomRole';
```

### Permissions of all custom roles
Explore the permissions assigned to all custom roles in your Azure environment. This can help you understand access controls and identify potential security risks.

```sql+postgres
select
  name,
  role_name,
  role_type,
  permission -> 'actions' as action,
  permission -> 'dataActions' as data_action,
  permission -> 'notActions' as no_action,
  permission -> 'notDataActions' as not_data_actions
from
  azure_role_definition
  cross join jsonb_array_elements(permissions) as permission
where
  role_type = 'CustomRole';
```

```sql+sqlite
select
  ard.name,
  ard.role_name,
  ard.role_type,
  json_extract(permission.value, '$.actions') as action,
  json_extract(permission.value, '$.dataActions') as data_action,
  json_extract(permission.value, '$.notActions') as no_action,
  json_extract(permission.value, '$.notDataActions') as not_data_actions
from
  azure_role_definition ard,
  json_each(ard.permissions) as permission
where
  ard.role_type = 'CustomRole';
```

### Permissions of all custom roles
Analyze the permissions assigned to all custom roles in your Azure environment. This can help in identifying roles with excessive permissions, thereby assisting in maintaining a principle of least privilege.

```sql_postgres
select
  name,
  role_name,
  role_type,
  permission -> 'actions' as action,
  permission -> 'dataActions' as data_action,
  permission -> 'notActions' as no_action,
  permission -> 'notDataActions' as not_data_actions
from
  azure_role_definition
  cross join jsonb_array_elements(permissions) as permission
where
  role_type = 'CustomRole';
```

```sql+sqlite
select
  ard.name,
  ard.role_name,
  ard.role_type,
  json_extract(permission.value, '$.actions') as action,
  json_extract(permission.value, '$.dataActions') as data_action,
  json_extract(permission.value, '$.notActions') as no_action,
  json_extract(permission.value, '$.notDataActions') as not_data_actions
from
  azure_role_definition ard,
  json_each(ard.permissions) as permission
where
  ard.role_type = 'CustomRole';
```