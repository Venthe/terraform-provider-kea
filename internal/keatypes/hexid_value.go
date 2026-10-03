package keatypes

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/venthe/terraform-provider-kea/kea"
)

var (
	_ basetypes.StringValuable                   = (*HexIDValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*HexIDValue)(nil)
	_ xattr.ValidateableAttribute                = (*HexIDValue)(nil)
	_ function.ValidateableParameter             = (*HexIDValue)(nil)
)

type HexIDValue struct {
	basetypes.StringValue
}
type HexID = HexIDValue

func (v HexIDValue) Equal(o attr.Value) bool {
	other, ok := o.(HexIDValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v HexIDValue) Type(ctx context.Context) attr.Type {
	return HexIDType{}
}

func (v HexIDValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(HexIDValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)
		return false, diags
	}

	// Already validated at this point, ignoring errors
	newHexID, _ := kea.ParseHexID(newValue.ValueString())
	currentHexID, _ := kea.ParseHexID(v.ValueString())

	return bytes.Equal(currentHexID, newHexID), diags
}

func (v HexIDValue) ValidateAttribute(_ context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}
	_, err := kea.ParseHexID(v.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Kea Hex Identifier String Value",
			"A string value was provided that is not valid Kea hex identifier string format.\n\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
	}
}

func (v HexIDValue) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}
	_, err := kea.ParseHexID(v.ValueString())
	if err != nil {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid Kea Hex Identifier String Value: "+
				"A string value was provided that is not valid Kea hex identifier string format.\n\n"+
				fmt.Sprintf("Position: %d", req.Position)+"\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
	}
}

func (v HexIDValue) ValueHexID() (kea.HexID, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.AddError("HexID ValueHexID Error", "Hex identifier string value is null")
		return nil, diags
	}

	if v.IsUnknown() {
		diags.AddError("HexID ValueHexID Error", "Hex identifier string value is unknown")
		return nil, diags
	}

	hexID, err := kea.ParseHexID(v.ValueString())
	if err != nil {
		diags.AddError("HexID ValueHexID Error", err.Error())
		return nil, diags
	}

	return hexID, nil
}

func NewHexIDNull() HexIDValue {
	return HexIDValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewHexIDUnknown() HexIDValue {
	return HexIDValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewHexIDValue(value string) HexIDValue {
	return HexIDValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewHexIDPointerValue(value *string) HexIDValue {
	return HexIDValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
