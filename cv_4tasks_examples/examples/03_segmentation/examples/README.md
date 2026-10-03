# 분할(Segmentation) — Semantic vs Instance 실습예제 (Colab/Python)

"컴퓨터비전 ③ 분할(Segmentation) — (사례) 의료영상 장기·종양 윤곽 분할 /
자율주행 도로 차선영역 분할"에 대응하는 Colab 쥬피터(Python) 실습
소스코드 모음입니다. **Semantic 방식 2개(사례별 1개씩) + Instance 방식
2개(사례별 1개씩), 총 4개** 실습예제로 구성했습니다.

## 폴더 구조

```
examples/
├── README.md
├── verified_outputs/                                실제 실행 결과 이미지(검증 증빙)
│   ├── medical_semantic_result.png
│   ├── driving_semantic_result.png
│   ├── medical_instance_result.png
│   └── driving_instance_result.png
├── ex1_semantic_segmentation_medical_code.py         실습1 (Semantic·의료영상)
├── ex1_semantic_segmentation_medical.ipynb
├── ex2_semantic_segmentation_driving_code.py         실습2 (Semantic·자율주행)
├── ex2_semantic_segmentation_driving.ipynb
├── ex3_instance_segmentation_medical_code.py         실습3 (Instance·의료영상)
├── ex3_instance_segmentation_medical.ipynb
├── ex4_instance_segmentation_driving_code.py         실습4 (Instance·자율주행)
└── ex4_instance_segmentation_driving.ipynb
```

## 실습 목록

| 실습 | 방식 | 사례 | 핵심 알고리즘 | 보여주는 것 |
|---|---|---|---|---|
| 실습1 | **Semantic** | 의료영상 | Multi-Otsu 임계값 (3클래스: 배경/장기/종양) | 종양 2개를 구분 없이 같은 "종양" 클래스 1개로 표시 |
| 실습2 | **Semantic** | 자율주행 | HSV 색상 규칙 (3클래스: 하늘/도로/차선) | 차선 3줄(좌·우·점선)을 구분 없이 같은 "차선 도색" 클래스 1개로 표시 |
| 실습3 | **Instance** | 의료영상 | Multi-Otsu + Watershed(거리변환+분수계) | 종양 2개를 "종양 #1", "종양 #2"로 각각 구분 |
| 실습4 | **Instance** | 자율주행 | HSV + 형태학적 팽창 병합 라벨링 | 끊어진 점선 7조각을 하나의 "차선 #1"로 올바르게 묶고, 좌/우 실선과 함께 총 3개 차선으로 구분 |

## Semantic과 Instance의 차이 (이 실습예제가 보여주는 것)

**Semantic 분할**은 "이 픽셀이 어떤 클래스인가"만 구분합니다. 종양이 2개든,
차선이 3줄이든 같은 클래스는 모두 같은 라벨(같은 색)로 칠해지므로 "몇 개"
있는지는 알 수 없습니다.

**Instance 분할**은 같은 클래스 안에서도 "서로 다른 개체"를 구분해 개별
ID를 부여합니다. 실습3은 두 종양이 비교적 떨어져 있어도(또는 서로 맞닿아
있어도) 거리 변환의 극댓값을 씨앗으로 삼는 Watershed로 분리하고, 실습4는
점선처럼 "원래 하나의 물체가 여러 조각으로 끊어져 보이는" 더 까다로운
경우를 형태학적 팽창으로 먼저 이어붙인 뒤 라벨링해 올바르게 하나의
개체로 묶어냅니다.

## 사용한 테스트 데이터

실제 의료영상(CT/MRI)·주행영상 데이터셋은 저작권·용량 문제로 다운로드가
필요해 Colab에서 바로 실행하기 어렵습니다. 이 실습에서는 **코드로 직접
생성하는 합성(synthetic) 이미지**를 사용해 외부 다운로드 없이 누구나
바로 실행할 수 있게 했습니다(난수 seed 고정으로 재현 가능).

- 의료영상 사례: 배경 위에 타원형 "장기", 그 안에 원형 "종양" 2개를
  배치하고 가우시안 잡음을 더한 합성 CT 단면.
- 자율주행 사례: 하늘 그라데이션 + 사다리꼴 도로 + 좌/우 실선 차선 +
  중앙 점선 차선(7조각)으로 구성한 합성 도로 장면.

## 검증 현황 (투명하게 안내)

이 코드를 작성한 환경에 `opencv-python-headless`, `scikit-image`,
`scipy`, `matplotlib`를 실제로 설치하여 **4개 전부 end-to-end로 직접
실행**하고 결과를 육안으로 확인했습니다.

- **실습1 (Semantic·의료)**: 배경 38,848 / 장기 25,072 / 종양 1,616픽셀로
  3클래스 정상 분류. 결과 이미지에서 종양 2개가 정확히 같은 빨간색으로
  칠해짐(= semantic의 한계를 그대로 보여줌) 확인.
- **실습2 (Semantic·자율주행)**: 하늘/배경 82,215 / 도로 32,740 / 차선
  도색 5,045픽셀로 3클래스 정상 분류. 좌/우 실선 + 중앙 점선 7조각이
  모두 같은 노란색으로 칠해짐 확인.
- **실습3 (Instance·의료)**: Watershed로 종양 **2개**를 정확히 분리해
  "종양 #1"(빨강), "종양 #2"(파랑)로 서로 다른 색 표시 확인.
