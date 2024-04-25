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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Host's name (FQDN)",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "IP address of the host",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "A description of this host",
			},
			"locality": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Host locality (e.g. 'Baltimore, MD')",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Host location (e.g. 'Lab 2')",
			},
			"platform": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Host hardware platform (e.g. 'Lenovo T61')",
			},
			"operating_system": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Host operating system and version (e.g. 'Fedora 9')",
			},
			"user_certificates": {
				Type:        schema.TypeList,
				Required:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Base-64 encoded host certificate",
			},
			"mac_addresses": {
				Type:        schema.TypeList,
				Required:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Hardware MAC address(es) on this host",
			},
			"ipasshpubkeys": {
				Type:        schema.TypeList,
				Required:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "SSH public key",
			},
			"userclass": {
				Type:        schema.TypeList,
				Required:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Host category (semantics placed on this attribute are for local interpretation)",
			},
			"assigned_idview": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Assigned ID View",
			},
			"krb_auth_indicators": {
				Type:        schema.TypeList,
				Required:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Defines a whitelist for Authentication Indicators. Use 'otp' to allow OTP-based 2FA authentications. Use 'radius' to allow RADIUS-based 2FA authentications. Other values may be used for custom configurations.",
			},
			"krb_preauth": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Pre-authentication is required for the service",
			},
			"trusted_for_delegation": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Client credentials may be delegated to the service",
			},
			"trusted_to_auth_as_delegate": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "The service is allowed to authenticate on behalf of a client",
			},
			"force": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Skip host's DNS check (A/AAAA) before adding it",
			},
			"userpassword": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password used in bulk enrollment",
			},
			"random_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Generate a random password to be used in bulk enrollment",
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
	if _v, ok := d.GetOkExists("random_password"); ok {
		v := _v.(bool)
		optArgs.Random = &v
	}
	if _v, ok := d.GetOkExists("userpassword"); ok {
		v := _v.(string)
		optArgs.Userpassword = &v
	}
	if _v, ok := d.GetOkExists("force"); ok {
		v := _v.(bool)
		optArgs.Force = &v
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
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] Host not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa host: %s", err)
		}
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
	if d.HasChange("userpassword") {
		if _v, ok := d.GetOkExists("userpassword"); ok {
			v := _v.(string)
			optArgs.Userpassword = &v
			hasChange = true
		}
	}
	if d.HasChange("random_password") {
		if _v, ok := d.GetOkExists("random_password"); ok {
			v := _v.(bool)
			optArgs.Random = &v
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
