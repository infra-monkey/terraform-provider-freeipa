package freeipa

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeIPADNSZone_basic(t *testing.T) {
	testZone := map[string]string{
		"index":     "0",
		"zone_name": "ipa.example.lan",
	}
	testDS := map[string]string{
		"index":     "0",
		"zone_name": "ipa.example.lan.",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dns-zone-0", "computed_zone_name", "ipa.example.lan."),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSZone_datasource(testDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dns-zone-0", "computed_zone_name", "ipa.example.lan."),
				),
			},
		},
	})
}

func TestAccFreeIPADNSZone_reverse(t *testing.T) {
	testZone := map[string]string{
		"index":           "0",
		"zone_name":       "192.168.23.0",
		"is_reverse_zone": "true",
	}
	testDS := map[string]string{
		"index":     "0",
		"zone_name": "23.168.192.in-addr.arpa.",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dns-zone-0", "computed_zone_name", "23.168.192.in-addr.arpa."),
				),
			},
			{
				Config: testAccFreeIPAProvider() + testAccFreeIPADNSZone_resource(testZone) + testAccFreeIPADNSZone_datasource(testDS),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freeipa_dns_zone.dns-zone-0", "computed_zone_name", "23.168.192.in-addr.arpa."),
				),
			},
		},
	})
}
