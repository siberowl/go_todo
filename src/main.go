package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	if _, err := os.Stat("todo.db"); err != nil {
		createDatabase()
	}

	db, _ := sql.Open("sqlite3", "./todo.db")

	_, err := db.Query("SELECT * FROM todo")
	if err != nil {
		createTable(db)
	}

	addPtr := flag.String("add", "", "Todo to add")
	delPtr := flag.Int("del", 0, "Todo ID to remove")
	donePtr := flag.Int("done", 0, "Todo ID to mark done")
	showPtr := flag.Bool("show", false, "Flag to show todo list")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	if *addPtr != "" {
		addTodo(db, *addPtr)
	}
	if *delPtr != 0 {
		delTodo(db, *delPtr)
	}
	if *donePtr != 0 {
		doneTodo(db, *donePtr)
	}
	if *showPtr {
		showTodo(db)
	}
}

type entry struct {
	id     int
	task   string
	status bool
}

func getEntries(db *sql.DB) []entry {
	rows, err := db.Query("SELECT COUNT(*) FROM todo")
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	rows.Close()

	entries := make([]entry, count)

	rows, err = db.Query("SELECT * FROM todo")
	checkErr(err)
	var id int
	var task string
	var status bool
	index := 0
	for rows.Next() {
		rows.Scan(&id, &task, &status)
		entries[index] = entry{id, task, status}
		index++
	}
	rows.Close()
	return entries
}

func showTodo(db *sql.DB) {
	entries := getEntries(db)
	for i := 0; i < len(entries); i++ {
		fmt.Println("["+strconv.Itoa(entries[i].id) + "] " + entries[i].task + " | " + strconv.FormatBool(entries[i].status))
	}
}

func addTodo(db *sql.DB, task string) {
	entries := getEntries(db)
	isnew := true
	for i := 0; i < len(entries); i++ {
		if entries[i].task == task {
			isnew = false
		}

	}
	if isnew {
		statement, err := db.Prepare(`INSERT INTO todo(id, task, status) values(?,?,?)`)
		checkErr(err)
		statement.Exec(len(entries)+1, task, 0)
	}
}

func delTodo(db *sql.DB, id int) {
	entries := getEntries(db)
	exists := false
	for i := 0; i < len(entries); i++ {
		if entries[i].id == id {
			exists = true
		}
	}
	if exists {
		statement, err := db.Prepare(`DELETE FROM todo WHERE id=` + strconv.Itoa(id))
		checkErr(err)
		statement.Exec()
		showTodo(db)
	} else {
		fmt.Println("ID does not exist")
	}
}

func doneTodo(db *sql.DB, id int) {
	entries := getEntries(db)
	exists := false
	for i := 0; i < len(entries); i++ {
		if entries[i].id == id {
			exists = true
		}
	}
	if exists {
		statement, err := db.Prepare(`UPDATE todo SET status = true WHERE id=` + strconv.Itoa(id))
		checkErr(err)
		statement.Exec()
		showTodo(db)
	} else {
		fmt.Println("ID does not exist")
	}
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
