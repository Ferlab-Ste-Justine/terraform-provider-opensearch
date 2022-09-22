package provider

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type UserModel struct {
	Username string        `json:"-"`
	Password string        `json:"password"`
	SecurityRoles []string `json:"opendistro_security_roles"`
	BackendRoles  []string `json:"backend_roles"`
}

func (reqCon *RequestContext) UpsertUser(user UserModel) error {
	userStr, marErr := json.Marshal(user)
    if marErr != nil {
        return marErr
    }

	res, err := reqCon.Do(
		"PUT", 
		path.Join("_plugins/_security/api/internalusers/", user.Username),
		"",
		string(userStr),
		[]int64{},
	)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}

func (reqCon *RequestContext) GetUser(username string) (*UserModel, error) {
	res, err := reqCon.Do(
		"GET", 
		path.Join("_plugins/_security/api/internalusers/", username),
		"",
		"",
		[]int64{},
	)
	
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, bErr := ioutil.ReadAll(res.Body)
	if bErr != nil {
		return nil, bErr
	}

	userMap := make(map[string]UserModel)
	uErr := json.Unmarshal(b, &userMap)
	if uErr != nil {
		return nil, uErr
	}
	
	user := userMap[username]
	user.Username = username
	return &user, nil
}

func (reqCon *RequestContext) DeleteUser(username string) error {
	res, err := reqCon.Do(
		"DELETE", 
		path.Join("_plugins/_security/api/internalusers/", username),
		"",
		"",
		[]int64{}, 
	)
	
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}