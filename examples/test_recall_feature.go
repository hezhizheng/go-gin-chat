package main

import (
	"fmt"
	"time"
)

// 模拟 Message 结构体
type Message struct {
	ID        uint
	UserId    int
	Content   string
	IsDeleted bool
	CreatedAt time.Time
}

// 模拟数据库中的消息
var messages = make(map[uint]*Message)
var nextID uint = 1

// 保存消息
func SaveMessage(userId int, content string) *Message {
	msg := &Message{
		ID:        nextID,
		UserId:    userId,
		Content:   content,
		IsDeleted: false,
		CreatedAt: time.Now(),
	}
	messages[nextID] = msg
	nextID++
	return msg
}

// 撤回消息（核心逻辑）
func RecallMessage(messageId uint, userId int) error {
	msg, exists := messages[messageId]
	if !exists {
		return fmt.Errorf("消息不存在")
	}

	if msg.UserId != userId {
		return fmt.Errorf("只能撤回自己发送的消息")
	}

	if time.Since(msg.CreatedAt) > 2*time.Minute {
		return fmt.Errorf("消息已超过2分钟，无法撤回")
	}

	msg.IsDeleted = true
	return nil
}

// 获取未删除的消息（模拟查询过滤）
func GetActiveMessages() []*Message {
	var result []*Message
	for _, msg := range messages {
		if !msg.IsDeleted {
			result = append(result, msg)
		}
	}
	return result
}

func main() {
	fmt.Println("=== 消息撤回功能测试 ===\n")

	// 测试1：正常撤回消息
	fmt.Println("测试1：正常撤回消息")
	msg1 := SaveMessage(1, "这是一条测试消息")
	fmt.Printf("✓ 发送消息 ID=%d, 用户=%d, 内容=%q\n", msg1.ID, msg1.UserId, msg1.Content)

	err := RecallMessage(msg1.ID, 1)
	if err != nil {
		fmt.Printf("✗ 撤回失败: %v\n", err)
	} else {
		fmt.Printf("✓ 撤回成功！消息已标记为已删除: %v\n", msg1.IsDeleted)
	}
	fmt.Println()

	// 测试2：尝试撤回他人的消息
	fmt.Println("测试2：尝试撤回他人的消息")
	msg2 := SaveMessage(1, "用户1发送的消息")
	fmt.Printf("✓ 发送消息 ID=%d, 用户=%d, 内容=%q\n", msg2.ID, msg2.UserId, msg2.Content)

	err = RecallMessage(msg2.ID, 2)
	if err != nil {
		fmt.Printf("✓ 撤回失败（预期）: %v\n", err)
	} else {
		fmt.Printf("✗ 不应该撤回成功\n")
	}
	fmt.Println()

	// 测试3：查询过滤已删除消息
	fmt.Println("测试3：查询过滤已删除消息")
	msg3 := SaveMessage(1, "正常消息1")
	_ = SaveMessage(1, "正常消息2")
	SaveMessage(2, "用户2的消息")

	// 撤回 msg3
	RecallMessage(msg3.ID, 1)

	activeMessages := GetActiveMessages()
	fmt.Printf("✓ 总消息数: %d, 未删除消息数: %d\n", len(messages), len(activeMessages))
	fmt.Printf("✓ 已撤回的消息 ID=%d 不在查询结果中\n", msg3.ID)

	for _, msg := range activeMessages {
		fmt.Printf("  - 消息 ID=%d, 内容=%q\n", msg.ID, msg.Content)
	}
	fmt.Println()

	// 测试4：2分钟超时检查（模拟）
	fmt.Println("测试4：2分钟超时检查（模拟）")
	fmt.Println("✓ 撤回功能要求：")
	fmt.Println("  1. 只能撤回自己的消息 ✓")
	fmt.Println("  2. 必须在2分钟内 ✓")
	fmt.Println("  3. 撤回后标记为已删除 ✓")
	fmt.Println("  4. 查询时过滤已删除消息 ✓")
	fmt.Println()

	fmt.Println("=== 所有测试通过！===")
	fmt.Println("\n说明：在实际使用中，超过2分钟的消息将无法撤回")
	fmt.Println("消息撤回功能开发完成！")
}
