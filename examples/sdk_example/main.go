package main

import (
	"fmt"
	"math/rand"
	"time"

	"victorgo/index"
	"victorgo/victor"
)

// Constants for the index and test parameters
const (
	// Index parameters
	IndexType = 0  // Assuming 0 is a valid index type (like linear search)
	Method    = 0  // Assuming 0 is a valid method
	Dims      = 50 // Dimension of vectors

	// Test parameters
	NumVectors = 10000 // Number of vectors to insert
	NumQueries = 3     // Number of searches to perform
	TopN       = 3     // Number of nearest neighbors to find
)

// generateRandomVector creates a random float32 vector of the specified dimension
func generateRandomVector(dims int) []float32 {
	vector := make([]float32, dims)
	for i := range vector {
		vector[i] = rand.Float32()*2.0 - 1.0 // Random values between -1.0 and 1.0
	}
	return vector
}

// printStats displays index statistics in a readable format
func printStats(stats *victor.IndexStats) {
	fmt.Println("\n--- Index Statistics ---")

	fmt.Println("Insert Operations:")
	fmt.Printf("  Count: %d, Total Time: %.6fs, Last: %.6fs, Min: %.6fs, Max: %.6fs\n",
		stats.Insert.Count, stats.Insert.Total, stats.Insert.Last, stats.Insert.Min, stats.Insert.Max)

	fmt.Println("Delete Operations:")
	fmt.Printf("  Count: %d, Total Time: %.6fs, Last: %.6fs, Min: %.6fs, Max: %.6fs\n",
		stats.Delete.Count, stats.Delete.Total, stats.Delete.Last, stats.Delete.Min, stats.Delete.Max)

	fmt.Println("Search Operations:")
	fmt.Printf("  Count: %d, Total Time: %.6fs, Last: %.6fs, Min: %.6fs, Max: %.6fs\n",
		stats.Search.Count, stats.Search.Total, stats.Search.Last, stats.Search.Min, stats.Search.Max)

	fmt.Println("SearchN Operations:")
	fmt.Printf("  Count: %d, Total Time: %.6fs, Last: %.6fs, Min: %.6fs, Max: %.6fs\n",
		stats.SearchN.Count, stats.SearchN.Total, stats.SearchN.Last, stats.SearchN.Min, stats.SearchN.Max)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Victor Vector Search Engine Test")
	fmt.Println("--------------------------------")

	// Create a new index
	fmt.Println("Creating index...")
	idxExample := index.HNSWIndex{EfSearch: 10, EfConstruct: 15, M0: 3}
	idx, err := victor.AllocIndex(IndexType, Method, Dims, &idxExample)
	if err != nil {
		panic(err)
	}
	defer idx.DestroyIndex() // Ensure index is destroyed when done

	// Insert random vectors
	fmt.Printf("Inserting %d random vectors...\n", NumVectors)
	vectors := make([][]float32, NumVectors)
	for i := 0; i < NumVectors; i++ {
		vectors[i] = generateRandomVector(Dims)
		err = idx.Insert(uint64(i+1), vectors[i])
		if err != nil {
			panic(err)
		}
	}

	// Check index size
	size, err := idx.Size()
	if err != nil {
		fmt.Printf("Failed to get index size: %v\n", err)
		panic(err)
	}
	fmt.Printf("Index size: %d vectors\n", size)

	// Check if vectors exist
	fmt.Println("Checking if vectors exist in index...")
	for i := 0; i < 5; i++ {
		id := uint64(rand.Intn(NumVectors) + 1)
		exists, err := idx.Contains(id)
		if err != nil {
			fmt.Printf("Error checking if ID %d exists: %v\n", id, err)
			continue
		}
		fmt.Printf("Vector ID %d exists: %v\n", id, exists)
	}

	// Test non-existent vector
	exists, err := idx.Contains(uint64(NumVectors + 100))
	if err != nil {
		fmt.Printf("Error checking non-existent vector: %v\n", err)
	}
	if exists {
		panic("Non-existent vector exists")
	}

	// Perform single vector searches
	fmt.Println("\nPerforming single vector searches...")
	for i := 0; i < NumQueries; i++ {
		queryVector := generateRandomVector(Dims)
		result, err := idx.Search(queryVector, Dims)
		if err != nil {
			fmt.Printf("Search failed: %v\n", err)
			continue
		}
		fmt.Printf("Query %d - Nearest neighbor: ID=%d, Distance=%.6f\n",
			i+1, result.ID, result.Distance)
	}

	// Perform multi-vector searches
	fmt.Printf("\nPerforming multi-vector searches (top %d)...\n", TopN)
	for i := 0; i < NumQueries; i++ {
		queryVector := generateRandomVector(Dims)
		results, err := idx.SearchN(queryVector, TopN)
		if err != nil {
			fmt.Printf("SearchN failed: %v\n", err)
			continue
		}

		fmt.Printf("Query %d - Top %d nearest neighbors:\n", i+1, len(results))
		for j, result := range results {
			fmt.Printf("  %d. ID=%d, Distance=%.6f\n", j+1, result.ID, result.Distance)
		}
	}

	// Delete some vectors
	fmt.Println("\nDeleting some vectors...")
	for i := 0; i < 4; i++ {
		id := uint64(rand.Intn(NumVectors) + 1)
		err := idx.Delete(id)
		if err != nil {
			fmt.Printf("Failed to delete vector %d: %v\n", id, err)
			continue
		}
		fmt.Printf("Deleted vector ID %d\n", id)

		// Verify deletion
		exists, err := idx.Contains(id)
		if err != nil {
			fmt.Printf("Error checking if ID %d exists: %v\n", id, err)
			continue
		}
		fmt.Printf("Vector ID %d exists after deletion: %v (should be false)\n", id, exists)
	}

	// Get updated size
	newSize, err := idx.Size()
	if err != nil {
		fmt.Printf("Failed to get updated index size: %v\n", err)
	} else {
		fmt.Printf("Updated index size: %d vectors (deleted %d vectors)\n",
			newSize, size-newSize)
	}

	// Get and print statistics
	stats, err := idx.GetStats()
	if err != nil {
		fmt.Printf("Failed to get index statistics: %v\n", err)
	} else {
		printStats(stats)
	}

	fmt.Println("\nTest completed successfully!")
}
