// Package infra는 클린 아키텍처의 인프라/어댑터 계층이다.
// gocv(OpenCV의 Go 바인딩)를 실제로 호출하는 유일한 패키지이며, domain
// 패키지에 정의된 포트(인터페이스)들을 구현하는 구체 어댑터들을 담는다.
// 이 패키지 밖에서는 "gocv.io/x/gocv"를 import하지 않는다.
package infra

import "gocv.io/x/gocv"

// MatFrame은 gocv.Mat을 domain.Frame 인터페이스로 감싸는 어댑터다.
// 유스케이스/main 계층이 gocv.Mat을 직접 알 필요 없이 domain.Frame으로만
// 다룰 수 있게 해준다(DIP). 인프라 계층 내부(이 패키지)에서는 Mat 필드로
// 실제 gocv 연산을 수행한다.
type MatFrame struct {
	Mat gocv.Mat
}

// Empty는 내부 Mat이 비어 있는지 여부를 반환해 domain.Frame을 만족한다.
func (f MatFrame) Empty() bool {
	return f.Mat.Empty()
}
