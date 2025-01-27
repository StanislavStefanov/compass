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

package tests

import (
	"encoding/json"
	"testing"

	"github.com/kyma-incubator/compass/tests/pkg/certs"
	"github.com/kyma-incubator/compass/tests/pkg/clients"
	"github.com/kyma-incubator/compass/tests/pkg/config"

	"github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/model"
	directorSchema "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/ptr"
	"github.com/stretchr/testify/require"
)

func TestAppRegistry(t *testing.T) {
	appInput := directorSchema.ApplicationRegisterInput{
		Name:           TestApp,
		ProviderName:   ptr.String("provider name"),
		Description:    ptr.String("my application"),
		HealthCheckURL: ptr.String("http://mywordpress.com/health"),
		Labels: &directorSchema.Labels{
			"scenarios": []interface{}{"DEFAULT"},
		},
	}

	descr := "test"
	runtimeInput := directorSchema.RuntimeInput{
		Name:        TestRuntime,
		Description: &descr,
		Labels: &directorSchema.Labels{
			"scenarios": []interface{}{"DEFAULT"},
		},
	}

	cfg := config.ConnectivityAdapterTestConfig{}
	config.ReadConfig(&cfg)

	directorClient, err := clients.NewDirectorClient(
		cfg.DirectorUrl,
		cfg.DirectorHealthzUrl,
		cfg.Tenant,
		[]string{"application:read", "application:write", "runtime:write", "runtime:read", "eventing:manage"})
	require.NoError(t, err)

	appID, err := directorClient.CreateApplication(appInput)
	require.NoError(t, err)

	defer func() {
		err = directorClient.DeleteApplication(appID)
		require.NoError(t, err)
	}()

	runtimeID, err := directorClient.CreateRuntime(runtimeInput)
	require.NoError(t, err)

	defer func() {
		err = directorClient.DeleteRuntime(runtimeID)
		require.NoError(t, err)
	}()

	err = directorClient.SetDefaultEventing(runtimeID, appID, cfg.EventsBaseURL)
	require.NoError(t, err)

	t.Run("App Registry Service flow for Application", func(t *testing.T) {
		client := clients.NewConnectorClient(directorClient, appID, cfg.Tenant, cfg.SkipSslVerify)
		clientKey := certs.CreateKey(t)

		crtResponse, infoResponse := createCertificateChain(t, client, clientKey)
		require.NotEmpty(t, crtResponse.CRTChain)
		require.NotEmpty(t, infoResponse.Api.ManagementInfoURL)
		require.NotEmpty(t, infoResponse.Certificate)

		certificates := certs.DecodeAndParseCerts(t, crtResponse)
		adapterClient := clients.NewSecuredClient(cfg.SkipSslVerify, clientKey, certificates.ClientCRT.Raw, cfg.Tenant)

		mgmInfoResponse, errorResponse := adapterClient.GetMgmInfo(t, infoResponse.Api.ManagementInfoURL)
		require.Nil(t, errorResponse)
		require.NotEmpty(t, mgmInfoResponse.URLs.RenewCertURL)
		require.NotEmpty(t, mgmInfoResponse.Certificate)
		require.Equal(t, infoResponse.Certificate, mgmInfoResponse.Certificate)

		defer func() {
			errorResponse = adapterClient.RevokeCertificate(t, mgmInfoResponse.URLs.RevokeCertURL)
			require.Nil(t, errorResponse)
		}()

		metadataURL := infoResponse.Api.MetadataURL

		services, errorResponse := adapterClient.ListServices(t, metadataURL)
		require.Nil(t, errorResponse)
		require.Len(t, services, 0)

		service := model.ServiceDetails{
			Name:        "test-service",
			Provider:    "provider",
			Description: "description",
			Api: &model.API{
				TargetUrl: "http://target.com",
				Credentials: &model.CredentialsWithCSRF{
					OauthWithCSRF: &model.OauthWithCSRF{
						Oauth: model.Oauth{
							URL:          "http://test.com/token",
							ClientID:     "client",
							ClientSecret: "secret",
						},
					},
				},
			},
			Labels: &map[string]string{},
			Events: &model.Events{
				Spec: json.RawMessage(`{"asyncapi":"1.2.0"}`),
			},
		}

		createServiceResponse, errorResponse := adapterClient.CreateService(t, metadataURL, service)
		require.Nil(t, errorResponse)
		require.NotNil(t, createServiceResponse.ID)

		expectedService := service
		expectedService.Provider = ""

		serviceResponse, errorResponse := adapterClient.GetService(t, metadataURL, createServiceResponse.ID)
		require.Nil(t, errorResponse)
		require.Equal(t, &expectedService, serviceResponse)

		expectedService.Api.TargetUrl = service.Api.TargetUrl + "/test"

		updateServiceResponse, errorResponse := adapterClient.UpdateService(t, metadataURL, createServiceResponse.ID, service)
		require.Nil(t, errorResponse)
		require.Equal(t, &expectedService, updateServiceResponse)

		services, errorResponse = adapterClient.ListServices(t, metadataURL)
		require.Nil(t, errorResponse)
		require.Len(t, services, 1)
		require.Equal(t, expectedService.Name, services[0].Name)
		require.Equal(t, expectedService.Description, services[0].Description)

		errorResponse = adapterClient.DeleteService(t, metadataURL, createServiceResponse.ID)
		require.Nil(t, errorResponse)

		services, errorResponse = adapterClient.ListServices(t, metadataURL)
		require.Nil(t, errorResponse)
		require.Len(t, services, 0)
	})
}
