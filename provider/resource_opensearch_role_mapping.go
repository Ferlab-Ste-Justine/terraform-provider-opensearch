package provider

import (
	"errors"
	"fmt"
	
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

func roleMappingSchemaToModel(d *schema.ResourceData) RoleMappingModel {
	model := RoleMappingModel{
		Role: "", 
		BackendRoles: []string{}, 
		Hosts: []string{}, 
		Users: []string{},
	}

	role, _ := d.GetOk("role")
	model.Role = role.(string)

	backendRoles, backendRolesExist := d.GetOk("backend_roles")
	if backendRolesExist {
		for _, val := range (backendRoles.(*schema.Set)).List() {
			role := val.(string)
			model.BackendRoles = append(model.BackendRoles, role)
		}
	}

	hosts, hostsExist := d.GetOk("hosts")
	if hostsExist {
		for _, val := range (hosts.(*schema.Set)).List() {
			host := val.(string)
			model.Hosts = append(model.Hosts, host)
		}
	}

	users, usersExist := d.GetOk("users")
	if usersExist {
		for _, val := range (users.(*schema.Set)).List() {
			user := val.(string)
			model.Users = append(model.Users, user)
		}
	}

	return model
}

func resourceOpensearchRoleMappingCreate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	roleMapping := roleMappingSchemaToModel(d)

	err := cli.GetRequestContext().UpsertRoleMapping(roleMapping)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating role mapping for role '%s': %s", roleMapping.Role, err.Error()))
	}

	d.SetId(roleMapping.Role)
	return resourceOpensearchRoleMappingRead(d, meta)
}

func resourceOpensearchRoleMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	roleMapping := roleMappingSchemaToModel(d)

	err := cli.GetRequestContext().UpsertRoleMapping(roleMapping)
	if err != nil {
		return errors.New(fmt.Sprintf("Error updating role mapping for role '%s': %s", roleMapping.Role, err.Error()))
	}

	return resourceOpensearchRoleMappingRead(d, meta)
}

func resourceOpensearchRoleMappingRead(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	role := d.Id()

	roleMapping, err := cli.GetRequestContext().GetRoleMapping(role)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving role mapping for role '%s': %s", role, err.Error()))
	}

	d.Set("role", role)
	d.Set("backend_roles", roleMapping.BackendRoles)
	d.Set("hosts", roleMapping.Hosts)
	d.Set("users", roleMapping.Users)

	return nil
}

func resourceOpensearchRoleMappingDelete(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	role := d.Id()

	err := cli.GetRequestContext().DeleteRoleMapping(role)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting role mapping for role '%s': %s", role, err.Error()))
	}

	return nil
}