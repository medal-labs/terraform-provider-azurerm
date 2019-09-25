package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHealthcareService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHealthcareServiceCreateUpdate,
		Read:   resourceArmHealthcareServiceRead,
		Update: resourceArmHealthcareServiceCreateUpdate,
		Delete: resourceArmHealthcareServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fhir",
			},

			"cosmosdb_throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},

			"access_policy_object_ids": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.UUID,
				},
			},

			"authentication_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authority": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"audience": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smart_proxy_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"cors_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allowed_headers": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								Elem: &schema.Schema{
									Type: schema.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"DELETE",
										"GET",
										"HEAD",
										"MERGE",
										"POST",
										"OPTIONS",
										"PUT"}, false),
								},
							},
						},
						"max_age_in_seconds": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 2000000000),
						},
						"allow_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmHealthcareServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Healthcare Service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	kind := d.Get("kind").(string)
	cdba := int32(d.Get("cosmosdb_throughput").(int))

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Healthcare Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_healthcare_service", *existing.ID)
		}
	}

	healthcareServiceDescription := healthcareapis.ServicesDescription{
		Location: utils.String(location),
		Tags:     expandedTags,
		Kind:     healthcareapis.Kind(kind),
		Properties: &healthcareapis.ServicesProperties{
			AccessPolicies: expandAzureRMhealthcareapisAccessPolicyEntries(d),
			CosmosDbConfiguration: &healthcareapis.ServiceCosmosDbConfigurationInfo{
				OfferThroughput: &cdba,
			},
			CorsConfiguration:           expandAzureRMhealthcareapisCorsConfiguration(d),
			AuthenticationConfiguration: expandAzureRMhealthcareapisAuthentication(d),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, healthcareServiceDescription)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Retrieving Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Healthcare Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHealthcareServiceRead(d, meta)
}

func resourceArmHealthcareServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["services"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Healthcare Service %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Healthcare Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if kind := resp.Kind; string(kind) != "" {
		d.Set("kind", kind)
	}
	if properties := resp.Properties; properties != nil {
		if config := properties.AccessPolicies; config != nil {
			d.Set("access_policy_object_ids", flattenHealthcareAccessPolicies(config))
		}
		if config := properties.CosmosDbConfiguration; config != nil {
			d.Set("cosmosdb_throughput", config.OfferThroughput)
		}

		authOutput := make([]interface{}, 0)
		if authConfig := properties.AuthenticationConfiguration; authConfig != nil {
			output := make(map[string]interface{})
			if authConfig.Authority != nil {
				output["authority"] = *authConfig.Authority
			}
			if authConfig.Audience != nil {
				output["audience"] = *authConfig.Audience
			}
			if authConfig.SmartProxyEnabled != nil {
				output["smart_proxy_enabled"] = *authConfig.SmartProxyEnabled
			}
			authOutput = append(authOutput, output)
		}

		if err := d.Set("authentication_configuration", authOutput); err != nil {
			return fmt.Errorf("Error setting `authentication_configuration`: %+v", authOutput)
		}

		corsOutput := make([]interface{}, 0)
		if corsConfig := properties.CorsConfiguration; corsConfig != nil {
			output := make(map[string]interface{})
			if corsConfig.Origins != nil {
				output["allowed_origins"] = *corsConfig.Origins
			}
			if corsConfig.Headers != nil {
				output["allowed_headers"] = *corsConfig.Headers
			}
			if corsConfig.Methods != nil {
				output["allowed_methods"] = *corsConfig.Methods
			}
			if corsConfig.MaxAge != nil {
				output["max_age_in_seconds"] = *corsConfig.MaxAge
			}
			if corsConfig.AllowCredentials != nil {
				output["allow_credentials"] = *corsConfig.AllowCredentials
			}
			corsOutput = append(corsOutput, output)
		}

		if err := d.Set("cors_configuration", corsOutput); err != nil {
			return fmt.Errorf("Error setting `cors_configuration`: %+v", corsOutput)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmHealthcareServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID: %+v", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["services"]
	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deleting Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandAzureRMhealthcareapisAccessPolicyEntries(d *schema.ResourceData) *[]healthcareapis.ServiceAccessPolicyEntry {
	accessPolicyObjectIds := d.Get("access_policy_object_ids").([]interface{})
	svcAccessPolicyArray := make([]healthcareapis.ServiceAccessPolicyEntry, 0)

	for _, objectId := range accessPolicyObjectIds {
		objectIdsStr := objectId.(string)
		svcAccessPolicyObjectId := healthcareapis.ServiceAccessPolicyEntry{ObjectID: &objectIdsStr}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}

	return &svcAccessPolicyArray
}

func expandAzureRMhealthcareapisCorsConfiguration(d *schema.ResourceData) *healthcareapis.ServiceCorsConfigurationInfo {
	corsConfigRaw := d.Get("cors_configuration").([]interface{})

	if len(corsConfigRaw) == 0 {
		return &healthcareapis.ServiceCorsConfigurationInfo{}
	}

	allowedOrigins := make([]string, 0)
	allowedHeaders := make([]string, 0)
	allowedMethods := make([]string, 0)
	maxAgeInSeconds := int32(0)
	allowCredentials := true

	for _, attr := range corsConfigRaw {
		corsConfigAttr := attr.(map[string]interface{})

		allowedOrigins = *utils.ExpandStringSlice(corsConfigAttr["allowed_origins"].([]interface{}))
		allowedHeaders = *utils.ExpandStringSlice(corsConfigAttr["allowed_headers"].([]interface{}))
		allowedMethods = *utils.ExpandStringSlice(corsConfigAttr["allowed_methods"].([]interface{}))
		maxAgeInSeconds = int32(corsConfigAttr["max_age_in_seconds"].(int))
		allowCredentials = corsConfigAttr["allow_credentials"].(bool)
	}

	cors := &healthcareapis.ServiceCorsConfigurationInfo{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		MaxAge:           &maxAgeInSeconds,
		AllowCredentials: &allowCredentials,
	}
	return cors
}

func expandAzureRMhealthcareapisAuthentication(d *schema.ResourceData) *healthcareapis.ServiceAuthenticationConfigurationInfo {
	authConfigRaw := d.Get("authentication_configuration").([]interface{})

	if len(authConfigRaw) == 0 {
		return &healthcareapis.ServiceAuthenticationConfigurationInfo{}
	}

	authority := ""
	audience := ""
	smart_proxy_enabled := true

	for _, attr := range authConfigRaw {
		authConfigAttr := attr.(map[string]interface{})

		authority = authConfigAttr["authority"].(string)
		audience = authConfigAttr["audience"].(string)
		smart_proxy_enabled = authConfigAttr["smart_proxy_enabled"].(bool)
	}

	auth := &healthcareapis.ServiceAuthenticationConfigurationInfo{
		Authority:         &authority,
		Audience:          &audience,
		SmartProxyEnabled: &smart_proxy_enabled,
	}
	return auth
}