module github.com/deckarep/golang-set/x/marshal

go 1.18

require (
	github.com/deckarep/golang-set/v3 v3.0.0
	go.mongodb.org/mongo-driver/v2 v2.6.0
)

// Local development and CI tests against the current version of the core module.
replace github.com/deckarep/golang-set/v3 => ../..
