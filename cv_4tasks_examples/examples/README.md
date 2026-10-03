# 컴퓨터비전(CV) 4대 핵심 태스크 실습예제

"컴퓨터비전(CV)의 4대 핵심 태스크 — ①분류 ②탐지 ③분할 ④추적" 강의자료(PPT)에 대응하는
실습 소스코드 모음입니다. 각 태스크마다 **Colab/OpenCV(Python) 2개 + GoCV(Go) 2개**,
총 **16개 실습예제**로 구성했습니다.

## 폴더 구조

```
examples/
├── README.md                              (본 파일)
├── assets/                                 공용 테스트 이미지·영상
│   ├── vtest.avi      — 보행자 CCTV 영상(OpenCV 공식 샘플, 탐지/추적 실습에 사용)
│   ├── coins.png       — 동전 사진(분할 실습 Watershed용)
│   ├── cat.jpg          — 고양이 사진(분류·분할 GrabCut 실습용)
│   ├── astronaut.jpg, coffee.jpg, rocket.jpg — 분류 실습용 보조 이미지
│
├── verified_outputs/                       Colab/OpenCV 8개를 실제 실행해 얻은 결과 이미지
│                                            (아래 "검증 현황" 참고 — 코드가 실제로 동작함을 보여주는 증빙)
│
├── 01_classification/                      ① 분류 (Classification)
│   ├── colab_opencv/
│   │   ├── ex1_color_histogram_knn.ipynb   실습1: 색상 히스토그램 + k-NN
│   │   └── ex2_hog_svm.ipynb               실습2: HOG 특징 + 선형 SVM
│   └── gocv/
│       ├── ex1_histogram_knn/              실습1: 히스토그램 + 최근접이웃(Go)
│       └── ex2_dnn_classification/         실습2: 사전학습 DNN 분류(Go)
│
├── 02_detection/                            ② 탐지 (Detection)
│   ├── colab_opencv/
│   │   ├── ex1_haar_face_detection.ipynb   실습1: Haar Cascade 얼굴 탐지
│   │   └── ex2_hog_pedestrian_detection.ipynb 실습2: HOG 보행자 탐지
│   └── gocv/
│       ├── ex1_haar_face_detection/        실습1: Haar Cascade 얼굴 탐지(Go, 실시간)
│       └── ex2_hog_pedestrian_detection/   실습2: HOG 보행자 탐지(Go)
│
├── 03_segmentation/                          ③ 분할 (Segmentation)
│   ├── colab_opencv/
│   │   ├── ex1_watershed_segmentation.ipynb 실습1: Watershed 동전 분할
│   │   └── ex2_grabcut_segmentation.ipynb   실습2: GrabCut 전경/배경 분할
│   └── gocv/
│       ├── ex1_watershed_segmentation/      실습1: Watershed(Go)
│       └── ex2_grabcut_segmentation/        실습2: GrabCut(Go)
│
└── 04_tracking/                              ④ 추적 (Tracking)
    ├── colab_opencv/
    │   ├── ex1_meanshift_tracking.ipynb     실습1: MeanShift 추적
    │   └── ex2_optical_flow_tracking.ipynb  실습2: Lucas-Kanade Optical Flow
    └── gocv/
        ├── ex1_tracker_mil/                 실습1: 내장 Tracker(MIL)(Go)
        └── ex2_bgsub_centroid_tracking/     실습2: 배경차분+중심점매칭 추적(Go)
```

## 태스크 ↔ 실습 매핑

| 태스크 | Colab/OpenCV (Python) | GoCV (Go) |
|---|---|---|
| ① 분류 | 색상 히스토그램+k-NN / HOG+SVM | 히스토그램 최근접이웃 / 사전학습 DNN |
| ② 탐지 | Haar Cascade 얼굴 / HOG 보행자 | Haar Cascade 얼굴(실시간) / HOG 보행자 |
| ③ 분할 | Watershed(동전) / GrabCut(전경추출) | Watershed(Go) / GrabCut(Go) |
| ④ 추적 | MeanShift / Optical Flow(Lucas-Kanade) | 내장 Tracker(MIL) / 배경차분+ID매칭 |

