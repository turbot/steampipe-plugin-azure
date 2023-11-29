---
title: "Steampipe Table: azure_iothub - Query Azure IoT Hub using SQL"
description: "Allows users to query Azure IoT Hubs"
---

# Table: azure_iothub - Query Azure IoT Hub using SQL

Azure IoT Hub is a managed service, hosted in the cloud, that acts as a central message hub for bi-directional communication between your IoT application and the devices it manages. You can use Azure IoT Hub to build IoT solutions with reliable and secure communications between millions of IoT devices and a cloud-hosted solution backend. It supports communications both from the device to the cloud and from the cloud to the device.

## Table Usage Guide

The 'azure_iothub' table provides insights into IoT Hubs within Azure IoT Hub. As a DevOps engineer, explore hub-specific details through this table, including the status, SKU, tier, and associated metadata. Utilize it to uncover information about IoT Hubs, such as their location, the number of devices connected, and the verification of their properties. The schema presents a range of attributes of the IoT Hub for your analysis, like the resource group, subscription ID, public network access, and associated tags.

## Examples

### Basic info
Analyze the settings of your Azure IoT Hub to understand its geographical distribution and types. This can help in managing resources and improving the efficiency of IoT devices across different regions.

```sql
select
  name,
  id,
  region,
  type
from
  azure_iothub;
```

### List hubs which are not active
Determine the areas in which inactive IoT hubs exist within the Azure platform. This can be beneficial in identifying potential issues or inefficiencies related to unused resources.

```sql
select
  name,
  id,
  region,
  type
from
  azure_iothub
where
  state <> 'Active';
```