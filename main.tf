# main.tf
terraform {
  required_providers {
    mirror = {
      source = "andrew/property-mirror"
    }
    linode = {
      source = "linode/linode"
      version = "~> 3.13.0"
    }
  }
  backend "local" {
    path = "/container_shared/tfstate/property-mirror.tfstate"
  }
}

variable "LINODE_API_KEY" {
  description = "The key to the Linode API"
  type        = string
  sensitive   = true
}

provider "linode" {
  token = var.LINODE_API_KEY
}


provider "mirror" {
  # this is a private network address.
  vboxwebsrv_endpoint = "http://daddy.siwko.org:18083/" 
  username = ""
  password = ""
}

resource "mirror_variable" "rhel_test_1" {
  value = "Running on RHEL 9.7"
}

resource "mirror_variable" "rhel_test_2" {
  value = "Also Running on RHEL 9.7"
}

resource "mirror_variable" "rhel_test_3" {
  value = "Also Running on RHEL 9.7 - new!"
}

output "var_1" {
  value = mirror_variable.rhel_test_1.value
}

output "var_2" {
  value = mirror_variable.rhel_test_2.value
} 

data "mirror_coffee" "my_espresso" {
  id = 1
}

output "coffee_name" {
  value = data.mirror_coffee.my_espresso.name
}

output "coffee_price" {
  value = data.mirror_coffee.my_espresso.price
}

data "mirror_vms" "all" {

}
output "vm_map" {
  value = {
    for vm in data.mirror_vms.all.vms : vm.name => {
      state        = vm.state
      cpus         = vm.cpus
      memory_mb    = vm.memory
      ip_addresses = vm.ip_addresses
    }
  }
}

locals {
  # Filter for VMs that are running and have at least one valid IP address
  active_vbox_vms = {
    for vm in data.mirror_vms.all.vms : vm.name => vm.ip_addresses[0]
    if vm.state == "Running" && length(vm.ip_addresses) > 0
  }
}
output "active_vbox_vm_map" {
  value = local.active_vbox_vms
}

# terraform import linode_domain.dns_zone 3417841
resource "linode_domain" "dns_zone" {
  type        = "master"
  domain      = "siwko.org"
  soa_email   = "asiwko@siwko.org"
  refresh_sec = 30
  retry_sec   = 30
  ttl_sec     = 30
  lifecycle {
    prevent_destroy = true
  }
}

resource "linode_domain_record" "virtualbox_records" {
  for_each = local.active_vbox_vms

  domain_id   = linode_domain.dns_zone.id
  record_type = "A"
  ttl_sec     = 30 # Match your aggressive local cluster testing profile

  # each.key is the VM name (e.g., "ansible-control")
  name        = "vbox_${each.key}"
  
  # each.value is the resolved IP address (e.g., "192.168.51.48")
  target      = each.value
}