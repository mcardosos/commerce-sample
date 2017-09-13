package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/utils"
	"github.com/mcardosos/azure-sdk-for-go/arm/commerce"
)

// This example requires that the following environment vars are set:
//
// AZURE_TENANT_ID: contains your Azure Active Directory tenant ID or domain
// AZURE_CLIENT_ID: contains your Azure Active Directory Application Client ID
// AZURE_CLIENT_SECRET: contains your Azure Active Directory Application Secret
// AZURE_SUBSCRIPTION_ID: contains your Azure Subscription ID
//

var (
	rateCardClient commerce.RateCardClient
)

func init() {
	authorizer, err := utils.GetAuthorizer(azure.PublicCloud)
	onErrorFail(err, "GetAuthorizer failed")

	subscriptionID := utils.GetEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	createClients(subscriptionID, authorizer)
}

func main() {
	fmt.Println("Get rate card...")
	rc, err := rateCardClient.Get("OfferDurableId eq 'MS-AZR-0062P' and Currency eq 'USD' and Locale eq 'en-US' and RegionInfo eq 'US'")
	if err != nil {
		b, errInner := ioutil.ReadAll(rc.Body)
		if errInner != nil {
			panic(errInner)
		}
		fmt.Println(string(b))
		onErrorFail(err, "Get failed")
	}
	printRC(rc)
}

func printRC(rc commerce.ResourceRateCardInfo) {
	fmt.Println("Currency: ", *rc.Currency)
	fmt.Println("Locale: ", *rc.Locale)
	fmt.Println("IsTaxIncluded: ", *rc.IsTaxIncluded)
	fmt.Println("OfferTerms:")
	printOT(*rc.OfferTerms)
	fmt.Println("Meters: ")
	printMI(*rc.Meters)
}

func printOT(oti []commerce.OfferTermInfo) {
	for _, ot := range oti[:getLen(oti)] {
		switch v := ot.(type) {
		case commerce.MonetaryCredit:
			fmt.Println("\tName: ", v.Name)
			fmt.Println("\tDate: ", *v.EffectiveDate)
			fmt.Println("\tCredit: ", *v.Credit)
			fmt.Println("\tMeterIDs: ", (*v.ExcludedMeterIds)[:getLen(*v.ExcludedMeterIds)])
		case commerce.MonetaryCommitment:
			fmt.Println("\tName: ", v.Name)
			fmt.Println("\tDate: ", *v.EffectiveDate)
			fmt.Println("\tDiscount: ", *v.TieredDiscount)
			fmt.Println("\tMeterIDs: ", (*v.ExcludedMeterIds)[:getLen(*v.ExcludedMeterIds)])
		case commerce.RecurringCharge:
			fmt.Println("\tName: ", v.Name)
			fmt.Println("\tDate: ", *v.EffectiveDate)
			fmt.Println("\tCharge: ", *v.RecurringCharge)
		default:
			fmt.Println("Not supported")
		}
		fmt.Println("=====")
	}
}

func printMI(mi []commerce.MeterInfo) {
	for _, m := range mi[:getLen(mi)] {
		fmt.Println("\tName: ", *m.MeterName)
		fmt.Println("\tMeterID: ", *m.MeterID)
		fmt.Println("=====")
	}
}

func getLen(slice interface{}) int {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return 0
	}
	index := s.Len()
	if index > 3 {
		index = 3
	}
	return index
}

func createClients(subscriptionID string, authorizer *autorest.BearerAuthorizer) {
	rateCardClient = commerce.NewRateCardClient(subscriptionID)
	rateCardClient.Authorizer = authorizer
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string, a ...interface{}) {
	if err != nil {
		fmt.Printf("%s: %s\n", fmt.Sprintf(message, a), err)
		os.Exit(1)
	}
}
