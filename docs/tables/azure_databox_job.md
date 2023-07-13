# Table: azure_databox_job

Azure Data Box is a family of data transfer and migration services offered by Microsoft Azure. Azure Data Box Job is one of the services in this family. It is a physical appliance provided by Azure that enables offline data transfer between on-premises data sources and Azure storage.

With Azure Data Box Job, you can securely transfer large amounts of data to Azure by shipping a physical device (Data Box) to your location. You can then copy your data to the device and ship it back to Azure for upload. This offline transfer method is useful when you have a large volume of data that would take a significant amount of time to transfer over the network.

Azure Data Box Job offers a simple and efficient way to migrate or transfer data to Azure, especially in scenarios where network bandwidth is limited or unreliable.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  status,
  region,
  sku_Family,
  is_deletable
from
  azure_databox_job;
```

### List deletable databox job

```sql
select
  name,
  id,
  type,
  transfer_type,
  is_deletable,
  is_shipping_address_editable
from
  azure_databox_job
where is_deletable;
```

### List cancellable databox job

```sql
select
  name,
  id,
  type,
  transfer_type,
  is_deletable,
  is_cancellable
from
  azure_databox_job
where is_cancellable;
```

### List jobs that are scheduled for delivery

```sql
select
  name,
  id,
  type,
  sku_name,
  delivery_type
from
  azure_databox_job
where
  delivery_type = 'Scheduled';
```

### Get error details of each job

```sql
select
  name,
  id,
  error ->> 'code' as error_code,
  error ->> 'Message' as error_message,
  error ->> 'Target' as error_target,
  error ->> 'Details' as error_details,
  error ->> 'AdditionalInfo' as additional_info
from
  azure_databox_job;
```

### List jobs that are cancellable without any fee

```sql
select
  name,
  id,
  start_time,
  transfer_type,
  is_cancellable_without_fee
from
  azure_databox_job
where
  is_cancellable_without_fee;
```