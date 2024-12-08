package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"notifier/internal/entity"
)

func Test_buildCfg(t *testing.T) {
	cfg, _ := buildCfg(map[string]string{
		"holdTasks": "2h",
		"holdDelay": "1m",
		"swapFile":  "./file.bin",
		"swapDelay": "1h",
	})

	assert.Equal(t, 2*time.Hour, cfg.holdDoneTasks)
	assert.Equal(t, 1*time.Hour, cfg.swapDelay)
	assert.Equal(t, "./file.bin", cfg.swapFile)
	assert.Equal(t, time.Minute, cfg.holdDoneDelay)

	cfg2, _ := buildCfg(map[string]string{})

	assert.Equal(t, defaultHoldDoneTasks, cfg2.holdDoneTasks)
	assert.Equal(t, defaultSwapDelay, cfg2.swapDelay)
	assert.Equal(t, defaultSwapFile, cfg2.swapFile)
	assert.Equal(t, defaultHoldDoneDelay, cfg2.holdDoneDelay)
}

func TestStorage_Save(t *testing.T) {
	storage, storageErr := NewStorage(context.Background(), map[string]string{})
	if storageErr != nil {
		t.Error(storageErr)
	}

	id, _ := uuid.NewV7()

	msg := entity.MessageTask{
		Id:                 id,
		Provider:           "smtp",
		Message:            entity.Message{},
		Status:             entity.TaskStatusToProcess,
		MaxRetryCnt:        0,
		StartExecutionTime: time.Time{},
		LastExecutionTime:  time.Time{},
		CreatedAt:          time.Time{},
		RetryCnt:           0,
	}

	saveErr := storage.Save(msg)

	if saveErr != nil {
		t.Error(storageErr)
	}

	if row, exists := storage.data.Process[id]; exists {
		assert.Equal(t, id, row.Id)
		assert.Equal(t, msg.Provider, row.Provider)
		assert.Equal(t, msg.Status, row.Status)
	} else {
		t.Error(fmt.Errorf("message with id: %s does not exists", id.String()))
	}

	msg.Status = entity.TaskStatusDone

	saveErr = storage.Save(msg)

	if saveErr != nil {
		t.Error(storageErr)
	}

	if _, exists := storage.data.Process[id]; exists {
		t.Error(fmt.Errorf("message with id: %s exists but doesn't", id.String()))
	}

	if row, exists := storage.data.Closed[id]; exists {
		assert.Equal(t, id, row.Id)
		assert.Equal(t, msg.Provider, row.Provider)
		assert.Equal(t, msg.Status, row.Status)
	} else {
		t.Error(fmt.Errorf("message with id: %s does not exists", id.String()))
	}

}

func TestStorage_ClearExpired(t *testing.T) {
	storage, storageErr := NewStorage(context.Background(), map[string]string{
		"holdTasks": "1h",
		"holdDelay": "1h",
	})
	if storageErr != nil {
		t.Error(storageErr)
	}

	id, _ := uuid.NewV7()
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()

	msg := entity.MessageTask{
		Id:                 id,
		Provider:           "smtp",
		Message:            entity.Message{},
		Status:             entity.TaskStatusToProcess,
		MaxRetryCnt:        0,
		StartExecutionTime: time.Time{},
		LastExecutionTime:  time.Time{},
		CreatedAt:          time.Time{},
		RetryCnt:           0,
	}
	msg1 := entity.MessageTask{
		Id:                 id1,
		Provider:           "smtp",
		Message:            entity.Message{},
		Status:             entity.TaskStatusError,
		MaxRetryCnt:        0,
		StartExecutionTime: time.Time{},
		LastExecutionTime:  time.Now().Add(-20 * time.Minute),
		CreatedAt:          time.Time{},
		RetryCnt:           0,
	}
	msg2 := entity.MessageTask{
		Id:                 id2,
		Provider:           "smtp",
		Message:            entity.Message{},
		Status:             entity.TaskStatusDone,
		MaxRetryCnt:        0,
		StartExecutionTime: time.Time{},
		LastExecutionTime:  time.Now().Add(-3 * time.Hour),
		CreatedAt:          time.Time{},
		RetryCnt:           0,
	}

	_ = storage.Save(msg)
	_ = storage.Save(msg1)
	_ = storage.Save(msg2)

	if _, exists := storage.data.Process[msg.Id]; !exists {
		t.Error("message one must be in Process queue")
	}
	if _, exists := storage.data.Closed[msg1.Id]; !exists {
		t.Error("message two must be in Closed queue")
	}
	if _, exists := storage.data.Closed[msg2.Id]; !exists {
		t.Error("message three must be in Closed queue")
	}

	storage.ClearExpired()

	if _, exists := storage.data.Process[msg.Id]; !exists {
		t.Error("message one must be in Process queue after clear")
	}
	if _, exists := storage.data.Closed[msg1.Id]; !exists {
		t.Error("message two must be in Closed queue after clear")
	}
	if _, exists := storage.data.Closed[msg2.Id]; exists {
		t.Error("message three must be deleted from Closed queue after clear")
	}
}