## 검증 현황 (투명하게 안내)

- **Colab/OpenCV 8개 전부**: 이 코드를 작성한 환경에 `opencv-python-headless`,
  `scikit-learn`, `scikit-image`, `scipy`를 실제로 설치하여 **end-to-end로 직접 실행**하고
  결과 이미지를 육안으로 확인했습니다.
  - 분류 실습1(히스토그램+k-NN): 테스트 정확도 92~100% 확인 (실행마다 랜덤 증강이 달라
    소폭 변동 — `verified_outputs/`에 실제 실행 결과 포함)
  - 분류 실습2(HOG+SVM): 테스트 정확도 100% 확인
  - 탐지 실습1(얼굴): 샘플 이미지에서 얼굴 정상 탐지 확인
  - 탐지 실습2(보행자): vtest.avi 프레임에서 보행자 4명 탐지 확인
  - 분할 실습1(Watershed): 동전 24~26개로 정상 분할 확인(실제 동전 수 24개)
  - 분할 실습2(GrabCut): 고양이 전경 분리 결과 육안 확인
  - 추적 실습1(MeanShift): 40프레임에 걸친 이동 경로(trajectory) 정상 생성 확인
  - 추적 실습2(Optical Flow): 다중 특징점의 이동 궤적 정상 생성 확인
  - 노트북(.ipynb) 파일 8개 모두 JSON 형식 유효성 검증 완료
- **GoCV 8개**: 이 코드 작성 환경에는 OpenCV C++ 네이티브 라이브러리와 외부 Go 모듈
  네트워크(proxy.golang.org)가 없어 실제 `gocv` 패키지로 `go build`하는 완전한 컴파일
  검증은 하지 못했습니다. 대신:
  1. `gofmt -e`로 8개 파일 모두 Go 문법 유효성을 확인했습니다.
  2. 사용된 gocv API의 **실제 함수 시그니처를 흉내 낸 로컬 스텁(stub) 패키지**를 직접
     작성해 8개 파일 모두 `go build`·`go vet` 통과를 확인했습니다(변수명 오류, 타입
     불일치, 미사용 변수 등 실제 버그 1건을 이 과정에서 발견해 수정했습니다 —
     분류 실습1의 `float32`/`float64` 비교 오류).
  3. 다만 이 스텁은 실제 gocv 라이브러리의 동작을 재현하지 않으므로, **실제 OpenCV
     처리 결과(탐지 정확도 등)까지 검증된 것은 아닙니다.** 본인 환경에서
     `go build` 확인을 권장합니다.

## 아키텍처 — Clean Architecture + SOLID 원칙 적용

사용자 요청에 따라, 16개 예제 전부 "동작만 하는 코드"가 아니라 **계층을 분리하고
의존성이 추상화를 향하도록(DIP)** 재구성했습니다.

### Python (`colab_opencv/*_code.py`)

하나의 `.py` 파일(Colab 셀 변환 규칙 유지) 안에서 아래 4개 계층을 주석 배너로 명확히 구분합니다.

1. `# ===== [도메인 계층 / Domain Layer] =====` — `@dataclass` 결과 값 객체와
   `abc.ABC` 추상 인터페이스(예: `FeatureExtractor`, `Classifier`, `ObjectDetector`,
   `ImageSegmenter`, `ObjectTracker`)만 정의합니다. cv2/sklearn/scipy에 의존하지 않습니다.
2. `# ===== [유스케이스 계층 / Use Case Layer] =====` — 예: `ClassificationUseCase`,
   `DetectionUseCase`, `SegmentationUseCase`, `TrackingUseCase`. 도메인 인터페이스만
   생성자로 주입받아 흐름을 조립하며, cv2를 직접 호출하지 않습니다.
