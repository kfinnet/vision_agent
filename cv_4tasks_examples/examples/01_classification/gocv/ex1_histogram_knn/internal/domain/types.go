// Package domain은 ①분류 실습1(색상 히스토그램 + 최근접 이웃)의 핵심 데이터 구조를 정의한다.
// 도메인 계층(Domain Layer) — gocv 등 어떤 구체적인 구현 기술에도 의존하지 않는 순수 Go 타입만 담는다.
package domain

// Feature는 특징 추출 결과를 표현하는 1차원 벡터다. 색상 히스토그램, HOG 등 어떤 알고리즘으로
// 추출되었는지와 무관한 순수 Go 타입으로, 인프라(gocv) 구현 세부사항을 전혀 노출하지 않는다.
type Feature []float32

// LabeledFeature는 레이블이 붙은 기준(reference) 이미지의 특징을 나타내는 엔티티다.
type LabeledFeature struct {
	Label   string
	Feature Feature
}

// ClassificationResult는 최근접 이웃 분류 결과를 표현하는 값 객체(Value Object)다.
// 예측 레이블과 함께, 각 기준 클래스와의 유사도 점수도 보조 정보로 담는다.
type ClassificationResult struct {
	PredictedLabel string
	BestScore      float32
	Scores         map[string]float32
}
