module github.com/proge-software/uninapoli-csml-csbot

go 1.14

replace github.com/proge-software/uninapoli-csml-csbot => ./

require (
	github.com/Azure/azure-sdk-for-go v42.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/satori/go.uuid v1.2.0 // indirect
	gopkg.in/tucnak/telebot.v2 v2.0.0-20200426184946-59629fe0483e
)
