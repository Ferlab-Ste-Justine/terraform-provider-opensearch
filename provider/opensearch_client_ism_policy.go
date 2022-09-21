package provider

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type EmptyModel struct {}

type IsmPsaReplicaCountModel struct {
	ReplicaCount int64 `json:"number_of_replicas"`
}

type IsmPsaIndexPriorityModel struct {
	IndexPriority int64 `json:"priority"`
}

type IsmPsaRetryModel struct {
	Count   int64  `json:"count"`
	Backoff string `json:"backoff,omitempty"`
	Delay   string `json:"delay,omitempty"`
}

type IsmPsActionModel struct {
	Timeout       string                     `json:"timeout,omitempty"`
	Retry         *IsmPsaRetryModel          `json:"retry,omitempty"`
	ReadOnly      *EmptyModel                `json:"read_only,omitempty"`
	ReadWrite     *EmptyModel                `json:"read_write,omitempty"`
	Open          *EmptyModel                `json:"open,omitempty"`
	Close         *EmptyModel                `json:"close,omitempty"`
	Delete        *EmptyModel                `json:"delete,omitempty"`
	ReplicaCount  *IsmPsaReplicaCountModel   `json:"replica_count,omitempty"`
	IndexPriority *IsmPsaIndexPriorityModel  `json:"index_priority,omitempty"`

}

type IsmPstConditionModel struct {
	MinIndexAge string `json:"min_index_age,omitempty"`
	MinDocCount *int64  `json:"min_doc_count,omitempty"`
	MinSize     string `json:"min_size,omitempty"`
}

type IsmPsTransitionModel struct {
	StateName  string					`json:"state_name"`
	Conditions []IsmPstConditionModel	`json:"conditions"`
}

type IsmPolicyStateModel struct {
    Name        string                          `json:"name"`
	Actions     []IsmPsActionModel		        `json:"actions"`
	Transitions []IsmPsTransitionModel          `json:"transitions"`
}

type IsmTemplateModel struct {
	Priority      int64		`json:"priority"`
	IndexPatterns []string	`json:"index_patterns"`
} 

type IsmPolicyModel struct {
	PolicyId     string                  `json:"-"`
	Description  string                  `json:"description"`
	IsmTemplate  IsmTemplateModel        `json:"ism_template"`
	DefaultState string                  `json:"default_state"`
	States       []IsmPolicyStateModel	 `json:"states"`
}

type IsmPolicyUpdateInfoModel struct {
	PrimaryTerm int64 `json:"_primary_term"`
	SeqNo       int64 `json:"_seq_no"`
}

func (reqCon *RequestContext) getIsmPolicyUpdateInfo(policyId string) (*IsmPolicyUpdateInfoModel, error) {
	return nil, nil
}

func (reqCon *RequestContext) UpsertIsmPolicy(ismPolicy IsmPolicyModel) error {
	ismPolicyStr, marErr := json.Marshal(ismPolicy)
    if marErr != nil {
        return marErr
    }

	res, err := reqCon.Do(
		"PUT", 
		path.Join("_plugins/_ism/policies/", ismPolicy.PolicyId),
		string(ismPolicyStr),
		[]int64{},
	)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}

func (reqCon *RequestContext) GetIsmPolicy(policyId string) (*IsmPolicyModel, error) {
	res, err := reqCon.Do(
		"GET", 
		path.Join("_plugins/_ism/policies/", policyId),
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

	policy := IsmPolicyModel{}
	uErr := json.Unmarshal(b, &policy)
	if uErr != nil {
		return nil, uErr
	}
	
	policy.PolicyId = policyId
	return &policy, nil
}

func (reqCon *RequestContext) DeleteIsmPolicy(policyId string) error {
	res, err := reqCon.Do(
		"DELETE", 
		path.Join("_plugins/_ism/policies/", policyId),
		"",
		[]int64{},
	)
	
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	return nil
}