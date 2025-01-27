package graphqlizer

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/pkg/errors"
)

// Graphqlizer is responsible for converting Go objects to input arguments in graphql format
type Graphqlizer struct{}

func (g *Graphqlizer) ApplicationRegisterInputToGQL(in graphql.ApplicationRegisterInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .ProviderName }}
		providerName: "{{ .ProviderName }}",
		{{- end }}
		{{- if .Description }}
		description: "{{ .Description }}",
		{{- end }}
        {{- if .Labels }}
		labels: {{ LabelsToGQL .Labels}},
		{{- end }}
		{{- if .Webhooks }}
		webhooks: [
			{{- range $i, $e := .Webhooks }} 
				{{- if $i}}, {{- end}} {{ WebhookInputToGQL $e }}
			{{- end }} ],
		{{- end}}
		{{- if .HealthCheckURL }}
		healthCheckURL: "{{ .HealthCheckURL }}",
		{{- end }}
		{{- if .Bundles }} 
		bundles: [
			{{- range $i, $e := .Bundles }} 
				{{- if $i}}, {{- end}} {{- BundleCreateInputToGQL $e }}
			{{- end }} ],
		{{- end }}
		{{- if .IntegrationSystemID }}
		integrationSystemID: "{{ .IntegrationSystemID }}",
		{{- end }}
		{{- if .StatusCondition }}
		statusCondition: {{ .StatusCondition }}
		{{- end }}
	}`)
}

func (g *Graphqlizer) ApplicationUpdateInputToGQL(in graphql.ApplicationUpdateInput) (string, error) {
	return g.genericToGQL(in, `{
		{{- if .ProviderName }}
		providerName: "{{ .ProviderName }}",
		{{- end }}
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
		{{- if .HealthCheckURL }}
		healthCheckURL: "{{ .HealthCheckURL }}",
		{{- end }}
		{{- if .IntegrationSystemID }}
		integrationSystemID: "{{ .IntegrationSystemID }}",
		{{- end }}
		{{- if .StatusCondition }}
		statusCondition: {{ .StatusCondition }}
		{{- end }}
	}`)
}

func (g *Graphqlizer) ApplicationTemplateInputToGQL(in graphql.ApplicationTemplateInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
		applicationInput: {{ ApplicationRegisterInputToGQL .ApplicationInput}},
		{{- if .Placeholders }}
		placeholders: [
			{{- range $i, $e := .Placeholders }} 
				{{- if $i}}, {{- end}} {{ PlaceholderDefinitionInputToGQL $e }}
			{{- end }} ],
		{{- end }}
		accessLevel: {{.AccessLevel}},
		{{- if .Webhooks }}
		webhooks: [
			{{- range $i, $e := .Webhooks }} 
				{{- if $i}}, {{- end}} {{ WebhookInputToGQL $e }}
			{{- end }} ],
		{{- end}}	
	}`)
}

