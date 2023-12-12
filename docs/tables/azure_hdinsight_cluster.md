---
title: "Steampipe Table: azure_hdinsight_cluster - Query Azure HDInsight Clusters using SQL"
description: "Allows users to query Azure HDInsight Clusters, providing insights into the configurations, properties, and states of these clusters."
---

# Table: azure_hdinsight_cluster - Query Azure HDInsight Clusters using SQL

Azure HDInsight is a fully managed, open-source analytics service for enterprises. It provides big data cloud offerings and is built on Hadoop, Spark, R, and Hive, among others. It enables processing massive amounts of data and running big data workloads in the cloud.

## Table Usage Guide

The `azure_hdinsight_cluster` table provides insights into HDInsight clusters within Azure. As a data engineer or data scientist, you can use this table to explore cluster-specific details, including properties, configurations, and states. This table can be utilized to uncover information about clusters, such as their health, location, provisioning state, and type.

## Examples

### Basic info
Determine the status and details of your Azure HDInsight clusters to manage and optimize your big data analytics. This query is useful in understanding the configuration and operational state of your clusters, including their version details and creation date.

```sql+postgres
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

```sql+sqlite
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
Determine the clusters that have enhanced security measures in place, specifically, the encryption during data transit. This is useful for auditing and ensuring compliance with data security standards.

```sql+postgres
select
  name,
  id,
  encryption_in_transit_properties -> 'isEncryptionInTransitEnabled' as is_encryption_in_transit_enabled
from
  azure_hdinsight_cluster
where
  (encryption_in_transit_properties ->> 'isEncryptionInTransitEnabled')::boolean;
```

```sql+sqlite
select
  name,
  id,
  json_extract(encryption_in_transit_properties, '$.isEncryptionInTransitEnabled') as is_encryption_in_transit_enabled
from
  azure_hdinsight_cluster
where
  json_extract(encryption_in_transit_properties, '$.isEncryptionInTransitEnabled') = 'true';
```

### List disk encryption details
Explore the encryption details of your Azure HDInsight clusters. This can help you understand your security setup and ensure that the right encryption measures are in place.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  json_extract(disk_encryption_properties, '$.encryptionAlgorithm') as encryption_algorithm,
  json_extract(disk_encryption_properties, '$.encryptionAtHost') as encryption_at_host,
  json_extract(disk_encryption_properties, '$.keyName') as key_name,
  json_extract(disk_encryption_properties, '$.keyVersion') as key_version,
  json_extract(disk_encryption_properties, '$.msiResourceId') as msi_resource_id,
  json_extract(disk_encryption_properties, '$.vaultUri') as vault_uri
from
  azure_hdinsight_cluster;
```

### List connectivity endpoint details
Explore the connectivity details of your HDInsight clusters in Azure. This query helps to understand the location, name, port, protocol, and private IP address of each endpoint, allowing for efficient cluster management and troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  name,
  c.id,
  json_extract(endpoint.value, '$.location') as endpoint_location,
  json_extract(endpoint.value, '$.name') as endpoint_name,
  json_extract(endpoint.value, '$.port') as endpoint_port,
  json_extract(endpoint.value, '$.protocol') as endpoint_protocol,
  json_extract(endpoint.value, '$.privateIpAddress') as endpoint_private_ip_address
from
  azure_hdinsight_cluster as c,
  json_each(connectivity_endpoints) as endpoint;
```