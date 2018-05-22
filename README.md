# go-lucener

go-lucener is a Go library for building [Lucene compatable expressions](https://github.com/Stratio/cassandra-lucene-index) for [Cassandra driver](https://github.com/gocql/gocql)

**Documentation:** [![GoDoc](https://godoc.org/github.com/engin/go-lucener/travis?status.svg)](https://godoc.org/github.com/engin/go-lucener)

**Build Status:** [![Build Status](https://travis-ci.org/engin/go-lucener.svg?branch=master)](https://travis-ci.org/engin/go-lucener)

go-lucener requires Go version 1.7 or greater.

## Example

You can download the example user table and index from [Stratio/cassandra-lucene-index repo](https://github.com/Stratio/cassandra-lucene-index/blob/branch-3.0.10/doc/resources/test-users-create.cql)

```go
package main

import (
	"fmt"
	"log"

	"github.com/engin/go-lucener"
	"github.com/gocql/gocql"
)

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "test"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("couldn't connect to cassandra: %v", err)
	}
	defer session.Close()

	var name, gender, animal, food string
	var age int

	exp := lucener.NewExpr()
	exp.Filter(
		lucener.BooleanMust(
			lucener.Wildcard("name", "Ali*"), lucener.Wildcard("food", "tu*"),
		),
	)
	q := `SELECT name,gender,animal,age,food FROM users WHERE expr(users_index, ?)`

	// generated cql will be like:
	// SELECT * FROM users WHERE expr(users_index, '{
	// 	filter: {
	// 	   type: "boolean",
	// 	   must: [
	// 		  {type: "wildcard", field: "name", value: "*a"},
	// 		  {type: "wildcard", field: "food", value: "tu*"}
	// 	   ]
	// 	}
	//  }');
	fmt.Sprintln(exp.String())
	iter := session.Query(q, exp).Iter()
	for iter.Scan(&name, &gender, &animal, &age, &food) {
		fmt.Println("User:", name, gender, animal, age, food)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

```

### TODO
I'll try to add more test and examples
