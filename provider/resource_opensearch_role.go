package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpensearchRole() *schema.Resource {
	return &schema.Resource{
		Description: "Opensearch role.",
		Create: resourceOpensearchRoleCreate,
		Update: resourceOpensearchRoleUpdate,
		Read:   resourceOpensearchRoleRead,
		Delete: resourceOpensearchRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the role.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Description: "Description of the role.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"cluster_permissions": {
				Description: "Permissions for cluster wide actions the role has.",
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_permissions": {
				Description: "Permissions for tenant access the role has.",
                Type: schema.TypeSet,
                Optional: true,
                ForceNew: false,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "tenant_patterns": {
                            Type: schema.TypeSet,
                            Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
                        },
                        "allowed_actions": {
                            Type: schema.TypeString,
                            Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
                        },
                    },
                },
			},
			"index_permissions": {
				Description: "Permissions for index access the role has.",
                Type: schema.TypeSet,
                Optional: true,
                ForceNew: false,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "index_patterns": {
                            Type: schema.TypeSet,
                            Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
                        },
                        "allowed_actions": {
                            Type: schema.TypeString,
                            Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
                        },
						"masked_fields": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"document_level_security": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"field_level_security": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
                    },
                },
			},
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