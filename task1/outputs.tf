data "azurerm_public_ip" "vm_public_ip" {
  name                = azurerm_public_ip.vm_public_ip.name
  resource_group_name = azurerm_linux_virtual_machine.vm.resource_group_name
}

output "vm_public_ip" {
  value = data.azurerm_public_ip.vm_public_ip.ip_address
}

output "mysql_server_fqdn" {
  value = azurerm_mysql_flexible_server.mysql_server.fqdn
}
