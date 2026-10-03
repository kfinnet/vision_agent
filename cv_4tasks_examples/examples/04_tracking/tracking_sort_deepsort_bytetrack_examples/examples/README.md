# 추적(Tracking) — SORT/DeepSORT, ByteTrack 실습예제 (Colab/Python)

"컴퓨터비전 ④ 추적(Tracking) — (활용사례) 스포츠 선수 동선·기록 분석, 매장
고객 동선 분석(리테일 애널리틱스), 교통량·차량 속도 측정"에 대응하는 Colab
쥬피터(Python) 실습 소스코드 모음입니다. **SORT/DeepSORT 1개 + ByteTrack
1개, 총 2개** 실습예제로 구성했습니다.

## 폴더 구조

```
examples/
├── README.md
├── verified_outputs/                              실제 실행 결과 이미지(검증 증빙)
│   ├── sort_tracking_result.png
│   ├── deepsort_tracking_result.png
│   ├── baseline_tracking_result.png
│   └── bytetrack_tracking_result.png
├── ex1_sort_deepsort_tracking_code.py              실습1 (SORT / DeepSORT)
├── ex1_sort_deepsort_tracking.ipynb
├── ex2_bytetrack_tracking_code.py                  실습2 (ByteTrack)
└── ex2_bytetrack_tracking.ipynb
```

## 실습 목록

| 실습 | 알고리즘 | 핵심 메커니즘 | 보여주는 것 |
|---|---|---|---|
| 실습1 | **SORT → DeepSORT** | 탐지 + 칼만필터(등속 운동 모델) + 헝가리안 매칭(IoU), DeepSORT는 여기에 색상 외형 특징 + 재식별(Re-ID) 갤러리 추가 | 두 객체가 완전히 겹쳐 긴 가려짐(18프레임)이 발생할 때, SORT는 트랙이 끊겨 새 ID(3개)가 생기지만 DeepSORT는 외형 특징으로 원래 ID를 복구해 2개로 유지 |
| 실습2 | **ByteTrack** | 1단계: 고신뢰 탐지↔전체 트랙 매칭 → 2단계: 1단계에서 남은 트랙↔저신뢰 탐지 매칭(새 트랙 생성에는 사용 안 함) | 장애물에 가려 탐지 신뢰도가 낮아지는 구간에서, 고신뢰 박스만 쓰는 방식은 ID가 끊겨 2개가 되지만 ByteTrack은 저신뢰 박스로 트랙을 5회 "구조"해 1개 ID로 유지 |

## 왜 합성(synthetic) 영상을 사용했는가

실제 SORT/DeepSORT/ByteTrack은 YOLO 등 딥러닝 검출기의 출력(박스+신뢰도)을
입력으로 받고, DeepSORT는 추가로 사전학습된 ReID 임베딩 네트워크를 사용합니다.
이 가중치들은 인터넷에서 내려받아야 해 Colab에서 매번 바로 실행하기
어렵습니다. 이 실습에서는:

- **검출기**: 색상 기반 contour 탐지 — 객체가 가려질수록 보이는 면적이
  줄어 "신뢰도(confidence)"가 자연스럽게 낮아지도록 설계했습니다(실제
  딥러닝 검출기가 가려진 객체에 낮은 점수를 주는 현상을 재현).
- **DeepSORT의 외형 특징**: 사전학습 ReID 네트워크 대신 색상 히스토그램
  (`cv2.calcHist`)을 사용했습니다. "외형으로 같은 객체인지 구분한다"는
  알고리즘적 역할은 동일합니다.

대신 **SORT/DeepSORT/ByteTrack 논문이 실제로 제안하는 알고리즘 로직
(칼만필터 수식, 헝가리안 매칭, 2단계 매칭, 재식별 갤러리)은 전부 직접
구현**했으며, 어떤 외부 사전학습 모델도 다운로드하지 않고 Colab에서 바로
실행됩니다.

## 검증 현황 (투명하게 안내)

이 코드를 작성한 환경에 `opencv-python-headless`, `scipy`, `matplotlib`를
실제로 설치하여 **2개 전부 end-to-end로 직접 실행**하고 결과를 육안으로
확인했습니다.

- **실습1 (SORT vs DeepSORT)**: 두 원형 객체가 서로 다른 방향(위→아래,
  아래→위)으로 움직이며 화면 중앙에서 약 18프레임 동안 완전히 겹치는
  장면을 생성했습니다(`max_age=6`으로 의도적으로 짧게 설정).
  - **SORT**: 겹침 구간을 넘기지 못해 트랙이 삭제되고, 다시 나타난 객체에
    새 ID가 부여되어 **총 3개** 트랙 ID가 생성됨을 확인했습니다.
  - **DeepSORT**: 사라지기 직전 색상 특징을 재식별 갤러리에 기억해두었다가,
    다시 나타난 탐지와 코사인 유사도로 비교해 원래 ID를 복구 — **총 2개**
    트랙 ID로 정확히 유지됨을 확인했습니다("재식별 성공" 로그 출력 포함).
