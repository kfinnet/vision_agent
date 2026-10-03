package domain

// FeatureExtractor는 이미지 파일로부터 분류용 특징 벡터를 추출하는 포트(인터페이스)다.
//
// SRP: "특징을 어떻게 추출할지"만 책임지며, 유사도 계산이나 분류 판단은 담당하지 않는다.
// ISP: 다른 책임과 섞이지 않은 작은 인터페이스로 유지한다.
// 구체적인 구현(색상 히스토그램, HOG 등)은 인프라 계층(infra 패키지)에서만 담당한다.
type FeatureExtractor interface {
	// Extract는 이미지 파일 경로를 입력받아 특징 벡터(Feature)를 반환한다.
	Extract(imagePath string) (Feature, error)
}

// SimilarityScorer는 두 특징 벡터 사이의 유사도를 계산하는 포트다. 값이 클수록 더 유사함을 의미한다.
//
// OCP: Correlation 외의 다른 유사도 계산 방식(예: 교차검증, 카이제곱 등)이 필요하면
//
//	이 인터페이스를 구현하는 새 타입을 infra에 추가하면 되고, 유스케이스는 바뀌지 않는다.
type SimilarityScorer interface {
	// Score는 두 특징 벡터 a, b 사이의 유사도를 반환한다.
	Score(a, b Feature) float32
}
