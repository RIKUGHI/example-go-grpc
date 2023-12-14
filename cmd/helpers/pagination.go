package helpers

import (
	"go-grpc/pb"
	"math"

	"gorm.io/gorm"
)

func Pagination(sql *gorm.DB, page int64, pagination *pb.Pagination) (int64, int64) {
	var total int64
	var limit int64 = 1
	var offest int64

	sql.Count(&total)

	if page == 1 {
		offest = 0
	} else {
		offest = (page - 1) * limit
	}

	pagination.Total = uint64(total)
	pagination.PerPage = uint32(limit)
	pagination.CurrentPage = uint32(page)
	pagination.LastPage = uint32(math.Ceil(float64(total) / float64(limit)))

	return offest, limit
}
