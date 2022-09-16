package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpensearchRole() *schema.Resource {
	return &schema.Resource{
		Description: "Opensearch role.",
		Create: resourceOpensearchRoleCreate,
		Read:   resourceOpensearchRoleRead,
		Delete: resourceOpensearchRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
		},
	}
}

func resourceOpensearchRoleCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpensearchRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpensearchRoleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpensearchRoleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}