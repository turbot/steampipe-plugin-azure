# Table: azure_databox_edge_device

Azure Data Box Gateway is a storage solution that enables you to seamlessly send data to Azure. This article provides you an overview of the Azure Data Box Gateway solution, benefits, key capabilities, and the scenarios where you can deploy this device.

## Examples

### Basic info

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
