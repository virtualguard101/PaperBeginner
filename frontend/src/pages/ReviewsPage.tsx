import { useState } from 'react'
import { motion } from 'framer-motion'
import { 
  BookOpen, 
  Plus,
  FileText,
  Star,
  Clock,
  CheckCircle,
  Edit3,
  Trash2,
  Sparkles,
  ChevronDown
} from 'lucide-react'

// Mock data
const mockReviews = [
  {
    id: '1',
    title: 'Transformer 架构在 NLP 中的应用综述',
    abstract: '本文综述了 Transformer 架构自提出以来在自然语言处理领域的发展与应用...',
    status: 'scored',
    category: '人工智能',
    paperCount: 5,
    score: 85,
    createdAt: '2024-01-10',
  },
  {
    id: '2',
    title: '联邦学习隐私保护机制研究',
    abstract: '随着数据隐私法规的完善，联邦学习作为一种分布式机器学习范式...',
    status: 'generated',
    category: '网络与信息安全',
    paperCount: 8,
    score: null,
    createdAt: '2024-01-08',
  },
  {
    id: '3',
    title: '图神经网络综述',
    abstract: '图结构数据在社交网络、知识图谱等领域广泛存在...',
    status: 'draft',
    category: '人工智能',
    paperCount: 3,
    score: null,
    createdAt: '2024-01-05',
  },
]

const statusConfig = {
  draft: { label: '草稿', color: 'text-dark-400 bg-dark-700' },
  generated: { label: '已生成', color: 'text-blue-400 bg-blue-400/10' },
  scored: { label: '已评分', color: 'text-green-400 bg-green-400/10' },
}

