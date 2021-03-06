// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	_ "net/http"
	"testing"
)

func TestScan(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Match all should return all documents
	cursor, err := client.Scan(testIndexName).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 3 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 3, cursor.Results.Hits.TotalHits)
	}
	if len(cursor.Results.Hits.Hits) != 0 {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", 0, len(cursor.Results.Hits.Hits))
	}

	pages := 0
	numDocs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 3 {
		t.Errorf("expected to retrieve %d hits; got %d", 3, numDocs)
	}
}

func TestScanWithSort(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// We sort on a numerical field, because sorting on the 'message' string field would
	// raise the whole question of tokenizing and analyzing.
	cursor, err := client.Scan(testIndexName).Sort("retweets", true).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 3 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 3, cursor.Results.Hits.TotalHits)
	}
	if len(cursor.Results.Hits.Hits) != 1 {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", 1, len(cursor.Results.Hits.Hits))
	}

	if cursor.Results.Hits.Hits[0].Id != "3" {
		t.Fatalf("expected hitID = %v; got %v", "3", cursor.Results.Hits.Hits[0].Id)

	}

	numDocs := 1 // The cursor already gave us a result
	pages := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 3 {
		t.Errorf("expected to retrieve %d hits; got %d", 3, numDocs)
	}
}

func TestScanWithSearchSource(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	src := NewSearchSource().
		Query(NewTermQuery("user", "olivere")).
		FetchSourceContext(NewFetchSourceContext(true).Include("retweets"))
	cursor, err := client.Scan(testIndexName).SearchSource(src).Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 2 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 2, cursor.Results.Hits.TotalHits)
	}

	numDocs := 0
	pages := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			if _, found := item["message"]; found {
				t.Fatalf("expected to not see field %q; got: %#v", "message", item)
			}
			numDocs += 1
		}
	}

	if pages != 3 {
		t.Errorf("expected to retrieve %d pages; got %d", 2, pages)
	}

	if numDocs != 2 {
		t.Errorf("expected to retrieve %d hits; got %d", 2, numDocs)
	}
}

func TestScanWithQuery(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun."}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Return tweets from olivere only
	termQuery := NewTermQuery("user", "olivere")
	cursor, err := client.Scan(testIndexName).
		Size(1).
		Query(termQuery).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 2 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 2, cursor.Results.Hits.TotalHits)
	}
	if len(cursor.Results.Hits.Hits) != 0 {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", 0, len(cursor.Results.Hits.Hits))
	}

	pages := 0
	numDocs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 2 {
		t.Errorf("expected to retrieve %d hits; got %d", 2, numDocs)
	}
}

func TestScanAndScrollWithMissingIndex(t *testing.T) {
	client := setupTestClient(t) // does not create testIndexName

	cursor, err := client.Scan(testIndexName).Scroll("30s").Do()
	if err != nil {
		t.Fatal(err)
	}
	if cursor == nil {
		t.Fatalf("expected cursor; got: %v", cursor)
	}

	// First request immediately returns EOS
	res, err := cursor.Next()
	if err != EOS {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatalf("expected results == %v; got: %v", nil, res)
	}
}

func TestScanAndScrollWithEmptyIndex(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	if isTravis() {
		t.Skip("test on Travis failes regularly with " +
			"Error 503 (Service Unavailable): SearchPhaseExecutionException[Failed to execute phase [init_scan], all shards failed]")
	}

	_, err := client.Flush().Index(testIndexName).WaitIfOngoing(true).Do()
	if err != nil {
		t.Fatal(err)
	}

	cursor, err := client.Scan(testIndexName).Scroll("30s").Do()
	if err != nil {
		t.Fatal(err)
	}
	if cursor == nil {
		t.Fatalf("expected cursor; got: %v", cursor)
	}

	// First request returns no error, but no hits
	res, err := cursor.Next()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatalf("expected results != nil; got: nil")
	}
	if res.ScrollId == "" {
		t.Fatalf("expected scrollId in results; got: %q", res.ScrollId)
	}
	if res.TotalHits() != 0 {
		t.Fatalf("expected TotalHits() = %d; got %d", 0, res.TotalHits())
	}
	if res.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got: nil")
	}
	if res.Hits.TotalHits != 0 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 0, res.Hits.TotalHits)
	}
	if res.Hits.Hits == nil {
		t.Fatalf("expected results.Hits.Hits != nil; got: %v", res.Hits.Hits)
	}
	if len(res.Hits.Hits) != 0 {
		t.Fatalf("expected len(results.Hits.Hits) == %d; got: %d", 0, len(res.Hits.Hits))
	}

	// Subsequent requests return EOS
	res, err = cursor.Next()
	if err != EOS {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatalf("expected results == %v; got: %v", nil, res)
	}

	res, err = cursor.Next()
	if err != EOS {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatalf("expected results == %v; got: %v", nil, res)
	}
}

