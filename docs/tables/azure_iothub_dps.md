# Table: azure_iothub_dps

The IoT Hub Device Provisioning Service (DPS) is a helper service for IoT Hub that enables zero-touch, just-in-time provisioning to the right IoT hub without requiring human intervention, allowing customers to provision millions of devices in a secure and scalable manner. 

## Examples

### Basic info

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
