package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ismStateActionRetrySchemaToModel(d *schema.ResourceData) IsmPsaRetryModel {
	model := IsmPsaRetryModel{}

	count, _ := d.GetOk("count")
	model.Count = int64(count.(int))

	backoff, backoffExists := d.GetOk("backoff")
	if backoffExists {
		model.Backoff = backoff.(string)
	}

	delay, delayExists := d.GetOk("delay")
	if delayExists {
		model.Delay = delay.(string)
	}

	return model
}

func ismStateActionSchemaToModel(d *schema.ResourceData) IsmPsActionModel {
	model := IsmPsActionModel{}

	timeout, timeoutExists := d.GetOk("timeout")
	if timeoutExists {
		model.Timeout = timeout.(string)
	}

	retry, retryExists := d.GetOk("retry")
	if retryExists {
		for _, val := range (retry.(*schema.Set)).List() {
			retry := ismStateActionRetrySchemaToModel(val.(*schema.ResourceData))
			model.Retry = &retry
		}
	}

	action, _ := d.GetOk("action")
	actionStr := action.(string)
	switch actionStr {
	case "read_only":
		model.ReadOnly = &EmptyModel{}
	case "read_write":
		model.ReadWrite = &EmptyModel{}
	case "open":
		model.Open = &EmptyModel{}
	case "close":
		model.Close = &EmptyModel{}
	case "delete":
		model.Delete = &EmptyModel{}
	case "replica_count":
		replicaCountInt64 := int64(-1)
		replicaCount, replicaCountExists := d.GetOk("replica_count")
		if replicaCountExists {
			replicaCountInt64 = int64(replicaCount.(int))
		}
		model.ReplicaCount = &IsmPsaReplicaCountModel{
			ReplicaCount: replicaCountInt64,
		}
	case "index_priority":
		indexPriorityint64 := int64(-1)
		indexPriority, indexPriorityExists := d.GetOk("index_priority")
		if indexPriorityExists {
			indexPriorityint64 = int64(indexPriority.(int))
		}
		model.IndexPriority = &IsmPsaIndexPriorityModel{
			IndexPriority: indexPriorityint64,
		}
	}

	return model
}

func ismPstConditionSchemaToModel(d map[string]interface{}) IsmPstConditionModel {
	model := IsmPstConditionModel{}

	minIndexAge, minIndexAgeExists := d["min_index_age"]
	if minIndexAgeExists {
		model.MinIndexAge = minIndexAge.(string)
	}
	
	minDocCount, minDocCountExists := d["min_doc_count"]
	if minDocCountExists {
		model.MinDocCount = int64(minDocCount.(int))
	}

	minSize, minSizeExists := d["min_size"]
	if minSizeExists {
		model.MinSize = minSize.(string)
	}

	return model
}

func ismStateTransitionSchemaToModel(d map[string]interface{}) IsmPsTransitionModel {
	model := IsmPsTransitionModel{}

	stateName := d["state_name"]
	model.StateName = stateName.(string)

	conditions := d["conditions"]
	for _, val := range (conditions.(*schema.Set)).List() {
		model.Conditions = ismPstConditionSchemaToModel(val.(map[string]interface{}))
	}

	return model
}

func ismStateSchemaToModel(d map[string]interface{}) IsmPolicyStateModel {
	model := IsmPolicyStateModel{
		Actions:     []IsmPsActionModel{},
		Transitions: []IsmPsTransitionModel{},
	}

	name, _ := d["name"]
	model.Name = name.(string)

	actions, actionsExist := d["actions"]
	if actionsExist {
		for _, val := range actions.([]interface{}) {
			model.Actions = append(model.Actions, ismStateActionSchemaToModel(val.(*schema.ResourceData)))
		}
	}

	transitions, transitionsExist := d["transitions"]
	if transitionsExist {
		for _, val := range transitions.([]interface{}) {
			model.Transitions = append(model.Transitions, ismStateTransitionSchemaToModel(val.(map[string]interface{})))
		}
	}

	return model
}

func ismTemplateSchemaToModel(d map[string]interface{}) IsmTemplateModel {
	model := IsmTemplateModel{
		IndexPatterns: []string{},
	}

	priority, priorityExists := d["priority"]
	if priorityExists {
		priorityInt64 := int64(priority.(int))
		model.Priority = &priorityInt64
	}

	indexPatterns, _ := d["index_patterns"]
	for _, val := range (indexPatterns.(*schema.Set)).List() {
		model.IndexPatterns = append(model.IndexPatterns, val.(string))
	}
	
	return model
}

