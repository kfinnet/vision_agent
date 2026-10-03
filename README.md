# 고전 컴퓨터비전(Classical CV) 3대 기법 실습예제 (Colab/Python)

"고전 컴퓨터비전: Canny 엣지 / Haar Cascade / SIFT·ORB" 비교 설명에 대응하는
Colab 쥬피터(Python) 실습 소스코드 모음입니다. 기법별로 **1개씩, 총 3개**
실습예제로 구성했습니다.

## 폴더 구조

```
examples/
├── README.md                                  (본 파일)
├── verified_outputs/                           실제 실행해 얻은 결과 이미지(검증 증빙)
│   ├── canny_edge_result.png
│   ├── haar_cascade_result.png
│   ├── sift_matching_result.png
│   └── orb_matching_result.png
├── ex1_canny_edge_detection_code.py            실습1 소스코드 (노트북 변환 전 원본)
├── ex1_canny_edge_detection.ipynb              실습1 Colab 노트북
├── ex2_haar_cascade_code.py                    실습2 소스코드
├── ex2_haar_cascade.ipynb                      실습2 Colab 노트북
├── ex3_sift_orb_feature_matching_code.py       실습3 소스코드
└── ex3_sift_orb_feature_matching.ipynb         실습3 Colab 노트북
```

## 실습 목록

| 실습 | 기법 | 핵심 내용 |
|---|---|---|
| 실습1 | **Canny 엣지 검출** | 블러→그래디언트 계산→비최대 억제→2단계 임계값으로 경계선만 선별. 임계값을 바꿔가며 비교까지 포함 |
| 실습2 | **Haar Cascade** | CLAHE 전처리 + OpenCV 내장 사전학습 모델로 정면 얼굴 탐지. `minNeighbors` 파라미터 비교 포함 |
| 실습3 | **SIFT / ORB** | 같은 이미지를 회전+축소시킨 변형본과 특징점 매칭. SIFT·ORB를 같은 인터페이스의 교체 가능한 어댑터로 구현해 한 코드로 두 알고리즘 모두 실행·비교 |

## 사용한 테스트 이미지

다운로드 없이 바로 실행되도록, 모든 예제는 `scikit-image`에 내장된 샘플
이미지(`skimage.data.astronaut()` — 사람 얼굴이 포함된 우주비행사 사진)를
사용합니다.

## 검증 현황 (투명하게 안내)

이 코드를 작성한 환경에 `opencv-python-headless`(4.13), `scikit-image`,
`matplotlib`를 실제로 설치하여 **3개 전부 end-to-end로 직접 실행**하고 결과를
육안으로 확인했습니다.

- **실습1 (Canny)**: 임계값(low=50, high=150)에서 엣지 픽셀 16,211개 검출,
  임계값을 30/90·50/150·100/200으로 바꿔가며 비교 실행까지 정상 동작 확인.
- **실습2 (Haar Cascade)**: 우주비행사 사진에서 얼굴 1개를 바운딩박스
  `(x=176, y=65, w=99, h=99)`로 정확히 탐지 확인. `minNeighbors` 3/5/8
  비교 실행도 정상 동작 확인.
- **실습3 (SIFT/ORB)**: 원본과 25도 회전+0.8배 축소한 변형 이미지 사이에서
  SIFT는 500개 특징점 중 293쌍, ORB는 500개 특징점 중 245쌍이 올바르게
  매칭됨을 시각적으로 확인(결과 이미지에서 사람 얼굴·헬멧·우주왕복선 등
  같은 지점끼리 선으로 정확히 이어짐).
- 노트북(.ipynb) 파일 3개 모두 JSON 형식 유효성 검증 완료.
- 한글 그래프 제목·라벨이 깨지지 않도록 NanumGothic 폰트를 명시적으로
  등록하는 코드를 포함했습니다(Colab 환경에 한글 폰트가 없다면 셀 1의
  안내 주석대로 `apt-get install fonts-nanum` 실행 권장).

## 아키텍처 — Clean Architecture + SOLID 원칙 적용

