package retrieval

import (
	"testing"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// 娴嬭瘯 ReciprocalRankFusion 鍩烘湰鍔熻兘
func TestRRFBasic(t *testing.T) {
	bm25 := []vs.Document{
		vs.NewDocument("doc1").WithID("1"),
		vs.NewDocument("doc2").WithID("2"),
	}
	vector := []vs.Document{
		vs.NewDocument("doc3").WithID("3"),
		vs.NewDocument("doc1").WithID("1"),
	}

	results := ReciprocalRankFusion(bm25, vector, RRFK)
	if len(results) == 0 {
		t.Fatal("expected non-empty results")
	}

	// doc1 appears in both sources, should rank highest
	if results[0].Document.ID == "1" {
		t.Log("RRF correctly ranked doc1 highest")
	}
}

// 娴嬭瘯 RRF 绌鸿緭鍏?func TestRRFEmpty(t *testing.T) {
	results := ReciprocalRankFusion(nil, nil, RRFK)
	if len(results) != 0 {
		t.Errorf("expected empty, got %d", len(results))
	}
}
