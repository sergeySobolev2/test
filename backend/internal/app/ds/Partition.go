package ds

// Partition - тип перегородки для звукоизоляции
type Partition struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Title          string `gorm:"type:varchar(100);not null" json:"title"`    // Название: "Гипсокартон 12мм"
	Category       string `gorm:"type:varchar(100)" json:"category"`          // Категория: "Легкие", "Тяжелые"
	Description    string `gorm:"type:text" json:"description"`               // Описание конструкции
	NoiseReduction string `gorm:"type:varchar(50)" json:"noise_reduction"`    // Снижение шума: "30-35 дБ"
	Thickness      string `gorm:"type:varchar(50)" json:"thickness"`          // Толщина: "10-15 см"
	Material       string `gorm:"type:varchar(100)" json:"material"`          // Материал: "ГКЛ + минвата"
	PricePerSqm    string `gorm:"type:varchar(50)" json:"price_per_sqm"`      // Цена: "500-800 руб/м²"
	ImageURL       string `gorm:"type:varchar(200)" json:"image_url"`         // URL изображения
	IsActive       bool   `gorm:"type:boolean;default:true" json:"is_active"` // Активность
}
