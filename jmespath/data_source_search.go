package jmespath

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gj "github.com/jmespath/go-jmespath"
)

func dataSourceSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadContext,

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
					s := v.(string)
					if s != "" && !json.Valid([]byte(s)) {
						errors = append(errors, fmt.Errorf("%q contains invalid JSON", k))
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

func dataSourceReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var input interface{}
	if err := json.Unmarshal([]byte(d.Get("input").(string)), &input); err != nil {
		return diag.Errorf("error umarshaling JSON: %s", err)
	}

	result, err := gj.Search(d.Get("expression").(string), input)
	if err != nil {
		return diag.Errorf("error performing JMES search: %s", err)
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return diag.Errorf("error marshaling JSON: %s", err)
	}

	hash := crc32.ChecksumIEEE(jsonResult)
	d.SetId(strconv.FormatUint(uint64(hash), 16))

	d.Set("result", string(jsonResult))

	return nil
}
