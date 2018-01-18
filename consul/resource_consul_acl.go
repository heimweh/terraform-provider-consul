package consul

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceConsulACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLCreate,
		Update: resourceConsulACLUpdate,
		Read:   resourceConsulACLRead,
		Delete: resourceConsulACLDelete,
		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "client",
				ValidateFunc: validation.StringInSlice([]string{
					"client",
					"management",
				}, false),
			},

			"rules": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConsulACLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	entry := &consulapi.ACLEntry{
		Type: d.Get("type").(string),
	}

	if v, ok := d.GetOk("acl_id"); ok {
		entry.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		entry.Name = v.(string)
	}

	if v, ok := d.GetOk("rules"); ok {
		entry.Rules = v.(string)
	}

	resp, _, err := client.ACL().Create(entry, nil)
	if err != nil {
		return err
	}

	d.SetId(resp)

	return resourceConsulACLRead(d, meta)
}

func resourceConsulACLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	acl, _, err := client.ACL().Info(d.Id(), nil)
	if err != nil {
		return err
	}

	d.Set("name", acl.Name)
	d.Set("type", acl.Type)
	d.Set("rules", acl.Rules)

	return nil
}

func resourceConsulACLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	acl, _, err := client.ACL().Info(d.Id(), nil)
	if err != nil {
		return err
	}

	if d.HasChange("name") {
		acl.Name = d.Get("name").(string)
	}

	if d.HasChange("type") {
		acl.Type = d.Get("type").(string)
	}

	if d.HasChange("rules") {
		acl.Rules = d.Get("rules").(string)
	}

	if _, err := client.ACL().Update(acl, nil); err != nil {
		return err
	}

	return resourceConsulACLRead(d, meta)
}

func resourceConsulACLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	if _, err := client.ACL().Destroy(d.Id(), nil); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
