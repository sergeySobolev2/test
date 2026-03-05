package repository

import (
	"partitionlab/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetOrCreateDraftCalculation(userID uint) (*ds.Calculation, error) {
	var request ds.Calculation
	err := r.db.Where("user_id = ? AND status = ?", userID, "черновик").First(&request).Error
	if err == nil {
		return &request, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	request = ds.Calculation{
		UserID: userID,
		Status: "черновик",
	}
	err = r.db.Create(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *Repository) AddPartitionToCalculation(calculationID, partitionID uint) error {
	CalculationItem := ds.CalculationItem{
		CalculationID: calculationID,
		PartitionID:   partitionID,
	}
	err := r.db.Create(&CalculationItem).Error
	if err != nil {
		// Ignore duplicate key errors
		return nil
	}
	return nil
}

func (r *Repository) DeleteCalculation(requestID uint) error {

	result := r.db.Exec("UPDATE calculations SET status = 'удалён' WHERE id = ?", requestID)
	return result.Error
}

func (r *Repository) GetDraftCalculation(userID uint) (*ds.Calculation, error) {
	var request ds.Calculation
	err := r.db.Where("user_id = ? AND status = ?", userID, "черновик").First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *Repository) GetDraftPartitions(userID uint) ([]ds.Partition, *ds.Calculation, error) {
	request, err := r.GetDraftCalculation(userID)
	if err != nil {
		return nil, nil, err
	}

	var partitions []ds.Partition
	err = r.db.Table("partitions").
		Joins("JOIN calculation_items ON calculation_items.partition_id = partitions.id").
		Where("calculation_items.calculation_id = ? AND partitions.is_active = ?", request.ID, true).
		Find(&partitions).Error
	if err != nil {
		return nil, nil, err
	}

	return partitions, request, nil
}

func (r *Repository) CountDraftItems(userID uint) (int64, error) {
	request, err := r.GetDraftCalculation(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	var count int64
	err = r.db.Model(&ds.CalculationItem{}).Where("calculation_id = ?", request.ID).Count(&count).Error
	return count, err
}

func (r *Repository) ListRequestsByUser(userID uint) ([]ds.Calculation, error) {
	var requests []ds.Calculation
	err := r.db.Where("user_id = ? AND status NOT IN ?", userID, []string{"удален", "удалён", "черновик"}).Order("id DESC").Find(&requests).Error
	return requests, err
}

func (r *Repository) GetCalculationByID(id uint) (*ds.Calculation, error) {
	var request ds.Calculation
	err := r.db.Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *Repository) GetCalculationPartitions(calculationID uint) ([]ds.Partition, error) {
	var partitions []ds.Partition
	err := r.db.Table("partitions").
		Joins("JOIN calculation_items ON calculation_items.partition_id = partitions.id").
		Where("calculation_items.calculation_id = ? AND partitions.is_active = ?", calculationID, true).
		Find(&partitions).Error
	return partitions, err
}

type PartitionWithParams struct {
	ds.Partition
	Quantity *int    `json:"quantity"`
	IsMain   *bool   `json:"is_main"`
	Comment  *string `json:"comment"`
}

type CalculationWithPartitions struct {
	Calculation ds.Calculation        `json:"calculation"`
	Partitions  []PartitionWithParams `json:"partitions"`
}

func (r *Repository) GetCalculationWithPartitions(id uint) (*CalculationWithPartitions, error) {
	dr, err := r.GetCalculationByID(id)
	if err != nil {
		return nil, err
	}

	// Получаем перегородки с параметрами из calculation_items
	var results []struct {
		ds.Partition
		Quantity *int    `gorm:"column:quantity"`
		IsMain   *bool   `gorm:"column:is_main"`
		Comment  *string `gorm:"column:comment"`
	}

	err = r.db.Table("partitions").
		Select("partitions.*, calculation_items.quantity, calculation_items.is_main, calculation_items.comment").
		Joins("JOIN calculation_items ON calculation_items.partition_id = partitions.id").
		Where("calculation_items.calculation_id = ? AND partitions.is_active = ?", id, true).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// Преобразуем в PartitionWithParams
	partitions := make([]PartitionWithParams, len(results))
	for i, r := range results {
		partitions[i] = PartitionWithParams{
			Partition: r.Partition,
			Quantity:  r.Quantity,
			IsMain:    r.IsMain,
			Comment:   r.Comment,
		}
	}

	return &CalculationWithPartitions{Calculation: *dr, Partitions: partitions}, nil
}

func (r *Repository) CountCalculationPartitions(calculationID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&ds.CalculationItem{}).Where("calculation_id = ?", calculationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountCalculationPartitionsWithComment считает M-M записи с непустым комментарием
func (r *Repository) CountCalculationPartitionsWithComment(calculationID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&ds.CalculationItem{}).
		Where("calculation_id = ? AND TRIM(COALESCE(comment, '')) <> ''", calculationID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) SetFormed(id uint, when time.Time) error {
	return r.db.Model(&ds.Calculation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    "сформирован",
		"formed_at": when,
	}).Error
}

func (r *Repository) SetCompleted(id uint, moderatorID uint, status string, when time.Time, RoomArea *float64, NoiseReductionDB *float64, RequiredThickness *float64, ExpertComment *string) error {
	updates := map[string]interface{}{
		"status":       status,
		"completed_at": when,
		"moderator_id": moderatorID,
	}
	if RoomArea != nil {
		updates["room_area"] = *RoomArea
	}
	if NoiseReductionDB != nil {
		updates["noise_reduction_db"] = *NoiseReductionDB
	}
	if RequiredThickness != nil {
		updates["required_thickness"] = *RequiredThickness
	}
	if ExpertComment != nil {
		updates["expert_comment"] = *ExpertComment
	}
	return r.db.Model(&ds.Calculation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) IsCalculationOwner(userID, requestID uint) (bool, error) {
	var cnt int64
	if err := r.db.Model(&ds.Calculation{}).Where("id = ? AND user_id = ?", requestID, userID).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// ListRequestsWithFilters фильтрует по статусу и диапазону дат формирования (formed_at)
func (r *Repository) ListRequestsWithFilters(userID uint, status *string, from, to *string) ([]ds.Calculation, error) {
	var list []ds.Calculation
	tx := r.db.Model(&ds.Calculation{}).
		Where("user_id = ? AND status NOT IN ?", userID, []string{"удален", "удалён"}).
		Preload("User").
		Preload("Moderator")
	if status != nil && *status != "" {
		tx = tx.Where("status = ?", *status)
	}
	if from != nil && *from != "" {
		tx = tx.Where("formed_at >= ?", *from)
	}
	if to != nil && *to != "" {
		tx = tx.Where("formed_at <= ?", *to)
	}
	if err := tx.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListAllRequestsWithFilters фильтрует все заявки (для модератора)
func (r *Repository) ListAllRequestsWithFilters(status *string, from, to *string) ([]ds.Calculation, error) {
	var list []ds.Calculation
	tx := r.db.Model(&ds.Calculation{}).
		Where("status NOT IN ?", []string{"удален", "удалён"}).
		Preload("User").
		Preload("Moderator")
	if status != nil && *status != "" {
		tx = tx.Where("status = ?", *status)
	}
	if from != nil && *from != "" {
		tx = tx.Where("formed_at >= ?", *from)
	}
	if to != nil && *to != "" {
		tx = tx.Where("formed_at <= ?", *to)
	}
	if err := tx.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateCalculationFields изменяет разрешенные поля темы (например, room_area, expert_comment)
func (r *Repository) UpdateCalculationFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&ds.Calculation{}).Where("id = ?", id).Updates(fields).Error
}

// ChangeRequestStatus меняет статус с проверкой допустимых переходов вне этого слоя (проверять в handler)
func (r *Repository) ChangeRequestStatus(id uint, status string) error {
	return r.db.Model(&ds.Calculation{}).Where("id = ?", id).Update("status", status).Error
}
