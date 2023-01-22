package main

import (
	"context"
	"log"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFreeIPADNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSDNSRecordCreate,
		ReadContext:   resourceFreeIPADNSDNSRecordRead,
		UpdateContext: resourceFreeIPADNSDNSRecordUpdate,
		DeleteContext: resourceFreeIPADNSDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"set_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFreeIPADNSDNSRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa dns record")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	name := d.Get("name").(string)
	zone_name := d.Get("zone_name").(string)

	args := ipa.DnsrecordAddArgs{
		Idnsname: name,
	}

	optArgs := ipa.DnsrecordAddOptionalArgs{
		Dnszoneidnsname: &zone_name,
	}

	_type := d.Get("type").(string)
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	if _v, ok := d.GetOkExists("ttl"); ok {
		v := _v.(int)
		optArgs.Dnsttl = &v
	}

	_, err = client.DnsrecordAdd(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
		} else {
			return diag.Errorf("Error creating freeipa dns record: %s", err)
		}
	}

	// Generate an ID
	vars := []string{
		zone_name,
		strings.ToLower(name),
		_type,
	}
	if v, ok := d.GetOk("set_identifier"); ok {
		vars = append(vars, v.(string))
	}

	d.SetId(strings.Join(vars, "_"))

	return resourceFreeIPADNSDNSRecordRead(ctx, d, meta)
}

func resourceFreeIPADNSDNSRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa dns record")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.DnsrecordShowArgs{
		Idnsname: d.Get("name").(string),
	}

	zone_name := d.Get("zone_name").(string)
	all := true
	optArgs := ipa.DnsrecordShowOptionalArgs{
		Dnszoneidnsname: &zone_name,
		All:             &all,
	}

	res, err := client.DnsrecordShow(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa dns record: %s", err)
	}

	_type := d.Get("type")

	switch _type {
	case "A":
		if res.Result.Arecord != nil {
			d.Set("records", *res.Result.Arecord)
		}
	case "AAAA":
		if res.Result.Aaaarecord != nil {
			d.Set("records", *res.Result.Aaaarecord)
		}
	case "MX":
		if res.Result.Mxrecord != nil {
			d.Set("records", *res.Result.Mxrecord)
		}
	case "NS":
		if res.Result.Nsrecord != nil {
			d.Set("records", *res.Result.Nsrecord)
		}
	case "PTR":
		if res.Result.Ptrrecord != nil {
			d.Set("records", *res.Result.Ptrrecord)
		}
	case "SRV":
		if res.Result.Srvrecord != nil {
			d.Set("records", *res.Result.Srvrecord)
		}
	case "TXT":
		if res.Result.Txtrecord != nil {
			d.Set("records", *res.Result.Txtrecord)
		}
	case "SSHFP":
		if res.Result.Sshfprecord != nil {
			d.Set("records", *res.Result.Sshfprecord)
		}
	}

	if res.Result.Dnsttl != nil {
		d.Set("ttl", *res.Result.Dnsttl)
	}

	return nil
}

func resourceFreeIPADNSDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Update freeipa dns record")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.DnsrecordModArgs{
		Idnsname: d.Get("name").(string),
	}

	zone_name := d.Get("zone_name").(string)
	optArgs := ipa.DnsrecordModOptionalArgs{
		Dnszoneidnsname: &zone_name,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	if _v, ok := d.GetOkExists("ttl"); ok {
		v := _v.(int)
		optArgs.Dnsttl = &v
	}

	_, err = client.DnsrecordMod(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "EmptyModlist") {
			log.Printf("[DEBUG] EmptyModlist (4202): no modifications to be performed")
		} else {
			return diag.Errorf("Error update freeipa dns record: %s", err)
		}
	}
	return resourceFreeIPADNSDNSRecordRead(ctx, d, meta)
}

func resourceFreeIPADNSDNSRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa dns record")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}
	args := ipa.DnsrecordDelArgs{
		Idnsname: d.Get("name").(string),
	}

	zone_name := d.Get("zone_name").(string)
	optArgs := ipa.DnsrecordDelOptionalArgs{
		Dnszoneidnsname: &zone_name,
	}

	_type := d.Get("type")
	_records := d.Get("records").(*schema.Set).List()
	records := make([]string, len(_records))
	for i, d := range _records {
		records[i] = d.(string)
	}
	switch _type {
	case "A":
		optArgs.Arecord = &records
	case "AAAA":
		optArgs.Aaaarecord = &records
	case "CNAME":
		optArgs.Cnamerecord = &records
	case "MX":
		optArgs.Mxrecord = &records
	case "NS":
		optArgs.Nsrecord = &records
	case "PTR":
		optArgs.Ptrrecord = &records
	case "SRV":
		optArgs.Srvrecord = &records
	case "TXT":
		optArgs.Txtrecord = &records
	case "SSHFP":
		optArgs.Sshfprecord = &records
	}

	_, err = client.DnsrecordDel(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa dns record: %s", err)
	}

	d.SetId("")
	return nil
}