3. `# ===== [인프라/어댑터 계층 / Infrastructure Layer] =====` — cv2/sklearn/scipy를
   실제로 호출하는 구현체(`HSVHistogramExtractor`, `HOGPedestrianDetector`,
   `WatershedSegmenter`, `MeanShiftTracker` 등)만 둡니다. 계산과 시각화(색칠·배경
   제거·오버레이)는 서로 다른 클래스로 분리했습니다.
4. `# ===== [구성 루트 / Composition Root — Colab 실행 코드] =====` — 실제 Colab
   셀들: 구체 어댑터 생성 및 의존성 주입(DI), 실행, matplotlib 시각화, 파일 저장을
   담당합니다. 파일 I/O·`!pip` 설치가 허용되는 유일한 계층입니다.

### Go (`gocv/<example>/`)

각 예제 폴더를 아래처럼 계층별 패키지로 분리했습니다.

```
<example>/
  go.mod
  main.go              — 구성 루트: CLI 인자 파싱 + 의존성 조립(DI)만 수행
  internal/
    domain/
      types.go          — 값 타입(엔티티)
      ports.go           — 인터페이스(포트) — usecase가 의존하는 추상화
    usecase/
      *_usecase.go       — 생성자 주입 + Execute/Run만 포함, gocv 미사용
    infra/
      gocv_*.go           — gocv.io/x/gocv를 import하는 유일한 패키지
```

### SOLID 적용 방식 (공통)

- **SRP**: 특징 추출/탐지/분할/추적 "계산" 로직과 시각화·후처리·파일 입출력을
  각각 별도 클래스·구조체로 분리했습니다.
- **OCP**: 새 알고리즘(예: 다른 특징 추출기, 다른 탐지기)은 인프라 계층에 새
  클래스를 추가하는 것만으로 확장되며, 유스케이스·도메인 코드는 수정하지 않습니다.
- **LSP**: 같은 인터페이스를 구현하는 어댑터들은 서로 완전히 치환 가능하도록
  설계했습니다. 단, 추적(④) 태스크처럼 알고리즘 특성상(MeanShift=단일 영역 vs
  Optical Flow=다중 특징점) 억지로 하나의 인터페이스에 끼워 맞추면 LSP를 오히려
  해치는 경우, 책임별로 작은 인터페이스를 따로 두는 쪽을 선택했습니다(ISP 우선).
- **ISP**: 인터페이스를 책임 단위로 작게 유지합니다(예: 탐지와 NMS 후처리,
  분할 계산과 결과 시각화를 별도 인터페이스/클래스로 분리).
- **DIP**: 유스케이스·도메인 계층은 cv2/gocv 같은 구체 구현이 아니라 추상
  인터페이스에만 의존하며, 실제 구현체는 구성 루트(Composition Root)에서
  생성자 주입으로 연결됩니다.

각 태스크 폴더의 README.md에 더 구체적인 클래스/패키지 구성이 안내되어 있습니다.

## 공통 사전 준비 (GoCV)

```bash
# OS별 OpenCV C++ 라이브러리 설치가 먼저 필요합니다
# Windows: https://gocv.io/getting-started/windows/
# macOS:   https://gocv.io/getting-started/macos/
# Linux:   https://gocv.io/getting-started/linux/

cd <실습폴더>
go get gocv.io/x/gocv
```

## 공통 사전 준비 (Colab)

1. https://colab.research.google.com 접속 후 새 노트북 생성
2. 각 `.ipynb` 파일을 업로드하여 열기 (파일 → 노트 업로드)
3. 셀을 위에서부터 순서대로 실행 (Shift+Enter)

로컬 Jupyter에서 실행하려면: `pip install opencv-python-headless matplotlib numpy scikit-learn scikit-image scipy`

각 태스크 폴더의 README.md에서 더 자세한 실행 방법을 안내합니다.
