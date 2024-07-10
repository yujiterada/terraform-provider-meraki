terraform {
  required_providers {
    meraki = {
      version = "0.2.5-alpha"
      source  = "hashicorp.com/edu/meraki"
      # "hashicorp.com/edu/meraki" is the local built source, change to "cisco-en-programmability/meraki" to use downloaded version from registry
    }
  }
}

provider "meraki" {
  meraki_debug = "true"
}

resource "meraki_networks_vlan_profiles" "vlan_profiles" {
   network_id = "L_828099381482771185"
   iname = "Default 2"
   name = "Default Profile 2"
   vlan_names = [ {
      name = "default_2",
      vlan_id = "1"
   }, {
    name = "test_2",
    vlan_id = "2"
   }]
   vlan_groups = []
}

