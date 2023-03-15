# Table: azure_bastion_host

Azure Bastion is a service you deploy that lets you connect to a virtual machine using your browser and the Azure portal, or via the native SSH or RDP client already installed on your local computer. The Azure Bastion service is a fully platform-managed PaaS service that you provision inside your virtual network. It provides secure and seamless RDP/SSH connectivity to your virtual machines directly from the Azure portal over TLS. When you connect via Azure Bastion, your virtual machines don't need a public IP address, agent, or special client software.

## Examples

### Basic info

```sql
select
	name,
	dns_name,
	provisioning_state,
	region,
	resource_group
from
	azure_bastion_host;
```


### List bastion hosts that failed creation

```sql
select
	name,
	dns_name,
	provisioning_state,
	region,
	resource_group
from
	azure_bastion_host
where
	provisioning_state = 'Failed';
```


### Get subnet details associated with each host

```sql
select
	h.name as bastion_host_name,
	s.id as subnet_id,
	s.name as subnet_name,
	address_prefix
from
	azure_bastion_host h,
	jsonb_array_elements(ip_configurations) ip,
	azure_subnet s
where
	s.id = ip -> 'properties' -> 'subnet' ->> 'id';
```

### Get ip configuration details associated with each host

```sql
select
	h.name as bastion_host_name,
	i.name as ip_configuration_name,
	ip_configuration_id,
	ip_address,
	public_ip_allocation_method,
	sku_name as ip_configuration_sku
from
	azure_bastion_host h,
	jsonb_array_elements(ip_configurations) ip,
	azure_public_ip i
where
	i.id = ip -> 'properties' -> 'publicIPAddress' ->> 'id';
```