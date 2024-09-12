package services

import (
	"errors"
	"github.com/google/uuid"
	"zadanie_6105/src/models"
)

func (s *Service) CheckIfUserIsResponsibleForTenderByBidID(creatorUsername string, bidId string) (bool, error) {
	bid, err := s.getBidLastVersion(bidId)
	if err != nil {
		return false, err
	}
	tenderId := bid.TenderId.String()

	return s.CheckIfUserIsResponsibleForTender(creatorUsername, tenderId)
}

func (s *Service) getEmployeeByUsername(username string) (*models.Employee, error) {
	var employee models.Employee

	err := s.db.Where("username = ?", username).
		First(&employee).Error

	return &employee, err
}

func (s *Service) CheckIfUserIsResponsibleForBid(username string, bidId string) (bool, error) {
	bid, err := s.getBidLastVersion(bidId)
	if err != nil {
		return false, err
	}

	user, err := s.getEmployeeByUsername(username)
	if err != nil {
		return false, err
	}

	return bid.AuthorId == user.ID, nil
}

func (s *Service) getEmployee(id string) (*models.Employee, error) {
	var employee models.Employee

	employeeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid employee ID format")
	}

	err = s.db.Where("id = ?", employeeID).
		First(&employee).Error

	return &employee, err
}

func (s *Service) CheckEmployeeExistence(id string) (bool, error) {
	employee, err := s.getEmployee(id)
	if err != nil {
		return false, err
	}
	return employee != nil, nil
}

func (s *Service) CreateBid(bid *models.Bid) error {
	if err := s.db.Create(bid).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) GetBidsByUser(username string, limit, offset int) (*[]models.Bid, error) {
	var bids []models.Bid

	subQuery := s.db.Table("bids as b1").
		Select("MAX(version)").
		Where("b1.id = bids.id")
	query := s.db.Joins("JOIN employee ON employee.id = bids.author_id").
		Where("employee.username = ?", username).
		Where("version = (?)", subQuery)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&bids).Error

	return &bids, err
}

// TODO: figure it out

func (s *Service) GetBidsByTender(tenderId string, limit, offset int) (*[]models.Bid, error) {

	var bids []models.Bid

	subQuery := s.db.Table("bids as b1").
		Select("MAX(b1.version)").
		Where("b1.id = bids.id").
		Where("b1.status = ?", "Published")
	query := s.db.Where("version = (?)", subQuery).
		Where("status = ?", "Published").
		Where("tender_id = ?", tenderId)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&bids).Error

	return &bids, err
}

func (s *Service) getBidLastVersion(id string) (*models.Bid, error) {
	var bid models.Bid

	bidID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid bid ID format")
	}

	err = s.db.Where("id = ?", bidID).
		Order("version DESC").
		First(&bid).Error

	return &bid, err
}

func (s *Service) GetBidStatus(id string) (string, error) {
	bid, err := s.getBidLastVersion(id)
	if err != nil {
		return "", err
	}
	return bid.Status, nil
}

func (s *Service) UpdateBidStatus(id string, status string) (*models.Bid, error) {
	bid, err := s.getBidLastVersion(id)
	if err != nil {
		return nil, err
	}
	bid.Status = status
	if err := s.db.Save(bid).Error; err != nil {
		return nil, err
	}
	return bid, nil
}

func (s *Service) UpdateBid(id string, edit *models.BidEdit) (*models.Bid, error) {
	bid, err := s.getBidLastVersion(id)
	if err != nil {
		return nil, err
	}

	newBid := models.Bid{
		ID:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      bid.Status,
		TenderId:    bid.TenderId,
		AuthorType:  bid.AuthorType,
		AuthorId:    bid.AuthorId,
		Version:     bid.Version + 1,
	}
	if edit.Name != "" {
		newBid.Name = edit.Name
	}
	if edit.Description != "" {
		newBid.Description = edit.Description
	}

	if err := s.db.Create(&newBid).Error; err != nil {
		return nil, err
	}

	return &newBid, nil
}

// TODO: figure it out

func (s *Service) SubmitBid(bidId string, decision bool) (*models.Bid, error) {
	bid, err := s.getBidLastVersion(bidId)
	if err != nil {
		return nil, err
	}
	tenderId := bid.TenderId.String()
	if decision {
		tender, err := s.getTenderLastVersion(tenderId)
		if err != nil {
			return nil, err
		}
		tender.Status = "Closed"
		if err := s.db.Save(tender).Error; err != nil {
			return nil, err
		}
	} else {
		bid.Status = "Cancelled"
		if err := s.db.Save(bid).Error; err != nil {
			return nil, err
		}
	}
	return bid, nil
}

func (s *Service) RollbackBid(id string, version int32) (*models.Bid, error) {
	var bid models.Bid

	err := s.db.Where("id = ? AND version = ?", id, version).First(&bid).Error
	if err != nil {
		return nil, err
	}

	lastBid, err := s.getBidLastVersion(id)
	if err != nil {
		return nil, err
	}

	newBid := models.Bid{
		ID:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      bid.Status,
		TenderId:    bid.TenderId,
		AuthorType:  bid.AuthorType,
		AuthorId:    bid.AuthorId,
		Version:     lastBid.Version + 1,
	}

	if err := s.db.Create(&newBid).Error; err != nil {
		return nil, err
	}
	return &newBid, nil
}
