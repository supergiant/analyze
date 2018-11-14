package prices

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pricing"

	"github.com/supergiant/robot/builtin/plugins/underutilizednodes/models"
)

var awsPartitions = map[string]string{
	"ap-northeast-1": "Asia Pacific (Tokyo)",
	"ap-northeast-2": "Asia Pacific (Seoul)",
	"ap-south-1":     "Asia Pacific (Mumbai)",
	"ap-southeast-1": "Asia Pacific (Singapore)",
	"ap-southeast-2": "Asia Pacific (Sydney)",
	"ca-central-1":   "Canada (Central)",
	"eu-central-1":   "EU (Frankfurt)",
	"eu-west-1":      "EU (Ireland)",
	"eu-west-2":      "EU (London)",
	"eu-west-3":      "EU (Paris)",
	"sa-east-1":      "South America (Sao Paulo)",
	"us-east-1":      "US East (N. Virginia)",
	"us-east-2":      "US East (Ohio)",
	"us-west-1":      "US West (N. California)",
	"us-west-2":      "US West (Oregon)",
}

func Get(pricingService *pricing.Pricing, region string) map[string][]models.PriceItem {
	var computeInstancesPrices = make(map[string][]models.PriceItem, 0)

	productsInput := &pricing.GetProductsInput{
		Filters: []pricing.Filter{
			{
				Field: aws.String("ServiceCode"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("AmazonEC2"),
			},
			{
				Field: aws.String("productFamily"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("Compute Instance"),
			},
			{
				Field: aws.String("operatingSystem"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("Linux"),
			},
			{
				Field: aws.String("preInstalledSw"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String("NA"),
			},
			//TODO: FIRST PRIORITY FIX, to filter by usagetype "EC2: Running Hours"
			//https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/selectdim.html
			//{
			//	Field: aws.String("tenancy"),
			//	Type:  pricing.FilterTypeTermMatch,
			//	Value: aws.String("Shared"),
			//},
			{
				Field: aws.String("location"),
				Type:  pricing.FilterTypeTermMatch,
				Value: aws.String(awsPartitions[region]), //TODO region to location??? bug, add PR to lib?
			},
		},
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(100), //TODO: add pagination
		ServiceCode:   aws.String("AmazonEC2"),
	}

	productsRequest := pricingService.GetProductsRequest(productsInput)

	productsPager := productsRequest.Paginate()
	for productsPager.Next() {
		page := productsPager.CurrentPage()

		if page != nil {
			for _, productItem := range page.PriceList {
				//b, _ := json.Marshal(productItem)
				var newPriceItem = getProduct(productItem)
				//TODO: some prices even for usagetype HostBoxUsage equal to zero. need to fix it later
				if newPriceItem.InstanceType == "" || newPriceItem.Memory == "" || newPriceItem.Vcpu == "" || newPriceItem.ValuePerUnit == "0.0000000000" {
					//b, _ := json.Marshal(productItem)
					//fmt.Printf("%s\n", b)
				}
				_, exists := computeInstancesPrices[newPriceItem.InstanceType]
				if !exists {
					computeInstancesPrices[newPriceItem.InstanceType] = make([]models.PriceItem, 0, 0)
				}
				computeInstancesPrices[newPriceItem.InstanceType] = append(computeInstancesPrices[newPriceItem.InstanceType], newPriceItem)
			}
		}
	}

	if err := productsPager.Err(); err != nil {
		fmt.Printf("failed to describe products, %v", err)
	}

	fmt.Printf("found product prices: %v\n", len(computeInstancesPrices))
	return computeInstancesPrices
}

// TODO add checks and return error
func getProduct(productItem aws.JSONValue) models.PriceItem {
	var pi = models.PriceItem{}
	productInterface, exists := productItem["product"]
	if !exists {
		fmt.Printf("product elemnt doesn't exist")
		return pi
	}

	product, ok := productInterface.(map[string]interface{})
	if !ok {
		fmt.Printf("product elemnt is not map")
		return pi
	}

	attributes, exists := product["attributes"]
	if !exists {
		fmt.Printf("product elemnt doesn't exist")
		return pi
	}

	attrs, ok := attributes.(map[string]interface{})
	if !ok {
		fmt.Printf("attributes elemnt doesn't exist")
		return pi
	}

	value := attrs["instanceType"]
	pi.InstanceType, _ = value.(string)
	value = attrs["memory"]
	pi.Memory, _ = value.(string)
	value = attrs["vcpu"]
	pi.Vcpu, _ = value.(string)
	value = attrs["usagetype"]
	pi.UsageType, _ = value.(string)
	value = attrs["tenancy"]
	pi.Tenancy, _ = value.(string)

	termsInterface, exists := productItem["terms"]
	if !exists {
		fmt.Printf("terms elemnt doesn't exist")
		return pi
	}

	terms, ok := termsInterface.(map[string]interface{})
	if !ok {
		fmt.Printf("terms elemnt is not map")
		return pi
	}

	onDemandInterface, exists := terms["OnDemand"]
	if !exists {
		fmt.Printf("OnDemand elemnt doesn't exist")
		return pi
	}

	onDemand, ok := onDemandInterface.(map[string]interface{})
	if !ok {
		fmt.Printf("onDemand elemnt is not map")
		return pi
	}

	for _, skuValueInterface := range onDemand {
		skuValue, ok := skuValueInterface.(map[string]interface{})
		if !ok {
			fmt.Printf("skuValue elemnt is not map")
			return pi
		}

		priceDimensionsInterface, exists := skuValue["priceDimensions"]
		if !exists {
			fmt.Printf("priceDimensions elemnt doesn't exist")
			return pi
		}

		priceDimensions, ok := priceDimensionsInterface.(map[string]interface{})
		if !ok {
			fmt.Printf("priceDimensions elemnt is not map")
			return pi
		}

		for _, priceDimentionInterface := range priceDimensions {
			priceDimention, ok := priceDimentionInterface.(map[string]interface{})
			if !ok {
				fmt.Printf("priceDimention elemnt is not map")
				return pi
			}

			unitInterface, exists := priceDimention["unit"]
			if !exists {
				fmt.Printf("unit elemnt doesn't exist")
				return pi
			}

			pi.Unit, ok = unitInterface.(string)
			if !ok {
				fmt.Printf("unit elemnt is not string")
				return pi
			}

			pricePerUnitInterface, exists := priceDimention["pricePerUnit"]
			if !exists {
				fmt.Printf("pricePerUnit elemnt doesn't exist")
				return pi
			}

			pricePerUnit, ok := pricePerUnitInterface.(map[string]interface{})
			if !ok {
				fmt.Printf("pricePerUnit elemnt is not map")
				return pi
			}

			for k, v := range pricePerUnit {
				pi.Currency = k

				pi.ValuePerUnit, ok = v.(string)
				if !ok {
					fmt.Printf("valuePerUnit elemnt is not map")
					return pi
				}
				return pi
			}

		}
	}

	return pi
}
