// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0
package provider

// RESOURCE ACTION

import (
	"context"

	merakigosdk "github.com/meraki/dashboard-api-go/v3/sdk"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &NetworksSmDevicesRebootResource{}
	_ resource.ResourceWithConfigure = &NetworksSmDevicesRebootResource{}
)

func NewNetworksSmDevicesRebootResource() resource.Resource {
	return &NetworksSmDevicesRebootResource{}
}

type NetworksSmDevicesRebootResource struct {
	client *merakigosdk.Client
}

func (r *NetworksSmDevicesRebootResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client := req.ProviderData.(MerakiProviderData).Client
	r.client = client
}

// Metadata returns the data source type name.
func (r *NetworksSmDevicesRebootResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_sm_devices_reboot"
}

// resourceAction
func (r *NetworksSmDevicesRebootResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"network_id": schema.StringAttribute{
				MarkdownDescription: `networkId path parameter. Network ID`,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"item": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{

					"ids": schema.SetAttribute{
						MarkdownDescription: `The Meraki Ids of the set of endpoints.`,
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"parameters": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"ids": schema.SetAttribute{
						MarkdownDescription: `The ids of the endpoints to be rebooted.`,
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"kext_paths": schema.SetAttribute{
						MarkdownDescription: `The KextPaths of the endpoints to be rebooted. Available for macOS 11+`,
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"notify_user": schema.BoolAttribute{
						MarkdownDescription: `Whether or not to notify the user before rebooting the endpoint. Available for macOS 11.3+`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"rebuild_kernel_cache": schema.BoolAttribute{
						MarkdownDescription: `Whether or not to rebuild the kernel cache when rebooting the endpoint. Available for macOS 11+`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"request_requires_network_tether": schema.BoolAttribute{
						MarkdownDescription: `Whether or not the request requires network tethering. Available for macOS and supervised iOS or tvOS`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"scope": schema.SetAttribute{
						MarkdownDescription: `The scope (one of all, none, withAny, withAll, withoutAny, or withoutAll) and a set of tags of the endpoints to be rebooted.`,
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"serials": schema.SetAttribute{
						MarkdownDescription: `The serials of the endpoints to be rebooted.`,
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"wifi_macs": schema.SetAttribute{
						MarkdownDescription: `The wifiMacs of the endpoints to be rebooted.`,
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}
func (r *NetworksSmDevicesRebootResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var data NetworksSmDevicesReboot

	var item types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(item.As(ctx, &data, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)

	if resp.Diagnostics.HasError() {
		return
	}
	//Has Paths
	vvNetworkID := data.NetworkID.ValueString()
	dataRequest := data.toSdkApiRequestCreate(ctx)
	response, restyResp1, err := r.client.Sm.RebootNetworkSmDevices(vvNetworkID, dataRequest)

	if err != nil || response == nil {
		if restyResp1 != nil {
			resp.Diagnostics.AddError(
				"Failure when executing RebootNetworkSmDevices",
				err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failure when executing RebootNetworkSmDevices",
			err.Error(),
		)
		return
	}
	//Item
	data = ResponseSmRebootNetworkSmDevicesItemToBody(data, response)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *NetworksSmDevicesRebootResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.Diagnostics.AddWarning("Error deleting Resource", "This resource has no delete method in the meraki lab, the resource was deleted only in terraform.")
}

func (r *NetworksSmDevicesRebootResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("Error Update Resource", "This resource has no update method in the meraki lab, the resource was deleted only in terraform.")
}

func (r *NetworksSmDevicesRebootResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Error deleting Resource", "This resource has no delete method in the meraki lab, the resource was deleted only in terraform.")
	resp.State.RemoveResource(ctx)
}

// TF Structs Schema
type NetworksSmDevicesReboot struct {
	NetworkID  types.String                       `tfsdk:"network_id"`
	Item       *ResponseSmRebootNetworkSmDevices  `tfsdk:"item"`
	Parameters *RequestSmRebootNetworkSmDevicesRs `tfsdk:"parameters"`
}

type ResponseSmRebootNetworkSmDevices struct {
	IDs types.Set `tfsdk:"ids"`
}

type RequestSmRebootNetworkSmDevicesRs struct {
	IDs                          types.Set  `tfsdk:"ids"`
	KextPaths                    types.Set  `tfsdk:"kext_paths"`
	NotifyUser                   types.Bool `tfsdk:"notify_user"`
	RebuildKernelCache           types.Bool `tfsdk:"rebuild_kernel_cache"`
	RequestRequiresNetworkTether types.Bool `tfsdk:"request_requires_network_tether"`
	Scope                        types.Set  `tfsdk:"scope"`
	Serials                      types.Set  `tfsdk:"serials"`
	WifiMacs                     types.Set  `tfsdk:"wifi_macs"`
}

// FromBody
func (r *NetworksSmDevicesReboot) toSdkApiRequestCreate(ctx context.Context) *merakigosdk.RequestSmRebootNetworkSmDevices {
	re := *r.Parameters
	var iDs []string = nil
	re.IDs.ElementsAs(ctx, &iDs, false)
	var kextPaths []string = nil
	re.KextPaths.ElementsAs(ctx, &kextPaths, false)
	notifyUser := new(bool)
	if !re.NotifyUser.IsUnknown() && !re.NotifyUser.IsNull() {
		*notifyUser = re.NotifyUser.ValueBool()
	} else {
		notifyUser = nil
	}
	rebuildKernelCache := new(bool)
	if !re.RebuildKernelCache.IsUnknown() && !re.RebuildKernelCache.IsNull() {
		*rebuildKernelCache = re.RebuildKernelCache.ValueBool()
	} else {
		rebuildKernelCache = nil
	}
	requestRequiresNetworkTether := new(bool)
	if !re.RequestRequiresNetworkTether.IsUnknown() && !re.RequestRequiresNetworkTether.IsNull() {
		*requestRequiresNetworkTether = re.RequestRequiresNetworkTether.ValueBool()
	} else {
		requestRequiresNetworkTether = nil
	}
	var scope []string = nil
	re.Scope.ElementsAs(ctx, &scope, false)
	var serials []string = nil
	re.Serials.ElementsAs(ctx, &serials, false)
	var wifiMacs []string = nil
	re.WifiMacs.ElementsAs(ctx, &wifiMacs, false)
	out := merakigosdk.RequestSmRebootNetworkSmDevices{
		IDs:                          iDs,
		KextPaths:                    kextPaths,
		NotifyUser:                   notifyUser,
		RebuildKernelCache:           rebuildKernelCache,
		RequestRequiresNetworkTether: requestRequiresNetworkTether,
		Scope:                        scope,
		Serials:                      serials,
		WifiMacs:                     wifiMacs,
	}
	return &out
}

// ToBody
func ResponseSmRebootNetworkSmDevicesItemToBody(state NetworksSmDevicesReboot, response *merakigosdk.ResponseSmRebootNetworkSmDevices) NetworksSmDevicesReboot {
	itemState := ResponseSmRebootNetworkSmDevices{
		IDs: StringSliceToSet(response.IDs),
	}
	state.Item = &itemState
	return state
}
