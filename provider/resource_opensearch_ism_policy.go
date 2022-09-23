package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpensearchIsmPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Opensearch ism policy. Not all the options supported by the api are supported by the provider at this time.",
		Create: resourceOpensearchIsmPolicyCreate,
		Update: resourceOpensearchIsmPolicyUpdate,
		Read:   resourceOpensearchIsmPolicyRead,
		Delete: resourceOpensearchIsmPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			//Missing
			//Fields: error_notification
			"policy_id": {
				Description: "Unique identifier for the policy.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Description: "Description for the policy.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"ism_template": {
				Description: "Match of the indices to apply the policy on.",
				Type:     schema.TypeSet,
                Optional: true,
                ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"index_patterns": {
							Description: "Indexes to include with wildcard support.",
							Type: schema.TypeSet,
							Required:     true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"default_state": {
				Description: "Default states that indices will have.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"states": {
				Description: "Permissions for index access the role has.",
                Type:        schema.TypeSet,
				Required:    true,
                ForceNew:    false,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the state.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"actions": {
							Description: "Actions that should be run when an index reach the state.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout": {
										Description: "Time limit to perform the action",
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"retry": {
										Description: "Retry policy when the action fails",
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"count": {
													Description:  "Number of retries",
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
												"backoff": {
													Description: "Backoff policy when retrying. Can be: Exponential, Constant and Linear",
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(
														[]string{
															"Exponential",
															"Linear",
															"Constant",
														}, 
														false,
													),
												},
												"delay": {
													Description: "Base time to wait between retries",
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
									"action": {
										//Missing
										//Actions: allocation, snapshot, notification, rollover, shrink, force_merge
										Description: "The action to execute. Currently supports: read_only, read_write, replica_count, open, close, delete, index_priority",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"read_only",
												"read_write",
												"replica_count",
												"open",
												"close",
												"delete",
												"index_prioty",
											}, 
											false,
										),
									},
									"index_priority": {
										Description: "Priority to set for the index if the action is index_priority",
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"replica_count": {
										Description: "Replicat count to set for the index if the action is replica_count",
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"transitions": {
							Description: "Transition specifications for when an index should transition to another state.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state_name": {
										Description: "Name of the state to transition to.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"conditions": {
										Description: "Conditions that trigger the state change.",
										Type:        schema.TypeSet,
										Required:    true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												//Missing
												//fields: cron
												"min_index_age": {
													Description: "Minimum age at which the index will transition.",
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"min_doc_count": {
													Description: "Minimum number of documents after which the index will transition.",
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
												"min_size": {
													Description: "Minimum size (not counting replication) after which the index will transition.",
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},
							},
						},
                    },
                },
			},
		},
	}
}

func resourceOpensearchIsmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	policyId := d.Id()

	policy, err := cli.GetRequestContext().GetIsmPolicy(policyId)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving existing policy '%s': %s", policyId, err.Error()))
	}

	writeIsmPolicyModelToSchema(d, policy)

	return nil
}

func resourceOpensearchIsmPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	policy := ismPolicySchemaToModel(d)

	err := cli.GetRequestContext().UpsertIsmPolicy(policy)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating policy '%s': %s", policy.PolicyId, err.Error()))
	}

	d.SetId(policy.PolicyId)
	return resourceOpensearchIsmPolicyRead(d, meta)
}

func resourceOpensearchIsmPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	policy := ismPolicySchemaToModel(d)

	err := cli.GetRequestContext().UpsertIsmPolicy(policy)
	if err != nil {
		return errors.New(fmt.Sprintf("Error updating existing policy '%s': %s", policy.PolicyId, err.Error()))
	}

	return resourceOpensearchIsmPolicyRead(d, meta)
}

func resourceOpensearchIsmPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	cli := meta.(OpensearchClient)
	policyId := d.Id()

	err := cli.GetRequestContext().DeleteIsmPolicy(policyId)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting existing policy '%s': %s", policyId, err.Error()))
	}

	return nil
}