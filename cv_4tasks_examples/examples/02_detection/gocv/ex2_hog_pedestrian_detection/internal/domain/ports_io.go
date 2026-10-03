package domain

// FrameDrawer는 탐지 결과를 프레임 위에 그려 새 프레임을 만드는 포트다.
// 탐지(ObjectDetector)와 시각화 책임을 분리해(SRP) 박스 색상/라벨 형식을
// 바꾸더라도 탐지 로직이나 유스케이스는 영향을 받지 않는다(OCP).
type FrameDrawer interface {
	// Draw는 frame 위에 detections를 그린 새 Frame을 반환한다.
	Draw(frame Frame, detections []Detection) Frame
}

// ImageReader는 정지 이미지 파일을 읽어 Frame으로 반환하는 포트다.
type ImageReader interface {
	Read(path string) (Frame, error)
}

// ImageWriter는 Frame을 이미지 파일로 저장하는 포트다.
type ImageWriter interface {
	Write(path string, frame Frame) (bool, error)
}

// VideoSource는 비디오 파일에서 프레임을 순차적으로 읽어오는 포트다.
type VideoSource interface {
	// Read는 다음 프레임을 읽어 ok로 성공 여부를 반환한다.
	Read() (frame Frame, ok bool)
	// Close는 사용 중인 자원을 해제한다.
	Close() error
}

// WindowRenderer는 프레임을 화면 창에 표시하는 포트다.
// 박스를 "그리는" 책임은 FrameDrawer가 담당하므로, 이 포트는 순수하게
// "보여주기"만 책임진다(SRP/ISP).
type WindowRenderer interface {
	// Show는 frame을 창에 출력한다.
	Show(frame Frame)
	// WaitKey는 지정한 시간(ms) 동안 키 입력을 기다리고 키 코드를 반환한다.
	WaitKey(delayMs int) int
	// Close는 창 자원을 해제한다.
	Close() error
}
