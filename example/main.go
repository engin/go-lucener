package main

import (
	"fmt"
	"log"

	"github.com/engin/go-lucener"
	"github.com/gocql/gocql"
)

func main() {
	// connect to the cassandra cluster
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
