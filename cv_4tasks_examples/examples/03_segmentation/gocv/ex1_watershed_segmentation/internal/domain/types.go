// Package domain은 ===== [도메인 계층 / Domain Layer] =====.
// "분할(Segmentation)이란 무엇인가"만 정의하는 계층으로, gocv 등 구체적인
// CV 라이브러리에는 전혀 의존하지 않는다. 다른 모든 계층이 참조하는
// 최하위(가장 안정적인) 계층이며, 그 반대 방향으로는 의존하지 않는다(DIP).
package domain

// SegmentationResult는 분할 알고리즘 1회 실행의 결과를 표현하는 값 타입이다.
//
//   - Markers: 각 픽셀이 속한 세그먼트 번호를 담은 라벨 맵(행 우선, Rows*Cols).
//     배경=1, 경계선=-1, 그 외 2 이상의 정수가 개별 물체(동전) 번호.
//   - Rows, Cols: 라벨 맵의 크기.
//   - NumSegments: 배경(1)을 제외하고 최종적으로 분리된 물체 개수.
type SegmentationResult struct {
	Markers     []int32
	Rows        int
	Cols        int
	NumSegments int
}

// At는 (y, x) 위치의 라벨 값을 반환한다.
func (r SegmentationResult) At(y, x int) int32 {
	return r.Markers[y*r.Cols+x]
}
