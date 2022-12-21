module postcard

go 1.18

require (
	github.com/ThreeDotsLabs/esja v0.0.0-20221208191400-8fbb493947e7
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.6
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/ThreeDotsLabs/pii v0.0.0-20221221144555-f2186024b30d // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ThreeDotsLabs/esja => ../../
