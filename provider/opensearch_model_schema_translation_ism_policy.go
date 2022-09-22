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

func ismPstConditionSchemaToModel(d *schema.ResourceData) IsmPstConditionModel {
	model := IsmPstConditionModel{}

	minIndexAge, minIndexAgeExists := d.GetOk("min_index_age")
	if minIndexAgeExists {
		model.MinIndexAge = minIndexAge.(string)
	}
	
	minDocCount, minDocCountExists := d.GetOk("min_doc_count")
	if minDocCountExists {
		minDocCountInt64 := int64(minDocCount.(int))
		model.MinDocCount = &minDocCountInt64
	}

	minSize, minSizeExists := d.GetOk("min_size")
	if minSizeExists {
		model.MinSize = minSize.(string)
	}

	return model
}

func ismStateTransitionSchemaToModel(d *schema.ResourceData) IsmPsTransitionModel {
	model := IsmPsTransitionModel{
		Conditions: []IsmPstConditionModel{},
	}

	stateName, _ := d.GetOk("name")
	model.StateName = stateName.(string)

	conditions, _ := d.GetOk("conditions")
	for _, val := range (conditions.(*schema.Set)).List() {
		model.Conditions = append(model.Conditions, ismPstConditionSchemaToModel(val.(*schema.ResourceData)))
	}

	return model
}

func ismStateSchemaToModel(d *schema.ResourceData) IsmPolicyStateModel {
	model := IsmPolicyStateModel{
		Actions:     []IsmPsActionModel{},
		Transitions: []IsmPsTransitionModel{},
	}

	name, _ := d.GetOk("name")
	model.Name = name.(string)

	actions, _ := d.GetOk("actions")
	for _, val := range (actions.(*schema.Set)).List() {
		model.Actions = append(model.Actions, ismStateActionSchemaToModel(val.(*schema.ResourceData)))
	}

	transitions, _ := d.GetOk("transitions")
	for _, val := range (transitions.(*schema.Set)).List() {
		model.Transitions = append(model.Transitions, ismStateTransitionSchemaToModel(val.(*schema.ResourceData)))
	}

	return model
}

func ismTemplateSchemaToModel(d *schema.ResourceData) IsmTemplateModel {
	model := IsmTemplateModel{
		IndexPatterns: []string{},
	}

	priority, priorityExists := d.GetOk("priority")
	if priorityExists {
		priorityInt64 := int64(priority.(int))
		model.Priority = &priorityInt64
	}

	indexPatterns, _ := d.GetOk("index_patterns")
	for _, val := range (indexPatterns.(*schema.Set)).List() {
		model.IndexPatterns = append(model.IndexPatterns, val.(string))
	}
	
	return model
}

func ismPolicySchemaToModel(d *schema.ResourceData) IsmPolicyModel {
	model := IsmPolicyModel{
		States: []IsmPolicyStateModel{},
	}

	policyId, _ := d.GetOk("policy_id")
	model.PolicyId = policyId.(string)

	description, descriptionExists := d.GetOk("description")
	if descriptionExists {
		model.Description = description.(string)
	}

	ismTemplate, ismTemplateExist := d.GetOk("ism_template")
	if ismTemplateExist {
		for _, val := range (ismTemplate.(*schema.Set)).List() {
			ismTemplateModel := ismTemplateSchemaToModel(val.(*schema.ResourceData))
			model.IsmTemplate = &ismTemplateModel
		}
	}

	defaultState, _ := d.GetOk("default_state")
	model.DefaultState = defaultState.(string)

	states, _ := d.GetOk("states")
	for _, val := range (states.(*schema.Set)).List() {
		model.States = append(model.States, ismStateSchemaToModel(val.(*schema.ResourceData)))
	}

	return model
}

func writeIsmPolicyModelToSchema(d *schema.ResourceData, m *IsmPolicyModel) {
	d.Set("description", m.Description)
	d.Set("default_state", m.DefaultState)

	if m.IsmTemplate != nil {
		ismTemplateElem := map[string]interface{}{
			"index_patterns": m.IsmTemplate.IndexPatterns,
		}
		
		if m.IsmTemplate.Priority != nil {
			ismTemplateElem["priority"] = (*m.IsmTemplate.Priority)
		}

		ismTemplate := make([]map[string]interface{}, 0)
		ismTemplate = append(ismTemplate, ismTemplateElem)
		d.Set("ism_template", ismTemplate)
	} else {
		d.Set("ism_template", nil)
	}

	states := make([]map[string]interface{}, 0)
	for _, v := range m.States {
		stateElem := map[string]interface{}{
			"name": v.Name,
		}

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

		transitions := make([]map[string]interface{}, 0)
		for _, t := range v.Transitions {
			transitionElem := map[string]interface{}{}

			transitionElem["state_name"] = t.StateName
			
			conditions := make([]map[string]interface{}, 0)
			for _, c := range t.Conditions {
				conditionElem := map[string]interface{}{}

				if c.MinIndexAge != "" {
					conditionElem["min_index_age"] = c.MinIndexAge
				}

				if c.MinDocCount != nil {
					conditionElem["min_doc_count"] = (*c.MinDocCount)
				}

				if c.MinSize != "" {
					conditionElem["min_size"] = c.MinSize
				}

				conditions = append(conditions, conditionElem)
			}
			transitionElem["conditions"] = conditions

			transitions = append(transitions, transitionElem)
		}
		stateElem["transitions"] = transitions

		states = append(states, stateElem)
	}
	d.Set("states", states)
}