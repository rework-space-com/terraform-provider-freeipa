package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAHBACServiceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAHBACServiceGroupCreate,
		ReadContext:   resourceFreeIPAHBACServiceGroupRead,
		UpdateContext: resourceFreeIPAHBACServiceGroupUpdate,
		DeleteContext: resourceFreeIPAHBACServiceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HBAC service group's name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "A description of this HBAC service group",
			},
		},
	}
}

func resourceFreeIPAHBACServiceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa_hbac_servicegroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HbacsvcgroupAddOptionalArgs{}

	args := ipa.HbacsvcgroupAddArgs{
		Cn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	_, err = client.HbacsvcgroupAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa_hbac_servicegroup: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPAHBACServiceGroupRead(ctx, d, meta)
}

func resourceFreeIPAHBACServiceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa_hbac_servicegroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	args := ipa.HbacsvcgroupShowArgs{
		Cn: d.Get("name").(string),
	}
	optArgs := ipa.HbacsvcgroupShowOptionalArgs{
		All: &all,
	}

	res, err := client.HbacsvcgroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] HBAC service group not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa_hbac_servicegroup: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa_hbac_servicegroup %s", res.Result.Cn)

	return nil
}

func resourceFreeIPAHBACServiceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa_hbac_servicegroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	var hasChange = false
	args := ipa.HbacsvcgroupModArgs{
		Cn: d.Get("name").(string),
	}
	optArgs := ipa.HbacsvcgroupModOptionalArgs{}

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			if v != "" {
				optArgs.Description = &v
				hasChange = true
			}
		}
	}
	if hasChange {
		_, err = client.HbacsvcgroupMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa hbac service group: %s", err)
			}
		}
	}

	return resourceFreeIPAHBACServiceGroupRead(ctx, d, meta)
}

func resourceFreeIPAHBACServiceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa_hbac_servicegroup")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.HbacsvcgroupDelArgs{
		Cn: []string{d.Get("name").(string)},
	}
	optArgs := ipa.HbacsvcgroupDelOptionalArgs{}

	_, err = client.HbacsvcgroupDel(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa_hbac_servicegroup: %s", err)
	}

	d.SetId("")

	return nil
}
