---
title: "Steampipe Table: azure_lighthouse_definition - Query Azure Lighthouse Definitions using SQL"
description: "Allows users to query Azure Lighthouse definitions, providing insights into cross-tenant management and governance."
---

# Table: azure_lighthouse_definition - Query Azure Lighthouse Definitions using SQL

Azure Lighthouse is a service within Microsoft Azure that enables cross-tenant management, allowing service providers to manage resources across multiple tenants while maintaining control and visibility. Azure Lighthouse provides greater automation, scalability, and enhanced governance across resources.

## Table Usage Guide

The `azure_lighthouse_definition` table provides insights into Azure Lighthouse definitions. As a Network Administrator or Service Provider, you can explore details about each Lighthouse definition, including its configuration, associated resources, and authorization details. Use this table to ensure your cross-tenant management and governance are correctly configured and to quickly identify any potential issues.

## Examples

### Basic info
Explore the status and details of your Azure Lighthouse definitions to understand their current state and type. This is beneficial for auditing and managing your cross-tenant resources effectively.

```sql+postgres
select
  name,
  id,
  managed_by_tenant_id,
  managed_by_tenant_name,
  managed_tenant_name,
  type
from
  azure_lighthouse_definition;
```

```sql+sqlite
select
  name,
  id,
  managed_by_tenant_id,
  managed_by_tenant_name,
  managed_tenant_name,
  type
from
  azure_lighthouse_definition;
```

### List authorization details for each Lighthouse definition
Identify the authorization details linked with each Lighthouse definition. This can help in managing access control and understanding the roles assigned to different Azure Active Directory principals.

```sql+postgres
select
  name,
  a ->> 'principalId' as principal_id,
  a ->> 'roleDefinitionId' as role_definition_id,
  a ->> 'principalIdDisplayName' as principal_id_display_name,
  a -> 'delegatedRoleDefinitionIds' as delegated_role_definition_i_ds
from
  azure_lighthouse_definition,
  jsonb_array_elements(authorizations) as a;
```

```sql+sqlite
select
  name,
  json_extract(a.value, '$.principalId') as principal_id,
  json_extract(a.value, '$.roleDefinitionId') as role_definition_id,
  json_extract(a.value, '$.principalIdDisplayName') as principal_id_display_name,
  json_extract(a.value, '$.delegatedRoleDefinitionIds') as delegated_role_definition_i_ds
from
  azure_lighthouse_definition,
  json_each(authorizations) as a;
```

### List eligible authorization details for each Lighthouse definition
Explore the eligible authorization details associated with each Lighthouse definition. This helps in understanding the just-in-time access Azure Active Directory principals will receive on the delegated resources.

```sql+postgres
select
  name,
  a ->> 'principalId' as principal_id,
  a ->> 'roleDefinitionId' as role_definition_id,
  a ->> 'principalIdDisplayName' as principal_id_display_name,
  a -> 'justInTimeAccessPolicy' as just_in_time_access_policy
from
  azure_lighthouse_definition,
  jsonb_array_elements(eligible_authorizations) as a;
```

```sql+sqlite
select
  name,
  json_extract(a.value, '$.principalId') as principal_id,
  json_extract(a.value, '$.roleDefinitionId') as role_definition_id,
  json_extract(a.value, '$.principalIdDisplayName') as principal_id_display_name,
  json_extract(a.value, '$.justInTimeAccessPolicy') as just_in_time_access_policy
from
  azure_lighthouse_definition,
  json_each(eligible_authorizations) as a;
```

### List plan details for each Lighthouse definition
Get an overview of the plan details for the managed services associated with each Lighthouse definition. This can assist in understanding the service plans and ensuring they align with your management requirements.

```sql+postgres
select
  name,
  plan ->> 'name' as plan_name,
  plan ->> 'product' as plan_product,
  plan ->> 'publisher' as plan_publisher,
  plan ->> 'version' as plan_version
from
  azure_lighthouse_definition;
```

```sql+sqlite
select
  name,
  json_extract(plan, '$.name') as plan_name,
  json_extract(plan, '$.product') as plan_product,
  json_extract(plan, '$.publisher') as plan_publisher,
  json_extract(plan, '$.version') as plan_version
from
  azure_lighthouse_definition;
```
