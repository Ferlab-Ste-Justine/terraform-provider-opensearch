package provider

import (
	"errors"
	"fmt"

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
                            Type: schema.TypeSet,
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
                            Type: schema.TypeSet,
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

func tenantPermissionSchemaToModel(d map[string]interface{}) TenantPermissionModel {
	model := TenantPermissionModel{
		TenantPatterns: []string{},
		AllowedActions: []string{},
	}

	tenantPatterns, tenantPatternsExist := d["tenant_patterns"]
	if tenantPatternsExist {
		for _, val := range (tenantPatterns.(*schema.Set)).List() {
			tenantPattern := val.(string)
			model.TenantPatterns = append(model.TenantPatterns, tenantPattern)
		}
	}

	allowedActions, allowedActionsExist := d["allowed_actions"]
	if allowedActionsExist {
		for _, val := range (allowedActions.(*schema.Set)).List() {
			allowedAction := val.(string)
			model.AllowedActions = append(model.AllowedActions, allowedAction)
		}
	}

	return model
}

func indexPermissionSchemaToModel(d map[string]interface{}) IndexPermissionModel {
	model := IndexPermissionModel{
		IndexPatterns:         []string{},
		AllowedActions:        []string{},
		MaskedFields:          []string{},
		DocumentLevelSecurity: "",
		FieldLevelSecurity:    []string{},
	}

	indexPatterns, indexPatternsExist := d["index_patterns"]
	if indexPatternsExist {
		for _, val := range (indexPatterns.(*schema.Set)).List() {
			indexPattern := val.(string)
			model.IndexPatterns = append(model.IndexPatterns, indexPattern)
		}
	}

	allowedActions, allowedActionsExist := d["allowed_actions"]
	if allowedActionsExist {
		for _, val := range (allowedActions.(*schema.Set)).List() {
			allowedAction := val.(string)
			model.AllowedActions = append(model.AllowedActions, allowedAction)
		}
	}

	maskedFields, maskedFieldsExist := d["masked_fields"]
	if maskedFieldsExist {
		for _, val := range (maskedFields.(*schema.Set)).List() {
			maskedField := val.(string)
			model.MaskedFields = append(model.MaskedFields, maskedField)
		}
	}

	documentLevelSecurity, documentLevelSecurityExist := d["document_level_security"]
	if documentLevelSecurityExist {
		model.DocumentLevelSecurity = documentLevelSecurity.(string)
	}

	fieldLevelSecurity, fieldLevelSecurityExist := d["field_level_security"]
	if fieldLevelSecurityExist {
		for _, val := range (fieldLevelSecurity.(*schema.Set)).List() {
			field := val.(string)
			model.FieldLevelSecurity = append(model.FieldLevelSecurity, field)
		}
	}

	return model
}

func roleSchemaToModel(d *schema.ResourceData) RoleModel {
	model := RoleModel{
		Name:               "",
		ClusterPermissions: []string{},
		TenantPermissions:  []TenantPermissionModel{},
		IndexPermissions:   []IndexPermissionModel{},
	}

	name, _ := d.GetOk("name")
	model.Name = name.(string)

	clusterPermissions, clusterPermissionsExist := d.GetOk("cluster_permissions")
	if clusterPermissionsExist {
		for _, val := range (clusterPermissions.(*schema.Set)).List() {
			permission := val.(string)
			model.ClusterPermissions = append(model.ClusterPermissions, permission)
		}
	}

	tenantPermissions, tenantPermissionsExist := d.GetOk("tenant_permissions")
	if tenantPermissionsExist {
		for _, val := range (tenantPermissions.(*schema.Set)).List() {
			model.TenantPermissions = append(model.TenantPermissions, tenantPermissionSchemaToModel(val.(map[string]interface{})))
		}
	}

	indexPermissions, indexPermissionsExist := d.GetOk("index_permissions")
	if indexPermissionsExist {
		for _, val := range (indexPermissions.(*schema.Set)).List() {
			model.IndexPermissions = append(model.IndexPermissions, indexPermissionSchemaToModel(val.(map[string]interface{})))
		}
	}

	return model
}

func resourceOpensearchRoleCreate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	role := roleSchemaToModel(d)

	err := cli.GetRequestContext().UpsertRole(role)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating role '%s': %s", role.Name, err.Error()))
	}

	d.SetId(role.Name)
	return resourceOpensearchRoleRead(d, meta)
}

func resourceOpensearchRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	role := roleSchemaToModel(d)

	err := cli.GetRequestContext().UpsertRole(role)
	if err != nil {
		return errors.New(fmt.Sprintf("Error updating existing role '%s': %s", role.Name, err.Error()))
	}

	return resourceOpensearchRoleRead(d, meta)
}

func resourceOpensearchRoleRead(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	name := d.Id()

	role, err := cli.GetRequestContext().GetRole(name)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving existing role '%s': %s", name, err.Error()))
	}

	d.Set("name", name)
	d.Set("cluster_permissions", role.ClusterPermissions)
	
    tenantPermissions := make([]map[string]interface{}, 0)
    for _, v := range role.TenantPermissions {
        tenantPermissions = append(tenantPermissions, map[string]interface{}{
            "tenant_patterns": v.TenantPatterns,
            "allowed_actions": v.AllowedActions,
        })
    }
    d.Set("tenant_permissions", tenantPermissions)

    indexPermissions := make([]map[string]interface{}, 0)
    for _, v := range role.IndexPermissions {
        indexPermissions = append(indexPermissions, map[string]interface{}{
            "index_patterns": v.IndexPatterns,
            "allowed_actions": v.AllowedActions,
			"masked_fields": v.MaskedFields,
			"document_level_security": v.DocumentLevelSecurity,
			"field_level_security": v.FieldLevelSecurity,
        })
    }
    d.Set("index_permissions", indexPermissions)

	return nil
}

func resourceOpensearchRoleDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	cli := meta.(OpensearchClient)

	err := cli.GetRequestContext().DeleteRole(name)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting existing role '%s': %s", name, err.Error()))
	}

	return nil
}