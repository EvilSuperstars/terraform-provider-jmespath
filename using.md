# `jmespath` Provider

The Terraform [jmespath](https://github.com/EvilSuperstars/terraform-provider-jmespath) provider evaluates a [JMESPath](http://jmespath.org/) expression against input data and returns the result.

This provider requires no configuration.

### Example Usage

```hcl
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
```

## Data Sources

### `jmespath_search`

#### Argument Reference

The following arguments are supported:

* `expression` - (Required, string) The JMESPath query expression.
* `input` - (Required, string) The string that the query is evaluated against.

#### Attributes Reference

The following attributes are exported in addition to the above configuration:

* `result` - (string) The query result as a JSON string.
