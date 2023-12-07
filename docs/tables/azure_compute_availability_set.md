---
title: "Steampipe Table: azure_compute_availability_set - Query Azure Compute Availability Sets using SQL"
description: "Allows users to query Azure Compute Availability Sets, providing insights into the availability of resources within an Azure Resource Group."
---

# Table: azure_compute_availability_set - Query Azure Compute Availability Sets using SQL

Azure Compute Availability Sets are a high-availability feature for providing redundant compute resources in Azure. They enable you to ensure that your application is available during planned and unplanned maintenance. Availability Sets are a strategy for achieving high availability and fault tolerance in Azure by ensuring that VM resources are located across multiple isolated hardware nodes in a cluster.

## Table Usage Guide

The `azure_compute_availability_set` table provides insights into the availability sets within Azure Compute. As a system administrator, explore availability set-specific details through this table, including fault domains, update domains, and associated metadata. Use it to uncover information about availability sets, such as those with high fault tolerance and the verification of update policies.

## Examples

### Basic info
Analyze the settings to understand the count of fault domains and update domains within your Azure Compute availability sets across different regions. This is beneficial for managing and optimizing your cloud resources, ensuring balanced workloads and high availability.

```sql+postgres
select
  name,
  platform_fault_domain_count,
  platform_update_domain_count,
  region
from
  azure_compute_availability_set;
```

```sql+sqlite
select
  name,
  platform_fault_domain_count,
  platform_update_domain_count,
  region
from
  azure_compute_availability_set;
```


### List of availability sets which does not use managed disks configuration
Determine the areas in which Azure availability sets are not utilizing the managed disks configuration. This can be useful in identifying potential opportunities for optimization and cost reduction.

```sql+postgres
select
  name,
  sku_name
from
  azure_compute_availability_set
where
  sku_name = 'Classic';
```

```sql+sqlite
select
  name,
  sku_name
from
  azure_compute_availability_set
where
  sku_name = 'Classic';
```


### List of availability sets without application tag key
Explore which Azure Compute availability sets are missing an 'application' tag. This is useful for identifying areas in your infrastructure that may lack important metadata, potentially impacting resource management and organization.

```sql+postgres
select
  name,
  tags
from
  azure_compute_availability_set
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  azure_compute_availability_set
where
  json_extract(tags, '$.application') is null;
```