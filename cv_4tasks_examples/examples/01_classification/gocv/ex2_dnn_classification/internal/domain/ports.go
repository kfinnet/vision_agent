package domain

// ImageClassifier는 입력 이미지에 대해 클래스별 점수(소프트맥스 점수 등)를 산출하는 포트다.
// 사전학습된 딥러닝 분류 모델의 추론을 추상화한다.
//
// SRP: "이미지를 모델에 넣어 점수를 얻는다"는 책임만 가지며, 상위 K개 선정이나
//
//	레이블 이름 조회는 담당하지 않는다.
type ImageClassifier interface {
	// ClassifyScores는 이미지 파일 경로를 입력받아 클래스 인덱스 순서의 점수 배열을 반환한다.
	ClassifyScores(imagePath string) ([]float32, error)
}

// LabelRepository는 클래스 인덱스에 대응하는 사람이 읽을 수 있는 레이블 이름을 제공하는 포트다.
//
// ISP: 레이블 조회라는 책임만 가지며, 모델 추론 책임(ImageClassifier)과는 분리된다.
// OCP: 레이블을 텍스트 파일이 아닌 다른 소스(DB, JSON 등)에서 가져오고 싶다면
//
//	이 인터페이스를 구현하는 새 어댑터만 추가하면 된다.
type LabelRepository interface {
	// Label은 클래스 인덱스에 해당하는 이름을 반환한다.
	Label(index int) string
}
