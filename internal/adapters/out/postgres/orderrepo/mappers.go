package orderrepo

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
)

func DomainToDTO(aggregate *order.Order) OrderDTO {
	var orderDTO OrderDTO
	orderDTO.ID = aggregate.ID()
	orderDTO.CourierID = aggregate.CourierID()
	orderDTO.Location = LocationDTO{
		X: aggregate.Location().X(),
		Y: aggregate.Location().Y(),
	}
	orderDTO.Volume = aggregate.Volume()
	orderDTO.Status = aggregate.Status()
	return orderDTO
}

func DtoToDomain(dto OrderDTO) *order.Order {
	var aggregate *order.Order
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	aggregate = order.RestoreOrder(dto.ID, dto.CourierID, location, dto.Volume, dto.Status)
	return aggregate
}
