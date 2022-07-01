package service

import (
	"bookman/entity"
	"bookman/events"
	"bookman/repository"
	"time"
)

type RentalService interface {
	StartRentBook(bookID, userID int) (*entity.RentalStatus, error)
	GetRentStatus(bookID int) (*entity.RentalStatus, error)
	EndRentBook(bookID, userID int) error
}

type rentalService struct {
	rentalRepo  repository.RentalRepo
	eventSender events.EventSender
}

func NewRentalService(rentalRepo repository.RentalRepo, eventSender events.EventSender) RentalService {
	return &rentalService{
		rentalRepo:  rentalRepo,
		eventSender: eventSender,
	}
}

func (s *rentalService) StartRentBook(bookID, userID int) (*entity.RentalStatus, error) {
	rental := entity.RentalStatus{
		BookID: bookID,
		UserID: userID,
		Status: "대여중",
		Start:  time.Now().Truncate(time.Second),
	}
	saveRental, err := s.rentalRepo.SaveRental(&rental)
	if err != nil {
		return nil, err
	}
	saveRental.Start = saveRental.Start.UTC()

	defer func() {
		go s.eventSender(events.EventStartRental, saveRental)
	}()

	return saveRental, nil
}

func (s *rentalService) GetRentStatus(bookID int) (*entity.RentalStatus, error) {
	return s.rentalRepo.FetchRentalStatusByBookID(bookID)
}

func (s *rentalService) EndRentBook(bookID, userID int) error {
	rental := entity.RentalStatus{
		BookID: bookID,
		UserID: userID,
	}

	foundRentalStatus, err := s.rentalRepo.FetchRentalStatus(&rental)
	if err != nil {
		return err
	}
	foundRentalStatus.Start = foundRentalStatus.Start.UTC()

	err = s.rentalRepo.DeleteRental(&rental)
	if err != nil {
		return err
	}

	defer func() {
		go s.eventSender(events.EventEndRental, foundRentalStatus)
	}()

	return nil
}
