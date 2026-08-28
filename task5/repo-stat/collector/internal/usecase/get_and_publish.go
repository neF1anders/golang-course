package usecase

import "context"

type GetAndPublishUseCase struct {
	getSubInfoUC *GetSubInfoUseCase
	publisher    Publisher
}

func NewGetAndPublishUseCase(
	get *GetSubInfoUseCase,
	publisher Publisher,
) *GetAndPublishUseCase {
	return &GetAndPublishUseCase{
		getSubInfoUC: get,
		publisher:    publisher,
	}
}

func (uc *GetAndPublishUseCase) Execute(ctx context.Context) error {
	data, err := uc.getSubInfoUC.Execute(ctx)
	if err != nil {
		return err
	}
	return uc.publisher.Publish(ctx, data)
}
