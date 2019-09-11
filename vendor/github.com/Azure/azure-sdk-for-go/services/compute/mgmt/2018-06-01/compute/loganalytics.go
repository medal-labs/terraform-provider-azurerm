package compute

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// LogAnalyticsClient is the compute Client
type LogAnalyticsClient struct {
	BaseClient
}

// NewLogAnalyticsClient creates an instance of the LogAnalyticsClient client.
func NewLogAnalyticsClient(subscriptionID string) LogAnalyticsClient {
	return NewLogAnalyticsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewLogAnalyticsClientWithBaseURI creates an instance of the LogAnalyticsClient client.
func NewLogAnalyticsClientWithBaseURI(baseURI string, subscriptionID string) LogAnalyticsClient {
	return LogAnalyticsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ExportRequestRateByInterval export logs that show Api requests made by this subscription in the given time window to
// show throttling activities.
// Parameters:
// parameters - parameters supplied to the LogAnalytics getRequestRateByInterval Api.
// location - the location upon which virtual-machine-sizes is queried.
func (client LogAnalyticsClient) ExportRequestRateByInterval(ctx context.Context, parameters RequestRateByIntervalInput, location string) (result LogAnalyticsExportRequestRateByIntervalFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/LogAnalyticsClient.ExportRequestRateByInterval")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: location,
			Constraints: []validation.Constraint{{Target: "location", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("compute.LogAnalyticsClient", "ExportRequestRateByInterval", err.Error())
	}

	req, err := client.ExportRequestRateByIntervalPreparer(ctx, parameters, location)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.LogAnalyticsClient", "ExportRequestRateByInterval", nil, "Failure preparing request")
		return
	}

	result, err = client.ExportRequestRateByIntervalSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.LogAnalyticsClient", "ExportRequestRateByInterval", result.Response(), "Failure sending request")
		return
	}

	return
}

// ExportRequestRateByIntervalPreparer prepares the ExportRequestRateByInterval request.
func (client LogAnalyticsClient) ExportRequestRateByIntervalPreparer(ctx context.Context, parameters RequestRateByIntervalInput, location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/logAnalytics/apiAccess/getRequestRateByInterval", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ExportRequestRateByIntervalSender sends the ExportRequestRateByInterval request. The method will close the
// http.Response Body if it receives an error.
func (client LogAnalyticsClient) ExportRequestRateByIntervalSender(req *http.Request) (future LogAnalyticsExportRequestRateByIntervalFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// ExportRequestRateByIntervalResponder handles the response to the ExportRequestRateByInterval request. The method always
// closes the http.Response Body.
func (client LogAnalyticsClient) ExportRequestRateByIntervalResponder(resp *http.Response) (result LogAnalyticsOperationResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ExportThrottledRequests export logs that show total throttled Api requests for this subscription in the given time
// window.
// Parameters:
// parameters - parameters supplied to the LogAnalytics getThrottledRequests Api.
// location - the location upon which virtual-machine-sizes is queried.
func (client LogAnalyticsClient) ExportThrottledRequests(ctx context.Context, parameters ThrottledRequestsInput, location string) (result LogAnalyticsExportThrottledRequestsFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/LogAnalyticsClient.ExportThrottledRequests")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: location,
			Constraints: []validation.Constraint{{Target: "location", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("compute.LogAnalyticsClient", "ExportThrottledRequests", err.Error())
	}

	req, err := client.ExportThrottledRequestsPreparer(ctx, parameters, location)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.LogAnalyticsClient", "ExportThrottledRequests", nil, "Failure preparing request")
		return
	}

	result, err = client.ExportThrottledRequestsSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.LogAnalyticsClient", "ExportThrottledRequests", result.Response(), "Failure sending request")
		return
	}

	return
}

// ExportThrottledRequestsPreparer prepares the ExportThrottledRequests request.
func (client LogAnalyticsClient) ExportThrottledRequestsPreparer(ctx context.Context, parameters ThrottledRequestsInput, location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/logAnalytics/apiAccess/getThrottledRequests", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ExportThrottledRequestsSender sends the ExportThrottledRequests request. The method will close the
// http.Response Body if it receives an error.
func (client LogAnalyticsClient) ExportThrottledRequestsSender(req *http.Request) (future LogAnalyticsExportThrottledRequestsFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// ExportThrottledRequestsResponder handles the response to the ExportThrottledRequests request. The method always
// closes the http.Response Body.
func (client LogAnalyticsClient) ExportThrottledRequestsResponder(resp *http.Response) (result LogAnalyticsOperationResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
