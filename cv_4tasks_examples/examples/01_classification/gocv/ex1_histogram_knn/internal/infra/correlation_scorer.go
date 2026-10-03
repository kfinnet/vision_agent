package infra

import "math"

import "ex1_histogram_knn/internal/domain"

// CorrelationScorer는 domain.SimilarityScorer를 피어슨 상관계수(Correlation) 방식으로
// 구현하는 어댑터다. OpenCV의 gocv.CompareHist(..., HistCmpCorrel)과 동일한 수식이며,
// 1.0에 가까울수록 두 히스토그램(특징 벡터)이 유사함을 의미한다.
//
// LSP: Score()의 시그니처와 동작이 domain.SimilarityScorer 계약을 그대로 만족하므로,
// SimilarityScorer를 기대하는 어떤 코드에도 안전하게 대입할 수 있다.
type CorrelationScorer struct{}

// NewCorrelationScorer는 CorrelationScorer를 생성한다.
func NewCorrelationScorer() *CorrelationScorer {
	return &CorrelationScorer{}
}

// Score는 두 특징 벡터 a, b 사이의 피어슨 상관계수를 계산해 반환한다 (값이 클수록 유사).
func (s *CorrelationScorer) Score(a, b domain.Feature) float32 {
	n := len(a)
	if n == 0 || len(b) != n {
		return -1.0
	}
	var sumA, sumB float64
	for i := 0; i < n; i++ {
		sumA += float64(a[i])
		sumB += float64(b[i])
	}
	meanA, meanB := sumA/float64(n), sumB/float64(n)

	var num, denA, denB float64
	for i := 0; i < n; i++ {
		da := float64(a[i]) - meanA
		db := float64(b[i]) - meanB
		num += da * db
		denA += da * da
		denB += db * db
	}
	if denA == 0 || denB == 0 {
		return 0
	}
	return float32(num / math.Sqrt(denA*denB))
}

// 컴파일 타임에 CorrelationScorer가 domain.SimilarityScorer 포트를 만족하는지 확인한다.
var _ domain.SimilarityScorer = (*CorrelationScorer)(nil)
