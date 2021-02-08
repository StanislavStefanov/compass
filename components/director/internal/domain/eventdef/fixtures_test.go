package eventdef_test

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	event "github.com/kyma-incubator/compass/components/director/internal/domain/eventdef"

	"github.com/kyma-incubator/compass/components/director/internal/domain/version"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
)

const (
	eventID          = "ddddddddd-dddd-dddd-dddd-dddddddddddd"
	specID           = "sssssssss-ssss-ssss-ssss-ssssssssssss"
	tenantID         = "ttttttttt-tttt-tttt-tttt-tttttttttttt"
	externalTenantID = "eeeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	bundleID         = "bbbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	packageID        = "ppppppppp-pppp-pppp-pppp-pppppppppppp"
	ordID            = "com.compass.ord.v1"
)

func fixEventDefinitionModel(id string, bndlID string, name string) *model.EventDefinition {
	return &model.EventDefinition{
		ID:       id,
		BundleID: &bndlID,
		Name:     name,
	}
}

func fixFullEventDefinitionModel(placeholder string) (model.EventDefinition, model.Spec) {
	eventType := model.EventSpecTypeAsyncAPI
	spec := model.Spec{
		ID:         specID,
		Data:       str.Ptr("spec_data_" + placeholder),
		Format:     model.SpecFormatYaml,
		ObjectType: model.EventSpecReference,
		ObjectID:   eventID,
		EventType:  &eventType,
	}

	deprecated := false
	forRemoval := false

	v := &model.Version{
		Value:           "v1.1",
		Deprecated:      &deprecated,
		DeprecatedSince: str.Ptr("v1.0"),
		ForRemoval:      &forRemoval,
	}

	boolVar := false
	return model.EventDefinition{
		ID:                  eventID,
		BundleID:            str.Ptr(bundleID),
		PackageID:           str.Ptr(packageID),
		Tenant:              tenantID,
		Name:                placeholder,
		Description:         str.Ptr("desc_" + placeholder),
		Group:               str.Ptr("group_" + placeholder),
		OrdID:               str.Ptr(ordID),
		ShortDescription:    str.Ptr("shortDescription"),
		SystemInstanceAware: &boolVar,
		Tags:                json.RawMessage("[]"),
		Countries:           json.RawMessage("[]"),
		Links:               json.RawMessage("[]"),
		ReleaseStatus:       str.Ptr("releaseStatus"),
		SunsetDate:          str.Ptr("sunsetDate"),
		Successor:           str.Ptr("successor"),
		ChangeLogEntries:    json.RawMessage("[]"),
		Labels:              json.RawMessage("[]"),
		Visibility:          str.Ptr("visibility"),
		Disabled:            &boolVar,
		PartOfProducts:      json.RawMessage("[]"),
		LineOfBusiness:      json.RawMessage("[]"),
		Industry:            json.RawMessage("[]"),
		Version:             v,
	}, spec
}

func fixFullGQLEventDefinition(placeholder string) *graphql.EventDefinition {
	data := graphql.CLOB("spec_data_" + placeholder)

	spec := &graphql.EventSpec{
		Data:         &data,
		Format:       graphql.SpecFormatYaml,
		Type:         graphql.EventSpecTypeAsyncAPI,
		DefinitionID: eventID,
	}

	deprecated := false
	forRemoval := false

	v := &graphql.Version{
		Value:           "v1.1",
		Deprecated:      &deprecated,
		DeprecatedSince: str.Ptr("v1.0"),
		ForRemoval:      &forRemoval,
	}

	return &graphql.EventDefinition{
		ID:          eventID,
		BundleID:    bundleID,
		Name:        placeholder,
		Description: str.Ptr("desc_" + placeholder),
		Spec:        spec,
		Group:       str.Ptr("group_" + placeholder),
		Version:     v,
	}
}

func fixModelEventDefinitionInput(name, description string, group string) (*model.EventDefinitionInput, *model.SpecInput) {
	data := "data"
	eventType := model.EventSpecTypeAsyncAPI

	spec := &model.SpecInput{
		Data:         &data,
		EventType:    &eventType,
		Format:       model.SpecFormatYaml,
		FetchRequest: &model.FetchRequestInput{},
	}

	deprecated := false
	deprecatedSince := ""
	forRemoval := false

	v := &model.VersionInput{
		Value:           "1.0.0",
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}

	return &model.EventDefinitionInput{
		Name:        name,
		Description: &description,
		Group:       &group,
		VersionInput:     v,
	}, spec
}

func fixGQLEventDefinitionInput(name, description string, group string) *graphql.EventDefinitionInput {
	data := graphql.CLOB("data")

	spec := &graphql.EventSpecInput{
		Data:         &data,
		Type:         graphql.EventSpecTypeAsyncAPI,
		Format:       graphql.SpecFormatYaml,
		FetchRequest: &graphql.FetchRequestInput{},
	}

	deprecated := false
	deprecatedSince := ""
	forRemoval := false

	v := &graphql.VersionInput{
		Value:           "1.0.0",
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}

	return &graphql.EventDefinitionInput{
		Name:        name,
		Description: &description,
		Group:       &group,
		Spec:        spec,
		Version:     v,
	}
}

