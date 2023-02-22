module aad-auth-proxy

go 1.19

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.1.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.1.0
	github.com/sirupsen/logrus v1.8.1
	software.sslmate.com/src/go-pkcs12 v0.0.0-20210415151418-c5206de65a78
)

require (
	github.com/google/uuid v1.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.0.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.5.1 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4 // indirect
	github.com/t-tomalak/logrus-easy-formatter v0.0.0-20190827215021-c074f06c5816
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

replace golang.org/x/text => golang.org/x/text v0.3.8
