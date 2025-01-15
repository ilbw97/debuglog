# **debuglog**

`debuglog`는 Go 애플리케이션에서 파일 기반 로깅을 간편하게 설정하고 관리할 수 있도록 도와주는 패키지입니다. Logrus와 Lumberjack을 기반으로 작성되었으며, 파일 롤링, 압축, 콘솔 출력 등 다양한 기능을 제공합니다.

---

## **특징**

-   **파일 기반 로깅**: 로그를 지정된 디렉토리에 파일로 저장.
-   **파일 롤링 지원**: 로그 파일 크기와 보관 기간을 설정하여 자동 관리.
-   **로그 압축**: 오래된 로그 파일을 압축하여 저장 공간 절약.
-   **콘솔 출력 옵션**: 로그를 파일과 콘솔에 동시에 출력 가능.
-   **로그 파일 이름 옵션**: PID 또는 날짜/시간 기반의 파일 이름 설정.
-   **사용자 정의 필드 정렬**: 로그 필드 순서를 커스터마이즈 가능.

---

## **설치**

Go 모듈을 사용하는 프로젝트에서 다음 명령어를 실행하여 패키지를 설치하세요:

```bash
go get github.com/ilbw97/debuglog
```

---

## **사용법**

### **1. Import 및 초기화**

패키지를 import하고 로거를 초기화합니다:

```go
package main
import (
    debuglog "github.com/ilbw97/debuglog"

)
func main() {
    // 로거 초기화logger := debuglog.DebugLogInit("myapp", true, true, true)
    // 로그 메시지 출력
    logger.Info("This is an info message.")
    logger.Error("This is an error message.")
}
```

### **2. 함수 설명**

#### `DebugLogInit`

`DebugLogInit` 함수는 로거를 초기화하며 다음과 같은 매개변수를 제공합니다:

```go
func DebugLogInit(options *LogConfig) *logrus.Logger
```

| 매개변수         | 타입     | 설명                                                                           |
| ---------------- | -------- | ------------------------------------------------------------------------------ |
| `logname`        | `string` | 로그 파일 이름의 기본값 (확장자는 자동으로 추가됨).                            |
| `makedir`        | `bool`   | 로그 디렉토리를 생성할지 여부 (`true`: 생성, `false`: 현재 디렉토리 사용).     |
| `usePID`         | `bool`   | 로그 파일 이름에 PID를 포함할지 여부 (`true`: 포함, `false`: 미포함).          |
| `useMultiWriter` | `bool`   | 로그를 콘솔과 파일에 동시에 출력할지 여부 (`true`: 활성화, `false`: 비활성화). |

---

## **구조체 (Structs)**

`debuglog` 패키지는 다음과 같은 주요 구조체를 제공합니다:

### `LogConfig` 구조체

로깅 동작을 세부적으로 제어하기 위한 구성 구조체입니다:

```go
type LogConfig struct {
	LogName     string  `json:"name"`                               // process name
	MakeDir     bool	`json:"make_dir" example:"true"`            // make directory
	UsePID      bool	`json:"use_pid" exmaple:"true"`             // use pid
	UseMultiWriter bool	`json:"use_multi_writer" example:"false"`   // use multiwriter
	LogRotateConfig
}
```

```go
type LogRotateConfig struct {
    MaxSize    int  // 최대 로그 파일 크기 (MB)
    MaxBackups int  // 최대 백업 파일 개수
    MaxAge     int  // 로그 파일 최대 보관 기간 (일)
    Compress   bool // 오래된 로그 파일 압축 여부
}
```

#### 구조체 사용 예시

```go
logConfig := &debuglog.LogConfig{
        LogName:     "myapp",
        MakeDir:     true,      // make directory
        UsePID:      true,      // use pid
        UseMultiWriter: true,   // use multiwriter
        LogRotateConfig: debuglog.LogRotateConfig{  // log rotate config
            MaxSize:    100,    // 100MB마다 새 로그 파일 생성
            MaxBackups: 5,      // 최대 5개의 백업 파일 유지
            MaxAge:     7,      // 7일 지난 로그 파일 삭제
            Compress:   true,   // 오래된 로그 파일 압축
        }
}
```

이 구조체를 사용하여 로깅 동작을 세밀하게 조정할 수 있습니다.

---

