---
title: "Steampipe Table: azure_hdinsight_cluster - Query Azure HDInsight Clusters using SQL"
description: "Allows users to query Azure HDInsight Clusters."
---

# Table: azure_hdinsight_cluster - Query Azure HDInsight Clusters using SQL

Azure HDInsight is a cloud distribution of the Hadoop components from the Hortonworks Data Platform (HDP). Azure HDInsight makes it easy, fast, and cost-effective to process massive amounts of data. You can use the most popular open-source frameworks such as Hadoop, Spark, Hive, LLAP, Kafka, Storm, R, and more.

## Table Usage Guide

The 'azure_hdinsight_cluster' table provides insights into HDInsight Clusters within Azure HDInsight. As a DevOps engineer, you can explore cluster-specific details through this table, including the cluster type, version, state, and associated metadata. Utilize it to uncover information about clusters, such as the number of worker nodes, the type of storage used, and the networking configurations. The schema presents a range of attributes of the HDInsight Cluster for your analysis, like the cluster ID, creation date, tier, and associated tags.

## Examples

### Basic info
Explore the status and details of your Azure HDInsight clusters to understand their configuration and performance. This can help in maintaining optimal cluster health and efficiency.

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
Assess the elements within your Azure HDInsight clusters to identify those with enabled encryption in transit. This can be useful to ensure data security and compliance with your organization's security policies.

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
Explore the encryption specifics of your disk resources to better understand your data's security. This query could be used to assess the encryption methods and algorithms in place, helping to identify potential vulnerabilities or areas for improvement.

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
Explore the connectivity details of your Azure HDInsight clusters to understand their communication protocols, locations, and private IP addresses. This information can be useful in managing network configurations and optimizing data transfer between various clusters.

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