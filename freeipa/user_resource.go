// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *ipa.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	FirstName              types.String `tfsdk:"first_name"`
	LastName               types.String `tfsdk:"last_name"`
	UID                    types.String `tfsdk:"name"`
	FullName               types.String `tfsdk:"full_name"`
	DisplayName            types.String `tfsdk:"display_name"`
	Initials               types.String `tfsdk:"initials"`
	HomeDirectory          types.String `tfsdk:"home_directory"`
	Gecos                  types.String `tfsdk:"gecos"`
	LoginShell             types.String `tfsdk:"login_shell"`
	KrbPrincipalName       types.List   `tfsdk:"krb_principal_name"`
	KrbPrincipalExpiration types.String `tfsdk:"krb_principal_expiration"`
	KrbPasswordExpiration  types.String `tfsdk:"krb_password_expiration"`
	UserPassword           types.String `tfsdk:"userpassword"`
	EmailAddress           types.List   `tfsdk:"email_address"`
	TelephoneNumbers       types.List   `tfsdk:"telephone_numbers"`
	MobileNumbers          types.List   `tfsdk:"mobile_numbers"`
	RandomPassword         types.Bool   `tfsdk:"random_password"`
	UidNumber              types.Int32  `tfsdk:"uid_number"`
	GidNumber              types.Int32  `tfsdk:"gid_number"`
	StreetAddress          types.String `tfsdk:"street_address"`
	City                   types.String `tfsdk:"city"`
	Province               types.String `tfsdk:"province"`
	PostalCode             types.String `tfsdk:"postal_code"`
	OrganisationUnit       types.String `tfsdk:"organisation_unit"`
	JobTitle               types.String `tfsdk:"job_title"`
	Manager                types.String `tfsdk:"manager"`
	EmployeeNumber         types.String `tfsdk:"employee_number"`
	EmployeeType           types.String `tfsdk:"employee_type"`
	PreferredLanguage      types.String `tfsdk:"preferred_language"`
	AccountDisabled        types.Bool   `tfsdk:"account_disabled"`
	SshPublicKeys          types.List   `tfsdk:"ssh_public_key"`
	CarLicense             types.List   `tfsdk:"car_license"`
	UserClass              types.List   `tfsdk:"userclass"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FreeIPA User resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "UID or login",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "Full name",
				Optional:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Optional:            true,
			},
			"initials": schema.StringAttribute{
				MarkdownDescription: "Initials",
				Optional:            true,
			},
			"home_directory": schema.StringAttribute{
				MarkdownDescription: "Home Directory",
				Optional:            true,
			},
			"gecos": schema.StringAttribute{
				MarkdownDescription: "GECOS",
				Optional:            true,
			},
			"login_shell": schema.StringAttribute{
				MarkdownDescription: "Login Shell",
				Optional:            true,
			},
			"krb_principal_name": schema.ListAttribute{
				MarkdownDescription: "Principal alias",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"krb_principal_expiration": schema.StringAttribute{
				MarkdownDescription: "Kerberos principal expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Optional: true,
			},
			"krb_password_expiration": schema.StringAttribute{
				MarkdownDescription: "User password expiration " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`)",
				Optional: true,
			},
			"userpassword": schema.StringAttribute{
				MarkdownDescription: "Prompt to set the user password",
				Optional:            true,
				Sensitive:           true,
			},
			"email_address": schema.ListAttribute{
				MarkdownDescription: "Email address",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"telephone_numbers": schema.ListAttribute{
				MarkdownDescription: "Telephone Number",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"mobile_numbers": schema.ListAttribute{
				MarkdownDescription: "Mobile Number",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"random_password": schema.BoolAttribute{
				MarkdownDescription: "Generate a random user password",
				Optional:            true,
			},
			"uid_number": schema.Int32Attribute{
				MarkdownDescription: "User ID Number (system will assign one if not provided)",
				Optional:            true,
			},
			"gid_number": schema.Int32Attribute{
				MarkdownDescription: "Group ID Number",
				Optional:            true,
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "Street address",
				Optional:            true,
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "City",
				Optional:            true,
			},
			"province": schema.StringAttribute{
				MarkdownDescription: "Province",
				Optional:            true,
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "Postal code",
				Optional:            true,
			},
			"organisation_unit": schema.StringAttribute{
				MarkdownDescription: "Org. Unit",
				Optional:            true,
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "Job Title",
				Optional:            true,
			},
			"manager": schema.StringAttribute{
				MarkdownDescription: "Manager",
				Optional:            true,
			},
			"employee_number": schema.StringAttribute{
				MarkdownDescription: "Employee Number",
				Optional:            true,
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "Employee Type",
				Optional:            true,
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "Preferred Language",
				Optional:            true,
			},
			"account_disabled": schema.BoolAttribute{
				MarkdownDescription: "Account disabled",
				Optional:            true,
			},
			"ssh_public_key": schema.ListAttribute{
				MarkdownDescription: "List of SSH public keys",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"car_license": schema.ListAttribute{
				MarkdownDescription: "Car Licenses",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"userclass": schema.ListAttribute{
				MarkdownDescription: "User category (semantics placed on this attribute are for local interpretation)",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ipa.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	optArgs := ipa.UserAddOptionalArgs{}

	args := ipa.UserAddArgs{
		Givenname: string(data.FirstName.ValueString()),
		Sn:        string(data.LastName.ValueString()),
	}
	optArgs.UID = data.UID.ValueStringPointer()
	if !data.FullName.IsUnknown() {
		optArgs.Cn = data.FullName.ValueStringPointer()
	}
	if !data.DisplayName.IsUnknown() {
		optArgs.Displayname = data.DisplayName.ValueStringPointer()
	}
	if !data.Initials.IsNull() {
		optArgs.Initials = data.Initials.ValueStringPointer()
	}
	if !data.HomeDirectory.IsNull() {
		optArgs.Homedirectory = data.HomeDirectory.ValueStringPointer()
	}
	if !data.Gecos.IsNull() {
		optArgs.Gecos = data.Gecos.ValueStringPointer()
	}
	if !data.LoginShell.IsNull() {
		optArgs.Loginshell = data.LoginShell.ValueStringPointer()
	}
	if len(data.KrbPrincipalName.Elements()) > 0 {
		var v []string
		for _, value := range data.KrbPrincipalName.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Krbprincipalname = &v
	}
	if !data.UserPassword.IsNull() {
		optArgs.Userpassword = data.UserPassword.ValueStringPointer()
	}
	if len(data.EmailAddress.Elements()) > 0 {
		var v []string
		for _, value := range data.EmailAddress.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mail = &v
	}
	if len(data.TelephoneNumbers.Elements()) > 0 {
		var v []string
		for _, value := range data.TelephoneNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Telephonenumber = &v
	}
	if len(data.MobileNumbers.Elements()) > 0 {
		var v []string
		for _, value := range data.MobileNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mobile = &v
	}
	if !data.RandomPassword.IsUnknown() {
		optArgs.Random = data.RandomPassword.ValueBoolPointer()
	}
	if !data.UidNumber.IsNull() {
		uid := int(data.UidNumber.ValueInt32())
		optArgs.Uidnumber = &uid
	}
	if !data.GidNumber.IsNull() {
		gid := int(data.GidNumber.ValueInt32())
		optArgs.Gidnumber = &gid
	}
	if !data.StreetAddress.IsNull() {
		optArgs.Street = data.StreetAddress.ValueStringPointer()
	}
	if !data.City.IsNull() {
		optArgs.L = data.City.ValueStringPointer()
	}
	if !data.Province.IsNull() {
		optArgs.St = data.Province.ValueStringPointer()
	}
	if !data.PostalCode.IsNull() {
		optArgs.Postalcode = data.PostalCode.ValueStringPointer()
	}
	if !data.OrganisationUnit.IsNull() {
		optArgs.Ou = data.OrganisationUnit.ValueStringPointer()
	}
	if !data.JobTitle.IsNull() {
		optArgs.Title = data.JobTitle.ValueStringPointer()
	}
	if !data.Manager.IsNull() {
		optArgs.Manager = data.Manager.ValueStringPointer()
	}
	if !data.EmployeeNumber.IsNull() {
		optArgs.Employeenumber = data.EmployeeNumber.ValueStringPointer()
	}
	if !data.EmployeeType.IsNull() {
		optArgs.Employeetype = data.EmployeeType.ValueStringPointer()
	}
	if !data.PreferredLanguage.IsNull() {
		optArgs.Preferredlanguage = data.PreferredLanguage.ValueStringPointer()
	}
	if !data.AccountDisabled.IsUnknown() {
		optArgs.Nsaccountlock = data.AccountDisabled.ValueBoolPointer()
	}
	if len(data.SshPublicKeys.Elements()) > 0 {
		var v []string
		for _, value := range data.SshPublicKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipasshpubkey = &v
	}
	if len(data.CarLicense.Elements()) > 0 {
		var v []string
		for _, value := range data.CarLicense.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Carlicense = &v
	}
	if !data.KrbPrincipalExpiration.IsNull() {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPrincipalExpiration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		optArgs.Krbprincipalexpiration = &timestamp
	}
	if !data.KrbPasswordExpiration.IsNull() {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPasswordExpiration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_password_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		optArgs.Krbpasswordexpiration = &timestamp
	}
	if len(data.UserClass.Elements()) > 0 {
		var v []string
		for _, value := range data.UserClass.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Userclass = &v
	}

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.UserAdd(&args, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating freeipa user group: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Create freeipa user %s returns %s", data.UID.String(), res.String()))

	data.Id = types.StringValue(data.UID.ValueString())

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	all := true
	optArgs := ipa.UserShowOptionalArgs{
		All: &all,
	}

	optArgs.UID = data.UID.ValueStringPointer()

	res, err := r.client.UserShow(&ipa.UserShowArgs{}, &optArgs)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			tflog.Debug(ctx, fmt.Sprintf("[DEBUG] User %s not found", data.Id.ValueString()))
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading user %s: %s", data.Id.ValueString(), err))
		}
	}

	if res.Result.Cn != nil && !data.FullName.IsNull() {
		data.FullName = types.StringValue(*res.Result.Cn)
	}
	if res.Result.Displayname != nil && !data.DisplayName.IsNull() {
		data.DisplayName = types.StringValue(*res.Result.Displayname)
	}
	if res.Result.Initials != nil && !data.Initials.IsNull() {
		data.Initials = types.StringValue(*res.Result.Initials)
	}
	if res.Result.Homedirectory != nil && !data.HomeDirectory.IsNull() {
		data.HomeDirectory = types.StringValue(*res.Result.Homedirectory)
	}
	if res.Result.Gecos != nil && !data.Gecos.IsNull() {
		data.Gecos = types.StringValue(*res.Result.Gecos)
	}
	if res.Result.Loginshell != nil && !data.LoginShell.IsNull() {
		data.LoginShell = types.StringValue(*res.Result.Loginshell)
	}
	if res.Result.Krbprincipalname != nil && !data.KrbPrincipalName.IsNull() {
		data.KrbPrincipalName, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Krbprincipalname)
	}
	//
	// if res.Result.Userpassword != nil {
	// 	data.UserPassword = types.StringValue(*res.Result.Userpassword)
	// }
	if res.Result.Mail != nil && !data.EmailAddress.IsNull() {
		data.EmailAddress, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mail)
	}
	if res.Result.Telephonenumber != nil {
		data.TelephoneNumbers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Telephonenumber)
	}
	if res.Result.Mobile != nil {
		data.MobileNumbers, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Mobile)
	}
	if res.Result.Random != nil {
		data.RandomPassword = types.BoolValue(*res.Result.Random)
	}
	if res.Result.Uidnumber != nil && !data.UidNumber.IsNull() {
		data.UidNumber = types.Int32Value(int32(*res.Result.Uidnumber))
	}
	if res.Result.Gidnumber != nil && !data.GidNumber.IsNull() {
		data.GidNumber = types.Int32Value(int32(*res.Result.Gidnumber))
	}
	if res.Result.Street != nil {
		data.StreetAddress = types.StringValue(*res.Result.Street)
	}
	if res.Result.L != nil {
		data.City = types.StringValue(*res.Result.L)
	}
	if res.Result.St != nil {
		data.Province = types.StringValue(*res.Result.St)
	}
	if res.Result.Postalcode != nil {
		data.PostalCode = types.StringValue(*res.Result.Postalcode)
	}
	if res.Result.Ou != nil {
		data.OrganisationUnit = types.StringValue(*res.Result.Ou)
	}
	if res.Result.Title != nil {
		data.JobTitle = types.StringValue(*res.Result.Title)
	}
	if res.Result.Manager != nil {
		data.Manager = types.StringValue(*res.Result.Manager)
	}
	if res.Result.Employeenumber != nil {
		data.EmployeeNumber = types.StringValue(*res.Result.Employeenumber)
	}
	if res.Result.Employeetype != nil {
		data.EmployeeType = types.StringValue(*res.Result.Employeetype)
	}
	if res.Result.Preferredlanguage != nil {
		data.PreferredLanguage = types.StringValue(*res.Result.Preferredlanguage)
	}
	// if res.Result.Nsaccountlock != nil {
	// 	data.AccountDisabled = types.BoolValue(*res.Result.Nsaccountlock)
	// }
	if res.Result.Ipasshpubkey != nil {
		data.SshPublicKeys, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Ipasshpubkey)
	}
	if res.Result.Carlicense != nil {
		data.CarLicense, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Carlicense)
	}
	if res.Result.Krbprincipalexpiration != nil && !data.KrbPrincipalExpiration.IsUnknown() {
		timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", res.Result.Krbprincipalexpiration.String())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		data.KrbPrincipalExpiration = types.StringValue(timestamp.Format(time.RFC3339))
	}
	if res.Result.Krbpasswordexpiration != nil && !data.KrbPasswordExpiration.IsUnknown() {
		timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", res.Result.Krbpasswordexpiration.String())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
			return
		}
		data.KrbPasswordExpiration = types.StringValue(timestamp.Format(time.RFC3339))
	}
	if res.Result.Userclass != nil {
		data.UserClass, _ = types.ListValueFrom(ctx, types.StringType, res.Result.Userclass)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	optArgs := ipa.UserModOptionalArgs{}

	if !data.UID.IsUnknown() {
		optArgs.UID = data.UID.ValueStringPointer()
	}
	if !data.FullName.IsUnknown() {
		optArgs.Cn = data.FullName.ValueStringPointer()
	}
	if !data.FirstName.IsUnknown() {
		optArgs.Givenname = data.FirstName.ValueStringPointer()
	}
	if !data.LastName.IsUnknown() {
		optArgs.Sn = data.LastName.ValueStringPointer()
	}
	if !data.DisplayName.IsUnknown() {
		optArgs.Displayname = data.DisplayName.ValueStringPointer()
	}
	if !data.Initials.IsUnknown() {
		optArgs.Initials = data.Initials.ValueStringPointer()
	}
	if !data.HomeDirectory.IsUnknown() {
		optArgs.Homedirectory = data.HomeDirectory.ValueStringPointer()
	}
	if !data.Gecos.IsUnknown() {
		optArgs.Gecos = data.Gecos.ValueStringPointer()
	}
	if !data.LoginShell.IsUnknown() {
		optArgs.Loginshell = data.LoginShell.ValueStringPointer()
	}
	if !data.UserPassword.IsUnknown() {
		optArgs.Userpassword = data.UserPassword.ValueStringPointer()
	}
	if !data.RandomPassword.IsUnknown() {
		optArgs.Random = data.RandomPassword.ValueBoolPointer()
	}
	if !data.UidNumber.IsNull() {
		uid := int(data.UidNumber.ValueInt32())
		optArgs.Uidnumber = &uid
	}
	if !data.GidNumber.IsNull() {
		gid := int(data.GidNumber.ValueInt32())
		optArgs.Gidnumber = &gid
	}
	if !data.StreetAddress.IsUnknown() {
		optArgs.Street = data.StreetAddress.ValueStringPointer()
	}
	if !data.City.IsUnknown() {
		optArgs.L = data.City.ValueStringPointer()
	}
	if !data.Province.IsUnknown() {
		optArgs.St = data.Province.ValueStringPointer()
	}
	if !data.PostalCode.IsUnknown() {
		optArgs.Postalcode = data.PostalCode.ValueStringPointer()
	}
	if !data.OrganisationUnit.IsUnknown() {
		optArgs.Ou = data.OrganisationUnit.ValueStringPointer()
	}
	if !data.JobTitle.IsUnknown() {
		optArgs.Title = data.JobTitle.ValueStringPointer()
	}
	if !data.EmployeeNumber.IsUnknown() {
		optArgs.Employeenumber = data.EmployeeNumber.ValueStringPointer()
	}
	if !data.EmployeeType.IsUnknown() {
		optArgs.Employeetype = data.EmployeeType.ValueStringPointer()
	}
	if !data.PreferredLanguage.IsUnknown() {
		optArgs.Preferredlanguage = data.PreferredLanguage.ValueStringPointer()
	}
	if !data.AccountDisabled.IsUnknown() {
		optArgs.Nsaccountlock = data.AccountDisabled.ValueBoolPointer()
	}
	if len(data.TelephoneNumbers.Elements()) > 0 {
		var v []string
		for _, value := range data.TelephoneNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Telephonenumber = &v
	}
	if len(data.MobileNumbers.Elements()) > 0 {
		var v []string
		for _, value := range data.MobileNumbers.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mobile = &v
	}
	if len(data.KrbPrincipalName.Elements()) > 0 {
		var v []string
		for _, value := range data.KrbPrincipalName.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Krbprincipalname = &v
	}
	if len(data.SshPublicKeys.Elements()) > 0 {
		var v []string
		for _, value := range data.SshPublicKeys.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Ipasshpubkey = &v
	}
	if len(data.CarLicense.Elements()) > 0 {
		var v []string
		for _, value := range data.CarLicense.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Carlicense = &v
	}
	if len(data.EmailAddress.Elements()) > 0 {
		var v []string
		for _, value := range data.EmailAddress.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Mail = &v
	}
	if !data.KrbPrincipalExpiration.IsNull() {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPrincipalExpiration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_principal_expiration timestamp could not be parsed as RFC3339: %s", err))
		}
		optArgs.Krbprincipalexpiration = &timestamp
	}
	if !data.KrbPasswordExpiration.IsNull() {
		timestamp, err := time.Parse(time.RFC3339, data.KrbPasswordExpiration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Attribute format", fmt.Sprintf("The krb_password_expiration timestamp could not be parsed as RFC3339: %s", err))
		}
		optArgs.Krbpasswordexpiration = &timestamp
	}
	if len(data.UserClass.Elements()) > 0 {
		var v []string
		for _, value := range data.UserClass.Elements() {
			val, _ := strconv.Unquote(value.String())
			v = append(v, val)
		}
		optArgs.Userclass = &v
	}

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.UserMod(&ipa.UserModArgs{}, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating freeipa user: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Update freeipa user %s returns %s", data.UID.String(), res.String()))

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel

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

	optArgs := ipa.UserDelOptionalArgs{}
	optArgs.UID = &[]string{data.UID.ValueString()}

	_, err := r.client.UserDel(&ipa.UserDelArgs{}, &optArgs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("User %s deletion failed: %s", data.Id.ValueString(), err))
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
