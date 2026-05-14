# main.tf
terraform {
  required_providers {
    mirror = {
      source = "andrew/property-mirror"
    }
  }
  backend "local" {
    path = "/container_shared/tfstate/property-mirror.tfstate"
  }
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

data "vboxdata_vms" "all" {
  
}

output "vm_list" {
  value = data.vboxdata_vms.all.names
}