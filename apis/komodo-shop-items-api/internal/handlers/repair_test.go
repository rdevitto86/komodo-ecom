package handlers

import (
	"net/http"
	"testing"
)

func Test_parsePaginationParams_defaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		query     string
		wantPage  int
		wantLimit int
	}{
		{
			name:      "no params uses defaults",
			query:     "",
			wantPage:  1,
			wantLimit: defaultRepairPageLimit,
		},
		{
			name:      "explicit page and limit",
			query:     "page=3&limit=10",
			wantPage:  3,
			wantLimit: 10,
		},
		{
			name:      "limit clamped to max",
			query:     "limit=999",
			wantPage:  1,
			wantLimit: maxRepairPageLimit,
		},
		{
			name:      "invalid page falls back to default",
			query:     "page=abc",
			wantPage:  1,
			wantLimit: defaultRepairPageLimit,
		},
		{
			name:      "zero page falls back to default",
			query:     "page=0",
			wantPage:  1,
			wantLimit: defaultRepairPageLimit,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			target := "/services/repair"
			if tc.query != "" {
				target += "?" + tc.query
			}
			req, err := http.NewRequest(http.MethodGet, target, nil)
			if err != nil {
				t.Fatalf("failed to build request: %v", err)
			}
			page, limit := parsePaginationParams(req)
			if page != tc.wantPage {
				t.Errorf("page: got %d, want %d", page, tc.wantPage)
			}
			if limit != tc.wantLimit {
				t.Errorf("limit: got %d, want %d", limit, tc.wantLimit)
			}
		})
	}
}

func Test_paginateSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		total     int
		page      int
		limit     int
		wantStart int
		wantEnd   int
	}{
		{"first page full", 50, 1, 20, 0, 20},
		{"second page full", 50, 2, 20, 20, 40},
		{"last partial page", 45, 3, 20, 40, 45},
		{"page beyond end returns empty range", 10, 5, 20, 10, 10},
		{"empty total", 0, 1, 20, 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			start, end := paginateSlice(tc.total, tc.page, tc.limit)
			if start != tc.wantStart || end != tc.wantEnd {
				t.Errorf("got [%d:%d], want [%d:%d]", start, end, tc.wantStart, tc.wantEnd)
			}
		})
	}
}
