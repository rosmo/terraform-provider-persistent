package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &PersistentCounterResource{}
var _ resource.ResourceWithImportState = &PersistentCounterResource{}

func NewPersistentCounterResource() resource.Resource {
	return &PersistentCounterResource{}
}

type PersistentCounterResource struct {
	client *http.Client
}

type PersistentCounterResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Keys         types.List   `tfsdk:"keys"`
	InitialValue types.Int64  `tfsdk:"initial_value"`
	LastValue    types.Int64  `tfsdk:"last_value"`
	Values       types.Map    `tfsdk:"values"`
}

func (r *PersistentCounterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "persistent_counter"
}

func (r *PersistentCounterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Persistent counter",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keys": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"initial_value": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The initial value to use for the counter.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					Int64DefaultValue(types.Int64Value(0)),
				},
			},
			"last_value": schema.Int64Attribute{
				Computed:    true,
				Description: "The last value that was used for the counter.",
			},
			"values": schema.MapAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				Description: "A map of strings that will cause a change to the counter when any of the values change.",
			},
		},
	}
}

func (r *PersistentCounterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PersistentCounterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PersistentCounterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue("persistent_counter")

	// Generate new set of keys
	if data.Values.IsUnknown() {
		keys := data.Keys.Elements()
		keysLength := len(keys)
		values := make(map[string]int64, keysLength)
		currentValue := data.InitialValue.ValueInt64()
		for i := 0; i < keysLength; i++ {
			keyValue := keys[i].(types.String).ValueString()
			values[keyValue] = currentValue
			currentValue += 1
		}
		data.LastValue = types.Int64Value(currentValue - 1)

		_values, diags := types.MapValueFrom(ctx, types.Int64Type, values)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Values = _values
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistentCounterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PersistentCounterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistentCounterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *PersistentCounterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Keys.Equal(state.Keys) {
		lastValue := state.LastValue.ValueInt64()

		keys := data.Keys.Elements()
		keysLength := len(keys)
		values := make(map[string]int64, keysLength)

		stateValues := state.Values.Elements()
		for _, key := range keys {
			keyValue := key.(types.String).ValueString()
			if val, ok := stateValues[keyValue]; ok {
				values[keyValue] = val.(types.Int64).ValueInt64()
			} else {
				lastValue += 1
				values[keyValue] = lastValue
			}
		}
		data.LastValue = types.Int64Value(lastValue)
		state.LastValue = data.LastValue
		_values, diags := types.MapValueFrom(ctx, types.Int64Type, values)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Values = _values
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistentCounterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PersistentCounterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *PersistentCounterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
