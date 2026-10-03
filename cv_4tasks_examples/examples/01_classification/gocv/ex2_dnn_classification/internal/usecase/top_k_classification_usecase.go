// Package usecase는 ①분류 실습2의 "이미지를 분류해 상위 K개 클래스를 뽑는다"는 응용
// 비즈니스 로직을 담당한다. 유스케이스 계층(Use Case Layer) — domain 패키지의 추상
// 인터페이스에만 의존하며(DIP), gocv/DNN 관련 코드를 직접 다루지 않는다.
package usecase

import (
	"sort"

	"ex2_dnn_classification/internal/domain"
)

// TopKClassificationUseCase는 ImageClassifier로 점수를 얻고, LabelRepository로 레이블
// 이름을 붙인 뒤, 점수가 높은 순으로 상위 K개를 골라내는 유스케이스를 담당한다 (SRP).
// 구체적인 모델/레이블 조회 구현은 생성자 주입(Constructor Injection)으로 전달받는다.
type TopKClassificationUseCase struct {
	classifier domain.ImageClassifier
	labels     domain.LabelRepository
}

// NewTopKClassificationUseCase는 ImageClassifier와 LabelRepository를 주입받아
// TopKClassificationUseCase를 생성하는 생성자다.
func NewTopKClassificationUseCase(classifier domain.ImageClassifier, labels domain.LabelRepository) *TopKClassificationUseCase {
	return &TopKClassificationUseCase{classifier: classifier, labels: labels}
}

// Execute는 이미지를 분류하고, 점수가 높은 순서로 상위 k개의 ScoredClass를 반환한다.
func (u *TopKClassificationUseCase) Execute(imagePath string, k int) ([]domain.ScoredClass, error) {
	scores, err := u.classifier.ClassifyScores(imagePath)
	if err != nil {
		return nil, err
	}

	ranked := make([]domain.ScoredClass, len(scores))
	for i, s := range scores {
		ranked[i] = domain.ScoredClass{Index: i, Label: u.labels.Label(i), Score: s}
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].Score > ranked[j].Score })

	if k > len(ranked) {
		k = len(ranked)
	}
	return ranked[:k], nil
}
