package jmespath

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	gj "github.com/jmespath/go-jmespath"
)

func dataSourceSearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if _, err := gj.Compile(v.(string)); err != nil {
						errors = append(errors, fmt.Errorf("%q contains an invalid JMESPath expression: %s", k, err))
					}

					return
				},
			},

			"input": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v == nil {
						return
					}

					s := v.(string)
					if s == "" {
						return
					}

					var i interface{}
					err := json.Unmarshal([]byte(s), &i)
					if err != nil {
						errors = append(errors, fmt.Errorf("%q contains invalid JSON: %s", k, err))
					}

					return
				},
			},

			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {
	var input interface{}
	if err := json.Unmarshal([]byte(d.Get("input").(string)), &input); err != nil {
		return err
	}

	result, err := gj.Search(d.Get("expression").(string), input)
	if err != nil {
		return err
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil
	}

	d.SetId(time.Now().UTC().String())
	d.Set("result", string(jsonResult))

	return nil
}
