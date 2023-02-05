package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func resourceFreeIPAHBACPolicyHostMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSHBACPolicyHostMembershipCreate,
		ReadContext:   resourceFreeIPADNSHBACPolicyHostMembershipRead,
		DeleteContext: resourceFreeIPADNSHBACPolicyHostMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hostgroup"},
			},
			"hostgroup": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"host"},
			},
		},
	}
}

func resourceFreeIPADNSHBACPolicyHostMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the HBAC policy host membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	hostmember_id := "h"
	optArgs := ipa.HbacruleAddHostOptionalArgs{}

	args := ipa.HbacruleAddHostArgs{
		Cn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("host"); ok {
		v := []string{_v.(string)}
		optArgs.Host = &v
		hostmember_id = "h"
	}
	if _v, ok := d.GetOkExists("hostgroup"); ok {
		v := []string{_v.(string)}
		optArgs.Hostgroup = &v
		hostmember_id = "hg"
	}

	_, err = client.HbacruleAddHost(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the HBAC policy host membership: %s", err)
	}
	switch hostmember_id {
	case "hg":
		id := fmt.Sprintf("%s/hg/%s", d.Get("name").(string), d.Get("hostgroup").(string))
		d.SetId(id)
	case "h":
		id := fmt.Sprintf("%s/h/%s", d.Get("name").(string), d.Get("host").(string))
		d.SetId(id)
	}

	return resourceFreeIPADNSHBACPolicyHostMembershipRead(ctx, d, meta)
}

func resourceFreeIPADNSHBACPolicyHostMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the HBAC policy host membership")

	name, typeId, hostId, err := parseHBACPolicyHostMembershipID(d.Id())
	all := true
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleShowArgs{
		Cn: name,
	}
	optArgs := ipa.HbacruleShowOptionalArgs{
		All: &all,
	}
	res, err := client.HbacruleShow(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error show freeipa HBAC policy host membership: %s", err)
	}

	switch typeId {
	case "hg":
		if res.Result.MemberhostHostgroup == nil || !slices.Contains(*res.Result.MemberhostHostgroup, hostId) {
			log.Printf("[DEBUG] Warning! Hostgroup membership does not exist")
			d.Set("host", "")
			d.Set("hostgroup", "")
			d.SetId("")
			return diag.Errorf("Error configuring freeipa HBAC policy, hostgroup not assigned: %s", hostId)
		}
	case "h":
		if res.Result.MemberhostHost == nil || !slices.Contains(*res.Result.MemberhostHost, hostId) {
			log.Printf("[DEBUG] Warning! Host membership does not exist")
			d.Set("host", "")
			d.Set("hostgroup", "")
			d.SetId("")
			return diag.Errorf("Error configuring freeipa HBAC policy, host not assigned: %s", hostId)
		}
	}

	return nil
}

func resourceFreeIPADNSHBACPolicyHostMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the HBAC policy host membership")

	name, typeId, hostId, err := parseHBACPolicyHostMembershipID(d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleRemoveHostArgs{
		Cn: name,
	}
	optArgs := ipa.HbacruleRemoveHostOptionalArgs{}

	if typeId == "h" {
		v := []string{hostId}
		optArgs.Host = &v
	}
	if typeId == "hg" {
		v := []string{hostId}
		optArgs.Hostgroup = &v
	}

	_, err = client.HbacruleRemoveHost(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa the HBAC policy host membership: %s", err)
	}

	d.SetId("")

	return nil
}

func parseHBACPolicyHostMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine host membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	host := idParts[2]

	return name, _type, host, nil
}
