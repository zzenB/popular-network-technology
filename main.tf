resource "azurerm_resource_group" "rg" {
  name     = "rg-task1"
  location = "East Asia"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-task1"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  address_space       = ["10.0.0.0/18"]
}

resource "azurerm_subnet" "vm_subnet" {
  name                 = "vm-subnet-task1"
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "db_subnet" {
  name                 = "db-subnet-task1"
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]

  delegation {
    name = "delegation-task1"
    service_delegation {
      name    = "Microsoft.DBforMySQL/flexibleServers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_security_group" "vm_nsg" {
  name                = "vm-nsg-task1"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "vm_nsg_association" {
  subnet_id                 = azurerm_subnet.vm_subnet.id
  network_security_group_id = azurerm_network_security_group.vm_nsg.id
}

resource "azurerm_public_ip" "vm_public_ip" {
  name                = "vm-public-ip-task1"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  allocation_method   = "Dynamic"
}

resource "azurerm_network_interface" "vm_nic" {
  name                = "vm-nic-task1"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.vm_subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm_public_ip.id
  }
}

resource "azurerm_linux_virtual_machine" "vm" {
  name                  = "vm-task1"
  resource_group_name   = azurerm_resource_group.rg.name
  location              = azurerm_resource_group.rg.location
  size                  = "Standard_D2s_v3"
  admin_username        = "adminuser"
  network_interface_ids = [azurerm_network_interface.vm_nic.id]
  priority              = "Spot" // so gk pricy ;)
  eviction_policy       = "Deallocate"

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
    # disk_size_gb         = 10 // minimum untuk image ubuntu itu semuanya 30gb jadi kalau ini di uncomment akan error
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }
}

resource "azurerm_private_dns_zone" "mysql_private_dns_zone" {
  name                = "task1.mysql.database.azure.com"
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "mysql_dns_link" {
  name                  = "dns-link-task1"
  private_dns_zone_name = azurerm_private_dns_zone.mysql_private_dns_zone.name
  virtual_network_id    = azurerm_virtual_network.vnet.id
  resource_group_name   = azurerm_resource_group.rg.name
}

resource "azurerm_mysql_flexible_server" "mysql_server" {
  name                   = "mysql-server-task1"
  resource_group_name    = azurerm_resource_group.rg.name
  location               = azurerm_resource_group.rg.location
  sku_name               = "B_Standard_B1s"
  administrator_login    = "mysqladmin"
  administrator_password = "password123!"
  delegated_subnet_id    = azurerm_subnet.db_subnet.id
  private_dns_zone_id    = azurerm_private_dns_zone.mysql_private_dns_zone.id

  depends_on = [azurerm_private_dns_zone_virtual_network_link.mysql_dns_link]
}

resource "azurerm_mysql_flexible_database" "mysql_db" {
  name                = "mysql-db-task1"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_mysql_flexible_server.mysql_server.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

resource "azurerm_mysql_flexible_server_firewall_rule" "mysql_firewall" {
  name                = "allow-vm-task1"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_mysql_flexible_server.mysql_server.name
  start_ip_address    = azurerm_network_interface.vm_nic.private_ip_address
  end_ip_address      = azurerm_network_interface.vm_nic.private_ip_address
}
