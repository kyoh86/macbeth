package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitCoordinates(t *testing.T) {
	want := []string{
		"fill", "~", "~-1", "1", "-1", "~1", "~",
	}
	got := splitCoordinates([]string{"fill", "~", "~-1", "1", "-1~1~"})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("invalid result (-want, +got):\n%s", diff)
	}
}

func TestChunk32Blocks(t *testing.T) {
	for _, c := range []struct {
		blocks    int
		wantSize  int
		wantCount int
		wantRem   int
	}{{
		blocks:    0,
		wantSize:  0,
		wantCount: 0,
		wantRem:   0,
	}, {
		blocks:    1,
		wantSize:  1,
		wantCount: 1,
		wantRem:   0,
	}, {
		blocks:    2,
		wantSize:  2,
		wantCount: 1,
		wantRem:   0,
	}, {
		blocks:    31,
		wantSize:  31,
		wantCount: 1,
		wantRem:   0,
	}, {
		blocks:    32,
		wantSize:  32,
		wantCount: 1,
		wantRem:   0,
	}, {
		blocks:    33,
		wantSize:  17,
		wantCount: 1,
		wantRem:   16,
	}, {
		blocks:    34,
		wantSize:  17,
		wantCount: 2,
		wantRem:   0,
	}, {
		blocks:    64,
		wantSize:  32,
		wantCount: 2,
		wantRem:   0,
	}, {
		blocks:    65,
		wantSize:  22,
		wantCount: 2,
		wantRem:   21,
	}, {
		blocks:    66,
		wantSize:  22,
		wantCount: 3,
		wantRem:   0,
	}, {
		blocks:    67,
		wantSize:  23,
		wantCount: 2,
		wantRem:   21,
	}} {
		t.Run(fmt.Sprintf("%d blocks", c.blocks), func(t *testing.T) {
			gotSize, gotCount, gotRem := chunk32Blocks(c.blocks)
			if c.wantSize != gotSize {
				t.Errorf("size missmatched (want: %d, got: %d)", c.wantSize, gotSize)
			}
			if c.wantCount != gotCount {
				t.Errorf("count missmatched (want: %d, got: %d)", c.wantCount, gotCount)
			}
			if c.wantRem != gotRem {
				t.Errorf("remain missmatched (want: %d, got: %d)", c.wantRem, gotRem)
			}
			if gotBlocks := gotSize*gotCount + gotRem; gotBlocks != c.blocks {
				t.Errorf("blocks missmatched (want: %d, got: %d)", c.blocks, gotBlocks)
			}
		})
	}
}
