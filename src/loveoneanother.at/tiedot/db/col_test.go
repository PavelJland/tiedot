package db

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const (
	COL_BENCH_SIZE    = 1000000 // Number of documents made available for collection benchmark
	COL_BENCH_THREADS = 16      // Number of threads for collection benchmark
)

func TestInsertRead(t *testing.T) {
	tmp := "/tmp/tiedot_col_test"
	os.RemoveAll(tmp)
	defer os.RemoveAll(tmp)
	col, err := OpenCol(tmp)
	if err != nil {
		t.Errorf("Failed to open: %v", err)
	}
	docs := []string{`{"a": 1}`, `{"b": 2}`}
	var jsonDoc [2]interface{}
	json.Unmarshal([]byte(docs[0]), &jsonDoc[0])
	json.Unmarshal([]byte(docs[1]), &jsonDoc[1])

	ids := [2]uint64{}
	if ids[0], err = col.Insert(jsonDoc[0]); err != nil {
		t.Errorf("Failed to insert: %v", err)
	}
	if ids[1], err = col.Insert(jsonDoc[1]); err != nil {
		t.Errorf("Failed to insert: %v", err)
	}

	if col.Read(ids[0]).(map[string]interface{})[string('a')].(float64) != 1.0 {
		t.Errorf("Failed to read back document, got %v", col.Read(ids[0]))
	}
	if col.Read(ids[1]).(map[string]interface{})[string('b')].(float64) != 2.0 {
		t.Errorf("Failed to read back document, got %v", col.Read(ids[1]))
	}
}

func TestInsertUpdateRead(t *testing.T) {
	tmp := "/tmp/tiedot_col_test"
	os.RemoveAll(tmp)
	defer os.RemoveAll(tmp)
	col, err := OpenCol(tmp)
	if err != nil {
		t.Errorf("Failed to open: %v", err)
	}

	docs := []string{`{"a": 1}`, `{"b": 2}`}
	var jsonDoc [2]interface{}
	json.Unmarshal([]byte(docs[0]), &jsonDoc[0])
	json.Unmarshal([]byte(docs[1]), &jsonDoc[1])

	updatedDocs := []string{`{"a": 2}`, `{"b": "abcdefghijklmnopqrstuvwxyz"}`}
	var updatedJsonDoc [2]interface{}
	json.Unmarshal([]byte(updatedDocs[0]), &updatedJsonDoc[0])
	json.Unmarshal([]byte(updatedDocs[1]), &updatedJsonDoc[1])

	ids := [2]uint64{}
	if ids[0], err = col.Insert(jsonDoc[0]); err != nil {
		t.Errorf("Failed to insert: %v", err)
	}
	if ids[1], err = col.Insert(jsonDoc[1]); err != nil {
		t.Errorf("Failed to insert: %v", err)
	}

	if ids[0], err = col.Update(ids[0], updatedJsonDoc[0]); err != nil {
		t.Errorf("Failed to update: %v", err)
	}
	if ids[1], err = col.Update(ids[1], updatedJsonDoc[1]); err != nil {
		t.Errorf("Failed to update: %v", err)
	}

	if col.Read(ids[0]).(map[string]interface{})[string('a')].(float64) != 2.0 {
		t.Errorf("Failed to read back document, got %v", col.Read(ids[0]))
	}
	if col.Read(ids[1]).(map[string]interface{})[string('b')].(string) != string("abcdefghijklmnopqrstuvwxyz") {
		t.Errorf("Failed to read back document, got %v", col.Read(ids[1]))
	}
}

func TestInsertDeleteRead(t *testing.T) {
	tmp := "/tmp/tiedot_col_test"
	os.RemoveAll(tmp)
	defer os.RemoveAll(tmp)
	col, err := OpenCol(tmp)
	if err != nil {
		t.Errorf("Failed to open: %v", err)
	}

	docs := []string{`{"a": 1}`, `{"b": 2}`}
	var jsonDoc [2]interface{}
	json.Unmarshal([]byte(docs[0]), &jsonDoc[0])
	json.Unmarshal([]byte(docs[1]), &jsonDoc[1])

	ids := [2]uint64{}
	if ids[0], err = col.Insert(jsonDoc[0]); err != nil {
		t.Errorf("Failed to insert: %v", err)
	}
	if ids[1], err = col.Insert(jsonDoc[1]); err != nil {
		t.Errorf("Failed to insert: %v", err)
	}
	col.Delete(ids[0])
	fmt.Println("Please ignore the next error message, it is intended to happen.")
	if col.Read(ids[0]) != nil {
		t.Errorf("Did not delete document, still read %v", col.Read(ids[0]))
	}
	if col.Read(ids[1]).(map[string]interface{})[string('b')].(float64) != 2 {
		t.Errorf("Failed to read back document, got %v", col.Read(ids[1]))
	}
}

func BenchmarkInsert(b *testing.B) {
	tmp := "/tmp/tiedot_col_bench"
	os.RemoveAll(tmp)
	defer os.RemoveAll(tmp)
	col, err := OpenCol(tmp)
	if err != nil {
		t.Errorf("Failed to open: %v", err)
	}
	var jsonDoc interface{}
	json.Unmarshal([]byte(`{"a"`), &jsonDoc)
}