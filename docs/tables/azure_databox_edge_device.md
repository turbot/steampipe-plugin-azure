---
title: "Steampipe Table: azure_databox_edge_device - Query Azure Databox Edge Devices using SQL"
description: "Allows users to query Azure Databox Edge Devices."
---

# Table: azure_databox_edge_device - Query Azure Databox Edge Devices using SQL

Azure Databox Edge is a physical network appliance, shipped by Microsoft, that brings compute, storage, and intelligence to the edge. It is designed to analyze, transform, and filter data at the edge, before it is transferred to Azure. This device is ideal for locations with limited or no network connectivity, and for reducing data transfer costs.

## Table Usage Guide

The 'azure_databox_edge_device' table provides insights into Databox Edge Devices within Azure. As a DevOps engineer, explore device-specific details through this table, including the device model, status, and associated metadata. Utilize it to uncover information about devices, such as those with high capacity, the network connectivity between devices, and the verification of transfer costs. The schema presents a range of attributes of the Databox Edge Device for your analysis, like the device name, serial number, model description, and associated tags.

## Examples

### Basic info
Explore the status and regional distribution of your Azure Databox Edge devices to gain insights into their operational efficiency and geographical spread. This can help in managing resources and enhancing data access performance.

```sql
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
Discover the segments that are offline in your Azure Data Box Edge devices. This helps in identifying devices that may require attention or troubleshooting for connectivity issues.

```sql
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