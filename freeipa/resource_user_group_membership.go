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

func resourceFreeIPAUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFreeIPAUserGroupMembershipCreate,
		ReadContext:   resourceFreeIPAUserGroupMembershipRead,
		DeleteContext: resourceFreeIPAUserGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Group name",
			},
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group", "external_member"},
				Description:   "User to add",
			},
			"group": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user", "external_member"},
				Description:   "Group to add",
			},
			"external_member": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user", "group"},
				Description:   "External member to add. name must refer to an external group. (Requires a valid AD Trust configuration).",
			},
		},
	}
}

func resourceFreeIPAUserGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating freeipa the user group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	user_id := "u"
	name := d.Get("name").(string)
	optArgs := ipa.GroupAddMemberOptionalArgs{}

	args := ipa.GroupAddMemberArgs{
		Cn: name,
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

	if _v, ok := d.GetOkExists("external_member"); ok {
		v := []string{_v.(string)}
		optArgs.Ipaexternalmember = &v
		user_id = "e"
	}

	_v, err := client.GroupAddMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error creating freeipa the user group membership: %s", err)
	}
	log.Printf("[DEBUG] Group add member response for group %s is %v", args.Cn, _v.Result)
	log.Printf("[DEBUG] Group add member failure for group %s is %v", args.Cn, _v.Failed)
	log.Printf("[DEBUG] Group add member number added for group %s is %d", args.Cn, _v.Completed)
	if _v.Completed == 0 {
		return diag.Errorf("Error creating freeipa the user group membership: %v", _v.Failed)
	}

	if _v, ok := d.GetOkExists("external_member"); ok {
		v := _v.(string)
		z := new(bool)
		*z = true
		groupRes, err := client.GroupShow(&ipa.GroupShowArgs{Cn: name}, &ipa.GroupShowOptionalArgs{All: z})
		if err != nil {
			return diag.Errorf("Error looking up freeipa user group membership: %s", err)
		}
		if !slices.Contains(*groupRes.Result.Ipaexternalmember, v) {
			_, err = client.GroupRemoveMember(&ipa.GroupRemoveMemberArgs{Cn: name}, &ipa.GroupRemoveMemberOptionalArgs{Ipaexternalmember: &[]string{v}})
			if err != nil {
				return diag.Errorf("Error deleting invalid freeipa user group membership: %s", err)
			}
			return diag.Errorf("Error, external member is not using the correct format. Use the lowercase upn format (ie: 'domain users@domain.net'): %s", v)
		} else {
			log.Printf("[DEBUG] group show %s is %v", name, groupRes.Result.String())
		}
	}

	switch user_id {
	case "g":
		id := fmt.Sprintf("%s/g/%s", name, d.Get("group").(string))
		d.SetId(id)
	case "u":
		id := fmt.Sprintf("%s/u/%s", name, d.Get("user").(string))
		d.SetId(id)
	case "e":
		id := fmt.Sprintf("%s/e/%s", name, d.Get("external_member").(string))
		d.SetId(id)
	}

	return resourceFreeIPAUserGroupMembershipRead(ctx, d, meta)
}

func resourceFreeIPAUserGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read freeipa the user group membership")

	name, typeId, userId, err := parseUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_user_group_membership: %s", err)
	}

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	reqArgs := ipa.GroupShowArgs{
		Cn: name,
	}
	z := new(bool)
	*z = true
	optArgs := ipa.GroupShowOptionalArgs{
		All: z,
	}

	res, err := client.GroupShow(&reqArgs, &optArgs)
	if err != nil {
		return diag.Errorf("Error find freeipa the user group membership: %s", err)
	}

	log.Printf("[DEBUG] group show %s is %v", name, res.Result.String())

	switch typeId {
	case "g":
		v := []string{userId}
		groups := *res.Result.MemberGroup
		log.Printf("[DEBUG] Group list in group %s is %v", name, groups)
		if slices.Contains(groups, v[0]) {
			return nil
		}
	case "u":
		v := []string{userId}
		users := *res.Result.MemberUser
		log.Printf("[DEBUG] User list in group %s is %v", name, users)
		if slices.Contains(users, v[0]) {
			return nil
		}
	case "e":
		v := []string{userId}
		extmembers := *res.Result.Ipaexternalmember
		log.Printf("[DEBUG] External member list in group %s is %v", name, extmembers)
		if slices.Contains(extmembers, v[0]) {
			return nil
		}
	}
	log.Printf("[DEBUG] Warning! Group or User membership not exist")
	d.Set("user", "")
	d.Set("group", "")
	d.Set("external_member", "")
	d.SetId("")

	return nil
}

func resourceFreeIPAUserGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Delete freeipa the user group membership")

	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.Errorf("Error creating freeipa identity client: %s", err)
	}

	optArgs := ipa.GroupRemoveMemberOptionalArgs{}

	nameId, typeId, userId, err := parseUserMembershipID(d.Id())

	if err != nil {
		return diag.Errorf("Error parsing ID of freeipa_user_group_membership: %s", err)
	}

	args := ipa.GroupRemoveMemberArgs{
		Cn: nameId,
	}

	switch typeId {
	case "g":
		v := []string{userId}
		optArgs.Group = &v
	case "u":
		v := []string{userId}
		optArgs.User = &v
	case "e":
		v := []string{userId}
		optArgs.Ipaexternalmember = &v
	}

	_, err = client.GroupRemoveMember(&args, &optArgs)
	if err != nil {
		return diag.Errorf("Error delete freeipa the user group membership: %s", err)
	}

	d.SetId("")

	return nil
}

func parseUserMembershipID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine user membership ID %s", id)
	}

	name := idParts[0]
	_type := idParts[1]
	user := idParts[2]

	return name, _type, user, nil
}
