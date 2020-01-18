package main

import (
	"flag"
	"fmt"
	"sort"
	"time"
)
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var P = make(map[string]int)
var count int
var Last_que int
var Last_select int
var Last_commit int
var Last_rollback int
var Last_slow int
var Last_ir int
var Last_iu int
var Last_id int
var Last_ii int

func Questions(conn string) {
	//db, err := sql.Open("mysql", "txc:123456@tcp(192.168.99.105:3309)/txc?charset=utf8")
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err)
	}
	//	defer db.Close()
	fmt.Println()
	var_name := []string{"Innodb_buffer_pool_pages_dirty", "Innodb_buffer_pool_pages_total", "Innodb_buffer_pool_read_requests", "Com_commit", "Com_rollback", "Slow_queries",
		"Innodb_buffer_pool_reads", "Innodb_data_reads", "Innodb_data_writes", "Innodb_log_waits", "Innodb_buffer_pool_write_requests", "Innodb_rows_deleted",
		"Innodb_rows_inserted", "Innodb_rows_read", "Innodb_rows_updated", "Innodb_log_writes", "Innodb_os_log_written", "Innodb_os_log_fsyncs", "Created_tmp_files",
		"Created_tmp_disk_tables", "Com_update", "Com_select", "Com_insert", "Com_delete", "Questions", "Threads_running", "Threads_connected", "Innodb_row_lock_current_waits"}
	rows, err := db.Query("show global status")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var Name string
		var Value int
		rows.Scan(&Name, &Value)
		for _, v := range var_name {
			if Name == v {
				P[Name] = Value

			}
		}

	}
}

func main() {
	user := flag.String("u", "", "mysql user")
	passwd := flag.String("p", "", "mysql user password")
	host := flag.String("h", "127.0.0.1", "mysql host")
	port := flag.String("port", "3306", "mysql port")
	flag.Parse()
	cn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", *user, *passwd, *host, *port)
	fmt.Print("----------|------------|----- Innodb row operation ----|---Buffer Pool---|--lock--|----------|----Thread---\n")
	fmt.Println("---Time---|--QPS--TPS--|  read inserted updated deleted|   Hit    dirty  |  lock  | slow logs|  conn   run")
	for  {
		Questions(cn)
		var s []string
		for k, _ := range P {
			s = append(s, k)
		}
		sort.Strings(s)
		a := float32(P["Innodb_buffer_pool_read_requests"] - P["Innodb_buffer_pool_reads"])
		Hit := (a / float32(P["Innodb_buffer_pool_read_requests"])) * 100
		b := (float32(P["Innodb_buffer_pool_pages_dirty"]) / float32(P["Innodb_buffer_pool_pages_total"])) * 100
		ts := time.Now().Format("15:04:05")
		//fmt.Println("innodb_buffer_pool_hit: ", Hit, "%", "           innodb_buffer_pool_dirty_percent: ", b, "%")
		//fmt.Println("Threads_running: ", P["Threads_running"], "     Threads_connected: ", P["Threads_connected"])
		if count < 1 {
			Last_que = P["Questions"]
			Last_select = P["Com_select"]
			Last_commit = P["Com_commit"]
			Last_rollback = P["Com_rollback"]
			Last_slow = P["Slow_queries"]
			Last_ir = P["Innodb_rows_read"]
			Last_id = P["Innodb_rows_deleted"]
			Last_ii = P["Innodb_rows_inserted"]
			Last_iu = P["Innodb_rows_updated"]

		}
		qps := P["Questions"] - Last_que
		tps := (P["Com_commit"] - Last_commit) + (P["Com_rollback"] - Last_rollback)
		slow := P["Slow_queries"] - Last_slow
		//sel :=P["Com_select"]-Last_select
		lock := P["Innodb_row_lock_current_waits"]
		conn := P["Threads_connected"]
		run := P["Threads_running"]
		ir := P["Innodb_rows_read"] - Last_ir
		iu := P["Innodb_rows_updated"] - Last_iu
		id := P["Innodb_rows_deleted"] - Last_id
		ii := P["Innodb_rows_inserted"] - Last_ii
		//tmp :=P["Created_tmp_disk_tables"]
		//fmt.Println("TPS:",((P["Com_commit"]-Last_commit)+(P["Com_rollback"]-Last_rollback)))
		//fmt.Println("slow log: ",P["Slow_queries"]-Last_slow)
		//fmt.Println("com_select: ",P["Com_select"]-Last_select)
		//fmt.Println("Innodb_row_lock_current_waits: ",P["Innodb_row_lock_current_waits"])
		//fmt.Println("Created_tmp_disk_tables: ",P["Created_tmp_disk_tables"])
		if count < 1 {

			fmt.Printf(" %s |   %d     %d  | %6d  %6d  %6d  %6d| %.2f     %.2f | %5d  | %6d   |%6d%6d \n", ts, qps, tps, ir, ii, iu, id, Hit, b, lock, slow, conn, run)
		} else {
			fmt.Printf(" %s |%d   %d|%6d  %6d  %6d  %6d| %.2f     %.2f | %5d  | %6d   |%6d%6d \n", ts, qps, tps, ir, ii, iu, id, Hit, b, lock, slow, conn, run)
		}

		count++
		Last_que = P["Questions"]
		Last_select = P["Com_select"]
		Last_commit = P["Com_commit"]
		Last_rollback = P["Com_rollback"]
		Last_slow = P["Slow_queries"]
		Last_ir = P["Innodb_rows_read"]
		Last_id = P["Innodb_rows_deleted"]
		Last_ii = P["Innodb_rows_inserted"]
		Last_iu = P["Innodb_rows_updated"]
		time.Sleep(10 * time.Second)
	}
}
