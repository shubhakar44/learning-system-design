package main

import (
	"context"
	db "learning/systems/internal/platform/database"
	"log"
	"sync"
	"time"
)

type Tasks struct {
	ID         int
	title      string
	status     string
	worker_id  int
	created_at string
}

func calculateWithoutPool(wg *sync.WaitGroup) {
	ctx := context.Background()
	conn := db.CreateInstance(ctx)
	defer conn.Close(ctx)
	defer wg.Done()

	rows, err := conn.Query(ctx, "SELECT title FROM tasks;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var u Tasks
		err := rows.Scan(&u.title)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		//fmt.Printf(" Name: %s\n", u.title)

	}
	defer rows.Close()
}

func calculateWithPool(pool *db.Pool, wg *sync.WaitGroup) {
	ctx := context.Background()
	conn, _ := db.GetPoolInstance(pool, ctx)

	defer wg.Done()

	rows, err := conn.Query(ctx, "SELECT title FROM tasks;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var u Tasks
		err := rows.Scan(&u.title)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		//fmt.Printf(" Name: %s\n", u.title)
	}
	defer rows.Close()
	defer db.PutPoolInstance(pool, conn)
}

func main() {
	var wg sync.WaitGroup
	maxConns := 10
	start := time.Now()
	for i := 0; i <= 400; i++ {
		wg.Add(1)
		go calculateWithoutPool(&wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	log.Print(elapsed)

	ctx := context.Background()
	pool := db.CreateConnectionPool(ctx, maxConns)

	start = time.Now()
	for i := 0; i <= 400; i++ {
		wg.Add(1)
		go calculateWithPool(pool, &wg)
	}
	wg.Wait()
	elapsed = time.Since(start)
	log.Print(elapsed)
}
