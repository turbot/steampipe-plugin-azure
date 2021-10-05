resource "null_resource" "named_resource" {
  provisioner "local-exec" {
    command = "az iot dps delete --resource-group {{ output.resource_name.value }} --name {{ output.resource_name.value }}"
  }
}
