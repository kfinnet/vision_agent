# ① 분류 (Classification) 실습

이미지 전체를 보고 "이게 무엇인가" 하나의 레이블을 고르는 가장 기본적인 CV 태스크입니다.

## Colab/OpenCV (Python)

### 실습1 — `ex1_color_histogram_knn.ipynb`
색상(HSV Hue·Saturation) 히스토그램을 특징으로 사용해, k-최근접이웃(k-NN)으로
고양이/사람/커피잔 3개 클래스를 분류합니다. 가장 가볍고 직관적인 고전적 접근입니다.
`scikit-image`의 표준 샘플 이미지에 밝기·플립·크롭 변형(augmentation)을 가해
학습/테스트셋을 만들기 때문에 별도 다운로드 없이 어디서나 동일하게 재현됩니다.
(검증 결과: 테스트 정확도 100%)

### 실습2 — `ex2_hog_svm.ipynb`
HOG(Histogram of Oriented Gradients, 형태/윤곽 기반) 특징 + 선형 SVM으로
고양이/사람/로켓 3개 클래스를 분류합니다. 색이 아니라 "모양"을 보기 때문에
조명·색상 변화에 더 강인합니다. 실습1과 같은 증강 방식을 사용합니다.
(검증 결과: 테스트 정확도 100%)

두 실습 모두 Colab에서 `.ipynb` 파일을 업로드해 셀을 순서대로 실행하면 됩니다.

## GoCV (Go)

### 실습1 — `gocv/ex1_histogram_knn/`
Python 실습1과 동일한 "색상 히스토그램" 특징을 OpenCV의 `CompareHist`로 비교해
가장 유사한 기준 이미지의 레이블을 선택하는 최근접이웃(k=1) 분류기입니다.

```bash
cd gocv/ex1_histogram_knn
go get gocv.io/x/gocv
go run main.go ../../../assets/cat.jpg=cat ../../../assets/astronaut.jpg=person ../../../assets/coffee.jpg=coffee -- ../../../assets/rocket.jpg
```

### 실습2 — `gocv/ex2_dnn_classification/`
사전학습된 딥러닝 분류 모델(ONNX)을 OpenCV DNN 모듈로 구동해 추론하는 실무형 예제입니다.
GoCV에는 scikit-learn 같은 전통 ML 라이브러리가 없으므로, Go 환경에서 분류를 실제로
서비스에 적용할 때 가장 흔히 쓰이는 "사전학습 모델 추론" 방식을 보여줍니다.

```bash
cd gocv/ex2_dnn_classification
go get gocv.io/x/gocv
# 모델/레이블 파일은 저작권 문제로 포함하지 않음 — 코드 상단 주석의 다운로드 링크 참고
go run main.go squeezenet.onnx classification_classes_ILSVRC2012.txt test.jpg
```

> GoCV 2개는 `gofmt`/스텁 기반 `go build`·`go vet`까지 확인했으나, 실제 `gocv`
> 라이브러리로 완전한 컴파일·실행 검증은 하지 못했습니다 — 상위 README의
> "검증 현황"을 참고하세요.

## 코드 구조 — 클린 아키텍처 & SOLID 원칙 적용

모든 예제는 알고리즘 로직(도메인 지식)과 cv2/gocv 같은 구체적인 구현 기술을 분리하는
클린 아키텍처(Clean Architecture)로 재구성되어 있습니다.

### Python (`colab_opencv/*_code.py`)

하나의 `.py` 파일 안에서 아래 4개 계층을 주석 배너로 명확히 구분합니다.

1. **도메인 계층** (`# ===== [도메인 계층 / Domain Layer] =====`) — `@dataclass` 값 객체와
   `abc.ABC` 추상 인터페이스(`FeatureExtractor`, `Classifier`)만 정의. cv2/sklearn에 의존하지 않습니다.
2. **유스케이스 계층** (`# ===== [유스케이스 계층 / Use Case Layer] =====`) — `ClassificationUseCase`가
   도메인 인터페이스만 생성자로 주입받아 학습/분류 흐름을 조립합니다. cv2/sklearn을 직접 호출하지 않습니다.
3. **인프라 계층** (`# ===== [인프라/어댑터 계층 / Infrastructure Layer] =====`) — `HSVHistogramExtractor`,
   `KNNClassifierAdapter`(실습1) / `HOGFeatureExtractor`, `LinearSVMClassifierAdapter`(실습2)처럼
   cv2·sklearn을 실제로 호출하는 구현체만 둡니다.
4. **구성 루트** (`# ===== [구성 루트 / Composition Root — Colab 실행 코드] =====`) — 샘플 준비,
   데이터 증강, 구체 어댑터 생성 및 의존성 주입(DI), 학습/테스트 실행, 시각화·파일 저장을 담당합니다.

### GoCV (`gocv/*/`)

각 예제 폴더를 `main.go` + `internal/{domain,usecase,infra}` 패키지로 분리합니다.

- `internal/domain` — `ports.go`(인터페이스)와 `types.go`(값 타입). 외부 라이브러리 의존 없음.
- `internal/usecase` — 도메인 인터페이스에만 의존하는 오케스트레이션 로직 (`NearestNeighborUseCase`,
  `TopKClassificationUseCase`). gocv를 import하지 않습니다.
- `internal/infra` — `gocv.io/x/gocv`를 import하는 유일한 패키지. 도메인 인터페이스의 실제 구현체
  (`HSVHistogramExtractor`, `CorrelationScorer`, `ONNXImageClassifier`, `TextFileLabelRepository`)를 둡니다.
- `main.go` — 구성 루트. CLI 인자 파싱, 구체 어댑터 생성, 유스케이스에 생성자 주입(DI), 결과 출력만 수행합니다.

### SOLID 원칙 적용 방식

- **SRP (단일 책임)**: 특징 추출, 분류 판단, 시각화/출력이 각각 별도 클래스·파일로 분리되어 있습니다.
- **OCP (개방-폐쇄)**: 새 알고리즘(예: 다른 특징 추출기)은 인프라 계층에 새 클래스를 추가하는 것만으로
  확장되며, 유스케이스·도메인 코드는 수정하지 않습니다.
- **LSP (리스코프 치환)**: 모든 어댑터는 자신이 구현하는 인터페이스의 입출력 규약을 그대로 지키므로
  서로 안전하게 치환할 수 있습니다.
- **ISP (인터페이스 분리)**: `FeatureExtractor`/`Classifier`(Python), `FeatureExtractor`/`SimilarityScorer`/
  `ImageClassifier`/`LabelRepository`(Go)처럼 책임별로 작은 인터페이스를 유지합니다.
- **DIP (의존성 역전)**: 유스케이스는 구체 구현(cv2, sklearn, gocv)이 아닌 추상 인터페이스에만 의존하며,
  실제 구현체는 구성 루트에서 생성자 주입으로 연결됩니다.
