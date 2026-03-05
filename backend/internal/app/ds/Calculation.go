package ds

import (
	"database/sql"
	"time"
)

// Calculation - расчет звукоизоляции помещения
type Calculation struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"not null" json:"user_id"`
	Status            string         `gorm:"type:varchar(20);not null" json:"status"` // черновик, сформирован, завершен, отклонен
	CreatedAt         time.Time      `gorm:"not null" json:"created_at"`
	FormedAt          *time.Time     `json:"formed_at"`
	CompletedAt       *time.Time     `json:"completed_at"`
	ModeratorID       *uint          `json:"moderator_id"`
	RoomArea          *float64       `gorm:"type:decimal(6,2)" json:"room_area"`               // Площадь помещения (м²)
	NoiseReductionDB  *float64       `gorm:"type:decimal(5,2)" json:"noise_reduction_db"`      // Требуемое снижение шума (дБ)
	RequiredThickness *float64       `gorm:"type:decimal(5,2)" json:"required_thickness"`      // Рекомендуемая толщина (см)
	ExpertComment     sql.NullString `gorm:"type:text" json:"expert_comment"`                  // Комментарий эксперта
	CalculationsCount int            `gorm:"type:integer;default:0" json:"calculations_count"` // Количество выполненных расчетов

	User      User  `gorm:"foreignKey:UserID" json:"user"`
	Moderator *User `gorm:"foreignKey:ModeratorID" json:"moderator"`
}
