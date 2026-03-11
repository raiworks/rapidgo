package database

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// PageResult holds metadata for offset-based pagination.
type PageResult struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Paginate performs offset-based pagination on the given query.
// page is clamped to ≥ 1, perPage is clamped to 1–100 (default 15).
// dest must be a pointer to a slice.
func Paginate(db *gorm.DB, page, perPage int, dest interface{}) (*PageResult, error) {
	if perPage <= 0 {
		perPage = 15
	}
	if perPage > 100 {
		perPage = 100
	}
	if page < 1 {
		page = 1
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	offset := (page - 1) * perPage
	if err := db.Offset(offset).Limit(perPage).Find(dest).Error; err != nil {
		return nil, err
	}

	return &PageResult{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// CursorResult holds metadata for cursor-based pagination.
type CursorResult struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	PerPage    int    `json:"per_page"`
	HasMore    bool   `json:"has_more"`
}

// CursorPaginate performs cursor-based pagination.
//
// Parameters:
//   - db: a *gorm.DB query (pre-scoped with any WHERE clauses).
//   - cursor: base64-encoded value of the last seen item (empty string = first page).
//   - orderCol: column to order and paginate by (must be unique and indexed, e.g. "id").
//   - perPage: items per page, clamped to 1–100 (default 15).
//   - direction: "next" (default, items after cursor) or "prev" (items before cursor).
//   - dest: pointer to a slice of model structs.
//
// The function sets NextCursor/PrevCursor on the result based on the last/first
// item's orderCol value. The caller can pass these cursors back for the next page.
func CursorPaginate(db *gorm.DB, cursor, orderCol string, perPage int, direction string, dest interface{}) (*CursorResult, error) {
	if perPage <= 0 {
		perPage = 15
	}
	if perPage > 100 {
		perPage = 100
	}
	if orderCol == "" {
		orderCol = "id"
	}
	if direction != "prev" {
		direction = "next"
	}

	query := db.Session(&gorm.Session{})

	if cursor != "" {
		decoded, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, fmt.Errorf("pagination: invalid cursor: %w", err)
		}
		cursorVal := string(decoded)
		if direction == "next" {
			query = query.Where(fmt.Sprintf("%s > ?", orderCol), cursorVal)
		} else {
			query = query.Where(fmt.Sprintf("%s < ?", orderCol), cursorVal)
		}
	}

	// Fetch one extra row to detect whether more pages exist.
	if direction == "next" {
		query = query.Order(fmt.Sprintf("%s ASC", orderCol))
	} else {
		query = query.Order(fmt.Sprintf("%s DESC", orderCol))
	}

	if err := query.Limit(perPage + 1).Find(dest).Error; err != nil {
		return nil, err
	}

	// Trim the extra row via reflection and determine hasMore.
	slicePtr := reflect.ValueOf(dest)
	slice := slicePtr.Elem()
	n := slice.Len()

	result := &CursorResult{PerPage: perPage}

	if n > perPage {
		result.HasMore = true
		slice.Set(slice.Slice(0, perPage))
		n = perPage
	}

	// Set cursors from the boundary items.
	if n > 0 {
		last := slice.Index(n - 1)
		result.NextCursor = encodeCursorField(last, orderCol)

		first := slice.Index(0)
		result.PrevCursor = encodeCursorField(first, orderCol)
	}

	return result, nil
}

// encodeCursor base64-encodes a cursor value.
func encodeCursor(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// encodeCursorField extracts a field from a struct by GORM column name and
// base64-encodes it. Falls back to the "ID" field if the column is not found.
func encodeCursorField(v reflect.Value, col string) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("gorm")
		if tagContainsColumn(tag, col) {
			return encodeCursor(fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}
	// Fallback: try matching field name case-insensitively.
	for i := 0; i < t.NumField(); i++ {
		if fmt.Sprintf("%s", t.Field(i).Name) == "ID" && col == "id" {
			return encodeCursor(fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}
	return ""
}

// tagContainsColumn checks if a GORM struct tag contains column:name.
func tagContainsColumn(tag, col string) bool {
	// Match "column:col" or "primaryKey" for "id"
	target := "column:" + col
	for _, part := range splitTag(tag) {
		if part == target {
			return true
		}
	}
	return false
}

func splitTag(tag string) []string {
	var parts []string
	current := ""
	for _, c := range tag {
		if c == ';' {
			if current != "" {
				parts = append(parts, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
