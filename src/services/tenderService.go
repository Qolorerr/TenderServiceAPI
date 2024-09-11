package services

import (
	"errors"
	"zadanie_6105/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenderService struct {
	db *gorm.DB
}

func NewTenderService(db *gorm.DB) *TenderService {
	return &TenderService{db: db}
}

func (s *TenderService) CheckIfUserIsResponsible(creatorUsername string, organizationID string) (bool, error) {
	var count int64

	err := s.db.Table("organization_responsible").
		Select("COUNT(*)").
		Joins("JOIN employee ON employee.id = organization_responsible.user_id").
		Joins("JOIN organization ON organization.id = organization_responsible.organization_id").
		Where("employee.username = ?", creatorUsername).
		Where("organization.id = ?", organizationID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, errors.New("the user is not responsible for the organization")
}

func (s *TenderService) CheckIfUserIsResponsibleForTender(creatorUsername string, tenderID string) (bool, error) {
	tender, err := s.getTenderLastVersion(tenderID)
	if err != nil {
		return false, err
	}
	organizationID := tender.OrganizationId

	return s.CheckIfUserIsResponsible(creatorUsername, organizationID)
}

func (s *TenderService) GetTenders(serviceTypes []string, limit, offset int) (*[]models.Tender, error) {
	var tenders []models.Tender

	subQuery := s.db.Table("tenders as t1").
		Select("MAX(t1.version)").
		Where("t1.id = tenders.id").
		Where("t1.status = ?", "Published")
	query := s.db.Where("version = (?)", subQuery).
		Where("status = ?", "Published")

	if len(serviceTypes) > 0 {
		query = query.Where("service_type IN ?", serviceTypes)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&tenders).Error

	return &tenders, err
}

func (s *TenderService) CreateTender(tender *models.Tender) error {
	if err := s.db.Create(tender).Error; err != nil {
		return err
	}
	return nil
}

func (s *TenderService) GetTendersByUser(username string, limit, offset int) (*[]models.Tender, error) {
	var tenders []models.Tender

	subQuery := s.db.Table("tenders as t1").
		Select("MAX(version)").
		Where("t1.id = tenders.id")
	query := s.db.Joins("JOIN organization_responsible ON organization_responsible.organization_id = tenders.organization_id").
		Joins("JOIN employee ON employee.id = organization_responsible.user_id").
		Where("employee.username = ?", username).
		Where("version = (?)", subQuery)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&tenders).Error

	return &tenders, err
}

func (s *TenderService) getTenderLastVersion(id string) (*models.Tender, error) {
	var tender models.Tender

	tenderID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid tender ID format")
	}

	err = s.db.Where("id = ?", tenderID).
		Order("version DESC").
		First(&tender).Error

	return &tender, err
}

func (s *TenderService) GetTenderStatus(id string) (string, error) {
	tender, err := s.getTenderLastVersion(id)
	if err != nil {
		return "", err
	}
	return tender.Status, nil
}

func (s *TenderService) UpdateTenderStatus(id string, status string) (*models.Tender, error) {
	tender, err := s.getTenderLastVersion(id)
	if err != nil {
		return nil, err
	}
	tender.Status = status
	if err := s.db.Save(tender).Error; err != nil {
		return nil, err
	}
	return tender, nil
}

func (s *TenderService) UpdateTender(id string, edit *models.TenderEdit) (*models.Tender, error) {
	tender, err := s.getTenderLastVersion(id)
	if err != nil {
		return nil, err
	}

	newTender := models.Tender{
		ID:             tender.ID,
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tender.Status,
		OrganizationId: tender.OrganizationId,
		Version:        tender.Version + 1,
	}
	if edit.Name != "" {
		newTender.Name = edit.Name
	}
	if edit.Description != "" {
		newTender.Description = edit.Description
	}
	if edit.ServiceType != "" {
		newTender.ServiceType = edit.ServiceType
	}

	if err := s.db.Create(&newTender).Error; err != nil {
		return nil, err
	}

	return &newTender, nil
}

func (s *TenderService) RollbackTender(id string, version int32) (*models.Tender, error) {
	var tender models.Tender

	err := s.db.Where("id = ? AND version = ?", id, version).First(&tender).Error
	if err != nil {
		return nil, err
	}

	lastTender, err := s.getTenderLastVersion(id)
	if err != nil {
		return nil, err
	}

	newTender := models.Tender{
		ID:             tender.ID,
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tender.Status,
		OrganizationId: tender.OrganizationId,
		Version:        lastTender.Version + 1,
	}

	if err := s.db.Create(&newTender).Error; err != nil {
		return nil, err
	}
	return &newTender, nil
}
