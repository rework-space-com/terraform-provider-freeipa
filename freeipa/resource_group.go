package freeipa

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
				Description: "Group name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description",
			},
			"gid_number": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"nonposix", "external"},
				Description:   "GID (use this option to set it manually)",
			},
			"nonposix": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"gid_number", "external"},
				Description:   "Create as a non-POSIX group",
			},
			"external": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"gid_number", "nonposix"},
				Description:   "Allow adding external non-IPA members from trusted domains",
			},
			"addattr": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Add an attribute/value pair. Format is attr=value. The attribute must be part of the schema.",
			},
			"setattr": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set an attribute to a name/value pair. Format is attr=value.",
			},
		},
	}
}

func resourceFreeIPADNSGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa group")

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
	if _v, ok := d.GetOk("addattr"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Addattr = &v
	}
	if _v, ok := d.GetOk("setattr"); ok {
		v := make([]string, len(_v.([]interface{})))
		for i, value := range _v.([]interface{}) {
			v[i] = value.(string)
		}
		optArgs.Setattr = &v
	}
	_, err = client.GroupAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa group: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSGroupRead(ctx, d, meta)
}

func resourceFreeIPADNSGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa group")

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

	log.Printf("[DEBUG] Read freeipa group %s", d.Id())
	res, err := client.GroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Group not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa group: %s", err)
		}
	}

	log.Printf("[DEBUG] Read freeipa group %s", res.Result.Cn)
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
	if d.HasChange("addattr") {
		if _v, ok := d.GetOkExists("addattr"); ok {
			v := make([]string, len(_v.([]interface{})))
			for i, value := range _v.([]interface{}) {
				v[i] = value.(string)
			}
			optArgs.Addattr = &v
			hasChange = true
		}
	}
	if d.HasChange("setattr") {
		if _v, ok := d.GetOkExists("setattr"); ok {
			v := make([]string, len(_v.([]interface{})))
			for i, value := range _v.([]interface{}) {
				v[i] = value.(string)
			}
			optArgs.Setattr = &v
			hasChange = true
		}
	}

	// TODO: Change No-Posix, Posix, External

	if hasChange {
		_, err = client.GroupMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa group: %s", err)
			}
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSGroupRead(ctx, d, meta)
}

func resourceFreeIPADNSGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa group")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.GroupDelArgs{
		Cn: []string{d.Id()},
	}
	_, err = client.GroupDel(&args, &ipa.GroupDelOptionalArgs{})
	if err != nil {
		return diag.Errorf("Error delete freeipa group: %s", err)
	}

	d.SetId("")
	return nil
}
