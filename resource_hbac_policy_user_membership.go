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

func resourceFreeIPAHBACPolicyUserMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPADNSHBACPolicyUserMembershipCreate,
		ReadContext:   resourceFreeIPADNSHBACPolicyUserMembershipRead,
		DeleteContext: resourceFreeIPADNSHBACPolicyUserMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group"},
			},
			"group": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user"},
			},
		},
	}
}

func resourceFreeIPADNSHBACPolicyUserMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the HBAC policy user membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	user_id := "u"
	optArgs := ipa.HbacruleAddUserOptionalArgs{}

	args := ipa.HbacruleAddUserArgs{
		Cn: d.Get("name").(string),
	}

	if _v, ok := d.GetOkExists("user"); ok {
		v := []string{_v.(string)}
		optArgs.User = &v
		user_id = "u"
	}
	if _v, ok := d.GetOkExists("group"); ok {
		v := []string{_v.(string)}
		optArgs.Group = &v
		user_id = "g"
	}

	_, err = client.HbacruleAddUser(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the HBAC policy user membership: %s", err)
	}
	switch user_id {
	case "g":
		id := fmt.Sprintf("%s/g/%s", d.Get("name").(string), d.Get("group").(string))
		d.SetId(id)
	case "u":
		id := fmt.Sprintf("%s/u/%s", d.Get("name").(string), d.Get("user").(string))
		d.SetId(id)
	}

	return resourceFreeIPADNSHBACPolicyUserMembershipRead(ctx, d, meta)
}

func resourceFreeIPADNSHBACPolicyUserMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the HBAC policy user membership")

	name, typeId, userId, err := parseHBACPolicyUserMembershipID(d.Id())
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
		return diag.Errorf("Error show freeipa HBAC policy user membership: %s", err)
	}

	switch typeId {
	case "g":
		if res.Result.MemberuserGroup == nil || !slices.Contains(*res.Result.MemberuserGroup, userId) {
			log.Printf("[DEBUG] Warning! Group membership does not exist")
			d.Set("user", "")
			d.Set("group", "")
			d.SetId("")
			return diag.Errorf("Error configuring freeipa HBAC policy, group not assigned: %s", userId)
		}
	case "u":
		if res.Result.MemberuserUser == nil || !slices.Contains(*res.Result.MemberuserUser, userId) {
			log.Printf("[DEBUG] Warning! User membership does not exist")
			d.Set("user", "")
			d.Set("group", "")
			d.SetId("")
			return diag.Errorf("Error configuring freeipa HBAC policy, user not assigned: %s", userId)
		}
	}

	return nil
}

func resourceFreeIPADNSHBACPolicyUserMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the HBAC policy user membership")

	name, typeId, userId, err := parseHBACPolicyUserMembershipID(d.Id())

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	args := ipa.HbacruleRemoveUserArgs{
		Cn: name,
	}
	optArgs := ipa.HbacruleRemoveUserOptionalArgs{}

	if typeId == "u" {
		v := []string{userId}
		optArgs.User = &v
	}
	if typeId == "g" {
		v := []string{userId}
		optArgs.Group = &v
	}

	_, err = client.HbacruleRemoveUser(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa the HBAC policy user membership: %s", err)
	}

	d.SetId("")

	return nil
}

func parseHBACPolicyUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine user membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
