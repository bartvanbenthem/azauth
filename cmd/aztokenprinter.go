package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bartvanbenthem/azauth/pkg/token"
)

// generic printer for access tokens
func AccessTokenPrinter(t token.Requester) {
	fmt.Println()
	tk, err := t.GetToken()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v\n", string(tk.AccessToken))
}

func main() {
	// get credentials from environment variables
	appid := os.Getenv("AZURE_CLIENT_ID")
	tenantid := os.Getenv("AZURE_TENANT_ID")
	secret := os.Getenv("AZURE_CLIENT_SECRET")

	credentials := token.Credential{
		ApplicationID: appid,
		TenantID:      tenantid,
		ClientSecret:  secret,
	}

	// get azure resource manager api token
	rmclient := token.RMClient{
		Auth: credentials,
	}
	AccessTokenPrinter(&rmclient)

	// get azure graph api token
	gclient := token.GraphClient{
		Auth: credentials,
	}
	AccessTokenPrinter(&gclient)

}
