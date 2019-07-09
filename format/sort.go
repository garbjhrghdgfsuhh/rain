package format

import (
	"sort"

	"github.com/aws-cloudformation/rain/template"
)

var orders = map[string][]string{
	"Template": {
		"AWSTemplateFormatVersion",
		"Description",
		"Metadata",
		"Parameters",
		"Mappings",
		"Conditions",
		"Transform",
		"Resources",
		"Outputs",
	},
	"Parameter": {
		"Type",
		"Default",
	},
	"Transform": {
		"Name",
		"Parameters",
	},
	"Resource": {
		"Type",
	},
	"Outputs": {
		"Description",
		"Value",
		"Export",
	},
	"Policy": {
		"PolicyName",
		"PolicyDocument",
	},
	"PolicyDocument": {
		"Version",
		"Id",
		"Statement",
	},
	"PolicyStatement": {
		"Sid",
		"Effect",
		"Principal",
		"NotPrincipal",
		"Action",
		"NotAction",
		"Resource",
		"NotResource",
		"Condition",
	},
	"ResourceProperties": {
		"Name",
		"Description",
		"Type",
	},
}

func sortMapKeys(value map[string]interface{}) []string {
	keys := make([]string, 0)
	for key, _ := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func sortAs(keys []string, name string) []string {
	// Map keys to known and unknown list
	known := make([]string, 0)
	unknown := make([]string, 0)

	seen := make(map[string]bool)

	for _, o := range orders[name] {
		for _, key := range keys {
			if key == o {
				known = append(known, key)
				seen[key] = true
				break
			}
		}
	}

	for _, key := range keys {
		if !seen[key] {
			unknown = append(unknown, key)
		}
	}

	return append(known, unknown...)
}

func (p *encoder) sortKeys() []string {
	keys := sortMapKeys(p.currentValue.(map[string]interface{}))

	// Specific length paths
	if len(p.path) == 0 {
		return sortAs(keys, "Template")
	} else if len(p.path) == 1 {
		if p.path[0] == "Resources" {
			t := template.Template(p.data.Data.(map[string]interface{}))
			g := t.Graph()
			sort.Sort(g)

			output := make([]string, 0)
			for _, item := range g.Items() {
				el := item.(template.Element)

				if el.Type == "Resources" {
					output = append(output, el.Name)
				}
			}

			return output
		}
	} else if len(p.path) == 2 {
		switch p.path[0] {
		case "Parameters":
			return sortAs(keys, "Parameter")
		case "Resources":
			return sortAs(keys, "Resource")
		case "Outputs":
			return sortAs(keys, "Outputs")
		}
	} else if len(p.path) > 3 {
		if p.path[0] == "Resources" && p.path[2] == "Properties" {
			return sortAs(keys, "ResourceProperties")
		}
	} else if len(p.path) > 2 {
		if p.path[len(p.path)-2] == "Policies" {
			return sortAs(keys, "Policy")
		} else if p.path[len(p.path)-2] == "Statement" {
			return sortAs(keys, "PolicyStatement")
		}
	}

	// Paths that can live anywhere
	if p.path[0] == "Transform" || p.path[len(p.path)-1] == "Fn::Transform" {
		return sortAs(keys, "Transform")
	} else if p.path[len(p.path)-1] == "PolicyDocument" || p.path[len(p.path)-1] == "AssumeRolePolicyDocument" {
		return sortAs(keys, "PolicyDocument")
	}

	return keys
}