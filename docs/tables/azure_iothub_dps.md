---
title: "Steampipe Table: azure_iothub_dps - Query Azure IoT Hub Device Provisioning Services using SQL"
description: "Allows users to query Azure IoT Hub Device Provisioning Services, offering details about the status, registration, and configuration of each device."
folder: "IoT Hub"
---

# Table: azure_iothub_dps - Query Azure IoT Hub Device Provisioning Services using SQL

Azure IoT Hub Device Provisioning Service is a helper service for IoT Hub that enables zero-touch, just-in-time provisioning to the right IoT hub without requiring human intervention, enabling customers to provision millions of devices in a secure and scalable manner. It provides a seamless, highly scalable way to register and provision IoT devices with an IoT hub. It enables customers to automate the process of registering devices with IoT Hub, reducing the complexity of initial device setup.

## Table Usage Guide

The `azure_iothub_dps` table provides insights into Device Provisioning Services within Azure IoT Hub. As an IoT developer, explore device-specific details through this table, including status, registration, and configuration. Utilize it to uncover information about devices, such as their provisioning status, the IoT hub they are associated with, and the attestation mechanism used.

## Examples

### Basic info
Explore which Azure IoT Hub Device Provisioning Services (DPS) are available and where they are located to better manage and distribute your IoT devices across different regions. This helps in planning and optimizing the distribution of your IoT devices.

```sql+postgres
select
  name,
  id,
  region,
  type
from
  azure_iothub_dps;
```

```sql+sqlite
select
  name,
  id,
  region,
  type
from
  azure_iothub_dps;
```

### List iot hub dps which are not active
Determine the areas in which IoT Hub Device Provisioning Services within Azure are not currently active. This can be beneficial for identifying potential issues or areas of underutilization within your IoT network.

```sql+postgres
select
  name,
  id,
  region,
  type
from
  azure_iothub_dps
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
  azure_iothub_dps
where
  state != 'Active';
```