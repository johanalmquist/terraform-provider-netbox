package netbox

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/virtualization"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxVirtualizationVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVirtualizationVMCreate,
		Read:   resourceNetboxVirtualizationVMRead,
		Update: resourceNetboxVirtualizationVMUpdate,
		Delete: resourceNetboxVirtualizationVMDelete,
		Exists: resourceNetboxVirtualizationVMExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"text", "integer", "boolean",
								"date", "url", "selection", "multiple"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"disk": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"local_context_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"platform_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"offline", "active",
					"planned", "staged", "failed", "decommissioning"}, false),
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[0-9]+((.[0-9]){0,1}[0-9]{0,1})$"),
					"Must be like ^[0-9]+((.[0-9]){0,1}[0-9]{0,1})$"),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == new+".00" || old == new {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceNetboxVirtualizationVMCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	clusterID := int64(d.Get("cluster_id").(int))
	comments := d.Get("comments").(string)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	disk := int64(d.Get("disk").(int))
	localContextData := d.Get("local_context_data").(string)
	memory := int64(d.Get("memory").(int))
	name := d.Get("name").(string)
	platformID := int64(d.Get("platform_id").(int))
	roleID := int64(d.Get("role_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	vcpus := d.Get("vcpus").(string)

	if !strings.Contains(vcpus, ".") {
		vcpus = vcpus + ".00"
	}

	newResource := &models.WritableVirtualMachineWithConfigContext{
		Cluster:      &clusterID,
		Comments:     comments,
		CustomFields: &customFields,
		Name:         &name,
		Status:       status,
		Tags:         convertTagsToNestedTags(tags),
	}

	if disk != 0 {
		newResource.Disk = &disk
	}

	if memory != 0 {
		newResource.Memory = &memory
	}

	if localContextData != "" {
		var localContextDataMap map[string]*interface{}
		if err := json.Unmarshal([]byte(localContextData), &localContextDataMap); err != nil {
			return err
		}
		newResource.LocalContextData = localContextDataMap
	}

	if platformID != 0 {
		newResource.Platform = &platformID
	}

	if roleID != 0 {
		newResource.Role = &roleID
	}

	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	if vcpus != "" {
		vcpusFloat, _ := strconv.ParseFloat(vcpus, 32)
		newResource.Vcpus = &vcpusFloat
	}

	resource := virtualization.NewVirtualizationVirtualMachinesCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationVirtualMachinesCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationVMRead(d, m)
}

func resourceNetboxVirtualizationVMRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)

	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("cluster_id", resource.Cluster.ID); err != nil {
				return err
			}

			var comments interface{}
			if resource.Comments == "" {
				comments = nil
			} else {
				comments = resource.Comments
			}

			if err = d.Set("comments", comments); err != nil {
				return err
			}

			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return err
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return err
			}

			if err = d.Set("disk", resource.Disk); err != nil {
				return err
			}

			if resource.LocalContextData != nil {
				localContextDataJSON, err := json.Marshal(resource.LocalContextData)
				if err != nil {
					return err
				}
				if err = d.Set("local_context_data",
					string(localContextDataJSON)); err != nil {
					return err
				}
			} else {
				d.Set("local_context_data", "")
			}

			if err = d.Set("memory", resource.Memory); err != nil {
				return err
			}

			if err = d.Set("name", resource.Name); err != nil {
				return err
			}

			if resource.Platform == nil {
				if err = d.Set("platform_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("platform_id", resource.Platform.ID); err != nil {
					return err
				}
			}

			if resource.Role == nil {
				if err = d.Set("role_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("role_id", resource.Role.ID); err != nil {
					return err
				}
			}

			if err = d.Set("status", resource.Status.Value); err != nil {
				return err
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return err
			}

			if resource.Tenant == nil {
				if err = d.Set("tenant_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("tenant_id", resource.Tenant.ID); err != nil {
					return err
				}
			}

			if err = d.Set("vcpus", fmt.Sprintf("%v", *resource.Vcpus)); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxVirtualizationVMUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableVirtualMachineWithConfigContext{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name
	clusterID := int64(d.Get("cluster_id").(int))
	params.Cluster = &clusterID

	if d.HasChange("comments") {
		if comments, exist := d.GetOk("comments"); exist {
			params.Comments = comments.(string)
		} else {
			params.Comments = " "
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("disk") {
		disk := int64(d.Get("disk").(int))
		params.Disk = &disk
	}

	if d.HasChange("local_context_data") {
		localContextData := d.Get("local_context_data").(string)
		var localContextDataMap map[string]*interface{}
		if localContextData == "" {
			localContextDataMap = nil
		} else if err := json.Unmarshal([]byte(localContextData), &localContextDataMap); err != nil {
			return err
		}
		params.LocalContextData = localContextDataMap
	}

	if d.HasChange("memory") {
		memory := int64(d.Get("memory").(int))
		params.Memory = &memory
	}

	if d.HasChange("platform_id") {
		platformID := int64(d.Get("platform_id").(int))
		params.Platform = &platformID
	}

	if d.HasChange("role_id") {
		roleID := int64(d.Get("role_id").(int))
		params.Role = &roleID
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		params.Status = status
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		params.Tenant = &tenantID
	}

	if d.HasChange("vcpus") {
		vcpus := d.Get("vcpus").(string)

		if !strings.Contains(vcpus, ".") {
			vcpus = vcpus + ".00"
		}

		vcpusFloat, _ := strconv.ParseFloat(vcpus, 32)
		params.Vcpus = &vcpusFloat
	}

	resource := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxVirtualizationVMRead(d, m)
}

func resourceNetboxVirtualizationVMDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationVMExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	p := virtualization.NewVirtualizationVirtualMachinesDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationVirtualMachinesDelete(
		p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxVirtualizationVMExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationVirtualMachinesList(
		params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
