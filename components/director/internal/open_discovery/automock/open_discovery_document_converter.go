// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	open_discovery "github.com/kyma-incubator/compass/components/director/internal/open_discovery"
)

// OpenDiscoveryDocumentConverter is an autogenerated mock type for the OpenDiscoveryDocumentConverter type
type OpenDiscoveryDocumentConverter struct {
	mock.Mock
}

// DocumentToGraphQLInputs provides a mock function with given fields: _a0
func (_m *OpenDiscoveryDocumentConverter) DocumentToGraphQLInputs(_a0 *open_discovery.Document) ([]*graphql.PackageInput, []*graphql.BundleInput, error) {
	ret := _m.Called(_a0)

	var r0 []*graphql.PackageInput
	if rf, ok := ret.Get(0).(func(*open_discovery.Document) []*graphql.PackageInput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.PackageInput)
		}
	}

	var r1 []*graphql.BundleInput
	if rf, ok := ret.Get(1).(func(*open_discovery.Document) []*graphql.BundleInput); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*graphql.BundleInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*open_discovery.Document) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}