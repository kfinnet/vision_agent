// Package usecase는 ①분류 실습1의 "최근접 이웃 분류"라는 응용 비즈니스 로직을 담당한다.
// 유스케이스 계층(Use Case Layer) — domain 패키지의 추상 인터페이스에만 의존하며(DIP),
// gocv를 직접 import하거나 호출하지 않는다.
package usecase

import "ex1_histogram_knn/internal/domain"

// NearestNeighborUseCase는 "테스트 이미지와 가장 유사한 기준 이미지의 레이블을 찾는다"는
// 분류 유스케이스를 담당한다 (SRP). 구체적인 특징 추출/유사도 계산 알고리즘은
// 생성자 주입(Constructor Injection)으로 전달받는다 — 어떤 구현체를 쓸지는 이 구조체가
// 아니라 구성 루트(main 패키지)가 결정한다.
type NearestNeighborUseCase struct {
	extractor domain.FeatureExtractor
	scorer    domain.SimilarityScorer
}

// NewNearestNeighborUseCase는 FeatureExtractor와 SimilarityScorer를 주입받아
// NearestNeighborUseCase를 생성하는 생성자다.
func NewNearestNeighborUseCase(extractor domain.FeatureExtractor, scorer domain.SimilarityScorer) *NearestNeighborUseCase {
	return &NearestNeighborUseCase{extractor: extractor, scorer: scorer}
}

// BuildLabeledFeature는 레이블이 붙은 기준 이미지 1장의 특징을 추출해 LabeledFeature로 감싼다.
func (u *NearestNeighborUseCase) BuildLabeledFeature(label, imagePath string) (domain.LabeledFeature, error) {
	feature, err := u.extractor.Extract(imagePath)
	if err != nil {
		return domain.LabeledFeature{}, err
	}
	return domain.LabeledFeature{Label: label, Feature: feature}, nil
}

// Classify는 테스트 이미지의 특징을 추출한 뒤, 기준 특징 목록(refs) 중 가장 유사한 것을 찾아
// 최근접 이웃(k=1) 분류 결과를 반환한다.
func (u *NearestNeighborUseCase) Classify(refs []domain.LabeledFeature, testImagePath string) (domain.ClassificationResult, error) {
	testFeature, err := u.extractor.Extract(testImagePath)
	if err != nil {
		return domain.ClassificationResult{}, err
	}

	result := domain.ClassificationResult{Scores: make(map[string]float32, len(refs))}
	best := float32(-1.0)
	for _, r := range refs {
		score := u.scorer.Score(testFeature, r.Feature)
		result.Scores[r.Label] = score
		if score > best {
			best = score
			result.PredictedLabel = r.Label
			result.BestScore = score
		}
	}
	return result, nil
}
