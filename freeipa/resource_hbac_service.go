package freeipa

import (
	"context"
	"log"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAHBACService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAHBACServiceCreate,
		ReadContext:   resourceFreeIPAHBACServiceRead,
		DeleteContext: resourceFreeIPAHBACServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the HBAC service",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "HBAC service description",
			},
		},
	}
}

func resourceFreeIPAHBACServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the HBAC service")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HbacsvcAddOptionalArgs{}

	args := ipa.HbacsvcAddArgs{
		Cn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}

	_, err = client.HbacsvcAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the HBAC service: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPAHBACServiceRead(ctx, d, meta)
}

func resourceFreeIPAHBACServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa HBAC service")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacsvcModArgs{
		Cn: d.Id(),
	}
	optArgs := ipa.HbacsvcModOptionalArgs{}

	var hasChange = false

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			optArgs.Description = &v
			hasChange = true
		}
	}

	// TODO: Change No-Posix, Posix, External

	if hasChange {
		_, err = client.HbacsvcMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa HBAC service: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPASudocmdRead(ctx, d, meta)
}

func resourceFreeIPAHBACServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the HBAC service")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	args := ipa.HbacsvcShowArgs{
		Cn: d.Id(),
	}

	optArgs := ipa.HbacsvcShowOptionalArgs{
		All: &all,
	}

	res, err := client.HbacsvcShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] HBAC service not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa HBAC service: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa sudo command %s", *&res.Result.Cn)
	return nil
}

func resourceFreeIPAHBACServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the HBAC service")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacsvcDelArgs{
		Cn: []string{d.Id()},
	}

	_, err = client.HbacsvcDel(&args, &ipa.HbacsvcDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa the HBAC service: %s", err)
	}

	d.SetId("")

	return nil
}
