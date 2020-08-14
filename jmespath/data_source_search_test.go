package jmespath

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testDataSourceConfig_basic = `
data "jmespath_search" "test" {
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
const testDataSourceExceptedResult_basic = `
{
  "WashingtonCities": "Bellevue, Olympia, Seattle"
}
`

func TestDataSource_basic(t *testing.T) {
	dataSourceName := "data.jmespath_search.test"

	expectedResult, err := normalizeJsonString(testDataSourceExceptedResult_basic)
	if err != nil {
		t.Fatal(err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "result", expectedResult),
				),
			},
		},
	})
}

const testDataSourceConfig_projection = `
data "jmespath_search" "test" {
  expression = "people[*].first"

  input =<<EOS
  {
    "people": [
      {"first": "James", "last": "d"},
      {"first": "Jacob", "last": "e"},
      {"first": "Jayden", "last": "f"},
      {"missing": "different"}
    ],
    "foo": {"bar": "baz"}
  }
EOS
}
`
const testDataSourceExceptedResult_projection = `
[
"James",
"Jacob",
"Jayden"
]
`

func TestDataSource_projection(t *testing.T) {
	dataSourceName := "data.jmespath_search.test"

	expectedResult, err := normalizeJsonString(testDataSourceExceptedResult_projection)
	if err != nil {
		t.Fatal(err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfig_projection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "result", expectedResult),
				),
			},
		},
	})
}

const testDataSourceConfig_multiSelect = `
data "jmespath_search" "test" {
  expression = "people[].{Name: name, State: state.name}"

  input =<<EOS
  {
    "people": [
      {
        "name": "a",
        "state": {"name": "up"}
      },
      {
        "name": "b",
        "state": {"name": "down"}
      },
      {
        "name": "c",
        "state": {"name": "up"}
      }
    ]
  }
EOS
}
`
const testDataSourceExceptedResult_multiSelect = `
[
    {
      "Name": "a",
      "State": "up"
    },
    {
      "Name": "b",
      "State": "down"
    },
    {
      "Name": "c",
      "State": "up"
    }
  ]
`

func TestDataSource_multiSelect(t *testing.T) {
	dataSourceName := "data.jmespath_search.test"

	expectedResult, err := normalizeJsonString(testDataSourceExceptedResult_multiSelect)
	if err != nil {
		t.Fatal(err)
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfig_multiSelect,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "result", expectedResult),
				),
			},
		},
	})
}

func normalizeJsonString(s string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return "", err
	}

	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(j), nil
}
