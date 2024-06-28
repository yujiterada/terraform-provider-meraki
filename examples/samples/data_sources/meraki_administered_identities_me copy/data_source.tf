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
data "meraki_administered_identities_me" "example" {
  count = 1000

}
# output "meraki_administered_identities_me_example" {
#   value = data.meraki_administered_identities_me.example.item
# }