// CV 4대 태스크 실습 — ①분류(Classification) 실습1 (GoCV)
// 색상 히스토그램(Color Histogram) + 최근접 이웃(Nearest-Neighbor) 분류
//
// Python/Colab 실습1(ex1_color_histogram_knn.ipynb)과 동일한 개념을 순수 Go + GoCV로 구현합니다.
// scikit-learn의 k-NN 대신, OpenCV의 CompareHist로 "가장 가까운 기준 이미지"를 직접 찾는
// 가장 단순한 형태(k=1 최근접 이웃)로 구현합니다.
//
// 클린 아키텍처 계층 구성:
//
//	internal/domain  — 특징/유사도 포트(인터페이스)와 값 타입 (gocv 의존 없음)
//	internal/usecase — 최근접 이웃 분류 로직 (domain 인터페이스에만 의존, DIP)
//	internal/infra   — gocv 기반 히스토그램 추출기 / 상관계수 계산기 (gocv를 import하는 유일한 패키지)
//	main (이 파일)    — 구성 루트: CLI 인자 해석 + 어댑터 생성 + 유스케이스 주입(DI) + 결과 출력
//
// 사용법:
//
//	go run main.go <기준이미지1=레이블1> <기준이미지2=레이블2> ... -- <테스트이미지>
//
// 예시:
//
//	go run main.go cat.jpg=cat person.jpg=person coffee.jpg=coffee -- test.jpg
package main

import (
	"fmt"
	"os"
	"strings"

	"ex1_histogram_knn/internal/domain"
	"ex1_histogram_knn/internal/infra"
	"ex1_histogram_knn/internal/usecase"
)

func main() {
	args := os.Args[1:]
	sepIdx := -1
	for i, a := range args {
		if a == "--" {
			sepIdx = i
			break
		}
	}
	if sepIdx == -1 || sepIdx == len(args)-1 {
		fmt.Println("사용법: go run main.go <기준이미지1>=<레이블1> ... -- <테스트이미지>")
		fmt.Println("예시:   go run main.go cat.jpg=cat person.jpg=person coffee.jpg=coffee -- test.jpg")
		os.Exit(1)
	}
	refArgs := args[:sepIdx]
	testPath := args[sepIdx+1]

	// 구성 루트(Composition Root)에서만 구체 어댑터를 생성하고 유스케이스에 주입한다.
	extractor := infra.NewHSVHistogramExtractor(32)
	scorer := infra.NewCorrelationScorer()
	uc := usecase.NewNearestNeighborUseCase(extractor, scorer)

	// 1) 기준(레이블별) 이미지들의 특징(히스토그램) 계산
	refs := make([]domain.LabeledFeature, 0, len(refArgs))
	for _, ra := range refArgs {
		parts := strings.SplitN(ra, "=", 2)
		if len(parts) != 2 {
			fmt.Println("형식 오류 (경로=레이블이어야 함):", ra)
			os.Exit(1)
		}
		path, label := parts[0], parts[1]
		lf, err := uc.BuildLabeledFeature(label, path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		refs = append(refs, lf)
	}

	// 2) 테스트 이미지 특징 계산 + 3) 최근접 이웃(k=1) 분류
	result, err := uc.Classify(refs, testPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("각 클래스와의 유사도(Correlation, 1.0에 가까울수록 유사):")
	for _, r := range refs {
		fmt.Printf("  - %-10s : %.4f\n", r.Label, result.Scores[r.Label])
	}
	fmt.Printf("\n분류 결과: \"%s\" (유사도 %.4f)\n", result.PredictedLabel, result.BestScore)
}
