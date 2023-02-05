package main

import (
	"context"
	"log"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPAHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSHostCreate,
		ReadContext:   resourceFreeIPADNSHostRead,
		UpdateContext: resourceFreeIPADNSHostUpdate,
		DeleteContext: resourceFreeIPADNSHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"locality": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"operating_system": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"user_certificates": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mac_addresses": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipasshpubkeys": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"userclass": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"assigned_idview": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"krb_auth_indicators": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"krb_preauth": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			"trusted_for_delegation": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			"trusted_to_auth_as_delegate": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
		},
	}
}

func resourceFreeIPADNSHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa host")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HostAddOptionalArgs{}

	args := ipa.HostAddArgs{
		Fqdn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("description"); ok {
		v := _v.(string)
		optArgs.Description = &v
	}
	if _v, ok := d.GetOkExists("ip_address"); ok {
		v := _v.(string)
		optArgs.IPAddress = &v
	}
	if _v, ok := d.GetOkExists("locality"); ok {
		v := _v.(string)
		optArgs.L = &v
	}
	if _v, ok := d.GetOkExists("location"); ok {
		v := _v.(string)
		optArgs.Nshostlocation = &v
	}
	if _v, ok := d.GetOkExists("platform"); ok {
		v := _v.(string)
		optArgs.Nshardwareplatform = &v
	}
	if _v, ok := d.GetOkExists("operating_system"); ok {
		v := _v.(string)
		optArgs.Nsosversion = &v
	}
	if _v, ok := d.GetOkExists("user_certificates"); ok {
		v := utilsGetArry(_v.([]interface{}))
		v2 := make([]interface{}, len(v))
		for i := range v {
			v2[i] = v[i]
		}
		optArgs.Usercertificate = &v2
	}
	if _v, ok := d.GetOkExists("mac_addresses"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Macaddress = &v
	}
	if _v, ok := d.GetOkExists("ipasshpubkeys"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Ipasshpubkey = &v
	}
	if _v, ok := d.GetOkExists("userclass"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Userclass = &v
	}
	if _v, ok := d.GetOkExists("assigned_idview"); ok {
		v := _v.(string)
		optArgs.Ipaassignedidview = &v
	}
	if _v, ok := d.GetOkExists("krb_auth_indicators"); ok {
		v := utilsGetArry(_v.([]interface{}))
		optArgs.Krbprincipalauthind = &v
	}
	if _v, ok := d.GetOkExists("krb_preauth"); ok {
		v := _v.(bool)
		optArgs.Ipakrbrequirespreauth = &v
	}
	if _v, ok := d.GetOkExists("trusted_for_delegation"); ok {
		v := _v.(bool)
		optArgs.Ipakrbokasdelegate = &v
	}
	if _v, ok := d.GetOkExists("trusted_to_auth_as_delegation"); ok {
		v := _v.(bool)
		optArgs.Ipakrboktoauthasdelegate = &v
	}
	_, err = client.HostAdd(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa host: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFreeIPADNSHostRead(ctx, d, meta)
}

func resourceFreeIPADNSHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa host")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	all := true
	args := ipa.HostShowArgs{
		Fqdn: d.Get("name").(string),
	}
	optArgs := ipa.HostShowOptionalArgs{
		All: &all,
	}

	res, err := client.HostShow(&args, &optArgs)

	if err != nil {
		return diag.Errorf("Error show freeipa host: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa host %s", res.Result.Fqdn)

	return nil
}

func resourceFreeIPADNSHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa host")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	var hasChange = false
	args := ipa.HostModArgs{
		Fqdn: d.Get("name").(string),
	}
	optArgs := ipa.HostModOptionalArgs{}

	if d.HasChange("description") {
		if _v, ok := d.GetOkExists("description"); ok {
			v := _v.(string)
			if v != "" {
				optArgs.Description = &v
				hasChange = true
			}
		}
	}
	if d.HasChange("locality") {
		if _v, ok := d.GetOkExists("locality"); ok {
			v := _v.(string)
			optArgs.L = &v
			hasChange = true
		}
	}
	if d.HasChange("location") {
		if _v, ok := d.GetOkExists("location"); ok {
			v := _v.(string)
			optArgs.Nshostlocation = &v
			hasChange = true
		}
	}
	if d.HasChange("platform") {
		if _v, ok := d.GetOkExists("platform"); ok {
			v := _v.(string)
			optArgs.Nshardwareplatform = &v
			hasChange = true
		}
	}
	if d.HasChange("operating_system") {
		if _v, ok := d.GetOkExists("operating_system"); ok {
			v := _v.(string)
			optArgs.Nsosversion = &v
			hasChange = true
		}
	}
	if d.HasChange("user_certificates") {
		if _v, ok := d.GetOkExists("user_certificates"); ok {
			v := utilsGetArry(_v.([]interface{}))
			v2 := make([]interface{}, len(v))
			for i := range v {
				v2[i] = v[i]
			}
			optArgs.Usercertificate = &v2
			hasChange = true
		}
	}
	if d.HasChange("mac_addresses") {
		if _v, ok := d.GetOkExists("mac_addresses"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Macaddress = &v
			hasChange = true
		}
	}
	if d.HasChange("ipasshpubkeys") {
		if _v, ok := d.GetOkExists("ipasshpubkeys"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Ipasshpubkey = &v
			hasChange = true
		}
	}
	if d.HasChange("userclass") {
		if _v, ok := d.GetOkExists("userclass"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Userclass = &v
			hasChange = true
		}
	}
	if d.HasChange("assigned_idview") {
		if _v, ok := d.GetOkExists("assigned_idview"); ok {
			v := _v.(string)
			optArgs.Ipaassignedidview = &v
			hasChange = true
		}
	}
	if d.HasChange("krb_auth_indicators") {
		if _v, ok := d.GetOkExists("krb_auth_indicators"); ok {
			v := utilsGetArry(_v.([]interface{}))
			optArgs.Krbprincipalauthind = &v
			hasChange = true
		}
	}
	if d.HasChange("krb_preauth") {
		if _v, ok := d.GetOkExists("krb_preauth"); ok {
			v := _v.(bool)
			optArgs.Ipakrbrequirespreauth = &v
			hasChange = true
		}
	}
	if d.HasChange("trusted_for_delegation") {
		if _v, ok := d.GetOkExists("trusted_for_delegation"); ok {
			v := _v.(bool)
			optArgs.Ipakrbokasdelegate = &v
			hasChange = true
		}
	}
	if d.HasChange("trusted_to_auth_as_delegation") {
		if _v, ok := d.GetOkExists("trusted_to_auth_as_delegation"); ok {
			v := _v.(bool)
			optArgs.Ipakrboktoauthasdelegate = &v
			hasChange = true
		}
	}
	if hasChange {
		_, err = client.HostMod(&args, &optArgs)
		if err != nil {
			if strings.Contains(err.Error(), "EmptyModlist") {
				log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
			} else {
				return diag.Errorf("Error update freeipa host: %s", err)
			}
		}
	}

	return resourceFreeIPADNSHostRead(ctx, d, meta)
}

func resourceFreeIPADNSHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa host")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	continuedel := false
	updatednsdel := true
	var fqdns = []string{d.Get("name").(string)}
	args := ipa.HostDelArgs{
		Fqdn: fqdns,
	}
	optArgs := ipa.HostDelOptionalArgs{
		Continue:  &continuedel,
		Updatedns: &updatednsdel,
	}

	_, err = client.HostDel(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa host: %s", err)
	}

	d.SetId("")

	return nil
}
