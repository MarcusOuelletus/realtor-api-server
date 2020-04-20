package main

import (
	"github.com/MarcusOuelletus/rets-server/database"
	"github.com/MarcusOuelletus/rets-server/helpers"
)

func main() {
	generatedToken := helpers.GenerateToken(40)

	err := database.Insert(&database.InsertQuery{
		Collection: "accounts",
		Data: map[string]interface{}{
			"Token":     generatedToken,
			"Name":      "BROKER_NAME",
			"Brokerage": "BROKERAGE",
			"Email":     "BROKER_EMAIL",
			"Tables":    []string{"res", "con", "com"},
		},
	})

	if err != nil {
		panic(err.Error())
	}
}
