// Package domain은 ④추적 실습1(TrackerMIL)의 도메인 계층이다.
// 외부 라이브러리(gocv)에 의존하지 않는 순수 값 타입/엔티티만 둔다.
package domain

import "image"

// TrackedObject는 한 프레임에서 추적된 대상 하나(박스 + 중심점)를 나타내는
// 값 객체다. 어떤 추적 알고리즘(TrackerMIL, CSRT, KCF, ...)을 쓰든 동일한
// 형태로 결과를 표현할 수 있도록 gocv 타입이 아닌 표준 image 패키지 타입만
// 사용한다(DIP: 도메인은 특정 CV 라이브러리를 몰라야 한다).
type TrackedObject struct {
	Box    image.Rectangle
	Center image.Point
	Found  bool // 이번 프레임에서 추적에 성공했는지 여부
}

// NewTrackedObject는 사각형 박스로부터 중심점을 계산해 TrackedObject를 만든다.
func NewTrackedObject(box image.Rectangle, found bool) TrackedObject {
	center := image.Pt(box.Min.X+box.Dx()/2, box.Min.Y+box.Dy()/2)
	return TrackedObject{Box: box, Center: center, Found: found}
}

// Trail은 추적 과정에서 누적된 중심점들의 이동 경로다.
type Trail struct {
	Points []image.Point
}

// Add는 새 중심점을 경로에 추가한다.
func (t *Trail) Add(p image.Point) {
	t.Points = append(t.Points, p)
}