func (g *Graphqlizer) ApplicationTemplateUpdateInputToGQL(in graphql.ApplicationTemplateUpdateInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
		applicationInput: {{ ApplicationRegisterInputToGQL .ApplicationInput}},
		{{- if .Placeholders }}
		placeholders: [
			{{- range $i, $e := .Placeholders }} 
				{{- if $i}}, {{- end}} {{ PlaceholderDefinitionInputToGQL $e }}
			{{- end }} ],
		{{- end }}
		accessLevel: {{.AccessLevel}},
	}`)
}

func (g *Graphqlizer) DocumentInputToGQL(in *graphql.DocumentInput) (string, error) {
	return g.genericToGQL(in, `{
		title: "{{.Title}}",
		displayName: "{{.DisplayName}}",
		description: "{{.Description}}",
		format: {{.Format}},
		{{- if .Kind }}
		kind: "{{.Kind}}",
		{{- end}}
		{{- if .Data }}
		data: "{{.Data}}",
		{{- end}}
		{{- if .FetchRequest }}
		fetchRequest: {{- FetchRequesstInputToGQL .FetchRequest }},
		{{- end}}
}`)
}

func (g *Graphqlizer) FetchRequestInputToGQL(in *graphql.FetchRequestInput) (string, error) {
	return g.genericToGQL(in, `{
		url: "{{.URL}}",
		{{- if .Auth }}
		auth: {{- AuthInputToGQL .Auth }},
		{{- end }}
		{{- if .Mode }}
		mode: {{.Mode}},
		{{- end}}
		{{- if .Filter}}
		filter: "{{.Filter}}",
		{{- end}}
	}`)
}

func (g *Graphqlizer) CredentialRequestAuthInputToGQL(in *graphql.CredentialRequestAuthInput) (string, error) {
	return g.genericToGQL(in, `{
		{{- if .Csrf }}
		csrf: {{ CSRFTokenCredentialRequestAuthInputToGQL .Csrf }},
		{{- end }}
	}`)
}

func (g *Graphqlizer) CredentialDataInputToGQL(in *graphql.CredentialDataInput) (string, error) {
	return g.genericToGQL(in, ` {
			{{- if .Basic }}
			basic: {
				username: "{{ .Basic.Username }}",
				password: "{{ .Basic.Password }}",
			},
			{{- end }}
			{{- if .Oauth }}
			oauth: {
				clientId: "{{ .Oauth.ClientID }}",
				clientSecret: "{{ .Oauth.ClientSecret }}",
				url: "{{ .Oauth.URL }}",
			},
			{{- end }}
	}`)
}

func (g *Graphqlizer) CSRFTokenCredentialRequestAuthInputToGQL(in *graphql.CSRFTokenCredentialRequestAuthInput) (string, error) {
	in.AdditionalHeadersSerialized = quoteHTTPHeadersSerialized(in.AdditionalHeadersSerialized)
	in.AdditionalQueryParamsSerialized = quoteQueryParamsSerialized(in.AdditionalQueryParamsSerialized)

	return g.genericToGQL(in, `{
			tokenEndpointURL: "{{ .TokenEndpointURL }}",
			{{- if .Credential }}
			credential: {{ CredentialDataInputToGQL .Credential }},
			{{- end }}
			{{- if .AdditionalHeaders }}
			additionalHeaders: {{ HTTPHeadersToGQL .AdditionalHeaders }},
			{{- end }}
			{{- if .AdditionalHeadersSerialized }}
			additionalHeadersSerialized: {{ .AdditionalHeadersSerialized }},
			{{- end }}
			{{- if .AdditionalQueryParams }}
			additionalQueryParams: {{ QueryParamsToGQL .AdditionalQueryParams }},
			{{- end }}
			{{- if .AdditionalQueryParamsSerialized }}
			additionalQueryParamsSerialized: {{ .AdditionalQueryParamsSerialized }},
			{{- end }}
	}`)
}

func (g *Graphqlizer) AuthInputToGQL(in *graphql.AuthInput) (string, error) {
	in.AdditionalHeadersSerialized = quoteHTTPHeadersSerialized(in.AdditionalHeadersSerialized)
	in.AdditionalQueryParamsSerialized = quoteQueryParamsSerialized(in.AdditionalQueryParamsSerialized)

	return g.genericToGQL(in, `{
		{{- if .Credential }}
		credential: {{ CredentialDataInputToGQL .Credential }},
		{{- end }}
		{{- if .AdditionalHeaders }}
		additionalHeaders: {{ HTTPHeadersToGQL .AdditionalHeaders }},
		{{- end }}
		{{- if .AdditionalHeadersSerialized }}
		additionalHeadersSerialized: {{ .AdditionalHeadersSerialized }},
		{{- end }}
		{{- if .AdditionalQueryParams }}
		additionalQueryParams: {{ QueryParamsToGQL .AdditionalQueryParams}},
		{{- end }}
		{{- if .AdditionalQueryParamsSerialized }}
		additionalQueryParamsSerialized: {{ .AdditionalQueryParamsSerialized }},
		{{- end }}
		{{- if .RequestAuth }}
		requestAuth: {{ CredentialRequestAuthInputToGQL .RequestAuth }},
		{{- end }}
	}`)
}

func (g *Graphqlizer) LabelsToGQL(in graphql.Labels) (string, error) {
	return g.marshal(in), nil
}

func (g *Graphqlizer) HTTPHeadersToGQL(in graphql.HttpHeaders) (string, error) {
	return g.genericToGQL(in, `{
		{{- range $k,$v := . }}
			{{$k}}: [
				{{- range $i,$j := $v }}
					{{- if $i}},{{- end}}"{{$j}}"
				{{- end }} ],
		{{- end}}
	}`)
}

func (g *Graphqlizer) QueryParamsToGQL(in graphql.QueryParams) (string, error) {
	return g.genericToGQL(in, `{
		{{- range $k,$v := . }}
			{{$k}}: [
				{{- range $i,$j := $v }}
					{{- if $i}},{{- end}}"{{$j}}"
				{{- end }} ],
		{{- end}}
	}`)
}

func (g *Graphqlizer) WebhookInputToGQL(in *graphql.WebhookInput) (string, error) {
	return g.genericToGQL(in, `{
		type: {{.Type}},
		{{- if .URL }}
		url: "{{.URL }}",
		{{- end }}
		{{- if .Auth }} 
		auth: {{- AuthInputToGQL .Auth }},
		{{- end }}
		{{- if .Mode }} 
		mode: {{.Mode }},
		{{- end }}
		{{- if .CorrelationIDKey }} 
		correlationIdKey: "{{.CorrelationIDKey }}",
		{{- end }}
		{{- if .RetryInterval }} 
		retryInterval: "{{.RetryInterval }}",
		{{- end }}
		{{- if .Timeout }} 
		timeout: "{{.Timeout }}",
		{{- end }}
		{{- if .URLTemplate }} 
		urlTemplate: "{{.URLTemplate }}",
		{{- end }}
		{{- if .InputTemplate }} 
		inputTemplate: "{{.InputTemplate }}",
		{{- end }}
		{{- if .HeaderTemplate }} 
		headerTemplate: "{{.HeaderTemplate }}",
		{{- end }}
		{{- if .OutputTemplate }} 
		outputTemplate: "{{.OutputTemplate }}",
		{{- end }}
		{{- if .StatusTemplate }} 
		statusTemplate: "{{.StatusTemplate }}",
		{{- end }}
	}`)
}

func (g *Graphqlizer) APIDefinitionInputToGQL(in graphql.APIDefinitionInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{ .Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end}}
		targetURL: "{{.TargetURL}}",
		{{- if .Group }}
		group: "{{.Group}}",
		{{- end }}
		{{- if .Spec }}
		spec: {{- ApiSpecInputToGQL .Spec }},
		{{- end }}
		{{- if .Version }}
		version: {{- VersionInputToGQL .Version }},
		{{- end}}
	}`)
}

