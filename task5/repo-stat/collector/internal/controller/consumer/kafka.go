package consumer

import (
	"context"
	"encoding/json"
	"repo-stat/collector/internal/usecase"
)

type OrderMessageHandler struct {
	collectAndPublishUC *usecase.GetAndPublishUseCase
}

func NewOrderMessageHandler(collectAndPublishUC *usecase.GetAndPublishUseCase) *OrderMessageHandler {
	return &OrderMessageHandler{collectAndPublishUC: collectAndPublishUC}
}
func (h *OrderMessageHandler) Handle(ctx context.Context, raw []byte) error {
	var cmdDTO struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(raw, &cmdDTO); err != nil {
		return err
	}
	if err := h.collectAndPublishUC.Execute(ctx); err != nil {
		return err
	}
	return nil
}
