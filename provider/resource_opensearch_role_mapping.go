package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpensearchRoleMapping() *schema.Resource {
	return &schema.Resource{
		Description: "Opensearch role mapping.",
		Create: resourceOpensearchRoleMappingCreate,
		Read:   resourceOpensearchRoleMappingRead,
		Delete: resourceOpensearchRoleMappingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
		},
	}
}

func resourceOpensearchRoleMappingCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpensearchRoleMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpensearchRoleMappingRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpensearchRoleMappingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}