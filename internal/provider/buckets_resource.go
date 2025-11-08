package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ resource.Resource = &PersistentBucketsResource{}
var _ resource.ResourceWithImportState = &PersistentBucketsResource{}

var itemObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"weight": types.Int64Type,
		"item":   types.StringType,
	},
}

var bucketsType = types.MapType{
	ElemType: itemObjectType,
}

var nestedItem = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"weight": schema.Int64Attribute{
			Required:    true,
			Description: `Weight to the item in the bucket. Counts against the capacity of the bucket.`,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"item": schema.StringAttribute{
			Optional:    true,
			Description: "Data for the item",
		},
	},
}

type BucketItem struct {
	Weight int64
	Item   string
}

func NewPersistentBucketsResource() resource.Resource {
	return &PersistentBucketsResource{}
}

type PersistentBucketsResource struct {
	client *http.Client
}

type PersistentBucketsResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Items          types.Map    `tfsdk:"items"`
	MaximumBuckets types.Int64  `tfsdk:"maximum_buckets"`
	BucketCapacity types.Int64  `tfsdk:"bucket_capacity"`
	Buckets        types.Set    `tfsdk:"buckets"`
}

func (r *PersistentBucketsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "persistent_buckets"
}

func (r *PersistentBucketsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
			Persistent buckets. Provisions a number of buckets (lists) containing resources
			defined according to bucket capacity and item size. Once a bucket's capacity
			is exhausted, new buckets are created up until ` + "`maximum_buckets`" + `. If items are
			removed, they are also removed from buckets and new items may be placed in 
			the freed space. If maximum buckets are reached, an error is raised.
		`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier (always fixed)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"items": schema.MapNestedAttribute{
				Description:  "Items that are placed in the buckets.",
				NestedObject: nestedItem,
				Required:     true,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
				},
			},
			"maximum_buckets": schema.Int64Attribute{
				Required:    true,
				Description: "Maximum number of buckets to provision.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIf(func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifier.RequiresReplaceIfFuncResponse) {
						if req.PlanValue.ValueInt64() < req.ConfigValue.ValueInt64() {
							resp = &int64planmodifier.RequiresReplaceIfFuncResponse{
								RequiresReplace: true,
							}
						} else {
							resp = &int64planmodifier.RequiresReplaceIfFuncResponse{
								RequiresReplace: true,
							}
						}
					}, "Replace resource if number of buckets is shrunk.", "Replace resource if number of buckets is shrunk."),
				},
			},
			"bucket_capacity": schema.Int64Attribute{
				Required:    true,
				Description: "Capacity of a single bucket.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"buckets": schema.SetAttribute{
				ElementType: types.MapType{
					ElemType: itemObjectType,
				},
				Optional:    true,
				Computed:    true,
				Description: "Set of filled buckets.",
			},
		},
	}
}

func (r *PersistentBucketsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func createItem(weight int64, item string, diagnostics *diag.Diagnostics) *basetypes.ObjectValue {
	obj, diags := types.ObjectValue(itemObjectType.AttrTypes, map[string]attr.Value{
		"weight": types.Int64Value(weight),
		"item":   types.StringValue(item),
	})
	if diags.HasError() {
		return nil
	}
	diagnostics.Append(diags...)
	return &obj
}

func findCapacity(capacities *[]int64, weight int64, maxCapacity int64) *int {
	for k, cap := range *capacities {
		if (cap + weight) <= maxCapacity {
			return &k
		}
	}
	return nil
}

func workTheBuckets(data, state *PersistentBucketsResourceModel, diagnostics *diag.Diagnostics) {
	data.Buckets = basetypes.NewSetNull(bucketsType)
	capacities := make([]int64, data.MaximumBuckets.ValueInt64())
	allBuckets := make([]map[string]BucketItem, data.MaximumBuckets.ValueInt64())
	for idx := 0; idx < int(data.MaximumBuckets.ValueInt64()); idx++ {
		allBuckets[idx] = make(map[string]BucketItem, 0)
		capacities[idx] = 0
	}
	bucketCapacity := data.BucketCapacity.ValueInt64()
	keysInBuckets := make(map[string]int, 0)

	// Fill buckets from TF data
	if state != nil && !state.Buckets.IsUnknown() {
		for bidx, bucket := range state.Buckets.Elements() {
			if bucketItems, ok := bucket.(basetypes.MapValue); ok {
				for k, v := range bucketItems.Elements() {
					if vv, ok := v.(basetypes.ObjectValue); ok {
						objAttrs := vv.Attributes()
						allBuckets[bidx][k] = BucketItem{
							Weight: objAttrs["weight"].(basetypes.Int64Value).ValueInt64(),
							Item:   objAttrs["item"].(basetypes.StringValue).ValueString(),
						}
						capacities[bidx] += allBuckets[bidx][k].Weight
						keysInBuckets[k] = bidx
					}
				}
			}
		}
	}

	keysDefined := make([]string, 0)
	newItems := make(map[string]BucketItem, 0)
	for k, v := range data.Items.Elements() {
		keysDefined = append(keysDefined, k)
		if _, ok := keysInBuckets[k]; !ok {
			if vv, ok := v.(basetypes.ObjectValue); ok {
				objAttrs := vv.Attributes()
				newItems[k] = BucketItem{
					Weight: objAttrs["weight"].(basetypes.Int64Value).ValueInt64(),
					Item:   objAttrs["item"].(basetypes.StringValue).ValueString(),
				}
			}
		} else {
			// Adjust bucket capacities
			keyInBucket := keysInBuckets[k]
			previousWeight := allBuckets[keyInBucket][k].Weight

			if vv, ok := v.(basetypes.ObjectValue); ok {
				objAttrs := vv.Attributes()
				newItem := BucketItem{
					Weight: objAttrs["weight"].(basetypes.Int64Value).ValueInt64(),
					Item:   objAttrs["item"].(basetypes.StringValue).ValueString(),
				}
				newWeight := newItem.Weight

				capacities[keyInBucket] -= previousWeight
				// Check if new weight would require moving the item to a new bucket
				if (capacities[keyInBucket] + newWeight) > bucketCapacity {
					newBucket := findCapacity(&capacities, newWeight, bucketCapacity)
					if newBucket == nil {
						diagnostics.AddError(fmt.Sprintf("unable to find bucket capacity for: %s (previous weight %d, new weight %d)", k, previousWeight, newWeight), fmt.Sprintf("bucket capacities: %+v", capacities))
						return
					}
					delete(allBuckets[keyInBucket], k)
					keysInBuckets[k] = *newBucket
					allBuckets[*newBucket][k] = newItem
					capacities[*newBucket] += newWeight
				} else {
					allBuckets[keyInBucket][k] = newItem
					capacities[keyInBucket] += newWeight
				}
			}
		}
	}

	// Remove items
	for bidx, bucket := range allBuckets {
		for k, v := range bucket {
			if !slices.Contains(keysDefined, k) {
				capacities[bidx] -= v.Weight
				delete(allBuckets[bidx], k)
			}
		}
	}

	// Add new items in buckets with capacity
	for k, v := range newItems {
		targetBucket := findCapacity(&capacities, v.Weight, bucketCapacity)
		if targetBucket == nil {
			diagnostics.AddError(fmt.Sprintf("unable to find bucket capacity for: %s (weight %d)", k, v.Weight), fmt.Sprintf("bucket capacities: %+v", capacities))
			return
		}
		allBuckets[*targetBucket][k] = v
		capacities[*targetBucket] += v.Weight
	}

	// Generate output data
	tfBuckets := make([]attr.Value, 0)
	for _, items := range allBuckets {
		tfItems := make(map[string]attr.Value, 0)
		for k, v := range items {
			tfItem := createItem(v.Weight, v.Item, diagnostics)
			if tfItem == nil {
				diagnostics.AddError(fmt.Sprintf("failed to create a map item for: %s", k), fmt.Sprintf("item: %s", v.Item))
				return
			}
			tfItems[k] = *tfItem
		}
		bucketMap, diags := types.MapValue(itemObjectType, tfItems)
		if diags.HasError() {
			return
		}
		diagnostics.Append(diags...)

		tfBuckets = append(tfBuckets, bucketMap)
	}

	// Set output value
	bucketsValue, diags := types.SetValue(bucketsType, tfBuckets)
	if diags.HasError() {
		diagnostics.AddError("failed to create buckets map for output", "")
		return
	}
	data.Buckets = bucketsValue
	diagnostics.Append(diags...)
}

func (r *PersistentBucketsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PersistentBucketsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue("persistent_buckets")

	workTheBuckets(data, nil, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistentBucketsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PersistentBucketsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistentBucketsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *PersistentBucketsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	workTheBuckets(data, state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistentBucketsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PersistentBucketsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *PersistentBucketsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
