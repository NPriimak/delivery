package courierrepo

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
)

func DomainToDTO(aggregate *courier.Courier) CourierDTO {
	var courierDTO CourierDTO
	courierDTO.ID = aggregate.ID()
	courierDTO.Name = aggregate.Name()
	courierDTO.Speed = aggregate.Speed()
	courierDTO.StoragePlaces = make([]*StoragePlaceDTO, 0)
	for _, storagePlace := range aggregate.StoragePlaces() {
		storagePlaceDTO := &StoragePlaceDTO{
			ID:          storagePlace.ID(),
			OrderID:     storagePlace.OrderID(),
			Name:        storagePlace.Name(),
			TotalVolume: storagePlace.TotalVolume(),
			CourierID:   aggregate.ID(),
		}
		courierDTO.StoragePlaces = append(courierDTO.StoragePlaces, storagePlaceDTO)
	}
	courierDTO.Location = LocationDTO{
		X: aggregate.Location().X(),
		Y: aggregate.Location().Y(),
	}
	return courierDTO
}

func DtoToDomain(dto CourierDTO) *courier.Courier {
	var aggregate *courier.Courier
	var storagePlaces []*courier.StoragePlace
	for _, dtoStoragePlace := range dto.StoragePlaces {
		item := courier.RestoreStoragePlace(dtoStoragePlace.ID, dtoStoragePlace.Name,
			dtoStoragePlace.TotalVolume, dtoStoragePlace.OrderID)
		storagePlaces = append(storagePlaces, item)
	}
	location, _ := kernel.NewLocation(dto.Location.X, dto.Location.Y)
	aggregate = courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, location, storagePlaces)
	return aggregate
}
