package postgres

// makeFormSelect selects one column from specified table and returns it as slice
// func makeFormSelect(table, col, usr string) []string {
// 	db, err := sql.Open("postgres", dbconf())
// 	checkErr(err)
// 	defer db.Close()

// 	err = db.Ping()
// 	checkErr(err)

// 	q := fmt.Sprintf("SELECT %s FROM %s WHERE act='true' ORDER by %s", col, table, col)
// 	//log.Println("DEBUG: q =", q)

// 	rows, err := db.Query(q)
// 	checkErr(err)
// 	defer rows.Close()

// 	var option string
// 	swap := make([]string, 500)

// 	i := 0
// 	swap[i] = "--"
// 	for rows.Next() {
// 		i++
// 		err = rows.Scan(&option)
// 		checkErr(err)
// 		swap[i] = option
// 	}
// 	nonNullCount := 0

// 	for i := 0; i < len(swap); i++ {
// 		if swap[i] != "" {
// 			nonNullCount++
// 		}
// 	}

// 	out := make([]string, nonNullCount)
// 	copy(out, swap)

// 	err = rows.Err()
// 	checkErr(err)
// 	return out
// }
