---
title: "Steampipe Table: azure_tenant - Query Azure Tenants using SQL"
description: "Allows users to query Azure Tenants for comprehensive data on the tenant's details, including tenant ID, domains, and display name."
---

# Table: azure_tenant - Query Azure Tenants using SQL

Azure Tenant refers to an organization's dedicated and isolated instance of Microsoft Azure that is automatically created when an organization signs up for a Microsoft cloud service subscription. Azure Tenants serve as dedicated, isolated containers for all of an organization's Azure resources, and provide a secure environment where an organization can store and manage its resources.

## Table Usage Guide

The 'azure_tenant' table provides insights into Azure Tenants within Microsoft Azure. As a DevOps engineer, explore tenant-specific details through this table, including tenant ID, domains, and display name. Utilize it to uncover information about tenants, such as those with specific domains, the tenant's display name, and the verification of tenant IDs. The schema presents a range of attributes of the Azure Tenant for your analysis, like the tenant ID, domains, and display name.

## Examples

### Basic info
Explore the basic details of your Azure tenants, including their names, IDs, categories, locations, and associated domains. This can be useful for gaining a high-level overview of your Azure environment, and for identifying areas for potential optimization or consolidation.

```sql
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