package time

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

// APIOutputTime 是一个自定义时间类型，用于在API响应中统一时间格式
type APIOutputTime struct {
	time.Time  // 内嵌标准 time.Time 类型
	timeFormat string
}

// OutputTimeFormat 定义了API响应中期望的时间格式
const OutputTimeFormat = time.DateTime

// MarshalJSON 实现了 json.Marshaler 接口
func (aot *APIOutputTime) MarshalJSON() ([]byte, error) {
	if aot.Time.IsZero() {
		// INFO: 零值作为nil处理
		return []byte("null"), nil
	}
	// 将时间格式化为期望的字符串，并加上JSON字符串应有的双引号
	formattedString := fmt.Sprintf("\"%s\"", aot.Time.Format(aot.timeFormat))
	return []byte(formattedString), nil
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口 (可选，如果此类型也用于请求绑定)
func (aot *APIOutputTime) UnmarshalJSON(data []byte) error {
	// 去掉可能的双引号
	s := string(data)
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "null" || s == "" {
		aot.Time = time.Time{}
		return nil
	}
	// 尝试解析期望的格式
	t, err := time.Parse(OutputTimeFormat, s)
	if err == nil {
		aot.Time = t
		return nil
	}
	// 如果期望格式解析失败，尝试解析 RFC3339 (标准库 time.Time 的默认JSON格式)
	// 这样可以增加对输入格式的兼容性
	t, errRFC := time.Parse(time.RFC3339Nano, s)
	if errRFC == nil {
		aot.Time = t
		return nil
	}
	// 如果两种都失败，返回最初的错误
	return fmt.Errorf("unable to parse time string '%s' with format '%s' or RFC3339: %w", s, OutputTimeFormat, err)
}

// --- 在 Gin 中使用 ---

// ExampleResponse DTO 示例，其中时间字段使用了 APIOutputTime
type ExampleResponse struct {
	EventID     string        `json:"event_id"`
	Description string        `json:"description"`
	ScheduledAt APIOutputTime `json:"scheduled_at"` // 使用自定义时间类型
	UpdatedAt   APIOutputTime `json:"updated_at"`   // 另一个使用自定义时间类型的字段
	RawTime     time.Time     `json:"raw_time"`     // 对比：标准 time.Time 字段
}

func main() {
	r := gin.Default()

	r.GET(
		"/event/:id", func(c *gin.Context) {
			eventID := c.Param("id")

			// 模拟从业务逻辑中获取数据，这里仍然是标准的 time.Time
			now := time.Now()                    // 业务逻辑中获取的当前时间
			scheduled := now.Add(24 * time.Hour) // 业务逻辑中计算得到的调度时间

			response := ExampleResponse{
				EventID:     eventID,
				Description: "This is a test event.",
				// 将标准的 time.Time 包装进 APIOutputTime 用于响应
				ScheduledAt: APIOutputTime{Time: scheduled},
				UpdatedAt:   APIOutputTime{Time: now},
				RawTime:     now, // 这个字段会以RFC3339Nano格式输出
			}

			// 当 response 被序列化为 JSON 时:
			// - ScheduledAt 和 UpdatedAt 字段会调用 APIOutputTime.MarshalJSON()
			// - RawTime 字段会调用 time.Time.MarshalJSON() (标准行为)
			c.JSON(http.StatusOK, response)
		},
	)

	fmt.Println("Gin server running on :8080")
	fmt.Println("Try: curl http://localhost:8080/event/123")
	/*
	 预计输出示例:
	 {
	   "event_id": "123",
	   "description": "This is a test event.",
	   "scheduled_at": "2025-05-20 17:56:22",  // (当前时间 + 24小时，格式 YYYY-MM-DD HH:MM:SS)
	   "updated_at": "2025-05-19 17:56:22",    // (当前时间，格式 YYYY-MM-DD HH:MM:SS)
	   "raw_time": "2025-05-19T17:56:22.123456789Z" // (标准RFC3339Nano格式)
	 }
	 (注意: .123456789Z 部分会根据实际纳秒和时区变化)
	*/
	r.Run(":8080")
}
