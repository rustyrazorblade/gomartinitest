package main

import "fmt"
import "github.com/codegangsta/martini"
import (
	"github.com/gocql/gocql"
	"github.com/codegangsta/martini-contrib/render"
)

func main() {
	fmt.Printf("Hello world!")
	db := DB()

	m := martini.Classic()
	m.Use(db)
	m.Use(render.Renderer())

	m.Get("/", Root)
	m.Get("/json", Json)

	m.Run()
}

func DB() martini.Handler {
	cluster := gocql.NewCluster("localhost")
	cluster.CQLVersion = "3.0.0"
	cluster.ProtoVersion = 1
	cluster.Keyspace = "users"
	cluster.Consistency = gocql.One

	if session, err := cluster.CreateSession(); err == nil{
		session.Query("select user_id, first_name, last_name from user limit 1")
		return func(c martini.Context) {
			c.Map(session)
		}
	} else {
		fmt.Println(err)
		panic("FREAK OUT")
	}
}

type result struct {
	FirstName string `json:"first_name"`
	LastName string  `json:"last_name"`
}

func Root(session *gocql.Session, r render.Render) {
	var first_name, last_name string
	var id gocql.UUID

	iter := session.Query("select user_id, first_name, last_name from user limit 1000").Iter()

	response := make([]result, 0, 1000)

	for iter.Scan(&id, &first_name, &last_name) {
		response = append(response, result{FirstName:first_name, LastName:last_name})
	}
	r.JSON(200, response)

}

func Json(session *gocql.Session, r render.Render) {
	r.JSON(200, map[string]interface {}{"shut_up": "captain"})
}
