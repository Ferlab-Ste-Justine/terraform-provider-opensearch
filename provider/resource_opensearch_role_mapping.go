package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpensearchRoleMapping() *schema.Resource {
	return &schema.Resource{
		Description: "Opensearch role mapping to map backend roles, users and hosts to a given role.",
		Create: resourceOpensearchRoleMappingCreate,
		Update: resourceOpensearchRoleMappingUpdate,
		Read:   resourceOpensearchRoleMappingRead,
		Delete: resourceOpensearchRoleMappingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role": {
				Description:  "Role that things should be mapped to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"backend_roles": {
				Description: "Backend roles to map to the role.",
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hosts": {
				Description: "Hosts to map to the role.",
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": {
				Description: "Users to map to the role.",
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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