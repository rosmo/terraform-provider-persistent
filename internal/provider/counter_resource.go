package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	Reuse        types.Bool   `tfsdk:"reuse"`
	InitialValue types.Int64  `tfsdk:"initial_value"`
	LastValue    types.Int64  `tfsdk:"last_value"`
	Values       types.Map    `tfsdk:"values"`
}

func (r *PersistentCounterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "persistent_counter"
}

func (r *PersistentCounterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
			Persistent counter. Generates sequentially increasing number counters for the strings specified 
			in the ` + "`keys`" + ` variable. As long as a specified key exist, it will always receive the same counter
			value. No counter value is re-used even if a new key is added and counter values will only ever
			increase.
		`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier (always fixed)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keys": schema.ListAttribute{
				Description: "List of keys to generate counters for.",
				ElementType: types.StringType,
				Required:    true,
			},
			"reuse": schema.BoolAttribute{
				Description: "Allows reusing freed keys for new ones.",
				Optional:    true,
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
				Description: "A map of keys to counter values.",
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
		keys := convertKeys(data.Keys.Elements())
		last, values := assignKeys(
			keys, nil, data.Reuse.ValueBool(),
			// initial value
			data.InitialValue.ValueInt64(),
			// use initial value - 1 for last value
			data.InitialValue.ValueInt64()-1,
		)

		data.LastValue = types.Int64Value(last)
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

	if !data.Keys.Equal(state.Keys) || !data.InitialValue.Equal(state.Keys) {

		keys := convertKeys(data.Keys.Elements())
		stateVals := convertState(state.Values.Elements())

		last, values := assignKeys(
			keys, stateVals, data.Reuse.ValueBool(),
			data.InitialValue.ValueInt64(),
			state.LastValue.ValueInt64(),
		)

		data.LastValue = types.Int64Value(last)
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

// convertKeys generates a string slice from the terraform string list representation
func convertKeys(tfKeys []attr.Value) []string {
	keys := make([]string, 0, len(tfKeys))
	for _, k := range tfKeys {
		val, ok := k.(types.String)
		if !ok {
			// conversion failed
			continue
		}
		str := val.ValueString()
		if str == "" {
			// invalid string value or string empty
			continue
		}
		keys = append(keys, str)

	}
	return keys
}

// convertState converts the counter state from terraform format to map[string]int64
func convertState(tfState map[string]attr.Value) map[string]int64 {
	state := make(map[string]int64, len(tfState))
	for k, v := range tfState {
		intVal, ok := v.(types.Int64)
		if !ok {
			continue
		}
		state[k] = intVal.ValueInt64()
	}
	return state
}
