terraform {
  required_providers {
    meraki = {
      version = "0.2.3-alpha"
      source  = "hashicorp.com/edu/meraki"
      # "hashicorp.com/edu/meraki" is the local built source, change to "cisco-en-programmability/meraki" to use downloaded version from registry
    }
  }
}
provider "meraki" {
  meraki_debug = "true"
}
data "meraki_devices" "example" {

  #   organization_id = "828099381482762766"
  organization_id = "828099381482762270"
  # serial = "QBSB-AX45-LY9A"
  network_ids = []
  #   /api/v1/organizations/828099381482762270/devices?
  #   networkIds=%221HOLA1%22&networkIds=%224HOLA4%22
}
output "meraki_devices_example" {
  value = data.meraki_devices.example.items
}