func fixEntityEventDefinition(id string, bndlID string, name string) event.Entity {
	return event.Entity{
		ID:       id,
		BundleID: repo.NewValidNullableString(bndlID),
		Name:     name,
	}
}

func fixFullEntityEventDefinition(eventID, placeholder string) event.Entity {
	return event.Entity{
		ID:                  eventID,
		TenantID:            tenantID,
		BundleID:            repo.NewValidNullableString(bundleID),
		PackageID:           repo.NewValidNullableString(packageID),
		Name:                placeholder,
		Description:         repo.NewValidNullableString("desc_" + placeholder),
		GroupName:           repo.NewValidNullableString("group_" + placeholder),
		OrdID:               repo.NewValidNullableString(ordID),
		ShortDescription:    repo.NewValidNullableString("shortDescription"),
		SystemInstanceAware: repo.NewValidNullableBool(false),
		ChangeLogEntries:    repo.NewValidNullableString("[]"),
		Links:               repo.NewValidNullableString("[]"),
		Tags:                repo.NewValidNullableString("[]"),
		Countries:           repo.NewValidNullableString("[]"),
		ReleaseStatus:       repo.NewValidNullableString("releaseStatus"),
		SunsetDate:          repo.NewValidNullableString("sunsetDate"),
		Successor:           repo.NewValidNullableString("successor"),
		Labels:              repo.NewValidNullableString("[]"),
		Visibility:          repo.NewValidNullableString("visibility"),
		Disabled:            repo.NewValidNullableBool(false),
		PartOfProducts:      repo.NewValidNullableString("[]"),
		LineOfBusiness:      repo.NewValidNullableString("[]"),
		Industry:            repo.NewValidNullableString("[]"),
		Version: version.Version{
			Value:           repo.NewNullableString(str.Ptr("v1.1")),
			Deprecated:      repo.NewValidNullableBool(false),
			DeprecatedSince: repo.NewNullableString(str.Ptr("v1.0")),
			ForRemoval:      repo.NewValidNullableBool(false),
		},
	}
}

func fixEventDefinitionColumns() []string {
	return []string{"id", "tenant_id", "bundle_id", "package_id", "name", "description", "group_name", "ord_id",
		"short_description", "system_instance_aware", "changelog_entries", "links", "tags", "countries", "release_status",
		"sunset_date", "successor", "labels", "visibility", "disabled", "part_of_products", "line_of_business", "industry", "version_value", "version_deprecated", "version_deprecated_since",
		"version_for_removal"}
}

func fixEventDefinitionRow(id, placeholder string) []driver.Value {
	boolVar := false
	return []driver.Value{id, tenantID, bundleID, packageID, placeholder, "desc_" + placeholder, "group_" + placeholder, ordID, "shortDescription", &boolVar,
		repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), "releaseStatus", "sunsetDate", "successor", repo.NewValidNullableString("[]"), "visibility", &boolVar,
		repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), repo.NewValidNullableString("[]"), "v1.1", false, "v1.0", false}
}

func fixEventCreateArgs(id string, event *model.EventDefinition) []driver.Value {
	return []driver.Value{id, tenantID, bundleID, packageID, event.Name, event.Description, event.Group, event.OrdID, event.ShortDescription,
		event.SystemInstanceAware, repo.NewNullableStringFromJSONRawMessage(event.ChangeLogEntries), repo.NewNullableStringFromJSONRawMessage(event.Links),
		repo.NewNullableStringFromJSONRawMessage(event.Tags), repo.NewNullableStringFromJSONRawMessage(event.Countries), event.ReleaseStatus, event.SunsetDate, event.Successor,
		repo.NewNullableStringFromJSONRawMessage(event.Labels), event.Visibility,
		event.Disabled, repo.NewNullableStringFromJSONRawMessage(event.PartOfProducts), repo.NewNullableStringFromJSONRawMessage(event.LineOfBusiness), repo.NewNullableStringFromJSONRawMessage(event.Industry),
		event.Version.Value, event.Version.Deprecated, event.Version.DeprecatedSince, event.Version.ForRemoval}
}

func fixModelFetchRequest(id, url string, timestamp time.Time) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     id,
		Tenant: tenantID,
		URL:    url,
		Auth:   nil,
		Mode:   "SINGLE",
		Filter: nil,
		Status: &model.FetchRequestStatus{
			Condition: model.FetchRequestStatusConditionInitial,
			Timestamp: timestamp,
		},
		ObjectType: model.SpecFetchRequestReference,
		ObjectID:   specID,
	}
}

func fixGQLFetchRequest(url string, timestamp time.Time) *graphql.FetchRequest {
	return &graphql.FetchRequest{
		Filter: nil,
		Mode:   graphql.FetchModeSingle,
		Auth:   nil,
		URL:    url,
		Status: &graphql.FetchRequestStatus{
			Timestamp: graphql.Timestamp(timestamp),
			Condition: graphql.FetchRequestStatusConditionInitial,
		},
	}
}