## **옵션별 동작**

### **PID 사용 여부**

-   **PID 사용 (`usePID=true`)**
-   로그 파일 이름 형식: `myapp.<PID>.log`

```go
logger := debuglog.DebugLogInit("myapp", true, true, false)
```

-   **PID 미사용 (`usePID=false`)**
-   로그 파일 이름 형식: `myapp.<YYYYMMDD_HHMMSS>.log`

```go
logger := debuglog.DebugLogInit("myapp", true, false, false)
```

### **MultiWriter 사용 여부**

-   **MultiWriter 사용 (`useMultiWriter=true`)**
-   로그가 콘솔(`os.Stdout`)과 파일에 동시에 출력됩니다.

```go
logger := debuglog.DebugLogInit("myapp", true, false, true)
```

-   **MultiWriter 미사용 (`useMultiWriter=false`)**
-   로그가 파일에만 저장됩니다.

```go
logger := debuglog.DebugLogInit("myapp", true, false, false)
```

---

## **환경 변수**

다음 환경 변수를 통해 로깅 동작을 제어할 수 있습니다:

| 환경 변수       | 기본값                   | 설명                          |
| --------------- | ------------------------ | ----------------------------- |
| `LOG_BASE_PATH` | 현재 작업 디렉토리(`./`) | 로그 파일이 저장될 기본 경로. |

예를 들어 `/var/log/myapp` 경로에 로그를 저장하려면 다음 명령어를 실행하세요:

```bash
export LOG_BASE_PATH=/var/log/myapp
```

---

## **로그 설정**

### 기본 설정

-   최대 파일 크기: 500MB
-   최대 백업 개수: 3개
-   최대 보관 기간: 3일
-   압축 활성화: `true`

이 설정은 코드 내에서 변경 가능합니다. 필요하면 직접 수정하여 사용할 수 있습니다.

---

## **예제 코드**

### PID와 MultiWriter 모두 활성화

```go
package main
import (
    debuglog "github.com/ilbw97/debuglog"
)
func main() {
    // 새로운 LogConfig 구조체를 사용한 로거 초기화
    logConfig := &debuglog.LogConfig{
        LogName:     "myapp",
        MakeDir:     true,
        UsePID:      true,
        UseMultiWriter: true,
        LogRotateConfig: debuglog.LogRotateConfig{
            MaxSize:    100,
            MaxBackups: 5,
            MaxAge:     7,
            Compress:   true,
        },
    }
    logger := debuglog.DebugLogInit(logConfig)

    // 로그 메시지 출력
    logger.Info("This is an info message.")
}
```

### PID 비활성화 및 MultiWriter 비활성화

```go
package main
import (
    debuglog "github.com/ilbw97/debuglog"
)
func main() {
    logger := debuglog.DebugLogInit("example", true, false, false)
    logger.Info("This log will only be written to the log file.")
}
```

---

## **테스트**

패키지가 예상대로 작동하는지 확인하려면 다음 테스트 시나리오를 실행하세요:

1. PID 및 MultiWriter 활성화:

```go
logConfig := &debuglog.LogConfig{
    LogName:     "myapp",
    MakeDir:     true,
    UsePID:      true,
    UseMultiWriter: true,
}
logger := debuglog.DebugLogInit(logConfig)

```

2. PID 및 MultiWriter 비활성화:

```go
logConfig := &debuglog.LogConfig{
    LogName:     "myapp",
    MakeDir:     true,
    UsePID:      false,
    UseMultiWriter: false,
}
logger := debuglog.DebugLogInit(logConfig)
```

3. 환경 변수로 경로 변경:

```bash
export LOG_BASE_PATH=/tmp/logs
```

```go
logger := debuglog.DebugLogInit("test", true, false, false)
logger.Info("Test log with custom base path.")
```

---

## **의존성**

이 패키지는 다음 외부 라이브러리를 사용합니다:

1. [Logrus](https://github.com/sirupsen/logrus): 고급 로깅 라이브러리.
2. [Lumberjack](https://github.com/natefinch/lumberjack): 파일 롤링 및 관리 라이브러리.

의존성은 `go.mod`에서 자동 관리됩니다.

---

## **라이선스**

이 프로젝트는 MIT 라이선스를 따릅니다. 자유롭게 수정 및 배포할 수 있습니다.
