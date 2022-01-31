package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	createDatabase()

	db, err := sql.Open("sqlite3", "./todo.db")
	checkErr(err)
	createStmt, err := db.Prepare(`CREATE TABLE todo(
    "id" INT,
    "task" VARCHAR(255),
    "done" BOOL
  )`)
	createStmt.Exec()
	checkErr(err)
	statement, err := db.Prepare("INSERT INTO todo(id, task, done) values(?,?,?)")
	checkErr(err)
	statement.Exec(1, "test task", false)

	// query
	rows, err := db.Query("SELECT * FROM todo")
	checkErr(err)
	var id int
	var task string
	var done bool

	for rows.Next() {
		err = rows.Scan(&id, &task, &done)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(task)
		fmt.Println(done)
	}

	rows.Close()

	addPtr := flag.String("add", "", "Todo to add")
	removePtr := flag.Int("remove", 0, "Todo ID to remove")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	fmt.Println("add:", *addPtr)
	fmt.Println("remove:", *removePtr)

	if *addPtr != "" {

	}
}

func showTodo(db *sql.DB) {
	// query
	rows, err := db.Query("SELECT * FROM todo")
  checkErr(err)
	var id int
	var task string
	var status bool
	fmt.Println("id, task, status")
	for rows.Next() {
		rows.Scan(&id, &task, &status)
		fmt.Println(string(id) + ", " + task + ", " + strconv.FormatBool(status))
	}
	rows.Close()
}

func addTodo(db *sql.DB, task string) {
	rows, err := db.Query(`SELECT COUNT(*) FROM todo`)
	checkErr(err)
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	rows.Close()
  count = count + 1
	statement, err := db.Prepare(`INSERT INTO todo(id, task, status) values(?,?,?)`)
  checkErr(err)
	statement.Exec(count, task, 0)
}

func createDatabase() {
	fmt.Println("Creating database...")
	file, err := os.Create("todo.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
}

func createTable(db *sql.DB) {
	fmt.Println("Creating table...")
	createStmt, err := db.Prepare(`CREATE TABLE todo(
    "id" INT,
    "task" VARCHAR(255),
    "status" BOOL
  )`)
  checkErr(err)
	createStmt.Exec()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
