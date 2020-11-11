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

package osb

import (
	"context"

	schema "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/system-broker/internal/director"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/log"
	"github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/pivotal-cf/brokerapi/v7/domain/apiresponses"
	"github.com/pkg/errors"
)

type BindLastOperationEndpoint struct {
	credentialsGetter packageCredentialsFetcher
}

func (b *BindLastOperationEndpoint) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	log.C(ctx).Infof("LastBindingOperation instanceID: %s bindingID: %s details: %+v", instanceID, bindingID, details)

	opType := details.OperationData
	appID := details.ServiceID  // may be empty per OSB spec
	packageID := details.PlanID // may be empty per OSB spec
	authID := bindingID

	logger := log.C(ctx).WithFields(map[string]interface{}{
		"opType":     opType,
		"appID":      appID,
		"packageID":  packageID,
		"authID":     authID,
		"instanceID": instanceID,
		"bindingID":  bindingID,
	})

	logger.Info("Fetching package instance credentials")
	resp, err := b.credentialsGetter.FetchPackageInstanceAuth(ctx, &director.PackageInstanceInput{
		InstanceAuthID: authID,
		Context: map[string]string{
			"instance_id": instanceID,
			"binding_id":  bindingID,
		},
	})
	if err != nil && !IsNotFoundError(err) {
		return domain.LastOperation{}, errors.Wrapf(err, "while getting package instance credentials from director")
	}

	if IsNotFoundError(err) {
		if opType == string(UnbindOp) {
			return domain.LastOperation{
				State:       domain.Succeeded,
				Description: "credentials were successfully deleted",
			}, nil
		}
		logger.Error("Package instance credentials not found")
		return domain.LastOperation{}, apiresponses.ErrBindingNotFound
	}

	instanceAuth := resp.InstanceAuth

	logger.Infof("Found package credentials during poll last op with status %+v", *instanceAuth.Status)

	var state domain.LastOperationState
	var opErr error
	switch opType {
	case string(BindOp):
		switch instanceAuth.Status.Condition {
		case schema.PackageInstanceAuthStatusConditionSucceeded: // success
			state = domain.Succeeded
		case schema.PackageInstanceAuthStatusConditionPending: // in progress
			state = domain.InProgress
		case schema.PackageInstanceAuthStatusConditionFailed: // failed
			// this would trigger orphan mitigation
			state = domain.Failed
		case schema.PackageInstanceAuthStatusConditionUnused: // error
			fallthrough
		default:
			// this should force platform to continue polling, should be the more flexiable approach
			opErr = errors.Errorf("operation reached unexpected state: op %s, status %+v", opType, *instanceAuth.Status)
		}
	case string(UnbindOp):
		state = domain.InProgress
	}

	if opErr != nil {
		return domain.LastOperation{}, opErr
	}

	return domain.LastOperation{
		State:       state,
		Description: instanceAuth.Status.Message,
	}, nil
}