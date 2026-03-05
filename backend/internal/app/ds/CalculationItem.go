package ds

// CalculationItem - связь между расчетом и перегородкой (M:N)
type CalculationItem struct {
	CalculationID uint   `gorm:"primaryKey;autoIncrement:false" json:"calculation_id"`
	PartitionID   uint   `gorm:"primaryKey;autoIncrement:false" json:"partition_id"`
	Quantity      *int   `gorm:"type:integer" json:"quantity"`              // Количество элементов
	IsMain        bool   `gorm:"type:boolean;default:false" json:"is_main"` // Основной тип перегородки
	Comment       string `gorm:"type:text" json:"comment"`                  // Примечание

	Calculation Calculation `gorm:"foreignKey:CalculationID" json:"calculation"`
	Partition   Partition   `gorm:"foreignKey:PartitionID" json:"partition"`
}
