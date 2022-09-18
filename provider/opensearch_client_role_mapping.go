package provider

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type RoleMappingModel struct {
	Role         string                  `json:"-"`
	BackendRoles []string                `json:"backend_roles"`
    Hosts        []string
	Users        []string
}

func (reqCon *RequestContext) UpsertRoleMapping(roleMapping RoleMappingModel) error {
	roleMappingStr, marErr := json.Marshal(roleMapping)
    if marErr != nil {
        return marErr
    }

	res, err := reqCon.Do(
		"PUT", 
		path.Join("_plugins/_security/api/rolesmapping/", roleMapping.Role),
		string(roleMappingStr), 
	)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}

func (reqCon *RequestContext) GetRoleMapping(role string) (*RoleMappingModel, error) {
	res, err := reqCon.Do(
		"GET", 
		path.Join("_plugins/_security/api/rolesmapping/", role),
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

	roleMappingMap := make(map[string]RoleMappingModel)
	uErr := json.Unmarshal(b, &roleMappingMap)
	if uErr != nil {
		return nil, uErr
	}
	
	roleMapping := roleMappingMap[role]
	roleMapping.Role = role
	return &roleMapping, nil
}

func (reqCon *RequestContext) DeleteRoleMapping(role string) error {
	res, err := reqCon.Do(
		"DELETE", 
		path.Join("_plugins/_security/api/rolesmapping/", role),
		"", 
	)
	
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}