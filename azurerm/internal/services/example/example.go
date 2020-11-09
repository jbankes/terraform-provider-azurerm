package example

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
)

type ExampleResource struct {
}

func (r ExampleResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"number": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"networks": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"networks_set": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"float": {
			Type: schema.TypeFloat,
			Optional: true,
		},
		"list": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"set": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

// Computed Only
func (r ExampleResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"output": {
			Type: schema.TypeString,
		},
	}
}

func (r ExampleResource) ResourceType() string {
	return "azurerm_example"
}

// NOTE: i guess we could return schema object to ensure everything is mapped and valid idk
type ExampleObj struct {
	Name     string   `hcl:"name"`
	Number   int      `hcl:"number"`
	Output   string   `hcl:"output" computed:"true" ip_address:"true"`
	Enabled  bool     `hcl:"enabled"`
	Networks []string `hcl:"networks"`
	NetworksSet []string `hcl:"networks_set"`
	List []NetworkList `hcl:"list"`
	Set []NetworkList `hcl:"set"`
	Float float64 `hcl:"float"`
}

type NetworkList struct {
	Name string `hcl:"name"`
}

func (r ExampleResource) Create() ResourceFunc {
	return CreateUpdate()
}

func (r ExampleResource) Read() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			return metadata.Encode(&ExampleObj{
				Name:    "updated",
				Number:  123,
				Enabled: true,
				Networks: []string{"123", "124"},
				NetworksSet: []string{"asdf", "qwer"},
				List: []NetworkList{{
					Name: "test1232",
				}},
				Set: []NetworkList{{
					Name: "set1232",
				}},
				Float: float64(123),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

// copy pasta create
func (r ExampleResource) Update() ResourceFunc {
	return CreateUpdate()
}

func (r ExampleResource) Delete() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ExampleResource) IDValidationFunc() schema.SchemaValidateFunc {
	return validate.SubnetID
}

func CreateUpdate() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			//metadata.ResourceData
			//metadata.Logger.WarnF("OHHAI %d", 3)
			//metadata.Client.Account.SubscriptionId
			metadata.Logger.Info("HEYO")

			var obj ExampleObj
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			id := parse.SubnetId{
				ResourceGroup:      "production-resources",
				VirtualNetworkName: "production-network",
				Name:               obj.Name,
			}

			metadata.Logger.InfoF("Name is %s", obj.Name)
			metadata.Logger.InfoF("Number is %d", obj.Number)
			metadata.Logger.InfoF("Float is %d", obj.Float)
			metadata.Logger.InfoF("Networks are %+v", obj.Networks)
			metadata.Logger.InfoF("Networks Set is %+v", obj.NetworksSet)
			metadata.Logger.InfoF("List  is %+v", obj.List)
			metadata.Logger.InfoF("Set  is %+v", obj.Set)


			metadata.SetID(id)
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}
