# Table: azure_iothub

Azure IoT Hub is Microsoft’s Internet of Things connector to the cloud. It’s a fully managed cloud service that enables reliable and secure bi-directional communications between millions of IoT devices and a solution back end.

## Examples

### Basic info

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
