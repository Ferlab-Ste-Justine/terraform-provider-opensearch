package provider

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type TenantPermissionsModel struct {
	TenantPatterns []string `json:"tenant_patterns"`
	AllowedActions []string `json:"allowed_actions"`
}

type IndexPermissionsModel struct {
	IndexPatterns         []string      `json:"index_patterns"`
	AllowedActions        []string      `json:"allowed_actions"`
	MaskedFields          []string      `json:"masked_fields"`
	DocumentLevelSecurity string        `json:"dls"`
	FieldLevelSecurity    []string	    `json:"fls"`
}

type RoleModel struct {
	Name               string
	Description        string
	ClusterPermissions []string               `json:"cluster_permissions"`
    TenantPermissions  TenantPermissionsModel `json:"tenant_permissions"`
	IndexPermissions   IndexPermissionsModel `json:"index_permissions"`
}

func (reqCon *RequestContext) UpsertRole(role RoleModel) error {
	roleStr, marErr := json.Marshal(role)
    if marErr != nil {
        return marErr
    }

	res, err := reqCon.Do(
		"PUT", 
		path.Join("_plugins/_security/api/roles/", role.Name),
		string(roleStr), 
	)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}

func (reqCon *RequestContext) GetRole(name string) (*RoleModel, error) {
	res, err := reqCon.Do(
		"GET", 
		path.Join("_plugins/_security/api/roles/", name),
		"", 
	)
	
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, bErr := ioutil.ReadAll(res.Body)
	if bErr != nil {
		return nil, bErr
	}

	var role RoleModel
	uErr := json.Unmarshal(b, &role)
	if uErr != nil {
		return nil, uErr
	}
	
	return &role, nil
}

func (reqCon *RequestContext) DeleteRole(name string) error {
	res, err := reqCon.Do(
		"DELETE", 
		path.Join("_plugins/_security/api/roles/", name),
		"", 
	)
	
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}