---
title: "Steampipe Table: azure_mssql_virtual_machine - Query Azure SQL Server Virtual Machines using SQL"
description: "Allows users to query Azure SQL Server Virtual Machines."
---

# Table: azure_mssql_virtual_machine - Query Azure SQL Server Virtual Machines using SQL

Azure SQL Server Virtual Machine is a service that offers the full control and features of a fully managed SQL Server instance. It provides the flexibility to choose the version, edition, and OS of SQL Server. You can also manage the VM size to meet your performance requirements. 

## Table Usage Guide

The 'azure_mssql_virtual_machine' table provides insights into SQL Server Virtual Machines within Microsoft Azure. As a DevOps engineer, explore VM-specific details through this table, including the SQL Server version, edition, OS type, VM size, and associated metadata. Utilize it to uncover information about virtual machines, such as those with specific SQL Server versions or OS types, and the verification of SQL Server configurations. The schema presents a range of attributes of the SQL Server Virtual Machine for your analysis, like the VM ID, resource group, location, SQL Server license type, and associated tags.

## Examples

### Basic info
Explore the configuration and status of your Azure SQL virtual machines. This query is useful for gaining insights into the types of SQL images and licenses in use, as well as where these resources are located geographically.

```sql
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
Explore which virtual machines have failed to provision in your Azure MSSQL environment, helping you to identify potential issues and take corrective action.

```sql
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