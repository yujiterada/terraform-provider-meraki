package provider

// RESOURCE NORMAL
import (
	"context"
	"fmt"
	"net/url"
	"strings"

	merakigosdk "github.com/meraki/dashboard-api-go/v3/sdk"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &OrganizationsPolicyObjectsResource{}
	_ resource.ResourceWithConfigure = &OrganizationsPolicyObjectsResource{}
)

func NewOrganizationsPolicyObjectsResource() resource.Resource {
	return &OrganizationsPolicyObjectsResource{}
}

type OrganizationsPolicyObjectsResource struct {
	client *merakigosdk.Client
}

func (r *OrganizationsPolicyObjectsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client := req.ProviderData.(MerakiProviderData).Client
	r.client = client
}

// Metadata returns the data source type name.
func (r *OrganizationsPolicyObjectsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_policy_objects"
}

func (r *OrganizationsPolicyObjectsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"category": schema.StringAttribute{
				MarkdownDescription: `Category of a policy object (one of: adaptivePolicy, network)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SuppressDiffString(),
				},
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: `CIDR Value of a policy object (e.g. 10.11.12.1/24")`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"fqdn": schema.StringAttribute{
				MarkdownDescription: `Fully qualified domain name of policy object (e.g. "example.com")`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_ids": schema.SetAttribute{
				MarkdownDescription: `The IDs of policy object groups the policy object belongs to`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},

				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"ip": schema.StringAttribute{
				MarkdownDescription: `IP Address of a policy object (e.g. "1.2.3.4")`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mask": schema.StringAttribute{
				MarkdownDescription: `Mask of a policy object (e.g. "255.255.0.0")`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: `Name of a policy object, unique within the organization (alphanumeric, space, dash, or underscore characters only)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_ids": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: `organizationId path parameter. Organization ID`,
				Required:            true,
			},
			"policy_object_id": schema.StringAttribute{
				MarkdownDescription: `policyObjectId path parameter. Policy object ID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: `Type of a policy object (one of: adaptivePolicyIpv4Cidr, cidr, fqdn, ipAndMask)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SuppressDiffString(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

//path params to set ['policyObjectId']
//path params to assign NOT EDITABLE ['category', 'type']

func (r *OrganizationsPolicyObjectsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var data OrganizationsPolicyObjectsRs

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
	vvOrganizationID := data.OrganizationID.ValueString()
	// organization_id
	vvName := data.Name.ValueString()
	//Items
	responseVerifyItem, restyResp1, err := getAllItemsOrganizationsPolicyObjects(*r.client, vvOrganizationID)
	//Have Create
	if err != nil {
		if restyResp1 != nil {
			if restyResp1.StatusCode() != 404 {
				resp.Diagnostics.AddError(
					"Failure when executing GetOrganizationPolicyObjects",
					err.Error(),
				)
				return
			}
		}
	}
	if responseVerifyItem != nil {
		responseStruct := structToMap(responseVerifyItem)
		result := getDictResult(responseStruct, "Name", vvName, simpleCmp)
		if result != nil {
			result2 := result.(map[string]interface{})
			vvPolicyObjectID, ok := result2["ID"].(string)
			if !ok {
				resp.Diagnostics.AddError(
					"Failure when parsing path parameter PolicyObjectID",
					err.Error(),
				)
				return
			}
			r.client.Organizations.UpdateOrganizationPolicyObject(vvOrganizationID, vvPolicyObjectID, data.toSdkApiRequestUpdate(ctx))
			responseVerifyItem2, _, _ := r.client.Organizations.GetOrganizationPolicyObject(vvOrganizationID, vvPolicyObjectID)
			if responseVerifyItem2 != nil {
				data = ResponseOrganizationsGetOrganizationPolicyObjectItemToBodyRs(data, responseVerifyItem2, false)
				// Path params update assigned
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
				return
			}
		}
	}
	dataRequest := data.toSdkApiRequestCreate(ctx)
	_, restyResp2, err := r.client.Organizations.CreateOrganizationPolicyObject(vvOrganizationID, dataRequest)

	if err != nil || restyResp2 == nil {
		if restyResp1 != nil {
			resp.Diagnostics.AddError(
				"Failure when executing CreateOrganizationPolicyObject",
				err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failure when executing CreateOrganizationPolicyObject",
			err.Error(),
		)
		return
	}
	//Items
	responseGet, restyResp1, err := getAllItemsOrganizationsPolicyObjects(*r.client, vvOrganizationID)
	// Has item and has items

	if err != nil || responseGet == nil {
		if restyResp1 != nil {
			resp.Diagnostics.AddError(
				"Failure when executing GetOrganizationPolicyObjects",
				err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failure when executing GetOrganizationPolicyObjects",
			err.Error(),
		)
		return
	}
	responseStruct := structToMap(responseGet)
	result := getDictResult(responseStruct, "Name", vvName, simpleCmp)
	if result != nil {
		result2 := result.(map[string]interface{})
		vvPolicyObjectID, ok := result2["ID"].(string)
		if !ok {
			resp.Diagnostics.AddError(
				"Failure when parsing path parameter PolicyObjectID",
				err.Error(),
			)
			return
		}
		responseVerifyItem2, restyRespGet, err := r.client.Organizations.GetOrganizationPolicyObject(vvOrganizationID, vvPolicyObjectID)
		if responseVerifyItem2 != nil && err == nil {
			data = ResponseOrganizationsGetOrganizationPolicyObjectItemToBodyRs(data, responseVerifyItem2, false)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		} else {
			if restyRespGet != nil {
				resp.Diagnostics.AddError(
					"Failure when executing GetOrganizationPolicyObject",
					err.Error(),
				)
				return
			}
			resp.Diagnostics.AddError(
				"Failure when executing GetOrganizationPolicyObject",
				err.Error(),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Error in result.",
			"Error in result.",
		)
		return
	}
}

func (r *OrganizationsPolicyObjectsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationsPolicyObjectsRs

	var item types.Object

	resp.Diagnostics.Append(req.State.Get(ctx, &item)...)
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
	// Has Item2

	vvOrganizationID := data.OrganizationID.ValueString()
	// organization_id
	vvPolicyObjectID := data.PolicyObjectID.ValueString()
	// policy_object_id
	responseGet, restyRespGet, err := r.client.Organizations.GetOrganizationPolicyObject(vvOrganizationID, vvPolicyObjectID)
	if err != nil || restyRespGet == nil {
		if restyRespGet != nil {
			if restyRespGet.StatusCode() == 404 {
				resp.Diagnostics.AddWarning(
					"Resource not found",
					"Deleting resource",
				)
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				"Failure when executing GetOrganizationPolicyObject",
				err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failure when executing GetOrganizationPolicyObject",
			err.Error(),
		)
		return
	}

	data = ResponseOrganizationsGetOrganizationPolicyObjectItemToBodyRs(data, responseGet, true)
	diags := resp.State.Set(ctx, &data)
	//update path params assigned
	resp.Diagnostics.Append(diags...)
}
func (r *OrganizationsPolicyObjectsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: attr_one,attr_two. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_object_id"), idParts[1])...)
}

func (r *OrganizationsPolicyObjectsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationsPolicyObjectsRs
	merge(ctx, req, resp, &data)

	if resp.Diagnostics.HasError() {
		return
	}
	//Has Paths
	//Update

	//Path Params
	vvOrganizationID := data.OrganizationID.ValueString()
	// organization_id
	vvPolicyObjectID := data.PolicyObjectID.ValueString()
	dataRequest := data.toSdkApiRequestUpdate(ctx)
	_, restyResp2, err := r.client.Organizations.UpdateOrganizationPolicyObject(vvOrganizationID, vvPolicyObjectID, dataRequest)
	if err != nil || restyResp2 == nil {
		if restyResp2 != nil {
			resp.Diagnostics.AddError(
				"Failure when executing UpdateOrganizationPolicyObject",
				err.Error(),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failure when executing UpdateOrganizationPolicyObject",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(req.Plan.Set(ctx, &data)...)
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *OrganizationsPolicyObjectsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationsPolicyObjectsRs
	var item types.Object

	resp.Diagnostics.Append(req.State.Get(ctx, &item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(item.As(ctx, &state, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}

	vvOrganizationID := state.OrganizationID.ValueString()
	vvPolicyObjectID := state.PolicyObjectID.ValueString()
	_, err := r.client.Organizations.DeleteOrganizationPolicyObject(vvOrganizationID, vvPolicyObjectID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure when executing DeleteOrganizationPolicyObject", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)

}

// TF Structs Schema
type OrganizationsPolicyObjectsRs struct {
	OrganizationID types.String `tfsdk:"organization_id"`
	PolicyObjectID types.String `tfsdk:"policy_object_id"`
	Category       types.String `tfsdk:"category"`
	Cidr           types.String `tfsdk:"cidr"`
	CreatedAt      types.String `tfsdk:"created_at"`
	GroupIDs       types.Set    `tfsdk:"group_ids"`
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	NetworkIDs     types.Set    `tfsdk:"network_ids"`
	Type           types.String `tfsdk:"type"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
	Fqdn           types.String `tfsdk:"fqdn"`
	IP             types.String `tfsdk:"ip"`
	Mask           types.String `tfsdk:"mask"`
}

// FromBody
func (r *OrganizationsPolicyObjectsRs) toSdkApiRequestCreate(ctx context.Context) *merakigosdk.RequestOrganizationsCreateOrganizationPolicyObject {
	emptyString := ""
	category := new(string)
	if !r.Category.IsUnknown() && !r.Category.IsNull() {
		*category = r.Category.ValueString()
	} else {
		category = &emptyString
	}
	cidr := new(string)
	if !r.Cidr.IsUnknown() && !r.Cidr.IsNull() {
		*cidr = r.Cidr.ValueString()
	} else {
		cidr = &emptyString
	}
	fqdn := new(string)
	if !r.Fqdn.IsUnknown() && !r.Fqdn.IsNull() {
		*fqdn = r.Fqdn.ValueString()
	} else {
		fqdn = &emptyString
	}
	var groupIDs []string = nil
	r.GroupIDs.ElementsAs(ctx, &groupIDs, false)
	iP := new(string)
	if !r.IP.IsUnknown() && !r.IP.IsNull() {
		*iP = r.IP.ValueString()
	} else {
		iP = &emptyString
	}
	mask := new(string)
	if !r.Mask.IsUnknown() && !r.Mask.IsNull() {
		*mask = r.Mask.ValueString()
	} else {
		mask = &emptyString
	}
	name := new(string)
	if !r.Name.IsUnknown() && !r.Name.IsNull() {
		*name = r.Name.ValueString()
	} else {
		name = &emptyString
	}
	typeR := new(string)
	if !r.Type.IsUnknown() && !r.Type.IsNull() {
		*typeR = r.Type.ValueString()
	} else {
		typeR = &emptyString
	}
	out := merakigosdk.RequestOrganizationsCreateOrganizationPolicyObject{
		Category: *category,
		Cidr:     *cidr,
		Fqdn:     *fqdn,
		GroupIDs: groupIDs,
		IP:       *iP,
		Mask:     *mask,
		Name:     *name,
		Type:     *typeR,
	}
	return &out
}
func (r *OrganizationsPolicyObjectsRs) toSdkApiRequestUpdate(ctx context.Context) *merakigosdk.RequestOrganizationsUpdateOrganizationPolicyObject {
	emptyString := ""
	cidr := new(string)
	if !r.Cidr.IsUnknown() && !r.Cidr.IsNull() {
		*cidr = r.Cidr.ValueString()
	} else {
		cidr = &emptyString
	}
	fqdn := new(string)
	if !r.Fqdn.IsUnknown() && !r.Fqdn.IsNull() {
		*fqdn = r.Fqdn.ValueString()
	} else {
		fqdn = &emptyString
	}
	var groupIDs []string = nil
	r.GroupIDs.ElementsAs(ctx, &groupIDs, false)
	iP := new(string)
	if !r.IP.IsUnknown() && !r.IP.IsNull() {
		*iP = r.IP.ValueString()
	} else {
		iP = &emptyString
	}
	mask := new(string)
	if !r.Mask.IsUnknown() && !r.Mask.IsNull() {
		*mask = r.Mask.ValueString()
	} else {
		mask = &emptyString
	}
	name := new(string)
	if !r.Name.IsUnknown() && !r.Name.IsNull() {
		*name = r.Name.ValueString()
	} else {
		name = &emptyString
	}
	out := merakigosdk.RequestOrganizationsUpdateOrganizationPolicyObject{
		Cidr:     *cidr,
		Fqdn:     *fqdn,
		GroupIDs: groupIDs,
		IP:       *iP,
		Mask:     *mask,
		Name:     *name,
	}
	return &out
}

// From gosdk to TF Structs Schema
func ResponseOrganizationsGetOrganizationPolicyObjectItemToBodyRs(state OrganizationsPolicyObjectsRs, response *merakigosdk.ResponseOrganizationsGetOrganizationPolicyObject, is_read bool) OrganizationsPolicyObjectsRs {
	itemState := OrganizationsPolicyObjectsRs{
		Category:   types.StringValue(response.Category),
		Cidr:       types.StringValue(response.Cidr),
		CreatedAt:  types.StringValue(response.CreatedAt),
		GroupIDs:   StringSliceToSet(response.GroupIDs),
		ID:         types.StringValue(response.ID),
		Name:       types.StringValue(response.Name),
		NetworkIDs: StringSliceToSet(response.NetworkIDs),
		Type:       types.StringValue(response.Type),
		UpdatedAt:  types.StringValue(response.UpdatedAt),
	}
	if is_read {
		return mergeInterfacesOnlyPath(state, itemState).(OrganizationsPolicyObjectsRs)
	}
	return mergeInterfaces(state, itemState, true).(OrganizationsPolicyObjectsRs)
}

func getAllItemsOrganizationsPolicyObjects(client merakigosdk.Client, organizationId string) (merakigosdk.ResponseOrganizationsGetOrganizationPolicyObjects, *resty.Response, error) {
	var all_response merakigosdk.ResponseOrganizationsGetOrganizationPolicyObjects
	response, r2, er := client.Organizations.GetOrganizationPolicyObjects(organizationId, &merakigosdk.GetOrganizationPolicyObjectsQueryParams{
		PerPage: 1000,
	})
	count := 0
	all_response = append(all_response, *response...)
	for len(*response) >= 1000 {
		count += 1
		fmt.Println(count)
		links := strings.Split(r2.Header().Get("Link"), ",")
		var link string
		if count > 1 {
			link = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.Split(links[2], ";")[0], ">", ""), "<", ""), client.RestyClient().BaseURL, "")
		} else {
			link = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.Split(links[1], ";")[0], ">", ""), "<", ""), client.RestyClient().BaseURL, "")
		}
		myUrl, _ := url.Parse(link)
		params, _ := url.ParseQuery(myUrl.RawQuery)
		if params["endingBefore"] != nil {
			response, r2, er = client.Organizations.GetOrganizationPolicyObjects(organizationId, &merakigosdk.GetOrganizationPolicyObjectsQueryParams{
				PerPage:      1000,
				EndingBefore: params["endingBefore"][0],
			})
			all_response = append(all_response, *response...)
		} else {
			break
		}
	}

	return all_response, r2, er
}
