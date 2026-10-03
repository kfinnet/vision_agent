// Package domain은 ①분류 실습2(사전학습 DNN 분류)의 핵심 데이터 구조를 정의한다.
// 도메인 계층(Domain Layer) — gocv 등 어떤 구체적인 구현 기술에도 의존하지 않는 순수 Go 타입만 담는다.
package domain

// ScoredClass는 분류 결과 중 한 클래스의 (레이블, 점수) 쌍을 나타내는 값 객체(Value Object)다.
type ScoredClass struct {
	Index int
	Label string
	Score float32
}
