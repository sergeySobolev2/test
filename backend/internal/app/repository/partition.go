package repository

import (
	"partitionlab/internal/app/ds"
	"strings"
)

func (r *Repository) GetActivePartitions() ([]ds.Partition, error) {
	var partitions []ds.Partition
	err := r.db.Where("is_active = ?", true).Find(&partitions).Error
	if err != nil {
		return nil, err
	}
	return partitions, nil
}

func (r *Repository) SearchPartitions(query string) ([]ds.Partition, error) {
	var partitions []ds.Partition
	q := strings.ToLower(query)
	err := r.db.Where("is_active = ? AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ?)",
		true, "%"+q+"%", "%"+q+"%").Find(&partitions).Error
	if err != nil {
		return nil, err
	}
	return partitions, nil
}

func (r *Repository) GetPartition(id uint) (ds.Partition, error) {
	var Partition ds.Partition
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&Partition).Error
	if err != nil {
		return ds.Partition{}, err
	}
	return Partition, nil
}

func (r *Repository) GetPartitionAny(id uint) (ds.Partition, error) {
	var Partition ds.Partition
	err := r.db.Where("id = ?", id).First(&Partition).Error
	if err != nil {
		return ds.Partition{}, err
	}
	return Partition, nil
}

// FilterPartitions фильтрует по названию (LIKE по title/description) и активности (если задана)
func (r *Repository) FilterPartitions(title string, active *bool) ([]ds.Partition, error) {
	var partitions []ds.Partition
	tx := r.db.Model(&ds.Partition{})
	if active != nil {
		tx = tx.Where("is_active = ?", *active)
	}
	if title != "" {
		q := strings.ToLower(title)
		tx = tx.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	if err := tx.Find(&partitions).Error; err != nil {
		return nil, err
	}
	return partitions, nil
}

func (r *Repository) CreatePartition(s *ds.Partition) error {
	return r.db.Create(s).Error
}

func (r *Repository) UpdatePartition(id uint, upd ds.Partition) (ds.Partition, error) {
	var s ds.Partition
	if err := r.db.First(&s, id).Error; err != nil {
		return s, err
	}
	// системные поля не трогаем: ID, ImageURL меняется отдельным методом
	s.Title = upd.Title
	s.Category = upd.Category
	s.Description = upd.Description
	s.NoiseReduction = upd.NoiseReduction
	s.Thickness = upd.Thickness
	s.Material = upd.Material
	s.PricePerSqm = upd.PricePerSqm
	s.IsActive = upd.IsActive
	if err := r.db.Save(&s).Error; err != nil {
		return s, err
	}
	return s, nil
}

func (r *Repository) UpdatePartitionImage(id uint, key string) error {
	return r.db.Model(&ds.Partition{}).Where("id = ?", id).Update("image_url", key).Error
}

func (r *Repository) DeletePartition(id uint) error {
	// Мягкое удаление: помечаем симптом как неактивный, чтобы не нарушать FK в request_symptoms
	return r.db.Model(&ds.Partition{}).Where("id = ?", id).Update("is_active", false).Error
}
