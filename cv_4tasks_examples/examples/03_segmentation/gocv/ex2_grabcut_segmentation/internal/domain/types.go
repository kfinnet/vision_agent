// Package domain은 ===== [도메인 계층 / Domain Layer] =====.
// "전경/배경 분할이란 무엇인가"만 정의하는 계층으로, gocv 등 구체적인
// CV 라이브러리에는 전혀 의존하지 않는다. 다른 모든 계층이 참조하는
// 최하위(가장 안정적인) 계층이며, 그 반대 방향으로는 의존하지 않는다(DIP).
package domain

// Rect는 분할 초기화에 사용하는 사각형 힌트(x, y, width, height)를 표현한다.
// gocv의 image.Rectangle에 의존하지 않도록 domain 전용 타입으로 둔다.
type Rect struct {
	X, Y, W, H int
}

// SegmentationResult는 전경/배경 분할 1회 실행의 결과를 표현하는 값 타입이다.
//
//   - ForegroundMask: 전경으로 분류된 픽셀을 255, 배경을 0으로 담은 평면 배열.
//   - Rows, Cols: 마스크 크기.
//   - ForegroundRatio: 전체 픽셀 중 전경으로 분류된 비율(0.0~1.0).
//   - InitRect: 분할 초기화에 사용한 사각형 힌트.
type SegmentationResult struct {
	ForegroundMask  []uint8
	Rows            int
	Cols            int
	ForegroundRatio float64
	InitRect        Rect
}

// At는 (y, x) 위치의 마스크 값을 반환한다.
func (r SegmentationResult) At(y, x int) uint8 {
	return r.ForegroundMask[y*r.Cols+x]
}
