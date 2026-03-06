package helpers

import "gorm.io/gorm"

// PaginateResult holds pagination metadata returned by Paginate.
type PaginateResult struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Paginate executes a counted, offset/limit query on db and populates dest.
// page is clamped to a minimum of 1. perPage is clamped to 15 if outside [1,100].
func Paginate(db *gorm.DB, page, perPage int, dest interface{}) (*PaginateResult, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	var total int64
	db.Count(&total)

	offset := (page - 1) * perPage
	err := db.Offset(offset).Limit(perPage).Find(dest).Error

	var totalPages int
	if total > 0 {
		totalPages = int(total) / perPage
		if int(total)%perPage != 0 {
			totalPages++
		}
	}

	return &PaginateResult{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, err
}
