package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpensearchUser() *schema.Resource {
	return &schema.Resource{
		Description: "Opensearch user.",
		Create: resourceOpensearchUserCreate,
		Update: resourceOpensearchUserUpdate,
		Read:   resourceOpensearchUserRead,
		Delete: resourceOpensearchUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Description: "Username of the user.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"password": {
				Description: "Password of the user.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"opendistro_security_roles": {
				Description: "Prebuilt security roles to assign to the user.",
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backend_roles": {
				Description: "Custom roles to assign to the user.",
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

func userSchemaToModel(d *schema.ResourceData) UserModel {
	model := UserModel{Username: "", Password: "", SecurityRoles: []string{}, BackendRoles: []string{}}

	username, _ := d.GetOk("username")
	model.Username = username.(string)

	password, _ := d.GetOk("password")
	model.Password = password.(string)

	securityRoles, securityRolesExist := d.GetOk("opendistro_security_roles")
	if securityRolesExist {
		for _, val := range (securityRoles.(*schema.Set)).List() {
			role := val.(string)
			model.SecurityRoles = append(model.SecurityRoles, role)
		}
	}

	backendRoles, backendRolesExist := d.GetOk("backend_roles")
	if backendRolesExist {
		for _, val := range (backendRoles.(*schema.Set)).List() {
			role := val.(string)
			model.BackendRoles = append(model.BackendRoles, role)
		}
	}

	return model
}

func resourceOpensearchUserCreate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	user := userSchemaToModel(d)

	err := cli.GetRequestContext().UpsertUser(user)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating user '%s': %s", user.Username, err.Error()))
	}

	d.SetId(user.Username)
	return resourceOpensearchUserRead(d, meta)
}

func resourceOpensearchUserRead(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	username := d.Id()

	user, err := cli.GetRequestContext().GetUser(username)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving existing user '%s': %s", username, err.Error()))
	}

	d.Set("username", username)
	d.Set("opendistro_security_roles", user.SecurityRoles)
	d.Set("backend_roles", user.BackendRoles)

	return nil
}

func resourceOpensearchUserUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	user := userSchemaToModel(d)

	err := cli.GetRequestContext().UpsertUser(user)
	if err != nil {
		return errors.New(fmt.Sprintf("Error updating existing user '%s': %s", user.Username, err.Error()))
	}

	return resourceOpensearchUserRead(d, meta)
}

func resourceOpensearchUserDelete(d *schema.ResourceData, meta interface{}) error {
	username := d.Id()
	cli := meta.(OpensearchClient)

	err := cli.GetRequestContext().DeleteUser(username)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting existing user '%s': %s", username, err.Error()))
	}

	return nil
}