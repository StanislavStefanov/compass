/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package director

import (
	"context"
	"encoding/json"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	schema "github.com/kyma-incubator/compass/components/director/pkg/graphql"

	"strconv"
)

type ApplicationExt struct {
	schema.Application    `mapstructure:",squash"`
	Labels                schema.Labels                           `json:"labels"`
	Webhooks              []schema.Webhook                        `json:"webhooks"`
	Auths                 []*schema.SystemAuth                    `json:"auths"`
	EventingConfiguration schema.ApplicationEventingConfiguration `json:"eventingConfiguration"`

	Packages []schema.PackageExt `json:"packages"`
}

type ApplicationsOutput []ApplicationExt

// go:generate Page
type ApplicationResponse struct {
	Result struct {
		Apps ApplicationsOutput `json:"data"`
		Page graphql.PageInfo   `json:"pageInfo"`
	} `json:"result"`
}

func (ao *ApplicationResponse) PageInfo() *graphql.PageInfo {
	return &ao.Result.Page
}

func (ao *ApplicationResponse) ListAll(ctx context.Context, pager *Pager) (ApplicationsOutput, error) {
	appsResult := ApplicationsOutput{}

	for pager.HasNext() {
		apps := &ApplicationResponse{}
		if err := pager.Next(ctx, apps); err != nil {
			return nil, err
		}
		appsResult = append(appsResult, apps.Result.Apps...)
	}
	return appsResult, nil
}

type PackagessOutput []schema.PackageExt

type ApiDefinitionsOutput []schema.APIDefinitionExt

type EventDefinitionsOutput []schema.EventAPIDefinitionExt

type RequestPackageInstanceCredentialsInput struct {
	PackageID   string `valid:"required"`
	Context     Values
	InputSchema Values
}

type Values map[string]interface{}

func (r *Values) MarshalToQGLJSON() (string, error) {
	input, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	return strconv.Quote(string(input)), nil
}

type RequestPackageInstanceCredentialsOutput struct {
	InstanceAuth *schema.PackageInstanceAuth `json:"result"`
}

type FindPackageInstanceCredentialsByContextInput struct {
	ApplicationID string `valid:"required"`
	PackageID     string `valid:"required"`
	Context       map[string]string
}

type FindPackageInstanceCredentialsOutput struct {
	InstanceAuths []*schema.PackageInstanceAuth
	TargetURLs    map[string]string
}

type FindPackageInstanceCredentialInput struct {
	PackageID      string `valid:"required"`
	ApplicationID  string `valid:"required"`
	InstanceAuthID string `valid:"required"`
}

type FindPackageInstanceCredentialOutput struct {
	InstanceAuth *schema.PackageInstanceAuth `json:"result"`
}

type RequestPackageInstanceAuthDeletionInput struct {
	InstanceAuthID string `valid:"required"`
}

type RequestPackageInstanceAuthDeletionOutput struct {
	ID     string                           `json:"id"`
	Status schema.PackageInstanceAuthStatus `json:"status"`
}

type FindPackageSpecificationInput struct {
	ApplicationID string `valid:"required"`
	PackageID     string `valid:"required"`
	DefinitionID  string `valid:"required"`
}

type FindPackageSpecificationOutput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	Data    *schema.CLOB      `json:"data,omitempty"`
	Format  schema.SpecFormat `json:"format"`
	Type    string            `json:"type"`
	Version *schema.Version   `json:"version,omitempty"`
}
