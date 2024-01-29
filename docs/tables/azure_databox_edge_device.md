---
title: "Steampipe Table: azure_databox_edge_device - Query Azure Databox Edge Devices using SQL"
description: "Allows users to query Azure Databox Edge Devices, providing insights into the device's status, SKU, model description, and more."
---

# Table: azure_databox_edge_device - Query Azure Databox Edge Devices using SQL

Azure Databox Edge is a physical network appliance, shipped by Microsoft, that brings computation and storage capabilities to the edge of your network. It acts as a storage gateway, creating a link between your site and Azure storage. This device provides AI-enabled edge compute, network and storage capabilities.

## Table Usage Guide

The `azure_databox_edge_device` table provides insights into Azure Databox Edge Devices within Microsoft Azure. As an IT administrator, explore device-specific details through this table, including the device's status, SKU, model description, and more. Utilize it to uncover information about the devices, such as their capacity, serial numbers, and the verification of device settings.

## Examples

### Basic info
Explore the status and geographical distribution of your Azure Databox Edge devices. This allows for efficient management and monitoring of your devices across different regions.

```sql+postgres
select
  name,
  id,
  type,
  data_box_edge_device_status,
  region
from
  azure_databox_edge_device;
```

```sql+sqlite
select
  name,
  id,
  type,
  data_box_edge_device_status,
  region
from
  azure_databox_edge_device;
```

### List offline data box edge devices
Determine the areas in which Azure Databox Edge devices are currently offline. This can be useful for identifying potential network issues or maintenance needs in your infrastructure.

```sql+postgres
select
  name,
  id,
  type,
  data_box_edge_device_status
from
  azure_databox_edge_device
where
  data_box_edge_device_status = 'Offline';
```

```sql+sqlite
select
  name,
  id,
  type,
  data_box_edge_device_status
from
  azure_databox_edge_device
where
  data_box_edge_device_status = 'Offline';
```