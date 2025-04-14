---
title: "Steampipe Table: azure_resource - Query Azure Resources using SQL"
description: "Allows users to query Azure resources, providing insights into resource properties, identities, tags, and more."
folder: "Resource Manager"
---

# Table: azure_resource - Query Azure Resources using SQL

Azure resources are entities managed by Azure, such as virtual machines, storage accounts, and network interfaces. This table allows you to query various details about these resources, including their properties, identities, tags, and more.

## Table Usage Guide

The `azure_resource` table provides insights into resources within your Azure environment. As an Azure administrator, you can explore resource-specific details through this table, including resource IDs, names, types, provisioning states, and associated tags. Utilize it to manage and audit resources, verify configurations, and enforce governance policies.

**Important notes:**
- For improved performance, it is advised that you use the optional qual to limit the result set. Optional quals are supported for the following columns:
  - `region`
  - `type`
  - `name`
  - `identity_principal_id`
  - `plan_publisher`
  - `plan_name`
  - `plan_product`
  - `plan_promotion_code`
  - `plan_version`
  - `resource_group`
  - `filter`

## Examples

### Basic Information
Retrieve basic information about all Azure resources, including their names, types, and regions.

```sql+postgres
select
  name,
  type,
  region
from
  azure_resource;
```

```sql+sqlite
select
  name,
  type,
  region
from
  azure_resource;
```

### Filter by resource type
Get a list of all virtual networks.

```sql+postgres
select
  name,
  type,
  region
from
  azure_resource
where
  filter = 'resourceType eq ''Microsoft.Network/virtualNetworks''';
```

```sql+sqlite
select
  name,
  type,
  region
from
  azure_resource
where
  filter = 'resourceType eq ''Microsoft.Network/virtualNetworks''';
```

### Retrieve resources with specific tag key
Fetch resources that have specific tags associated with them.

```sql+postgres
select
  name,
  type,
  tags
from
  azure_resource
where
  filter = 'startswith(tagName, ''cost'')';
```

```sql+sqlite
select
  name,
  type,
  tags
from
  azure_resource
where
  filter = 'startswith(tagName, ''cost'')';
```

### Get Resource identity details
List resources along with their identity principal IDs and types.

```sql+postgres
select
  name,
  type,
  identity_principal_id,
  identity ->> 'type' as identity_type,
  identify -> 'userAssignedIdentities' as user_assigned_identities
from
  azure_resource;
```

```sql+sqlite
select
  name,
  type,
  identity_principal_id,
  json_extract(identity, '$.type') as identity_type,
  json_extract(identity, '$.userAssignedIdentities') as user_assigned_identities
from
  azure_resource;
```

### List managed resources
Identify resources that are managed by other resources.

```sql+postgres
select
  name,
  type,
  managed_by
from
  azure_resource
where
  managed_by is not null;
```

```sql+sqlite
select
  name,
  type,
  managed_by
from
  azure_resource
where
  managed_by is not null;
```