package http

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
)

var _ servers.StrictServerInterface = &serverHandlers{}

type serverHandlers struct {
	allCouriersQueryHandler      queries.AllCouriersQueryHandler
	incompleteOrdersQueryHandler queries.IncompleteOrdersQueryHandler
	createCourierCommandHandler  commands.CreateCourierCommandHandler
	createOrderCommandHandler    commands.CreateOrderCommandHandler
}

func NewServerHandlers(
	allCouriersQueryHandler queries.AllCouriersQueryHandler,
	incompleteOrdersQueryHandler queries.IncompleteOrdersQueryHandler,
	createCourierCommandHandler commands.CreateCourierCommandHandler,
	createOrderCommandHandler commands.CreateOrderCommandHandler,
) (servers.StrictServerInterface, error) {

	if allCouriersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("allCouriersQueryHandler")
	}

	if incompleteOrdersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("incompleteOrdersQueryHandler")
	}

	if createCourierCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createCourierCommandHandler")
	}

	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}

	return &serverHandlers{
		allCouriersQueryHandler:      allCouriersQueryHandler,
		incompleteOrdersQueryHandler: incompleteOrdersQueryHandler,
		createCourierCommandHandler:  createCourierCommandHandler,
		createOrderCommandHandler:    createOrderCommandHandler,
	}, nil
}

func (s serverHandlers) GetCouriers(ctx context.Context, _ servers.GetCouriersRequestObject) (servers.GetCouriersResponseObject, error) {
	couriers, err := s.allCouriersQueryHandler.Handle(ctx)
	if err != nil {
		return nil, err
	}

	responseCouriers := servers.GetCouriers200JSONResponse{}
	for _, courier := range couriers {
		responseCouriers = append(responseCouriers,
			servers.Courier{
				Id: courier.CourierID,
				Location: servers.Location{
					X: courier.LocationX,
					Y: courier.LocationY,
				},
				Name: courier.Name,
			})
	}

	return responseCouriers, nil
}

func (s serverHandlers) CreateCourier(ctx context.Context, request servers.CreateCourierRequestObject) (servers.CreateCourierResponseObject, error) {
	cmd, err := commands.NewCreateCourierCommand(request.Body.Name, request.Body.Speed)
	if err != nil {
		return nil, err
	}

	err = s.createCourierCommandHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return servers.CreateCourier201Response{}, nil
}

func (s serverHandlers) CreateOrder(ctx context.Context, _ servers.CreateOrderRequestObject) (servers.CreateOrderResponseObject, error) {

	// cmd, err := commands.NewCreateOrderCommand(request)
	// if err != nil {
	// 	return nil, err
	// }

	cmd, err := commands.NewCreateOrderCommand(uuid.New(), "Несуществующая", 5)
	if err != nil {
		return nil, err
	}

	err = s.createOrderCommandHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return servers.CreateOrder201Response{}, nil
}

func (s serverHandlers) GetOrders(ctx context.Context, _ servers.GetOrdersRequestObject) (servers.GetOrdersResponseObject, error) {
	orders, err := s.incompleteOrdersQueryHandler.Handle(ctx)
	if err != nil {
		return nil, err
	}

	responseOrders := servers.GetOrders200JSONResponse{}
	for _, order := range orders {
		responseOrders = append(responseOrders,
			servers.Order{
				Id: order.OrderID,
				Location: servers.Location{
					X: order.LocationX,
					Y: order.LocationY,
				},
			})
	}

	return responseOrders, nil
}
