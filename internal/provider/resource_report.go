package provider

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const reportExportTemplate = `<#list reports as report>
    <#list report.transactionList as expense>
        <#if expense.modifiedMerchant?has_content>
            <#assign merchant = expense.modifiedMerchant>
        <#else>
            <#assign merchant = expense.merchant>
        </#if>
        <#if expense.modifiedAmount?has_content>
            <#assign amount = expense.modifiedAmount>
        <#else>
            <#assign amount = expense.amount>
        </#if>
        <#if expense.modifiedCreated?has_content>
            <#assign created = expense.modifiedCreated>
        <#else>
            <#assign created = expense.created>
        </#if>
        ${merchant},<#t>
		${amount},<#t>
		${expense.currency},<#t>
        ${created}<#lt>
    </#list>
</#list>`

func resourceReport() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceReportRead,
		Delete: func(*schema.ResourceData, interface{}) error {
			return fmt.Errorf("delete of expense reports is not supported in the API")
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"expense": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"merchant": {
							Type:     schema.TypeString,
							Required: true,
						},
						"date": {
							Type:     schema.TypeString,
							Required: true,
						},
						"amount": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"currency": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "USD",
						},
					},
				},
			},
		},
	}
}

func resourceReportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	file, err := c.File(ctx, FileRequest{
		OnReceive: OnReceive{
			ImmediateResponse: []string{"returnRandomFileName"},
		},
		InputSettings: InputSettings{
			Type: "combinedReportData",
			Filters: InputSettingsFilters{
				ReportIDList: id,
			},
		},
		OutputSettings: OutputSettings{
			FileExtension: "txt",
		},
	}, reportExportTemplate)
	if err != nil {
		//return fmt.Errorf("")
		panic(err)
	}

	textData, err := c.Download(ctx, file, "")
	if err != nil {
		panic(err)
	}

	expenses := []interface{}{}
	r := csv.NewReader(strings.NewReader(textData))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		amount, err := strconv.Atoi(record[1])
		if err != nil {
			panic(err)
		}

		expenses = append(expenses, map[string]interface{}{
			"merchant": record[0],
			"amount":   amount,
			"currency": record[2],
			"date":     record[3],
		})
	}

	err = d.Set("expense", expenses)
	if err != nil {
		panic(err)
	}

	return nil
}
