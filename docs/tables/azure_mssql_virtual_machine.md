---
title: "Steampipe Table: azure_mssql_virtual_machine - Query Azure SQL Server Virtual Machines using SQL"
description: "Allows users to query Azure SQL Server Virtual Machines, specifically providing insights into the configuration, status, and operational aspects of the SQL Server instances running on Azure Virtual Machines."
folder: "SQL Server"
---

# Table: azure_mssql_virtual_machine - Query Azure SQL Server Virtual Machines using SQL

Azure SQL Server Virtual Machines are a fully managed service that provides the broadest SQL Server engine compatibility and native virtual network (VNET) support. This service offers a set of capabilities for enterprise-grade data workloads, enabling users to run their SQL Server workloads on a virtual machine in Azure. It is an ideal choice for applications requiring OS-level access.

## Table Usage Guide

The `azure_mssql_virtual_machine` table provides insights into SQL Server instances running on Azure Virtual Machines. As a database administrator or a DevOps engineer, explore instance-specific details through this table, including configurations, status, and operational aspects. Utilize it to manage and monitor your SQL Server workloads running on Azure Virtual Machines effectively.

## Examples

### Basic info
Analyze the settings of your Azure SQL virtual machines to gain insights into their current status and configurations. This can help you understand the provisioning state, image offer, license type, and geographical location of each machine, aiding in resource management and optimization.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  sql_image_offer,
  sql_server_license_type,
  region
from
  azure_mssql_virtual_machine;
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state,
  sql_image_offer,
  sql_server_license_type,
  region
from
  azure_mssql_virtual_machine;
```

### List failed virtual machines
Determine the areas in which virtual machines have failed to provision properly within your Azure SQL environment, allowing you to address and rectify these issues promptly.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state
from
  azure_mssql_virtual_machine
where
  provisioning_state = 'Failed';
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state
from
  azure_mssql_virtual_machine
where
  provisioning_state = 'Failed';
```