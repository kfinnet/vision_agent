// CV 4대 태스크 실습 — ①분류(Classification) 실습2 (GoCV)
// 사전학습 딥러닝 모델(DNN) 기반 이미지 분류
//
// Python/Colab 실습2(ex2_hog_svm.ipynb)는 HOG+SVM(고전 머신러닝)으로 분류했다면,
// 이 GoCV 실습은 실무에서 더 흔히 쓰이는 방식 — "사전학습된 딥러닝 분류 모델을 불러와
// 추론만 수행"하는 방식을 보여준다. GoCV는 scikit-learn 같은 ML 라이브러리가 없으므로,
// OpenCV의 DNN 모듈(gocv.ReadNet)로 ONNX 분류 모델을 직접 구동한다.
//
// 클린 아키텍처 계층 구성:
//
//	internal/domain  — 분류/레이블 조회 포트(인터페이스)와 값 타입 (gocv 의존 없음)
//	internal/usecase — 상위 K개 분류 결과 선정 로직 (domain 인터페이스에만 의존, DIP)
//	internal/infra   — gocv DNN 기반 ONNX 추론기 / 텍스트 레이블 저장소 (gocv를 import하는 유일한 패키지)
//	main (이 파일)    — 구성 루트: CLI 인자 해석 + 어댑터 생성 + 유스케이스 주입(DI) + 결과 출력
//
// 사전 준비 (필수 — 모델 파일은 저작권/용량 문제로 본 저장소에 포함하지 않음):
//  1. ONNX 형식의 이미지 분류 모델 (예: SqueezeNet, MobileNetV2 등)
//     예시 다운로드: https://github.com/onnx/models/tree/main/validated/vision/classification/squeezenet
//  2. ImageNet 클래스 레이블 텍스트 파일 (줄바꿈으로 구분된 1000개 클래스명)
//     예시: https://raw.githubusercontent.com/opencv/opencv/master/samples/data/dnn/classification_classes_ILSVRC2012.txt
//
// 사용법:
//
//	go run main.go <모델.onnx> <레이블.txt> <입력이미지>
package main

import (
	"fmt"
	"os"

	"ex2_dnn_classification/internal/infra"
	"ex2_dnn_classification/internal/usecase"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("사용법: go run main.go <모델.onnx> <레이블.txt> <입력이미지>")
		os.Exit(1)
	}
	modelPath, labelPath, imgPath := os.Args[1], os.Args[2], os.Args[3]

	// 구성 루트(Composition Root)에서만 구체 어댑터를 생성하고 유스케이스에 주입한다.
	labels, err := infra.NewTextFileLabelRepository(labelPath)
	if err != nil {
		fmt.Println("레이블 파일을 읽을 수 없습니다:", err)
		os.Exit(1)
	}

	classifier, err := infra.NewONNXImageClassifier(modelPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer classifier.Close()

	uc := usecase.NewTopKClassificationUseCase(classifier, labels)
	top5, err := uc.Execute(imgPath, 5)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("분류 결과 (상위 5개):")
	for i, sc := range top5 {
		fmt.Printf("  %d위. %-30s (점수 %.4f)\n", i+1, sc.Label, sc.Score)
	}
}
