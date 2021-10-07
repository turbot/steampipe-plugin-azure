# Table: azure_hdinsight_cluster

Azure HDInsight is a managed, full-spectrum, open-source analytics service in the cloud for enterprises. You can use open-source frameworks such as Hadoop, Apache Spark, Apache Hive, LLAP, Apache Kafka, Apache Storm, R, and more.

## Examples

### Basic info

```sql
select
  name,
  id,
  provisioning_state,
  type,
  cluster_hdp_version,
  cluster_id,
  cluster_state,
  cluster_version,
  created_date
from
  azure_hdinsight_cluster;
```

### List clusters with encryption in transit enabled

```sql
select
  name,
  id,
  encryption_in_transit_properties -> 'isEncryptionInTransitEnabled' as is_encryption_in_transit_enabled
from
  azure_hdinsight_cluster
where
  (encryption_in_transit_properties ->> 'isEncryptionInTransitEnabled')::boolean;
```

### List disk encryption details

```sql
select
  name,
  id,
  disk_encryption_properties ->> 'encryptionAlgorithm' as encryption_algorithm,
  disk_encryption_properties -> 'encryptionAtHost' as encryption_at_host,
  disk_encryption_properties ->> 'keyName' as key_name,
  disk_encryption_properties ->> 'keyVersion' as key_version,
  disk_encryption_properties ->> 'msiResourceId' as msi_resource_id,
  disk_encryption_properties ->> 'vaultUri' as vault_uri 
from
  azure_hdinsight_cluster;
```

### List connectivity endpoint details

```sql
select
  name,
  id,
  endpoint ->> 'location' as endpoint_location,
  endpoint ->> 'name' as endpoint_name,
  endpoint -> 'port' as endpoint_port,
  endpoint ->> 'protocol' as endpoint_protocol,
  endpoint ->> 'privateIpAddress' as endpoint_private_ip_address
from
  azure_hdinsight_cluster,
  jsonb_array_elements(connectivity_endpoints) as endpoint;
```
