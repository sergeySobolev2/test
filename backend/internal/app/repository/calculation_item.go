package repository

import "partitionlab/internal/app/ds"

func (r *Repository) DeleteCalculationItem(calculationID, partitionID uint) error {
	return r.db.Where("calculation_id = ? AND partition_id = ?", calculationID, partitionID).Delete(&ds.CalculationItem{}).Error
}

func (r *Repository) UpdateCalculationItem(calculationID, partitionID uint, quantity *int, comment *string, isMain *bool) error {
	updates := map[string]interface{}{}
	if quantity != nil {
		updates["quantity"] = *quantity
	}
	if comment != nil {
		updates["comment"] = *comment
	}
	if isMain != nil {
		updates["is_main"] = *isMain
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&ds.CalculationItem{}).Where("calculation_id = ? AND partition_id = ?", calculationID, partitionID).Updates(updates).Error
}