- **실습2 (ByteTrack)**: 객체 1개가 폭 28px(객체 지름보다 좁은) 고정
  장애물 뒤를 지나가며 신뢰도가 0.01~0.9 사이로 변동하는 장면을
  생성했습니다(`max_age=5`).
  - **고신뢰만 사용(비교 기준)**: 신뢰도 0.6 미만 구간을 모두 버려
    트랙이 끊기고 **총 2개** 트랙 ID가 생성됨을 확인했습니다.
  - **ByteTrack**: 2단계 매칭으로 저신뢰 탐지 5회를 활용해 트랙을
    계속 살려, **총 1개** 트랙 ID로 끊김 없이 유지됨을 확인했습니다.
- 두 실습 모두 결과를 궤적(trajectory) 그래프로 시각화해 `verified_outputs/`
  에 증빙으로 포함했습니다. 노트북(.ipynb) 2개 모두 JSON 형식 유효성
  검증을 완료했습니다.

## 아키텍처 — Clean Architecture + SOLID 원칙 적용

두 예제 모두 하나의 `.py` 파일(Colab 셀 변환 규칙 `# Colab 셀 N: 설명`
유지) 안에서 아래 계층을 주석 배너로 명확히 구분합니다.

1. **`[도메인 계층 / Domain Layer]`** — `@dataclass` 값 객체
   (`Detection`, `Track`)와 `abc.ABC` 추상 인터페이스
   (`ObjectDetector`, `MotionModel`, `Associator`, 실습1은 추가로
   `TrackGallery`)만 정의합니다. **cv2/scipy를 전혀 import하지 않습니다.**
   IoU 계산처럼 외부 라이브러리 없이 사칙연산만으로 가능한 순수 함수
   (`compute_iou`)도 도메인에 둡니다.
2. **`[유스케이스 계층 / Use Case Layer]`** — `MultiObjectTrackingUseCase`.
   "예측 → 탐지 → 매칭 → 갱신 → 생성 → 소멸"이라는 추적의 전체 흐름만
   조립하며, SORT/DeepSORT의 차이나 ByteTrack의 2단계 매칭 같은 세부
   전략은 전혀 모릅니다 — 전부 주입받은 `Associator`/`TrackGallery`
   구현체 안에 캡슐화되어 있습니다.
3. **`[인프라/어댑터 계층 / Infrastructure Layer]`** — cv2/scipy를 실제로
   호출하는 구현체(`ColorContourDetector`/`OcclusionAwareDetector`,
   `ConstantVelocityKalmanMotionModel`, `IoUHungarianAssociator`,
   `AppearanceIoUAssociator`, `ByteTrackAssociator`,
   `NullTrackGallery`/`AppearanceReidentificationGallery`)만 둡니다.
   "계산"과 "시각화"(`TrackingVisualizer`)는 분리했습니다.
4. **`[구성 루트 / Composition Root — Colab 실행 코드]`** — 합성 영상
   생성, 구체 어댑터 생성 및 의존성 주입(DI)으로 SORT 구성과 DeepSORT
   구성(또는 비교기준과 ByteTrack 구성)을 나란히 실행·비교합니다.

### SOLID 적용 방식

- **SRP**: 탐지, 운동 예측(칼만필터), 매칭, 재식별, 시각화를 각각 별도
  클래스로 분리했습니다.
- **OCP**: 실습1은 `IoUHungarianAssociator`(SORT) → `AppearanceIoUAssociator`
  (DeepSORT)로, 실습2는 `HighConfidenceOnlyAssociator` →
  `ByteTrackAssociator`로 **Associator 구현체 하나만 교체**하는 것으로
  알고리즘을 바꿨습니다. `MultiObjectTrackingUseCase`는 두 실습 모두
  단 한 줄도 수정하지 않았습니다.
- **LSP**: `Associator`의 모든 구현체(IoU만 쓰는 것, 외형을 더하는 것,
  2단계로 나누는 것)는 완전히 치환 가능합니다. `TrackGallery`의
  `NullTrackGallery`(아무것도 안 함)와 `AppearanceReidentificationGallery`
  (실제로 기억·재인식함)도 마찬가지입니다 — "아무것도 하지 않는 구현"도
  인터페이스 계약을 정직하게 지키는 한 올바른 LSP 사례입니다.
- **ISP**: "트랙과 탐지를 짝짓는 책임"(`Associator`)과 "사라진 트랙을
  나중에 재인식하는 책임"(`TrackGallery`)을 완전히 분리했습니다. 재식별
  전략만 바꾸고 싶을 때 매칭 코드를 건드릴 필요가 없습니다.
- **DIP**: 유스케이스 계층은 cv2/scipy 같은 구체 구현이 아니라 추상
  인터페이스에만 의존하며, 실제 구현체는 구성 루트에서 생성자 주입으로
  연결됩니다.

## 실행 방법

### Colab

1. https://colab.research.google.com 접속 후 새 노트북 생성
2. 각 `.ipynb` 파일을 업로드하여 열기 (파일 → 노트 업로드)
3. 셀을 위에서부터 순서대로 실행 (Shift+Enter)

### 로컬 Jupyter

```bash
pip install opencv-python-headless scipy matplotlib numpy
jupyter notebook ex1_sort_deepsort_tracking.ipynb
```
