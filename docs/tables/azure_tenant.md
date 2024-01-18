---
title: "Steampipe Table: azure_tenant - Query Azure Tenants using SQL"
description: "Allows users to query Azure Tenants, providing insights into the organizations associated with the Azure subscriptions."
---

# Table: azure_tenant - Query Azure Tenants using SQL

Azure Tenants represent an organization in Azure. Each organization has at least one tenant, and each Azure subscription is associated with a tenant. Tenants are used to manage access to Azure resources.

## Table Usage Guide

The `azure_tenant` table provides insights into the organizations associated with Azure subscriptions. As a Cloud Administrator, you can use this table to explore details such as tenant IDs and domains. This information can be useful for managing access to Azure resources and for understanding the organizational structure of your Azure subscriptions.

## Examples

### Basic info
Discover the segments that are part of your Azure tenant, including their geographical location and associated domains. This is useful for understanding the distribution and categorization of your Azure resources.

```sql+postgres
select
  name,
  id,
  tenant_id,
  tenant_category,
  country,
  country_code,
  display_name,
  domains
from
  azure_tenant;
```

```sql+sqlite
select
  name,
  id,
  tenant_id,
  tenant_category,
  country,
  country_code,
  display_name,
  domains
from
  azure_tenant;
```