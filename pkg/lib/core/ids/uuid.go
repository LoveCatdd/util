package ids

import (
	uuid "github.com/satori/go.uuid"
)

// 返回基于当前时间戳和MAC地址的UUID。
func UUIDV1() uuid.UUID {
	return uuid.NewV1()
}

// 返回基于POSIX UID/GID的DCE安全UUID。
func UUIDV2() uuid.UUID {

	return uuid.NewV2(byte(176))
}

// 返回基于命名空间UUID和名称的MD5哈希的UUID。
func UUIDV3(u uuid.UUID, name string) uuid.UUID {
	if isEmpty(u) {
		u = UUIDV2()
	}
	return uuid.NewV3(u, name)
}

// 返回随机生成的UUID。
func UUIDV4() uuid.UUID {
	return uuid.NewV4()
}

// 返回基于命名空间UUID和名称的SHA-1哈希的UUID。
func UUIDV5(u uuid.UUID, name string) uuid.UUID {
	if isEmpty(u) {
		u = UUIDV2()
	}
	return uuid.NewV5(u, name)
}

// 将字符串转换成UUID
func FromString(s string) (uuid.UUID, error) {
	return uuid.FromString(s)
}

func isEmpty(arr uuid.UUID) bool {
	for _, b := range arr {
		if b != 0 {
			return false
		}
	}
	return true
}
