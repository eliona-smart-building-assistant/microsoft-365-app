module microsoft-365

go 1.24

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.17.1
	github.com/eliona-smart-building-assistant/app-integration-tests v1.1.2
	github.com/eliona-smart-building-assistant/go-eliona v1.10.5
	github.com/eliona-smart-building-assistant/go-utils v1.1.4
	github.com/gorilla/mux v1.8.1
	github.com/microsoft/kiota-abstractions-go v1.9.1
	github.com/microsoft/kiota-authentication-azure-go v1.2.1
	github.com/microsoftgraph/msgraph-sdk-go v1.67.0
	github.com/microsoftgraph/msgraph-sdk-go-core v1.3.1
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/sqlboiler/v4 v4.18.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.4.2 // indirect
	github.com/cjlapao/common-go v0.0.48 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ericlagergren/decimal v0.0.0-20240411145413-00de7ca16731 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/microsoft/kiota-http-go v1.5.1 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.1.1 // indirect
	github.com/microsoft/kiota-serialization-json-go v1.1.1 // indirect
	github.com/microsoft/kiota-serialization-multipart-go v1.1.1 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.1.1 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.3 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/volatiletech/randomize v0.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

// Bugfix see: https://github.com/volatiletech/sqlboiler/blob/91c4f335dd886d95b03857aceaf17507c46f9ec5/README.md
// decimal library showing errors like: pq: encode: unknown type types.NullDecimal is a result of a too-new and broken version of the github.com/ericlargergren/decimal package, use the following version in your go.mod: github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
replace github.com/ericlagergren/decimal => github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5

require (
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.2
	github.com/eliona-smart-building-assistant/go-eliona-api-client/v2 v2.8.2
	github.com/friendsofgo/errors v0.9.2
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.8
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
)
