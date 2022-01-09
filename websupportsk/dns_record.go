package websupportsk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
)

type DnsClient struct {
	client *Client
}

type DnsZone struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	UpdateTime int    `json:"updateTime"`
}

type DnsRecord struct {
	Id      int     `json:"id"`
	Type    string  `json:"type"`
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Prio    int     `json:"prio"`
	Port    int     `json:"port"`
	Weight  int     `json:"weight"`
	Ttl     int     `json:"ttl"`
	Note    string  `json:"note"`
	Zone    DnsZone `json:"zone"`
}

type DnsRecordStatusWrapper struct {
	Status string    `json:"status"`
	Item   DnsRecord `json:"item"`
}

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDnsRecordImport,
		},
		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "MX", "ANAME", "CNAME", "NS", "TXT", "SRV"}, true),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"prio": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDnsRecordCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	body := GetDnsRecord(data)

	record, err := client.Dns.CreateDnsRecord(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(record.Id))

	return resourceDnsRecordRead(ctx, data, meta)
}

func resourceDnsRecordRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	record, err := client.Dns.ReadDnsRecord(ctx, data.Get("zone_name").(string), data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	record.SetResourceData(data)

	return nil
}

func resourceDnsRecordUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	body := GetDnsRecord(data)

	record, err := client.Dns.UpdateDnsRecord(ctx, data.Id(), body)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(record.Id))

	return resourceDnsRecordRead(ctx, data, meta)
}

func resourceDnsRecordDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	if err := client.Dns.DeleteDnsRecord(ctx, data.Get("zone_name").(string), data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDnsRecordImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Client)

	// split the id so we can lookup
	idAttr := strings.SplitN(d.Id(), "/", 2)
	var zoneName string
	var recordId string
	if len(idAttr) == 2 {
		zoneName = idAttr[0]
		recordId = idAttr[1]
	} else {
		return nil, fmt.Errorf("invalid id %q specified, should be in format \"zoneName/recordId\" for import", d.Id())
	}

	record, err := client.Dns.ReadDnsRecord(ctx, zoneName, recordId)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Found record: %s", record.Name)

	d.Set("zone_name", zoneName)
	d.SetId(recordId)

	resourceDnsRecordRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func (client *DnsClient) CreateDnsRecord(ctx context.Context, body DnsRecord) (DnsRecord, error) {
	bodyData, errJson := json.Marshal(body)
	if errJson != nil {
		return DnsRecord{}, errJson
	}

	request, errReq := client.client.NewRequest(ctx, "POST", fmt.Sprintf("/v1/user/self/zone/%s/record", body.Zone.Name), bytes.NewReader(bodyData))
	if errReq != nil {
		return DnsRecord{}, errReq
	}

	var response DnsRecordStatusWrapper
	if _, errDo := client.client.Do(request, &response); errDo != nil {
		return DnsRecord{}, errDo
	}

	return response.Item, nil
}

func (client *DnsClient) ReadDnsRecord(ctx context.Context, zone string, id string) (DnsRecord, error) {
	request, errReq := client.client.NewRequest(ctx, "GET", fmt.Sprintf("/v1/user/self/zone/%s/record/%s", zone, id), nil)
	if errReq != nil {
		return DnsRecord{}, errReq
	}

	var response DnsRecord
	if _, errDo := client.client.Do(request, &response); errDo != nil {
		return DnsRecord{}, errDo
	}

	return response, nil
}

func (client *DnsClient) UpdateDnsRecord(ctx context.Context, id string, body DnsRecord) (DnsRecord, error) {
	bodyData, errJson := json.Marshal(body)
	if errJson != nil {
		return DnsRecord{}, errJson
	}

	request, errReq := client.client.NewRequest(ctx, "PUT", fmt.Sprintf("/v1/user/self/zone/%s/record/%s", body.Zone.Name, id), bytes.NewReader(bodyData))
	if errReq != nil {
		return DnsRecord{}, errReq
	}

	var response DnsRecordStatusWrapper
	if _, errDo := client.client.Do(request, &response); errDo != nil {
		return DnsRecord{}, errDo
	}

	return response.Item, nil
}

func (client *DnsClient) DeleteDnsRecord(ctx context.Context, zone string, id string) error {
	request, errReq := client.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/v1/user/self/zone/%s/record/%s", zone, id), nil)
	if errReq != nil {
		return errReq
	}

	var response DnsRecordStatusWrapper
	if _, errDo := client.client.Do(request, &response); errDo != nil {
		return errDo
	}

	return nil
}

func (dnsRecord *DnsRecord) SetResourceData(data *schema.ResourceData) {
	data.Set("id", dnsRecord.Id)
	data.Set("type", dnsRecord.Type)
	data.Set("name", dnsRecord.Name)
	data.Set("content", dnsRecord.Content)
	data.Set("ttl", dnsRecord.Ttl)
	data.Set("note", dnsRecord.Note)

	switch dnsRecord.Type {
	case "MX":
		data.Set("prio", dnsRecord.Prio)
	case "SRV":
		data.Set("prio", dnsRecord.Prio)
		data.Set("port", dnsRecord.Port)
		data.Set("weight", dnsRecord.Weight)
	}
}

func GetDnsRecord(data *schema.ResourceData) DnsRecord {
	id, _ := strconv.Atoi(data.Id())
	body := DnsRecord{
		Id:      id,
		Type:    data.Get("type").(string),
		Name:    data.Get("name").(string),
		Content: data.Get("content").(string),
		Ttl:     data.Get("ttl").(int),
		Note:    data.Get("note").(string),

		Zone: DnsZone{
			Name: data.Get("zone_name").(string),
		},
	}

	switch body.Type {
	case "SRV":
		body.Prio = data.Get("prio").(int)
		body.Port = data.Get("port").(int)
		body.Weight = data.Get("weight").(int)
	case "MX":
		body.Prio = data.Get("prio").(int)
	}

	return body
}