export default function ReviewsPage() {
  const [reviews] = useState(mockReviews)
  const [showGenerator, setShowGenerator] = useState(false)
  const [expandedReview, setExpandedReview] = useState<string | null>(null)

  return (
    <div className="space-y-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex flex-col md:flex-row md:items-center justify-between gap-4"
      >
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-3">
            <BookOpen className="w-7 h-7 text-primary-400" />
            论文综述
          </h1>
          <p className="text-dark-400 mt-1">AI 辅助撰写论文综述，智能评分与改进建议</p>
        </div>
        <button 
          onClick={() => setShowGenerator(!showGenerator)}
          className="btn-primary flex items-center gap-2 self-start"
        >
          <Plus className="w-5 h-5" />
          新建综述
        </button>
      </motion.div>

      {/* Generator Panel */}
      {showGenerator && (
        <motion.div
          initial={{ opacity: 0, height: 0 }}
          animate={{ opacity: 1, height: 'auto' }}
          exit={{ opacity: 0, height: 0 }}
          className="glass-card p-6"
        >
          <h3 className="font-semibold mb-4 flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-accent-400" />
            生成论文综述
          </h3>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-2">综述标题</label>
              <input 
                type="text" 
                className="input-field"
                placeholder="例如：深度学习在医学影像分析中的应用综述"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-dark-300 mb-2">选择论文</label>
              <div className="p-4 rounded-xl bg-dark-800/50 border border-dark-700 border-dashed text-center">
                <FileText className="w-8 h-8 text-dark-500 mx-auto mb-2" />
                <p className="text-dark-400 text-sm">从我的论文中选择作为参考</p>
                <button className="btn-secondary mt-3 text-sm">
                  选择论文
                </button>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-dark-300 mb-2">综述风格</label>
              <div className="grid grid-cols-3 gap-3">
                {[
                  { id: 'academic', name: '学术风格', desc: '正式、严谨' },
                  { id: 'summary', name: '总结风格', desc: '简洁、精炼' },
                  { id: 'comprehensive', name: '全面风格', desc: '详尽、深入' },
                ].map((style) => (
                  <button
                    key={style.id}
                    className="p-3 rounded-lg bg-dark-800/50 border border-transparent hover:border-primary-500/30 transition-all text-left"
                  >
                    <p className="font-medium text-sm">{style.name}</p>
                    <p className="text-xs text-dark-500">{style.desc}</p>
                  </button>
                ))}
              </div>
            </div>

            <button className="btn-primary w-full flex items-center justify-center gap-2">
              <Sparkles className="w-5 h-5" />
              生成综述
            </button>
          </div>
        </motion.div>
      )}

      {/* Reviews List */}
      <div className="space-y-4">
        {reviews.map((review, index) => {
          const status = statusConfig[review.status as keyof typeof statusConfig]
          const isExpanded = expandedReview === review.id

          return (
            <motion.div
              key={review.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              className="glass-card overflow-hidden"
            >
              <div 
                className="p-6 cursor-pointer"
                onClick={() => setExpandedReview(isExpanded ? null : review.id)}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-2">
                      <span className={`px-2 py-1 rounded text-xs font-medium ${status.color}`}>
                        {status.label}
                      </span>
                      <span className="category-badge">{review.category}</span>
                      <span className="text-dark-500 text-sm flex items-center gap-1">
                        <FileText className="w-4 h-4" />
                        {review.paperCount} 篇论文
                      </span>
                    </div>

                    <h3 className="text-lg font-semibold mb-2">{review.title}</h3>
                    <p className="text-dark-400 text-sm line-clamp-2">{review.abstract}</p>

                    <div className="flex items-center gap-4 mt-4">
                      <span className="text-dark-500 text-sm flex items-center gap-1">
                        <Clock className="w-4 h-4" />
                        {review.createdAt}
                      </span>
                      {review.score && (
                        <span className="text-accent-400 text-sm flex items-center gap-1">
                          <Star className="w-4 h-4" />
                          评分: {review.score}/100
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <ChevronDown 
                      className={`w-5 h-5 text-dark-500 transition-transform ${isExpanded ? 'rotate-180' : ''}`} 
                    />
                  </div>
                </div>
              </div>

              {/* Expanded Content */}
              {isExpanded && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  className="px-6 pb-6 border-t border-dark-700/50"
                >
                  <div className="pt-4 space-y-4">
                    {review.status === 'scored' && (
                      <div className="p-4 rounded-xl bg-dark-800/50">
                        <h4 className="font-medium mb-3 flex items-center gap-2">
                          <Star className="w-5 h-5 text-accent-400" />
                          AI 评分详情
                        </h4>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                          {[
                            { name: '全面性', score: 82 },
                            { name: '组织结构', score: 88 },
                            { name: '批判性分析', score: 85 },
                            { name: '写作质量', score: 86 },
                          ].map((criterion) => (
                            <div key={criterion.name} className="text-center">
                              <div className="text-2xl font-bold text-primary-400">{criterion.score}</div>
                              <div className="text-sm text-dark-400">{criterion.name}</div>
                            </div>
                          ))}
                        </div>
                        <div className="space-y-2">
                          <div className="p-3 rounded-lg bg-green-500/10 border border-green-500/20">
                            <p className="text-sm text-green-400">
                              <span className="font-medium">优点：</span>论文选取全面，结构清晰，分析深入
                            </p>
                          </div>
                          <div className="p-3 rounded-lg bg-amber-500/10 border border-amber-500/20">
                            <p className="text-sm text-amber-400">
                              <span className="font-medium">建议：</span>可增加对未来研究方向的探讨
                            </p>
                          </div>
                        </div>
                      </div>
                    )}

                    <div className="flex items-center gap-3">
                      <button className="btn-secondary flex items-center gap-2 flex-1">
                        <Edit3 className="w-4 h-4" />
                        编辑综述
                      </button>
                      {review.status !== 'scored' && (
                        <button className="btn-primary flex items-center gap-2 flex-1">
                          <CheckCircle className="w-4 h-4" />
                          获取评分
                        </button>
                      )}
                      <button className="p-3 rounded-xl bg-dark-800 text-red-400 hover:bg-red-500/10 transition-colors">
                        <Trash2 className="w-5 h-5" />
                      </button>
                    </div>
                  </div>
                </motion.div>
              )}
            </motion.div>
          )
        })}
      </div>

      {reviews.length === 0 && !showGenerator && (
        <div className="text-center py-16">
          <BookOpen className="w-16 h-16 text-dark-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">还没有论文综述</h3>
          <p className="text-dark-400 mb-6">选择您已上传的论文，让 AI 帮您撰写高质量综述</p>
          <button 
            onClick={() => setShowGenerator(true)}
            className="btn-primary inline-flex items-center gap-2"
          >
            <Sparkles className="w-5 h-5" />
            新建综述
          </button>
        </div>
      )}
    </div>
  )
}

