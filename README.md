# go-ctechpay

A work-in-progress client library for [CTechPay APIs](http://docs.ctechpay.com/) in Go.

## Requirements

- Go 1.21+

## Usage

Add to your project

```
go get github.com/golang-malawi/go-ctechpay
```

Import in your code, the following is a rough example, more detailed examples will be provided later:

```go
package something


import (
    "github.com/golang-malawi/go-ctechpay"
)


func paymentHandler() {
    
    client := ctechpay.NewSandboxClient("token from ctechpay", 30*time.Duration)

    client.SetRedirectURL("http://localhost:8080/path/to/your/billing-complete-page")

    transactionID := ulid.Make().String()
    orderResponse, err := client.InitiateCardOrder(transactionID, *big.NewFloat(30.0), true)
    
    // handle error in case the call to CTechPay fails 

    // then render a redirect page
    // renderHtml("ctech_redirect.html", map[string]any {
	//	"transactionID":          transactionID,
	//	"internalTransactionRef": "TEST",
	//	"redirectURL":            orderResponse.PaymentPageURL,
	// })
}
```
