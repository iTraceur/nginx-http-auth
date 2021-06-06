package utils

import (
	"math"
)

type Page struct {
	PageNo     int
	PageSize   int
	TotalPage  int
	TotalCount int64
	FirstPage  bool
	LastPage   bool
	RangePrev  int
	Range      []int
	RangeNext  int
}

// 分页生成函数
func Paginate(count int64, pageNo int, pageSize int) Page {
	// 计算总页数
	tp := count / int64(pageSize)
	if count%int64(pageSize) > 0 {
		tp = tp + 1
	}

	// 根据当前页号计算出可跳转的开始页、最后页，向前快速跳转页、向后快速跳转页、可跳转页列表
	var start, end, rPrev, rNext int
	var ran []int
	if pageNo <= 2 {
		start = 1
		end = int(math.Min(5, float64(tp)))
		rPrev = 0
		if int(tp) < 6 {
			rNext = 0
		} else {
			rNext = pageNo + 5
		}
	} else if int64(pageNo) >= tp-2 {
		start = int(math.Max(1, float64(tp-4)))
		end = int(tp)
		if tp == 4 || tp == 5 {
			rPrev = 0
		} else if pageNo >= 4 && pageNo <= 5 {
			rPrev = 1
		} else {
			rPrev = int(math.Max(0, float64(pageNo-5)))
		}
		rNext = 0
	} else {
		start = pageNo - 2
		end = pageNo + 2
		if pageNo <= 3 {
			rPrev = 0
		} else if pageNo > 3 && pageNo <= 5 {
			rPrev = 1
		} else {
			rPrev = int(math.Max(0, float64(pageNo-5)))
		}
		if pageNo >= int(tp)-2 {
			rNext = 0
		} else if pageNo >= int(tp)-5 && pageNo < int(tp)-2 {
			rNext = int(tp)
		} else {
			rNext = pageNo + 5
		}
	}

	// 可跳转页列表
	for i := start; i <= end; i++ {
		ran = append(ran, i)
	}

	// 向前快速跳转页
	if rPrev != 0 {
		rPrev = int(math.Max(1, float64(rPrev)))
	}
	// 向后快速跳转页
	if rNext != 0 {
		rNext = int(math.Min(float64(tp), float64(rNext)))
	}

	return Page{
		PageNo:     pageNo,
		PageSize:   pageSize,
		TotalPage:  int(tp),
		TotalCount: count,
		FirstPage:  pageNo == 1,
		LastPage:   pageNo == int(tp),
		RangePrev:  rPrev,
		Range:      ran,
		RangeNext:  rNext,
	}
}
