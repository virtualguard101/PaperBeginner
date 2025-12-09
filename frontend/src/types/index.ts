export interface User {
  id: string
  email: string
  name: string
  avatar?: string
  role: 'user' | 'admin'
  is_active: boolean
  preferences?: UserPreferences
  created_at: string
  last_login_at?: string
}

export interface UserPreferences {
  preferred_llm?: string
  research_fields?: string[]
  notify_on_trending?: boolean
  language?: string
  openai_api_key?: string
  anthropic_api_key?: string
  deepseek_api_key?: string
}

export interface Paper {
  id: string
  title: string
  authors: string[]
  abstract: string
  keywords: string[]
  published_at?: string
  venue?: string
  doi?: string
  ccf_category?: CCFCategory
  ccf_rank?: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  analyses?: PaperAnalysis[]
  created_at: string
}

export interface PaperAnalysis {
  id: string
  paper_id: string
  analysis_type: AnalysisType
  content: string
  structured_data?: Record<string, unknown>
  llm_provider: string
  llm_model: string
  tokens_used: number
  created_at: string
}

export type AnalysisType = 
  | 'summary'
  | 'methodology'
  | 'contributions'
  | 'citations'
  | 'comparison'
  | 'strengths'
  | 'weaknesses'

export interface CCFCategory {
  id: number
  name: string
  name_en: string
}

export interface TrendingItem {
  id: string
  source: 'github' | 'ccf' | 'arxiv'
  source_id: string
  title: string
  description: string
  url: string
  stars?: number
  forks?: number
  language?: string
  topics: string[]
  ccf_category?: CCFCategory
  trend_score: number
  crawled_at: string
}

export interface TrendReport {
  id: string
  period: string
  category?: CCFCategory
  summary: string
  highlights: TrendHighlight[]
  generated_by: string
  created_at: string
}

export interface TrendHighlight {
  title: string
  description: string
  sources: string[]
  importance: number
}

export interface LearningPath {
  id: string
  title: string
  description: string
  category?: CCFCategory
  stages: LearningStage[]
  prerequisites: string[]
  estimated_time: string
  difficulty: 'beginner' | 'intermediate' | 'advanced'
  created_at: string
}

export interface LearningStage {
  order: number
  title: string
  description: string
  resources: LearningResource[]
  duration: string
  skills: string[]
}

export interface LearningResource {
  type: 'documentation' | 'course' | 'book' | 'video' | 'tutorial'
  title: string
  url: string
  provider: string
  language: 'en' | 'zh'
  is_free: boolean
  description?: string
}

export interface Review {
  id: string
  title: string
  abstract: string
  content: string
  paper_ids: string[]
  category?: CCFCategory
  status: 'draft' | 'generated' | 'scored'
  score?: ReviewScore
  generated_by: string
  created_at: string
  updated_at: string
}

export interface ReviewScore {
  overall_score: number
  criteria: ScoreCriterion[]
  strengths: string[]
  weaknesses: string[]
  suggestions: string[]
  generated_at: string
}

export interface ScoreCriterion {
  name: string
  score: number
  weight: number
  description: string
  feedback: string
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: string
  }
  meta?: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
}