func ismPolicySchemaToModel(d *schema.ResourceData) IsmPolicyModel {
	model := IsmPolicyModel{
		States:      []IsmPolicyStateModel{},
		IsmTemplate: []IsmTemplateModel{},
	}

	policyId, _ := d.GetOk("policy_id")
	model.PolicyId = policyId.(string)

	description, _ := d.GetOk("description")
	model.Description = description.(string)

	ismTemplate, ismTemplateExist := d.GetOk("ism_template")
	if ismTemplateExist {
		for _, val := range (ismTemplate.(*schema.Set)).List() {
			ismTemplateModel := ismTemplateSchemaToModel(val.(map[string]interface{}))
			model.IsmTemplate = append(model.IsmTemplate, ismTemplateModel)
		}
	}

	defaultState, _ := d.GetOk("default_state")
	model.DefaultState = defaultState.(string)

	states, _ := d.GetOk("states")
	for _, val := range (states.(*schema.Set)).List() {
		model.States = append(model.States, ismStateSchemaToModel(val.(map[string]interface{})))
	}

	return model
}

func writeIsmPolicyModelToSchema(d *schema.ResourceData, m *IsmPolicyModel) {
	d.Set("description", m.Description)
	d.Set("default_state", m.DefaultState)

	if len(m.IsmTemplate) > 0 {
		ismTemplate := make([]map[string]interface{}, 0)
		
		for _, val := range m.IsmTemplate {
			ismTemplateElem := map[string]interface{}{
				"index_patterns": val.IndexPatterns,
			}
			
			if val.Priority != nil {
				ismTemplateElem["priority"] = (*val.Priority)
			}
			
			ismTemplate = append(ismTemplate, ismTemplateElem)
		}

		d.Set("ism_template", ismTemplate)
	} else {
		d.Set("ism_template", nil)
	}

	states := make([]map[string]interface{}, 0)
	for _, v := range m.States {
		stateElem := map[string]interface{}{
			"name": v.Name,
		}

		if len(v.Actions) > 0 {
			actions := make([]map[string]interface{}, 0)
			for _, a := range v.Actions {
				actionElem := map[string]interface{}{}
				if a.Timeout != "" {
					actionElem["timeout"] = a.Timeout
				}
	
				if a.Retry != nil {
					retryElem := map[string]interface{}{}
	
					retryElem["count"] = a.Retry.Count
					if a.Retry.Backoff != "" {
						retryElem["backoff"] = a.Retry.Backoff
					}
					if a.Retry.Delay != "" {
						retryElem["delay"] = a.Retry.Delay
					}
	
					actionElem["retry"] = []map[string]interface{}{retryElem}
				}
	
				if a.ReadOnly != nil {
					actionElem["action"] = "read_only"
				} else if a.ReadWrite != nil {
					actionElem["action"] = "read_write"
				} else if a.Open != nil {
					actionElem["action"] = "open"
				} else if a.Close != nil {
					actionElem["action"] = "close"
				} else if a.Delete != nil {
					actionElem["action"] = "delete"
				} else if a.ReplicaCount != nil {
					actionElem["action"] = "replica_count"
					actionElem["replica_count"] = a.ReplicaCount.ReplicaCount
				} else if a.IndexPriority != nil {
					actionElem["action"] = "index_priority"
					actionElem["index_priority"] = a.IndexPriority.IndexPriority
				}
	
				actions = append(actions, actionElem)
			}
			stateElem["actions"] = actions
		}

		if len(v.Transitions) > 0 {
			transitions := make([]map[string]interface{}, 0)
			for _, t := range v.Transitions {
				transitionElem := map[string]interface{}{}
	
				transitionElem["state_name"] = t.StateName
				
				conditions := map[string]interface{}{}

				c := t.Conditions
				if c.MinIndexAge != "" {
					conditions["min_index_age"] = c.MinIndexAge
				}

				if c.MinDocCount > 0 {
					conditions["min_doc_count"] = c.MinDocCount
				}

				if c.MinSize != "" {
					conditions["min_size"] = c.MinSize
				}

				transitionElem["conditions"] = []map[string]interface{}{conditions}
	
				transitions = append(transitions, transitionElem)
			}
			stateElem["transitions"] = transitions
		}

		states = append(states, stateElem)
	}
	d.Set("states", states)
}