- **실습4 (Instance·자율주행)**: 형태학적 병합 라벨링으로 끊어진 점선
  7조각이 올바르게 **1개**의 중앙 차선으로 합쳐지고, 좌/우 실선과 함께
  총 **3개** 차선으로 정확히 구분됨을 결과 이미지로 확인(병합 없이
  단순 Connected Components만 썼다면 최대 9개로 잘못 나뉘었을 것이라는
  비교 설명도 코드 실행 결과에 포함).
- 노트북(.ipynb) 파일 4개 모두 JSON 형식 유효성 검증 완료.
- 한글 그래프 제목·범례가 깨지지 않도록 시스템에 설치된 한글 폰트
  (Noto Sans CJK KR 또는 나눔고딕)를 자동 탐색해 등록하는 코드를
  포함했습니다. Colab처럼 한글 폰트가 전혀 없는 환경이면 코드 상단
  주석의 `!apt-get install fonts-nanum` 안내를 먼저 실행하세요.

## 아키텍처 — Clean Architecture + SOLID 원칙 적용

4개 예제 모두 하나의 `.py` 파일(Colab 셀 변환 규칙 `# Colab 셀 N: 설명`
유지) 안에서 아래 계층을 주석 배너로 명확히 구분합니다.

1. **`[도메인 계층 / Domain Layer]`** — `@dataclass` 값 객체
   (`SemanticSegmentationResult`, `InstanceSegmentationResult`)와
   `abc.ABC` 추상 인터페이스(`SemanticSegmenter`, `InstanceSplitter`)만
   정의합니다. **cv2/skimage/scipy를 전혀 import하지 않습니다.**
2. **`[유스케이스 계층 / Use Case Layer]`** — `SemanticSegmentationUseCase`
   (실습1·2), `InstanceSegmentationUseCase`(실습3·4). 인스턴스 분할
   유스케이스는 "시맨틱 분할기로 관심 클래스 추출 → 인스턴스 분리기로
   개체별 분리"라는 2단계 파이프라인을 두 인터페이스만으로 조립하며
   (DIP), cv2를 직접 호출하지 않습니다.
3. **`[인프라/어댑터 계층 / Infrastructure Layer]`** — cv2/skimage/scipy를
   실제로 호출하는 구현체(`MultiOtsuMedicalSegmenter`,
   `HSVColorRoadSegmenter`, `WatershedInstanceSplitter`,
   `MorphologicalMergeInstanceSplitter`)만 둡니다. "계산"과 "시각화"
   (`*Visualizer`)는 서로 다른 클래스로 분리했습니다.
4. **`[구성 루트 / Composition Root — Colab 실행 코드]`** — 합성 영상
   생성, 구체 어댑터 생성 및 의존성 주입(DI), 실행, matplotlib 시각화·
   저장을 담당합니다. 파일 I/O·`!pip` 설치가 허용되는 유일한 계층입니다.

### SOLID 적용 방식

- **SRP**: 시맨틱 분류 계산, 인스턴스 분리 계산, 시각화를 각각 별도
  클래스로 분리했습니다. 예: `MultiOtsuMedicalSegmenter`(분류 계산) ≠
  `InstanceSegmentationVisualizer`(시각화).
- **OCP**: 새 분할 알고리즘(예: 다른 색상 규칙, 다른 임계값 전략)은
  `SemanticSegmenter`를 구현하는 새 클래스를 추가하는 것만으로
  확장되며, 유스케이스·도메인 코드는 수정하지 않습니다.
- **LSP**: 실습1의 `MultiOtsuMedicalSegmenter`와 실습2의
  `HSVColorRoadSegmenter`는 같은 `SemanticSegmenter` 인터페이스를
  완전히 치환 가능하게 구현합니다. 마찬가지로 실습3의
  `WatershedInstanceSplitter`와 실습4의
  `MorphologicalMergeInstanceSplitter`는 같은 `InstanceSplitter`
  인터페이스의 서로 다른 전략이지만 완전히 교체 가능합니다 —
  `InstanceSegmentationUseCase`는 둘 중 어떤 것을 주입받아도 동일하게
  동작합니다.
- **ISP**: "클래스를 분류하는 책임"(`SemanticSegmenter`)과 "같은
  클래스 안에서 개체를 나누는 책임"(`InstanceSplitter`)을 완전히
  분리된 작은 인터페이스로 두었습니다. 인스턴스 분리 전략만 바꾸고
  싶을 때 시맨틱 분류 코드를 건드릴 필요가 없습니다.
- **DIP**: 유스케이스 계층은 cv2/skimage/scipy 같은 구체 구현이 아니라
  추상 인터페이스(`SemanticSegmenter`, `InstanceSplitter`)에만
  의존하며, 실제 구현체는 구성 루트에서 생성자 주입으로 연결됩니다.

## 실행 방법

### Colab

1. https://colab.research.google.com 접속 후 새 노트북 생성
2. 각 `.ipynb` 파일을 업로드하여 열기 (파일 → 노트 업로드)
3. 셀을 위에서부터 순서대로 실행 (Shift+Enter)

### 로컬 Jupyter

```bash
pip install opencv-python-headless scikit-image scipy matplotlib numpy
jupyter notebook ex1_semantic_segmentation_medical.ipynb
```