func (g *Graphqlizer) EventDefinitionInputToGQL(in graphql.EventDefinitionInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
		{{- if .Spec }}
		spec: {{ EventAPISpecInputToGQL .Spec }},
		{{- end }}
		{{- if .Group }}
		group: "{{.Group}}", 
		{{- end }}
		{{- if .Version }}
		version: {{- VersionInputToGQL .Version }},
		{{- end}}
	}`)
}

func (g *Graphqlizer) EventAPISpecInputToGQL(in graphql.EventSpecInput) (string, error) {
	in.Data = quoteCLOB(in.Data)
	return g.genericToGQL(in, `{
		{{- if .Data }}
		data: {{.Data}},
		{{- end }}
		type: {{.Type}},
		{{- if .FetchRequest }}
		fetchRequest: {{- FetchRequesstInputToGQL .FetchRequest }},
		{{- end }}
		format: {{.Format}},
	}`)
}

func (g *Graphqlizer) ApiSpecInputToGQL(in graphql.APISpecInput) (string, error) {
	in.Data = quoteCLOB(in.Data)
	return g.genericToGQL(in, `{
		{{- if .Data}}
		data: {{.Data}},
		{{- end}}	
		type: {{.Type}},
		format: {{.Format}},
		{{- if .FetchRequest }}
		fetchRequest: {{- FetchRequesstInputToGQL .FetchRequest }},
		{{- end }}
	}`)
}

func (g *Graphqlizer) VersionInputToGQL(in graphql.VersionInput) (string, error) {
	return g.genericToGQL(in, `{
		value: "{{.Value}}",
		{{- if .Deprecated }}
		deprecated: {{.Deprecated}},
		{{- end}}
		{{- if .DeprecatedSince }}
		deprecatedSince: "{{.DeprecatedSince}}",
		{{- end}}
		{{- if .ForRemoval }}
		forRemoval: {{.ForRemoval }},
		{{- end }}
	}`)
}

func (g *Graphqlizer) RuntimeInputToGQL(in graphql.RuntimeInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
		{{- if .Labels }}
		labels: {{ LabelsToGQL .Labels}},
		{{- end }}
		{{- if .StatusCondition }}
		statusCondition: {{ .StatusCondition }},
		{{- end }}
	}`)
}

func (g *Graphqlizer) LabelDefinitionInputToGQL(in graphql.LabelDefinitionInput) (string, error) {
	return g.genericToGQL(in, `{
		key: "{{.Key}}",
		{{- if .Schema }}
		schema: {{.Schema}},
		{{- end }}
	}`)
}

func (g *Graphqlizer) LabelFilterToGQL(in graphql.LabelFilter) (string, error) {
	return g.genericToGQL(in, `{
		key: "{{.Key}}",
		{{- if .Query }}
		query: "{{.Query}}",
		{{- end }}
	}`)
}

