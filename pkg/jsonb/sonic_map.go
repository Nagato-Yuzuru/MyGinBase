package jsonb

import (
	"database/sql/driver"
	"fmt"
	"terraqt.io/colas/bedrock-go/pkg/errs"

	"github.com/bytedance/sonic"
)

// SonicMap 是一个 map[string]interface{} 的包装类型，
// 用于配合 sqlc 和 sonic 处理 PostgreSQL 的 JSONB 类型。
type SonicMap struct {
	Json any
}

// Value 实现了 driver.Valuer 接口。
// 当 Go 的 SonicMap 类型需要被写入数据库时，此方法被调用。
// 它使用 sonic 将 map[string]interface{} 序列化为 JSON []byte。
func (sm *SonicMap) Value() (driver.Value, error) {
	if sm == nil || sm.Json == nil {
		// 如果 map 为 nil，我们将其作为 SQL NULL 存入数据库。
		return nil, nil
	}
	// 使用 sonic 将 map 序列化为 JSON 字节数组。
	jsonBytes, err := sonic.Marshal(sm)
	if err != nil {
		return nil, errs.WrapCodeError(
			errs.ErrMarshalFailed,
			fmt.Errorf("sonicMap:  sonic can't marshal JSON: %w", err),
		)
	}
	return jsonBytes, nil
}

// Scan 实现了 sql.Scanner 接口。
// 当从数据库读取 JSONB 数据并赋值给 Go 的 SonicMap 类型时，此方法被调用。
// src 参数是数据库驱动传过来的原始数据，通常是 []byte (对于 jsonb)。
func (sm *SonicMap) Scan(src any) error {
	if src == nil {
		sm.Json = nil
		return nil
	}

	var sourceBytes []byte
	switch s := src.(type) {
	case []byte:
		sourceBytes = s
	case string:
		sourceBytes = []byte(s)
	default:
		return errs.WrapCodeError(
			errs.ErrUnmarshalFailed,
			fmt.Errorf("customtypes: can't support type %T to SonicMap，expect []byte or string", src),
		)
	}

	// 如果字节数组为空 (例如，PostgreSQL 的 'null'::jsonb 可能会被驱动视为空的 []byte)，
	// 我们也将其视作 nil map。
	if len(sourceBytes) == 0 {
		sm.Json = nil
		return nil
	}

	// sonic.Unmarshal 期望一个非 nil 的指针。
	// 我们创建一个临时的 map 来接收反序列化的结果，然后再赋值给 *sm。
	var tempMap any
	// 使用 sonic 将 JSON 字节数组反序列化到临时的 map 中。
	err := sonic.Unmarshal(sourceBytes, &tempMap)
	if err != nil {
		return errs.WrapCodeError(
			errs.ErrUnmarshalFailed, fmt.Errorf(
				"customtypes: can't sonic use JSONB (%s) Unmarshal to SonicMap: %w",
				string(sourceBytes), err,
			),
		)
	}
	sm.Json = tempMap // 将反序列化后的 map 赋值给指针指向的 SonicMap
	return nil
}
