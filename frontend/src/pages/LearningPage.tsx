import { useState } from 'react'
import { motion } from 'framer-motion'
import {
  GraduationCap,
  Plus,
  Book,
  Video,
  FileText,
  ExternalLink,
  Clock,
  Sparkles,
  Loader2,
} from 'lucide-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { learningApi } from '@/services/api'
import type { LearningPath } from '@/types'

const categories = [
  { id: 1, name: '计算机体系结构' },
  { id: 2, name: '计算机网络' },
  { id: 3, name: '网络与信息安全' },
  { id: 4, name: '软件工程' },
  { id: 5, name: '数据库/数据挖掘' },
  { id: 8, name: '人工智能' },
]

const difficulties = [
  { id: 'beginner', name: '入门级', desc: '适合零基础' },
  { id: 'intermediate', name: '进阶级', desc: '有一定基础' },
  { id: 'advanced', name: '高级', desc: '深入研究' },
]

const resourceIcons = {
  course: GraduationCap,
  book: Book,
  video: Video,
  documentation: FileText,
  tutorial: FileText,
}

export default function LearningPage() {
  const queryClient = useQueryClient()
  const [showGenerator, setShowGenerator] = useState(false)
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null)
  const [selectedDifficulty, setSelectedDifficulty] = useState<string | null>(null)
  const [activeId, setActiveId] = useState<string | null>(null)

  const { data: paths = [] } = useQuery({
    queryKey: ['learning-paths'],
    queryFn: async () => {
      const res = await learningApi.getPaths()
      return (res.data.data || []) as LearningPath[]
    },
  })

  const generateMutation = useMutation({
    mutationFn: () =>
      learningApi.generate({
        category_id: selectedCategory!,
        difficulty: selectedDifficulty!,
      }),
    onSuccess: (res) => {
      toast.success('学习路线已生成')
      queryClient.invalidateQueries({ queryKey: ['learning-paths'] })
      setActiveId(res.data.data.id)
      setShowGenerator(false)
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: { message?: string } } } }
      toast.error(err.response?.data?.error?.message || '生成失败')
    },
  })

  const current = paths.find((p) => p.id === activeId) || paths[0]

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-3">
            <GraduationCap className="w-7 h-7 text-primary-400" />
            学习路线
          </h1>
          <p className="text-dark-400 mt-1">AI 生成个性化学习路径，快速入门研究领域</p>
        </div>
        <button onClick={() => setShowGenerator(!showGenerator)} className="btn-primary flex items-center gap-2 self-start">
          <Plus className="w-5 h-5" />
          生成新路线
        </button>
      </motion.div>

      {showGenerator && (
        <div className="glass-card p-6">
          <h3 className="font-semibold mb-4 flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-accent-400" />
            生成学习路线
          </h3>
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-3">选择研究领域</label>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                {categories.map((cat) => (
                  <button
                    key={cat.id}
                    onClick={() => setSelectedCategory(cat.id)}
                    className={`p-3 rounded-lg text-left transition-all ${
                      selectedCategory === cat.id
                        ? 'bg-primary-500/20 border border-primary-500/30 text-primary-400'
                        : 'bg-dark-800/50 border border-transparent hover:bg-dark-800'
                    }`}
                  >
                    {cat.name}
                  </button>
                ))}
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-3">选择难度级别</label>
              <div className="grid grid-cols-3 gap-3">
                {difficulties.map((diff) => (
                  <button
                    key={diff.id}
                    onClick={() => setSelectedDifficulty(diff.id)}
                    className={`p-4 rounded-lg text-center transition-all ${
                      selectedDifficulty === diff.id ? 'bg-primary-500/20 border border-primary-500/30' : 'bg-dark-800/50 border border-transparent hover:bg-dark-800'
                    }`}
                  >
                    <p className="font-medium">{diff.name}</p>
                    <p className="text-sm text-dark-400">{diff.desc}</p>
                  </button>
                ))}
              </div>
            </div>
            <button
              className="btn-primary w-full flex items-center justify-center gap-2"
              disabled={!selectedCategory || !selectedDifficulty || generateMutation.isPending}
              onClick={() => generateMutation.mutate()}
            >
              {generateMutation.isPending ? <Loader2 className="w-5 h-5 animate-spin" /> : <Sparkles className="w-5 h-5" />}
              生成学习路线
            </button>
          </div>
        </div>
      )}

      {paths.length > 1 && (
        <div className="flex flex-wrap gap-2">
          {paths.map((p) => (
            <button
              key={p.id}
              onClick={() => setActiveId(p.id)}
              className={`px-3 py-1.5 rounded-lg text-sm ${(current && current.id === p.id) ? 'bg-primary-500 text-white' : 'bg-dark-800 text-dark-300'}`}
            >
              {p.title}
            </button>
          ))}
        </div>
      )}

      {current && (
        <div className="glass-card p-6">
          <div className="flex items-center gap-3 mb-2">
            {current.category && <span className="category-badge">{current.category.name}</span>}
            <span className="px-2 py-1 rounded text-xs bg-dark-700 text-dark-300">
              {current.difficulty === 'beginner' ? '入门级' : current.difficulty === 'intermediate' ? '进阶级' : '高级'}
            </span>
          </div>
          <h2 className="text-xl font-bold mb-2">{current.title}</h2>
          <p className="text-dark-400">{current.description}</p>
          <div className="flex items-center gap-2 mt-3 text-dark-500 text-sm">
            <Clock className="w-4 h-4" />
            预计时长: {current.estimated_time}
          </div>
          <div className="space-y-4 mt-6">
            {(current.stages || []).map((stage) => (
              <div key={stage.order} className="p-5 rounded-xl border bg-dark-800/50 border-dark-700/50">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold">{stage.order}. {stage.title}</h3>
                  <span className="text-dark-500 text-sm">{stage.duration}</span>
                </div>
                <p className="text-dark-400 text-sm mb-4">{stage.description}</p>
                <div className="space-y-2">
                  {(stage.resources || []).map((resource, rIndex) => {
                    const ResourceIcon = resourceIcons[resource.type as keyof typeof resourceIcons] || FileText
                    return (
                      <a key={rIndex} href={resource.url} target="_blank" rel="noopener noreferrer" className="flex items-center gap-3 p-3 rounded-lg bg-dark-900/50 hover:bg-dark-900 group">
                        <ResourceIcon className="w-5 h-5 text-dark-500" />
                        <div className="flex-1 min-w-0">
                          <p className="font-medium text-sm group-hover:text-primary-400">{resource.title}</p>
                          <p className="text-dark-500 text-xs">{resource.provider}</p>
                        </div>
                        {resource.is_free && <span className="px-2 py-0.5 rounded text-xs bg-green-500/10 text-green-400">免费</span>}
                        <ExternalLink className="w-4 h-4 text-dark-500" />
                      </a>
                    )
                  })}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {!current && !showGenerator && (
        <div className="text-center py-16">
          <GraduationCap className="w-16 h-16 text-dark-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">还没有学习路线</h3>
          <p className="text-dark-400 mb-6">让 AI 根据您的研究方向生成个性化学习路径</p>
          <button onClick={() => setShowGenerator(true)} className="btn-primary inline-flex items-center gap-2">
            <Sparkles className="w-5 h-5" />
            生成学习路线
          </button>
        </div>
      )}
    </div>
  )
}