3개 예제 모두 하나의 `.py` 파일(Colab 셀 변환 규칙 `# Colab 셀 N: 설명` 유지)
안에서 아래 4개 계층을 주석 배너로 명확히 구분합니다.

1. **`[도메인 계층 / Domain Layer]`** — `@dataclass` 값 객체(예:
   `EdgeDetectionResult`, `BoundingBox`, `FeatureMatchResult`)와 `abc.ABC`
   추상 인터페이스(예: `EdgeDetector`, `ObjectDetector`, `FeatureExtractor`,
   `FeatureMatcher`)만 정의합니다. **cv2를 전혀 import하지 않습니다.**
2. **`[유스케이스 계층 / Use Case Layer]`** — `EdgeDetectionUseCase`,
   `FaceDetectionUseCase`, `FeatureMatchingUseCase`. 도메인 인터페이스만
   생성자로 주입받아(DIP) 실행 흐름을 조립하며, cv2를 직접 호출하지 않습니다.
3. **`[인프라/어댑터 계층 / Infrastructure Layer]`** — cv2를 실제로 호출하는
   구현체(`CannyEdgeDetector`, `HaarCascadeFaceDetector`,
   `SIFTFeatureExtractor`/`ORBFeatureExtractor`, `RatioTestBFMatcher`
   등)만 둡니다. "계산" 로직과 "시각화"(`*Visualizer`)는 서로 다른
   클래스로 분리했습니다.
4. **`[구성 루트 / Composition Root — Colab 실행 코드]`** — 실제 Colab 셀:
   이미지 준비, 구체 어댑터 생성 및 의존성 주입(DI), 실행, matplotlib
   시각화·파일 저장을 담당합니다. 파일 I/O·`!pip` 설치가 허용되는
   유일한 계층입니다.

### SOLID 적용 방식

- **SRP**: 전처리/검출(또는 추출)/매칭 "계산" 로직과 시각화를 각각 별도
  클래스로 분리했습니다. 예: `CannyEdgeDetector`(계산) ≠
  `EdgeDetectionVisualizer`(시각화).
- **OCP**: 새 알고리즘(예: 다른 엣지 검출기, 다른 특징 추출기)은 인프라
  계층에 새 클래스를 추가하는 것만으로 확장되며, 유스케이스·도메인
  코드는 수정하지 않습니다. 실습1·2의 "구성 루트" 마지막 셀에서
  파라미터만 바꿔 새 어댑터 인스턴스를 만드는 비교 실행으로 이를
  직접 보여줍니다.
- **LSP**: 실습3의 `SIFTFeatureExtractor`와 `ORBFeatureExtractor`는 같은
  `FeatureExtractor` 인터페이스를 완전히 치환 가능하게 구현합니다 —
  `FeatureMatchingUseCase`는 둘 중 어떤 것을 주입받아도 동일하게 동작합니다.
- **ISP**: 인터페이스를 책임 단위로 작게 유지합니다. 예: 특징 "추출"
  (`FeatureExtractor`)과 "매칭"(`FeatureMatcher`)을 별도 인터페이스로
  분리해, 매칭 전략만 바꾸고 싶을 때 추출기 코드를 건드릴 필요가
  없습니다.
- **DIP**: 유스케이스 계층은 cv2 같은 구체 구현이 아니라 추상 인터페이스
  (`EdgeDetector`, `ObjectDetector`, `FeatureExtractor`, `FeatureMatcher`)
  에만 의존하며, 실제 구현체는 구성 루트에서 생성자 주입으로
  연결됩니다.

## 실행 방법

### Colab

1. https://colab.research.google.com 접속 후 새 노트북 생성
2. 각 `.ipynb` 파일을 업로드하여 열기 (파일 → 노트 업로드)
3. 셀을 위에서부터 순서대로 실행 (Shift+Enter)

### 로컬 Jupyter

```bash
pip install opencv-python-headless matplotlib numpy scikit-image
jupyter notebook ex1_canny_edge_detection.ipynb
```
