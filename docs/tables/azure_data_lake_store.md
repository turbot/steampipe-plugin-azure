---
title: "Steampipe Table: azure_data_lake_store - Query Azure Data Lake Store using SQL"
description: "Allows users to query Azure Data Lake Stores, providing insights into the data storage and analytics service in Azure."
---

# Table: azure_data_lake_store - Query Azure Data Lake Store using SQL

Azure Data Lake Store is a hyper-scale repository for big data analytic workloads in Azure. It allows you to store and analyze petabyte-size files and trillions of objects. Azure Data Lake Store offers high-speed integration with Azure HDInsight, Azure Data Factory, and Azure Machine Learning.

## Table Usage Guide

The `azure_data_lake_store` table provides insights into the data storage and analytics service in Azure. As a data engineer or data scientist, explore details about your Azure Data Lake Stores through this table, including their properties, encryption settings, and firewalls rules. Utilize it to manage and monitor your data lake stores, ensuring they are configured according to your organization's security and compliance policies.

## Examples

### Basic info
Explore the general information of your Azure Data Lake Store resources to understand their current state and type. This can help in monitoring the provisioning status and managing these resources effectively.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store;
```

### List data lake stores with encryption disabled
Explore which Azure data lake stores have disabled encryption, a potential security risk. This can be useful in auditing and improving your data security measures.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  encryption_state = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  encryption_state = 'Disabled';
```

### List data lake stores with firewall disabled
Discover the segments that utilize Azure Data Lake stores with disabled firewalls, enabling you to identify potential security risks and take necessary precautions. This is particularly useful for ensuring optimal security measures are in place and avoiding potential data breaches.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  firewall_state = 'Disabled';
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_data_lake_store
where
  firewall_state = 'Disabled';
```