// Package domain은 클린 아키텍처의 도메인 계층이다.
// 이 패키지는 어떤 구체 기술(gocv 등)에도 의존하지 않는 순수한 값 타입만
// 정의한다 — "탐지 결과"라는 개념을 애플리케이션 전반에서 공통으로 쓰기
// 위한 최소 단위다.
package domain

// Detection은 탐지된 객체 하나를 나타내는 값 타입이다.
// X, Y, W, H는 바운딩 박스(좌상단 좌표 + 너비/높이), Label은 객체 종류
// (예: "face"), Confidence는 탐지 신뢰도(0~1)를 뜻한다. Haar Cascade처럼
// 신뢰도를 제공하지 않는 알고리즘은 1.0으로 채워 다른 구현체와 동일한
// 형태를 유지한다(LSP: 모든 탐지기 구현이 동일하게 취급될 수 있어야 한다).
type Detection struct {
	X          int
	Y          int
	W          int
	H          int
	Label      string
	Confidence float64
}
