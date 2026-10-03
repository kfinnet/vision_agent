// Package infra는 ①분류 실습2의 도메인 인터페이스를 gocv(OpenCV DNN) 및 파일 시스템으로
// 구현하는 어댑터들을 담는다. 인프라/어댑터 계층(Infrastructure Layer) — 이 패키지만
// "gocv.io/x/gocv"를 import하며, 실제 OpenCV DNN API 호출은 모두 여기에 있다.
package infra

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"

	"ex2_dnn_classification/internal/domain"
)

// ONNXImageClassifier는 domain.ImageClassifier를 OpenCV DNN 모듈(gocv.ReadNet)의 ONNX
// 추론으로 구현하는 어댑터다. GoCV는 scikit-learn 같은 전통 ML 라이브러리가 없으므로,
// 실무에서 흔히 쓰이는 "사전학습 모델 추론" 방식을 보여준다.
//
// OCP: TensorFlow/Caffe 등 다른 포맷의 모델이 필요하면 ImageClassifier를 구현하는
//
//	새 어댑터를 추가하면 되고, 유스케이스는 전혀 바뀌지 않는다.
//
// LSP: ClassifyScores()의 입출력 규약이 domain.ImageClassifier와 동일하므로
//
//	ImageClassifier를 기대하는 어떤 코드에도 안전하게 대입할 수 있다.
type ONNXImageClassifier struct {
	net gocv.Net
}

// NewONNXImageClassifier는 ONNX 분류 모델 파일을 읽어 ONNXImageClassifier를 생성한다.
//
// 사전 준비 (필수 — 모델 파일은 저작권/용량 문제로 본 저장소에 포함하지 않음):
//  1. ONNX 형식의 이미지 분류 모델 (예: SqueezeNet, MobileNetV2 등)
//     예시 다운로드: https://github.com/onnx/models/tree/main/validated/vision/classification/squeezenet
//  2. ImageNet 클래스 레이블 텍스트 파일 (줄바꿈으로 구분된 1000개 클래스명)
//     예시: https://raw.githubusercontent.com/opencv/opencv/master/samples/data/dnn/classification_classes_ILSVRC2012.txt
func NewONNXImageClassifier(modelPath string) (*ONNXImageClassifier, error) {
	net := gocv.ReadNet(modelPath, "")
	if net.Empty() {
		return nil, fmt.Errorf("모델을 불러오지 못했습니다: %s", modelPath)
	}
	net.SetPreferableBackend(gocv.NetBackendDefault)
	net.SetPreferableTarget(gocv.NetTargetCPU)
	return &ONNXImageClassifier{net: net}, nil
}

// Close는 내부에 보유한 Net 리소스를 해제한다.
func (c *ONNXImageClassifier) Close() error {
	return c.net.Close()
}

// ClassifyScores는 이미지를 읽어 전처리(blob 변환) 후 순전파(forward)를 실행하고,
// 클래스별 점수(소프트맥스 등) 배열을 반환한다.
func (c *ONNXImageClassifier) ClassifyScores(imagePath string) ([]float32, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("이미지를 읽을 수 없습니다: %s", imagePath)
	}
	defer img.Close()

	// 대부분의 ImageNet 분류 모델 전처리 규약: 224x224 리사이즈, BGR 평균값 보정, 1/255 스케일
	blob := gocv.BlobFromImage(img, 1.0/255.0, image.Pt(224, 224),
		gocv.NewScalar(0.485*255, 0.456*255, 0.406*255, 0), true, false)
	defer blob.Close()

	c.net.SetInput(blob, "")
	probMat := c.net.Forward("")
	defer probMat.Close()

	data, err := probMat.DataPtrFloat32()
	if err != nil {
		return nil, err
	}
	scores := make([]float32, len(data))
	copy(scores, data)
	return scores, nil
}

// 컴파일 타임에 ONNXImageClassifier가 domain.ImageClassifier 포트를 만족하는지 확인한다.
var _ domain.ImageClassifier = (*ONNXImageClassifier)(nil)
