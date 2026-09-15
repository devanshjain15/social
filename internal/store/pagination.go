package store

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=25"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	q := r.URL.Query()

	limit := q.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}

		fq.Limit = l
	}

	offset := q.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}

		fq.Offset = l
	}

	sort := q.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := q.Get("tags")
	if tags != "" {
		parts := strings.Split(tags, ",")
		fq.Tags = make([]string, 0, len(parts))
		for _, t := range parts {
			t = strings.TrimSpace(t)
			if decoded, err := url.QueryUnescape(t); err == nil {
				t = decoded
			}
			if t != "" {
				fq.Tags = append(fq.Tags, t)
			}
		}
	}

	search := q.Get("search")
    if search != "" {
       fq.Search = search
    }

	since := q.Get("since")
	if since != "" {
		fq.Since = parseDateTime(since)
	}

	until := q.Get("search")
	if search != "" {
		fq.Until = parseDateTime(until)
	}

	return fq, nil
}

func parseDateTime(s string) string {
	parsed, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}
	return parsed.Format(time.DateTime)
}
