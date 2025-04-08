---
title: "Steampipe Table: azure_iothub - Query Azure IoT Hub using SQL"
description: "Allows users to query Azure IoT Hubs, providing insights into the configurations, properties, and status of each IoT hub deployed in Azure."
folder: "IoT Hub"
---

# Table: azure_iothub - Query Azure IoT Hub using SQL

Azure IoT Hub is a managed service, hosted in the cloud, that acts as a central message hub for bi-directional communication between your IoT application and the devices it manages. It provides reliable and secure communication between millions of IoT devices and a cloud-hosted solution backend. Azure IoT Hub supports communications both from the device to the cloud and from the cloud to the device.

## Table Usage Guide

The `azure_iothub` table provides insights into IoT Hubs within Azure. As a IoT developer or cloud solutions architect, explore IoT Hub-specific details through this table, including configurations, properties, and status. Utilize it to uncover information about each IoT Hub, such as its SKU, location, and the number of devices it can support, to ensure optimal performance and resource allocation.

## Examples

### Basic info
Explore the basic characteristics of your Azure IoT Hub resources, such as their names, IDs, regions, and types. This can help you manage your resources more effectively by understanding their distribution and categorization.

```sql+postgres
select
  name,
  id,
  region,
  type
from
  azure_iothub;
```

```sql+sqlite
select
  name,
  id,
  region,
  type
from
  azure_iothub;
```

### List hubs which are not active
Determine the areas in which certain hubs within the Azure IoT platform are not currently active. This can help prioritize troubleshooting efforts or identify opportunities for resource optimization.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  region,
  type
from
  azure_iothub
where
  state != 'Active';
```