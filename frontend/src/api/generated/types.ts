export type UpdateCalculationPayload = {
  status?: string;
  room_area?: number;
  noise_reduction_db?: number;
  expert_comment?: string;
};
// Generated from docs/swagger.json (simplified subset)

export type User = {
  id: number
  login: string
  is_moderator: boolean
}

export type partition = {
  id: number
  title: string
  description?: string
  is_active: boolean
  image_url?: string
  public_image_url?: string
  category?: string
  noise_reduction?: string
  thickness?: string
  material?: string
  price_per_sqm?: string
}

export type Calculation = {
  id: number
  user_id: number
  status: string
  created_at: string
  formed_at?: string
  completed_at?: string
  moderator_id?: number
  room_area?: number
  noise_reduction_db?: number
  required_thickness?: number
  expert_comment?: { string: string, valid: boolean } | null
}

export type RequestWithCommentsCount = Calculation & { comments_count: number }

export type RequestPartitionUpdate = {
  calculation_id: number
  partition_id: number
  quantity?: number
  comment?: string
  is_main?: boolean
}
