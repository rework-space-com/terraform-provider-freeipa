package freeipa

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

func resourceFreeIPAHBACServiceServiceGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAHBACServiceServiceGroupMembershipCreate,
		ReadContext:   resourceFreeIPAHBACServiceServiceGroupMembershipRead,
		DeleteContext: resourceFreeIPAHBACServiceServiceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the HBAC service group",
			},
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HBAC service to add to the group",
			},
		},
	}
}

func resourceFreeIPAHBACServiceServiceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa hbac service group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.HbacsvcgroupAddMemberOptionalArgs{}

	args := ipa.HbacsvcgroupAddMemberArgs{
		Cn: d.Get("name").(string),
	}
	if _v, ok := d.GetOkExists("service"); ok {
		v := []string{_v.(string)}
		optArgs.Hbacsvc = &v
	}

	_, err = client.HbacsvcgroupAddMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa hbac service group membership: %s", err)
	}

	id := fmt.Sprintf("%s/svc/%s", encodeSlash(d.Get("name").(string)), d.Get("service").(string))
	d.SetId(id)

	return resourceFreeIPAHBACServiceServiceGroupMembershipRead(ctx, d, meta)
}

func resourceFreeIPAHBACServiceServiceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa hbac service group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	hbacsvcgroupId, typeId, svcId, err := parseHBACServiceGroupMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_hbac_service_servicegroup_membership: %s", err)
	}

	all := true
	optArgs := ipa.HbacsvcgroupShowOptionalArgs{
		All: &all,
	}

	args := ipa.HbacsvcgroupShowArgs{
		Cn: hbacsvcgroupId,
	}

	res, err := client.HbacsvcgroupShow(&args, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			d.SetId("")
			log.Printf("[DEBUG] HBAC service group not found")
			return nil
		} else {
			return diag.Errorf("Error reading freeipa HBAC service group: %s", err)
		}
	}

	switch typeId {
	case "svc":
		if res.Result.MemberHbacsvc == nil || !slices.Contains(*res.Result.MemberHbacsvc, svcId) {
			log.Printf("[DEBUG] Warning! HBAC service group membership does not exist")
			d.Set("name", "")
			d.Set("hbacsvc", "")
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return diag.Errorf("Error show freeipa hbac service group membership: %s", err)
	}

	log.Printf("[DEBUG] Read freeipa hbac service group membership %s", res.Result.Cn)
	return nil
}

func resourceFreeIPAHBACServiceServiceGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa hbac service group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	hbacsvcgroupId, typeId, svcId, err := parseHBACServiceGroupMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_hbac_service_servicegroup_membership: %s", err)
	}

	optArgs := ipa.HbacsvcgroupRemoveMemberOptionalArgs{}

	args := ipa.HbacsvcgroupRemoveMemberArgs{
		Cn: hbacsvcgroupId,
	}

	switch typeId {
	case "svc":
		v := []string{svcId}
		optArgs.Hbacsvc = &v
	}

	_, err = client.HbacsvcgroupRemoveMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa hbac service group membership: %s", err)
	}

	d.SetId("")
	return nil
}

func parseHBACServiceGroupMembershipID(id string) (string, string, string, error) {
	idParts := strings.SplitN(id, "/", 3)
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine hbac service group membership ID %s", id)
	}

	name := decodeSlash(idParts[0])
	_type := idParts[1]
	hbacsvc := idParts[2]

	return name, _type, hbacsvc, nil
}
