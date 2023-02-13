package loadbalancer

import (
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandACLConditions(rawConditions []interface{}) []loadbalancer.ACLCondition {
	var conditions []loadbalancer.ACLCondition
	for _, rawCondition := range rawConditions {
		condition := rawCondition.(map[string]interface{})

		arguments := make(map[string]loadbalancer.ACLArgument)
		for _, argument := range expandACLArguments(condition["argument"].(*schema.Set).List()) {
			arguments[argument.Name] = argument
		}

		conditions = append(conditions, loadbalancer.ACLCondition{
			Name:      condition["name"].(string),
			Arguments: arguments,
		})
	}

	return conditions
}

func flattenACLConditions(conditions []loadbalancer.ACLCondition) []map[string]interface{} {
	var flattenedConditions []map[string]interface{}
	for _, condition := range conditions {
		flattenedCondition := make(map[string]interface{})
		flattenedCondition["name"] = condition.Name
		flattenedCondition["arguments"] = flattenACLArguments(condition.Arguments)

		flattenedConditions = append(flattenedConditions, flattenedCondition)
	}

	return flattenedConditions
}

func expandACLActions(rawActions []interface{}) []loadbalancer.ACLAction {
	var actions []loadbalancer.ACLAction
	for _, rawAction := range rawActions {
		action := rawAction.(map[string]interface{})

		arguments := make(map[string]loadbalancer.ACLArgument)
		for _, argument := range expandACLArguments(action["argument"].(*schema.Set).List()) {
			arguments[argument.Name] = argument
		}

		actions = append(actions, loadbalancer.ACLAction{
			Name:      action["name"].(string),
			Arguments: arguments,
		})
	}

	return actions
}

func flattenACLActions(actions []loadbalancer.ACLAction) []map[string]interface{} {
	var flattenedActions []map[string]interface{}
	for _, action := range actions {
		flattenedAction := make(map[string]interface{})
		flattenedAction["name"] = action.Name
		flattenedAction["arguments"] = flattenACLArguments(action.Arguments)

		flattenedActions = append(flattenedActions, flattenedAction)
	}

	return flattenedActions
}

func expandACLArguments(rawArguments []interface{}) []loadbalancer.ACLArgument {
	var arguments []loadbalancer.ACLArgument
	for _, rawArgument := range rawArguments {
		argument := rawArgument.(map[string]interface{})

		arguments = append(arguments, loadbalancer.ACLArgument{
			Name:  argument["name"].(string),
			Value: argument["value"].(string),
		})
	}

	return arguments
}

func flattenACLArguments(arguments map[string]loadbalancer.ACLArgument) []map[string]interface{} {
	var flattenedArguments []map[string]interface{}
	for _, argument := range arguments {
		flattenedArguments = append(flattenedArguments, map[string]interface{}{
			"name":  argument.Name,
			"value": argument.Value,
		})
	}

	return flattenedArguments
}

func flattenACLArgument(argument loadbalancer.ACLArgument) (string, map[string]interface{}) {
	arguments := make(map[string]interface{})
	arguments["name"] = argument.Name
	arguments["value"] = argument.Value
	return argument.Name, arguments
}
