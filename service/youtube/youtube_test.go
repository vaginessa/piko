package youtube

import (
	"flag"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/piko/service"
	"github.com/mlvzk/piko/service/testutil"
)

const base = "https://www.youtube.com"

var update = flag.Bool("update", false, "update .golden files")

func TestIsValidTarget(t *testing.T) {
	tests := map[string]bool{
		"https://www.youtube.com/watch?v=HOK0uF-Z0xM": true,
		"https://youtube.com/watch?v=HOK0uF-Z0xM":     true,
		"youtube.com/watch?v=HOK0uF-Z0xM":             true,
		"https://youtu.be/HOK0uF-Z0xM":                true,
		"https://imgur.com/":                          false,
	}

	for target, expected := range tests {
		if (Youtube{}).IsValidTarget(target) != expected {
			t.Errorf("Invalid result, target: %v, expected: %v", target, expected)
		}
	}
}

func TestIteratorNext(t *testing.T) {
	ts := testutil.CacheHttpRequest(t, base, *update)
	defer ts.Close()

	iterator := YoutubeIterator{
		url: ts.URL + "/watch?v=HOK0uF-Z0xM",
	}

	items, err := iterator.Next()
	if err != nil {
		t.Fatalf("iterator.Next() error: %v", err)
	}

	if len(items) < 1 {
		t.Fatalf("Items array is empty")
	}

	for k := range items[0].Meta {
		if k[0] == '_' {
			items[0].Meta[k] = "ignore"
		}
	}

	expected := []service.Item{
		{
			Meta: map[string]string{
				"_ytConfig": "ignore",
				"author":    "Veltnox",
				"title":     "2 hours of Bloomer Music",
				"ext":       "mkv",
			},
			DefaultName: "%[title].%[ext]",
			AvailableOptions: map[string]([]string){
				"quality":   []string{"best", "medium", "worst"},
				"useFfmpeg": []string{"yes", "no"},
				"onlyAudio": []string{"yes", "no"},
			},
			DefaultOptions: map[string]string{
				"quality":   "medium",
				"useFfmpeg": "yes",
				"onlyAudio": "no",
			},
		},
	}

	if diff := pretty.Compare(items, expected); diff != "" {
		t.Errorf("%s diff:\n%s", t.Name(), diff)
	}
}
