package main

import (
	"context"
	"log"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSGroupCreate,
		ReadContext:   resourceFreeIPADNSGroupRead,
		UpdateContext: resourceFreeIPADNSGroupUpdate,
		DeleteContext: resourceFreeIPADNSGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gid_number": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"nonposix", "external"},
			},
			"nonposix": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"gid_number", "external"},
			},
			"external": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"gid_number", "nonposix"},
			},
		},
	}
}

func resourceFreeIPADNSGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa grupe")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.GroupAddOptionalArgs{}

	args := ipa.GroupAddArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	if _v, ok := d.GetOkExists("gid_number"); ok {
		v := _v.(int)
		optArgs.Gidnumber = &v
	}
	if _v, ok := d.GetOkExists("nonposix"); ok {
		v := _v.(bool)
		optArgs.Nonposix = &v
	}
	if _v, ok := d.GetOkExists("external"); ok {
		v := _v.(bool)
		optArgs.External = &v
	}
	_, err = client.GroupAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa group: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSGroupRead(ctx, d, meta)
}

func resourceFreeIPADNSGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa grupe")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	optArgs := ipa.GroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.GroupShowArgs{
		Cn: d.Id(),
	}

	res, err := client.GroupShow(&args, &optArgs)

	if err != nil {
		return diag.Errorf("Error show freeipa grupe: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa grupe %s", res.Result.Cn)
	return nil
}

func resourceFreeIPADNSGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa group")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.GroupModArgs{
		Cn: d.Id(),
	}
	optArgs := ipa.GroupModOptionalArgs{}

	var hasChange = false

	if d.HasChange("name") {
		if _v, ok := d.GetOkExists("name"); ok {
			v := _v.(string)
			optArgs.Rename = &v
			hasChange = true
		}
	}
	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			optArgs.Description = &v
			hasChange = true
		}
	}
	if d.HasChange("gid_number") {
		if _v, ok := d.GetOkExists("gid_number"); ok {
			v := _v.(int)
			if v != 0 {
				optArgs.Gidnumber = &v
				hasChange = true
			}
		}
	}

	// TODO: Change No-Posix, Posix, External

	if hasChange {
		_, err = client.GroupMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa groupe: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSGroupRead(ctx, d, meta)
}

func resourceFreeIPADNSGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa grupe")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.GroupDelArgs{
		Cn: []string{d.Id()},
	}
	_, err = client.GroupDel(&args, &ipa.GroupDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa grupe: %s", err)
	}

	d.SetId("")
	return nil
}
