package main

import (
	"bloomFilters/simpleBloom"
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"time"
)

func main() {
	//test()
	dbPerfTest()
	//bloomPerfTest()

}

//Test Bloom Filter

func test() {
	bloomFilter := simpleBloom.NewBloom(100, nil)
	simpleBloom.AddElement(bloomFilter, "Sivani")
	simpleBloom.AddElement(bloomFilter, "Miray")
	fmt.Println(simpleBloom.CheckElement(bloomFilter, "Sivani"))
	fmt.Println(simpleBloom.CheckElement(bloomFilter, "Miray"))
	fmt.Println(simpleBloom.CheckElement(bloomFilter, "Mouli"))
}

func dbPerfTest() {
	//Calculate time taken for performance test
	startTime := time.Now()

	missingCount := 0
	wordsInserted := 0
	// Insert words from words_insert.txt into database
	file, err := os.Open("words_insert.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	//Create a postgres database connection
	conn, err := pgx.Connect(context.Background(), "postgresql://moulisanketh:password@localhost/bloomingo")

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Insert into postgres database
		_, err = conn.Exec(context.Background(), "INSERT INTO words (word) VALUES ($1)", scanner.Text())
		if err != nil {
			log.Fatalf("Unable to insert into database: %v\n", err)
		}
		wordsInserted++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	// Fetch words from words_test.txt and check if they exist in database
	file, err = os.Open("words_test.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		// Check if word exists in database
		var word string
		err = conn.QueryRow(context.Background(), "SELECT word FROM words WHERE word = $1", scanner.Text()).Scan(&word)
		if err != nil {
			missingCount++
		}
	}

	fmt.Printf("Words inserted: %d\n", wordsInserted)
	fmt.Printf("Words that are not definitely present: %d\n", 50)
	fmt.Printf("Words that are not present as reported by database: %d\n", missingCount)
	//Print Time taken for performance test
	fmt.Printf("Time taken for db performance test: %v\n", time.Since(startTime))
}

//Performance Test for Bloom Filter

func bloomPerfTest() {
	startTime := time.Now()
	bloomFilterSize := 100000
	bloomFilter := simpleBloom.NewBloom(bloomFilterSize, nil)
	missingCount := 0
	wordsInserted := 0

	// Fetch words from words_insert.txt and insert into bloom filter
	file, err := os.Open("words_insert.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		simpleBloom.AddElement(bloomFilter, scanner.Text())
		wordsInserted++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	// Fetch words from words_test.txt and check if they exist in bloom filter
	file, err = os.Open("words_test.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		if !simpleBloom.CheckElement(bloomFilter, scanner.Text()) {
			missingCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	fmt.Printf("Words inserted: %d\n", wordsInserted)
	fmt.Printf("Size of Bloom Filter: %d\n", bloomFilterSize)
	fmt.Printf("Words that are not definitely present: %d\n", 50)
	fmt.Printf("Words that are not present as reported by Bloom Filter: %d\n", missingCount)
	fmt.Printf("Time taken for bloom filter performance test: %v\n", time.Since(startTime))

}
