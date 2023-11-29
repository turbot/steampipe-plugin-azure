---
title: "Steampipe Table: azure_iothub_dps - Query Azure IoT Hub Device Provisioning Services using SQL"
description: "Allows users to query Azure IoT Hub Device Provisioning Services."
---

# Table: azure_iothub_dps - Query Azure IoT Hub Device Provisioning Services using SQL

Azure IoT Hub Device Provisioning Service is a helper service for IoT Hub that enables zero-touch, just-in-time provisioning to the right IoT hub without requiring human intervention, enabling customers to provision millions of devices in a secure and scalable manner. It brings the scalability, security, and reliability of Azure IoT Hub and Device Provisioning Service to your on-premises Internet of Things (IoT) applications. The service supports provisioning of both Azure IoT Edge devices and IoT devices running on other operating systems.

## Table Usage Guide

The 'azure_iothub_dps' table provides insights into Device Provisioning Services within Azure IoT Hub. As a DevOps engineer, explore service-specific details through this table, including the provisioning state, IoT Hub linked with the service, and associated metadata. Utilize it to uncover information about services, such as those with specific provisioning states, the IoT Hubs associated with the services, and the verification of the service operations monitoring level. The schema presents a range of attributes of the IoT Hub Device Provisioning Service for your analysis, like the service name, provisioning state, IoT Hub Device ID, and associated tags.

## Examples

### Basic info
Explore the basic details of your Azure IoT Hub Device Provisioning Services (DPS) to understand their locations and types. This can be useful to manage and organize your IoT devices across different regions.

```sql
select
  name,
  id,
  region,
  type
from
  azure_iothub_dps;
```

### List iot hub dps which are not active
Explore which IoT Hub Device Provisioning Services are not currently active, to potentially identify any issues or areas requiring attention. This could be beneficial in maintaining optimal network performance and avoiding potential service disruptions.

```sql
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