func TestIssue119(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch."}
	comment1 := comment{User: "nico", Comment: "You bet."}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic."}

	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Type("comment").Id("1").Parent("1").BodyJson(&comment1).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Match all should return all documents
	cursor, err := client.Scan(testIndexName).Fields("_source", "_parent").Size(1).Do()
	if err != nil {
		t.Fatal(err)
	}

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		for _, hit := range searchResult.Hits.Hits {
			if hit.Type == "tweet" {
				if _, ok := hit.Fields["_parent"].(string); ok {
					t.Errorf("Type `tweet` cannot have any parent...")

					toPrint, _ := json.MarshalIndent(hit, "", "    ")
					t.Fatal(string(toPrint))
				}
			}

			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestScanWithSorter(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndLog(t, SetTraceLog(log.New(os.Stdout, "", 0)))
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "user1", Message: "Welcome to Golang and Elasticsearch.", Retweets: 8}
	tweet2 := tweet{User: "user2", Message: "Another unrelated topic.", Retweets: 4}
	tweet3 := tweet{User: "user3", Message: "Cycling is fun.", Retweets: 1}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	// Return tweets from olivere only
	cursor, err := client.Scan(testIndexName).
		Size(1).
		SortBy(NewFieldSort("retweets").Desc()).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	res := cursor.Results

	if res == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if res.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if res.Hits.TotalHits != 3 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 3, cursor.Results.Hits.TotalHits)
	}

	pages := 0
	numDocs := 0

	for {
		pages += 1

		for _, hit := range res.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			var doc tweet
			err := json.Unmarshal(*hit.Source, &doc)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
			switch numDocs {
			case 1:
				if got, want := doc.User, "user1"; got != want {
					t.Fatalf("expected document from %q first; got: %q", want, got)
				}
			case 2:
				if got, want := doc.User, "user2"; got != want {
					t.Fatalf("expected document from %q first; got: %q", want, got)
				}
			case 3:
				if got, want := doc.User, "user3"; got != want {
					t.Fatalf("expected document from %q first; got: %q", want, got)
				}
			default:
				t.Fatalf("expected no more than 3 documents; got: %d", numDocs)
			}
		}

		res, err = cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 3 {
		t.Errorf("expected to retrieve %d hits; got %d", 2, numDocs)
	}
}

func TestScanWithRawBody(t *testing.T) {
	//client := setupTestClientAndCreateIndexAndLog(t)
	client := setupTestClientAndCreateIndex(t)

	tweet1 := tweet{User: "olivere", Message: "Welcome to Golang and Elasticsearch.", Retweets: 4}
	tweet2 := tweet{User: "olivere", Message: "Another unrelated topic.", Retweets: 10}
	tweet3 := tweet{User: "sandrae", Message: "Cycling is fun.", Retweets: 3}

	// Add all documents
	_, err := client.Index().Index(testIndexName).Type("tweet").Id("1").BodyJson(&tweet1).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("2").BodyJson(&tweet2).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Index().Index(testIndexName).Type("tweet").Id("3").BodyJson(&tweet3).Do()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Flush().Index(testIndexName).Do()
	if err != nil {
		t.Fatal(err)
	}

	cursor, err := client.
		Scan(testIndexName).
		Body(`{"query":{"term":{"user":"olivere"}}}`).
		Size(1).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	if cursor.Results == nil {
		t.Fatalf("expected results != nil; got nil")
	}
	if cursor.Results.Hits == nil {
		t.Fatalf("expected results.Hits != nil; got nil")
	}
	if cursor.Results.Hits.TotalHits != 2 {
		t.Fatalf("expected results.Hits.TotalHits = %d; got %d", 2, cursor.Results.Hits.TotalHits)
	}
	if len(cursor.Results.Hits.Hits) != 0 {
		t.Fatalf("expected len(results.Hits.Hits) = %d; got %d", 0, len(cursor.Results.Hits.Hits))
	}

	pages := 0
	numDocs := 0

	for {
		searchResult, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		pages += 1

		for _, hit := range searchResult.Hits.Hits {
			if hit.Index != testIndexName {
				t.Fatalf("expected SearchResult.Hits.Hit.Index = %q; got %q", testIndexName, hit.Index)
			}
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				t.Fatal(err)
			}
			numDocs += 1
		}
	}

	if pages <= 0 {
		t.Errorf("expected to retrieve at least 1 page; got %d", pages)
	}

	if numDocs != 2 {
		t.Errorf("expected to retrieve %d hits; got %d", 2, numDocs)
	}
}
