package jmespath

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testDataSourceConfig_basic = `
provider "jmespath" {}

data "jmespath_search" "foo" {
  expression = "locations[?state == 'WA'].name | sort(@) | {WashingtonCities: join(', ', @)}"

  input =<<EOS
{
  "locations": [
    {"name": "Seattle", "state": "WA"},
    {"name": "New York", "state": "NY"},
    {"name": "Bellevue", "state": "WA"},
    {"name": "Olympia", "state": "WA"}
  ]
}
EOS
}
`

func TestDataSource_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jmespath_search.foo", "result"),
				),
			},
		},
	})
}