func (g *Graphqlizer) IntegrationSystemInputToGQL(in graphql.IntegrationSystemInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
	}`)
}

func (g *Graphqlizer) PlaceholderDefinitionInputToGQL(in graphql.PlaceholderDefinitionInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{.Name}}",
		{{- if .Description }}
		description: "{{.Description}}",
		{{- end }}
	}`)
}

func (g *Graphqlizer) TemplateValueInputToGQL(in graphql.TemplateValueInput) (string, error) {
	return g.genericToGQL(in, `{
		placeholder: "{{.Placeholder}}"
		value: "{{.Value}}"
	}`)
}

func (g *Graphqlizer) ApplicationFromTemplateInputToGQL(in graphql.ApplicationFromTemplateInput) (string, error) {
	return g.genericToGQL(in, `{
		templateName: "{{.TemplateName}}"
		{{- if .Values }}
		values: [
			{{- range $i, $e := .Values }} 
				{{- if $i}}, {{- end}} {{ TemplateValueInput $e }}
			{{- end }} ],
		{{- end }},
	}`)
}

func (g *Graphqlizer) BundleCreateInputToGQL(in graphql.BundleCreateInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{ .Name }}"
		{{- if .Description }}
		description: "{{ .Description }}"
		{{- end }}
		{{- if .InstanceAuthRequestInputSchema }}
		instanceAuthRequestInputSchema: {{ .InstanceAuthRequestInputSchema }}
		{{- end }}
		{{- if .DefaultInstanceAuth }}
		defaultInstanceAuth: {{- AuthInputToGQL .DefaultInstanceAuth }}
		{{- end }}
		{{- if .APIDefinitions }}
		apiDefinitions: [
			{{- range $i, $e := .APIDefinitions }}
				{{- if $i}}, {{- end}} {{ APIDefinitionInputToGQL $e }}
			{{- end }}],
		{{- end }}
		{{- if .EventDefinitions }}
		eventDefinitions: [
			{{- range $i, $e := .EventDefinitions }}
				{{- if $i}}, {{- end}} {{ EventDefinitionInputToGQL $e }}
			{{- end }}],
		{{- end }}
		{{- if .Documents }} 
		documents: [
			{{- range $i, $e := .Documents }} 
				{{- if $i}}, {{- end}} {{- DocumentInputToGQL $e }}
			{{- end }} ],
		{{- end }}
	}`)
}

func (g *Graphqlizer) BundleUpdateInputToGQL(in graphql.BundleUpdateInput) (string, error) {
	return g.genericToGQL(in, `{
		name: "{{ .Name }}"
		{{- if .Description }}
		description: "{{ .Description }}"
		{{- end }}
		{{- if .InstanceAuthRequestInputSchema }}
		instanceAuthRequestInputSchema: {{ .InstanceAuthRequestInputSchema }}
		{{- end }}
		{{- if .DefaultInstanceAuth }}
		defaultInstanceAuth: {{- AuthInputToGQL .DefaultInstanceAuth }}
		{{- end }}
	}`)
}

func (g *Graphqlizer) BundleInstanceAuthStatusInputToGQL(in graphql.BundleInstanceAuthStatusInput) (string, error) {
	return g.genericToGQL(in, `{
		condition: {{ .Condition }}
		{{- if .Message }}
		message: "{{ .Message }}"
		{{- end }}
		{{- if .Reason }}
		reason: "{{ .Reason }}"
		{{- end }}
	}`)
}

func (g *Graphqlizer) BundleInstanceAuthRequestInputToGQL(in graphql.BundleInstanceAuthRequestInput) (string, error) {
	return g.genericToGQL(in, `{
		{{- if .ID }}
		id: "{{ .ID }}"
		{{- end }}		
		{{- if .Context }}
		context: {{ .Context }}
		{{- end }}
		{{- if .InputParams }}
		inputParams: {{ .InputParams }}
		{{- end }}
	}`)
}

func (g *Graphqlizer) BundleInstanceAuthSetInputToGQL(in graphql.BundleInstanceAuthSetInput) (string, error) {
	return g.genericToGQL(in, `{
		{{- if .Auth }}
		auth: {{- AuthInputToGQL .Auth}}
		{{- end }}
		{{- if .Status }}
		status: {{- BundleInstanceAuthStatusInputToGQL .Status }}
		{{- end }}
	}`)
}

func (g *Graphqlizer) LabelSelectorInputToGQL(in graphql.LabelSelectorInput) (string, error) {
	return g.genericToGQL(in, `{
		key: "{{ .Key }}"
		value: "{{ .Value }}"
	}`)
}

func (g *Graphqlizer) AutomaticScenarioAssignmentSetInputToGQL(in graphql.AutomaticScenarioAssignmentSetInput) (string, error) {
	return g.genericToGQL(in, `{
		scenarioName: "{{ .ScenarioName }}"
		selector: {{- LabelSelectorInputToGQL .Selector }}
	}`)
}

func (g *Graphqlizer) marshal(obj interface{}) string {
	var out string

	val := reflect.ValueOf(obj)

	switch val.Kind() {
	case reflect.Map:
		s, err := g.genericToGQL(obj, `{ {{- range $k, $v := . }}{{ $k }}:{{ marshal $v }},{{ end -}} }`)
		if err != nil {
			return ""
		}
		out = s
	case reflect.Slice, reflect.Array:
		s, err := g.genericToGQL(obj, `[{{ range $i, $e := . }}{{ if $i }},{{ end }}{{ marshal $e }}{{ end }}]`)
		if err != nil {
			return ""
		}
		out = s
	default:
		marshalled, err := json.Marshal(obj)
		if err != nil {
			return ""
		}
		out = string(marshalled)
	}

	return out
}

func (g *Graphqlizer) genericToGQL(obj interface{}, tmpl string) (string, error) {
	fm := sprig.TxtFuncMap()
	fm["marshal"] = g.marshal
	fm["ApplicationRegisterInputToGQL"] = g.ApplicationRegisterInputToGQL
	fm["DocumentInputToGQL"] = g.DocumentInputToGQL
	fm["FetchRequesstInputToGQL"] = g.FetchRequestInputToGQL
	fm["AuthInputToGQL"] = g.AuthInputToGQL
	fm["LabelsToGQL"] = g.LabelsToGQL
	fm["WebhookInputToGQL"] = g.WebhookInputToGQL
	fm["APIDefinitionInputToGQL"] = g.APIDefinitionInputToGQL
	fm["EventDefinitionInputToGQL"] = g.EventDefinitionInputToGQL
	fm["ApiSpecInputToGQL"] = g.ApiSpecInputToGQL
	fm["VersionInputToGQL"] = g.VersionInputToGQL
	fm["HTTPHeadersToGQL"] = g.HTTPHeadersToGQL
	fm["QueryParamsToGQL"] = g.QueryParamsToGQL
	fm["EventAPISpecInputToGQL"] = g.EventAPISpecInputToGQL
	fm["CredentialDataInputToGQL"] = g.CredentialDataInputToGQL
	fm["CSRFTokenCredentialRequestAuthInputToGQL"] = g.CSRFTokenCredentialRequestAuthInputToGQL
	fm["CredentialRequestAuthInputToGQL"] = g.CredentialRequestAuthInputToGQL
	fm["PlaceholderDefinitionInputToGQL"] = g.PlaceholderDefinitionInputToGQL
	fm["TemplateValueInput"] = g.TemplateValueInputToGQL
	fm["BundleInstanceAuthStatusInputToGQL"] = g.BundleInstanceAuthStatusInputToGQL
	fm["BundleCreateInputToGQL"] = g.BundleCreateInputToGQL
	fm["LabelSelectorInputToGQL"] = g.LabelSelectorInputToGQL

	t, err := template.New("tmpl").Funcs(fm).Parse(tmpl)
	if err != nil {
		return "", errors.Wrapf(err, "while parsing template")
	}

	var b bytes.Buffer

	if err := t.Execute(&b, obj); err != nil {
		return "", errors.Wrap(err, "while executing template")
	}
	return b.String(), nil
}

func quoteCLOB(in *graphql.CLOB) *graphql.CLOB {
	if in == nil {
		return nil
	}

	quoted := strconv.Quote(string(*in))
	return (*graphql.CLOB)(&quoted)
}

func quoteHTTPHeadersSerialized(in *graphql.HttpHeadersSerialized) *graphql.HttpHeadersSerialized {
	if in == nil {
		return nil
	}

	quoted := strconv.Quote(string(*in))
	return (*graphql.HttpHeadersSerialized)(&quoted)
}

func quoteQueryParamsSerialized(in *graphql.QueryParamsSerialized) *graphql.QueryParamsSerialized {
	if in == nil {
		return nil
	}

	quoted := strconv.Quote(string(*in))
	return (*graphql.QueryParamsSerialized)(&quoted)
}
