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

	"github.com/paultyng/terraform-provider-expensify/internal/sdk"
)

const reportExportTemplate = `
${reports?size}<#lt>
<#list reports as report>
	${report.transactionList?size},<#t>
	${report.accountEmail},<#t>
	${report.reportName},<#t>
	${report.policyID}<#lt>
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
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
						"amount_cents": {
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
	c := meta.(*sdk.Client)

	id := d.Id()

	file, err := c.File(ctx, sdk.FileRequest{
		OnReceive: sdk.OnReceive{
			ImmediateResponse: []string{"returnRandomFileName"},
		},
		InputSettings: sdk.InputSettings{
			Type: "combinedReportData",
			Filters: sdk.InputSettingsFilters{
				ReportIDList: id,
			},
		},
		OutputSettings: sdk.OutputSettings{
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
	r.FieldsPerRecord = 1

	// read report size (should be 1)
	record, err := r.Read()
	if err != nil {
		panic(err)
	}
	numReports, err := strconv.Atoi(record[0])
	if err != nil {
		panic(err)
	}
	if numReports > 1 {
		panic("expected 1 report, got %d")
	}
	if numReports == 0 {
		// not found
		d.SetId("")
		return nil
	}

	r.FieldsPerRecord = 4
	// read report header (1 row)
	record, err = r.Read()
	if err != nil {
		panic(err)
	}
	numExpenses, err := strconv.Atoi(record[0])
	if err != nil {
		panic(err)
	}
	email := record[1]
	title := record[2]
	policyID := record[3]

	d.Set("email", email)
	d.Set("title", title)
	d.Set("policy_id", policyID)

	for {
		r.FieldsPerRecord = 4
		record, err = r.Read()
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
			"merchant":     record[0],
			"amount_cents": amount,
			"currency":     record[2],
			"date":         record[3],
		})
	}

	if numExpenses != len(expenses) {
		panic("mismatch in expenses length")
	}

	err = d.Set("expense", expenses)
	if err != nil {
		panic(err)
	}

	return nil
}
