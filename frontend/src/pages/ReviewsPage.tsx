import { useState } from 'react'
import { motion } from 'framer-motion'
import {
  BookOpen,
  Plus,
  FileText,
  Star,
  Clock,
  CheckCircle,
  Trash2,
  Sparkles,
  ChevronDown,
  Loader2,
} from 'lucide-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import ReactMarkdown from 'react-markdown'
import toast from 'react-hot-toast'
import { paperApi, reviewApi } from '@/services/api'
import type { Paper, Review } from '@/types'

const statusConfig = {
  draft: { label: '草稿', color: 'text-dark-400 bg-dark-700' },
  generated: { label: '已生成', color: 'text-blue-400 bg-blue-400/10' },
  scored: { label: '已评分', color: 'text-green-400 bg-green-400/10' },
}

function errMsg(error: unknown) {
  const err = error as { response?: { data?: { error?: { message?: string } } } }
  return err.response?.data?.error?.message || '操作失败'
}

export default function ReviewsPage() {
  const queryClient = useQueryClient()
  const [showGenerator, setShowGenerator] = useState(false)
  const [expandedReview, setExpandedReview] = useState<string | null>(null)
  const [title, setTitle] = useState('')
  const [style, setStyle] = useState('academic')
  const [selectedIds, setSelectedIds] = useState<string[]>([])

  const { data: reviews = [] } = useQuery({
    queryKey: ['reviews'],
    queryFn: async () => {
      const res = await reviewApi.list({ page: 1, per_page: 50 })
      return (res.data.data || []) as Review[]
    },
  })

  const { data: papers = [] } = useQuery({
    queryKey: ['papers'],
    queryFn: async () => {
      const res = await paperApi.list({ page: 1, per_page: 50 })
      return (res.data.data || []) as Paper[]
    },
  })

  const generateMutation = useMutation({
    mutationFn: () => reviewApi.generate({ title, paper_ids: selectedIds, style }),
    onSuccess: () => {
      toast.success('综述已生成')
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
      setShowGenerator(false)
      setTitle('')
      setSelectedIds([])
    },
    onError: (e) => toast.error(errMsg(e)),
  })

  const scoreMutation = useMutation({
    mutationFn: (id: string) => reviewApi.score(id),
    onSuccess: () => {
      toast.success('评分完成')
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
    },
    onError: (e) => toast.error(errMsg(e)),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => reviewApi.delete(id),
    onSuccess: () => {
      toast.success('已删除')
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
    },
    onError: (e) => toast.error(errMsg(e)),
  })

  const togglePaper = (id: string) => {
    setSelectedIds((prev) => (prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]))
  }

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-3">
            <BookOpen className="w-7 h-7 text-primary-400" />
            论文综述
          </h1>
          <p className="text-dark-400 mt-1">AI 辅助撰写论文综述，智能评分与改进建议</p>
        </div>
        <button onClick={() => setShowGenerator(!showGenerator)} className="btn-primary flex items-center gap-2 self-start">
          <Plus className="w-5 h-5" />
          新建综述
        </button>
      </motion.div>

      {showGenerator && (
        <div className="glass-card p-6">
          <h3 className="font-semibold mb-4 flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-accent-400" />
            生成论文综述
          </h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-2">综述标题</label>
              <input type="text" className="input-field" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="例如：深度学习在医学影像分析中的应用综述" />
            </div>
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-2">选择论文</label>
              <div className="space-y-2 max-h-48 overflow-auto">
                {papers.length === 0 && <p className="text-dark-500 text-sm">请先在「我的论文」上传并分析论文</p>}
                {papers.map((p) => (
                  <label key={p.id} className="flex items-center gap-3 p-3 rounded-lg bg-dark-800/50 cursor-pointer">
                    <input type="checkbox" checked={selectedIds.includes(p.id)} onChange={() => togglePaper(p.id)} />
                    <span className="text-sm">{p.title}</span>
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-2">综述风格</label>
              <div className="grid grid-cols-3 gap-3">
                {[
                  { id: 'academic', name: '学术风格', desc: '正式、严谨' },
                  { id: 'summary', name: '总结风格', desc: '简洁、精炼' },
                  { id: 'comprehensive', name: '全面风格', desc: '详尽、深入' },
                ].map((s) => (
                  <button
                    key={s.id}
                    onClick={() => setStyle(s.id)}
                    className={`p-3 rounded-lg text-left ${style === s.id ? 'border border-primary-500/30 bg-primary-500/10' : 'bg-dark-800/50 border border-transparent'}`}
                  >
                    <p className="font-medium text-sm">{s.name}</p>
                    <p className="text-xs text-dark-500">{s.desc}</p>
                  </button>
                ))}
              </div>
            </div>
            <button
              className="btn-primary w-full flex items-center justify-center gap-2"
              disabled={!title || selectedIds.length === 0 || generateMutation.isPending}
              onClick={() => generateMutation.mutate()}
            >
              {generateMutation.isPending ? <Loader2 className="w-5 h-5 animate-spin" /> : <Sparkles className="w-5 h-5" />}
              生成综述
            </button>
          </div>
        </div>
      )}

      <div className="space-y-4">
        {reviews.map((review) => {
          const status = statusConfig[review.status] || statusConfig.generated
          const isExpanded = expandedReview === review.id
          return (
            <div key={review.id} className="glass-card overflow-hidden">
              <div className="p-6 cursor-pointer" onClick={() => setExpandedReview(isExpanded ? null : review.id)}>
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-2">
                      <span className={`px-2 py-1 rounded text-xs font-medium ${status.color}`}>{status.label}</span>
                      {review.category && <span className="category-badge">{review.category.name}</span>}
                      <span className="text-dark-500 text-sm flex items-center gap-1">
                        <FileText className="w-4 h-4" />
                        {review.paper_ids?.length || 0} 篇论文
                      </span>
                    </div>
                    <h3 className="text-lg font-semibold mb-2">{review.title}</h3>
                    <p className="text-dark-400 text-sm line-clamp-2">{review.abstract}</p>
                    <div className="flex items-center gap-4 mt-4">
                      <span className="text-dark-500 text-sm flex items-center gap-1">
                        <Clock className="w-4 h-4" />
                        {review.created_at?.slice(0, 10)}
                      </span>
                      {review.score && (
                        <span className="text-accent-400 text-sm flex items-center gap-1">
                          <Star className="w-4 h-4" />
                          评分: {review.score.overall_score}/100
                        </span>
                      )}
                    </div>
                  </div>
                  <ChevronDown className={`w-5 h-5 text-dark-500 transition-transform ${isExpanded ? 'rotate-180' : ''}`} />
                </div>
              </div>
              {isExpanded && (
                <div className="px-6 pb-6 border-t border-dark-700/50 pt-4 space-y-4">
                  <div className="prose prose-invert max-w-none text-sm text-dark-200">
                    <ReactMarkdown>{review.content}</ReactMarkdown>
                  </div>
                  {review.score && (
                    <div className="p-4 rounded-xl bg-dark-800/50">
                      <h4 className="font-medium mb-3 flex items-center gap-2">
                        <Star className="w-5 h-5 text-accent-400" />
                        AI 评分详情
                      </h4>
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                        {(review.score.criteria || []).map((criterion) => (
                          <div key={criterion.name} className="text-center">
                            <div className="text-2xl font-bold text-primary-400">{criterion.score}</div>
                            <div className="text-sm text-dark-400">{criterion.name}</div>
                          </div>
                        ))}
                      </div>
                      {(review.score.strengths || []).map((s) => (
                        <div key={s} className="p-3 rounded-lg bg-green-500/10 border border-green-500/20 mb-2">
                          <p className="text-sm text-green-400"><span className="font-medium">优点：</span>{s}</p>
                        </div>
                      ))}
                      {(review.score.suggestions || []).map((s) => (
                        <div key={s} className="p-3 rounded-lg bg-amber-500/10 border border-amber-500/20 mb-2">
                          <p className="text-sm text-amber-400"><span className="font-medium">建议：</span>{s}</p>
                        </div>
                      ))}
                    </div>
                  )}
                  <div className="flex items-center gap-3">
                    {review.status !== 'scored' && (
                      <button
                        className="btn-primary flex items-center gap-2 flex-1"
                        onClick={(e) => { e.stopPropagation(); scoreMutation.mutate(review.id) }}
                        disabled={scoreMutation.isPending}
                      >
                        {scoreMutation.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <CheckCircle className="w-4 h-4" />}
                        获取评分
                      </button>
                    )}
                    <button
                      className="p-3 rounded-xl bg-dark-800 text-red-400 hover:bg-red-500/10"
                      onClick={(e) => { e.stopPropagation(); deleteMutation.mutate(review.id) }}
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>

      {reviews.length === 0 && !showGenerator && (
        <div className="text-center py-16">
          <BookOpen className="w-16 h-16 text-dark-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">还没有论文综述</h3>
          <p className="text-dark-400 mb-6">选择您已上传的论文，让 AI 帮您撰写高质量综述</p>
          <button onClick={() => setShowGenerator(true)} className="btn-primary inline-flex items-center gap-2">
            <Sparkles className="w-5 h-5" />
            新建综述
          </button>
        </div>
      )}
    </div>
  )
}
