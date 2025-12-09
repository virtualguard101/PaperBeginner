import { useState } from 'react'
import { motion } from 'framer-motion'
import { 
  GraduationCap, 
  Plus,
  ChevronRight,
  Book,
  Video,
  FileText,
  ExternalLink,
  Clock,
  CheckCircle,
  Sparkles
} from 'lucide-react'

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

// Mock learning path
const mockLearningPath = {
  id: '1',
  title: '人工智能研究入门路线',
  description: '从数学基础到深度学习前沿，系统掌握 AI 研究必备知识',
  category: '人工智能',
  difficulty: 'beginner',
  estimatedTime: '3-6 个月',
  stages: [
    {
      order: 1,
      title: '数学基础',
      description: '线性代数、概率论、优化理论',
      duration: '4-6 周',
      completed: true,
      resources: [
        { type: 'course', title: 'MIT 18.06 线性代数', provider: 'MIT OCW', url: '#', isFree: true },
        { type: 'book', title: '深度学习的数学', provider: '图灵图书', url: '#', isFree: false },
      ]
    },
    {
      order: 2,
      title: '机器学习基础',
      description: '监督学习、无监督学习、评估方法',
      duration: '6-8 周',
      completed: false,
      resources: [
        { type: 'course', title: 'Stanford CS229', provider: 'Stanford', url: '#', isFree: true },
        { type: 'course', title: '吴恩达机器学习', provider: 'Coursera', url: '#', isFree: true },
      ]
    },
    {
      order: 3,
      title: '深度学习',
      description: 'CNN、RNN、Transformer',
      duration: '8-12 周',
      completed: false,
      resources: [
        { type: 'course', title: 'Stanford CS231n', provider: 'Stanford', url: '#', isFree: true },
        { type: 'video', title: '李沐动手学深度学习', provider: 'B站', url: '#', isFree: true },
      ]
    },
  ]
}

const resourceIcons = {
  course: GraduationCap,
  book: Book,
  video: Video,
  documentation: FileText,
  tutorial: FileText,
}

export default function LearningPage() {
  const [showGenerator, setShowGenerator] = useState(false)
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null)
  const [selectedDifficulty, setSelectedDifficulty] = useState<string | null>(null)
  const [hasPath] = useState(true) // Mock: user has a learning path

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
            <GraduationCap className="w-7 h-7 text-primary-400" />
            学习路线
          </h1>
          <p className="text-dark-400 mt-1">AI 生成个性化学习路径，快速入门研究领域</p>
        </div>
        <button 
          onClick={() => setShowGenerator(!showGenerator)}
          className="btn-primary flex items-center gap-2 self-start"
        >
          <Plus className="w-5 h-5" />
          生成新路线
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
            生成学习路线
          </h3>
          
          <div className="space-y-6">
            {/* Category Selection */}
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

            {/* Difficulty Selection */}
            <div>
              <label className="block text-sm font-medium text-dark-300 mb-3">选择难度级别</label>
              <div className="grid grid-cols-3 gap-3">
                {difficulties.map((diff) => (
                  <button
                    key={diff.id}
                    onClick={() => setSelectedDifficulty(diff.id)}
                    className={`p-4 rounded-lg text-center transition-all ${
                      selectedDifficulty === diff.id
                        ? 'bg-primary-500/20 border border-primary-500/30'
                        : 'bg-dark-800/50 border border-transparent hover:bg-dark-800'
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
              disabled={!selectedCategory || !selectedDifficulty}
            >
              <Sparkles className="w-5 h-5" />
              生成学习路线
            </button>
          </div>
        </motion.div>
      )}

      {/* Learning Path */}
      {hasPath && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="glass-card p-6"
        >
          <div className="flex items-start justify-between mb-6">
            <div>
              <div className="flex items-center gap-3 mb-2">
                <span className="category-badge">{mockLearningPath.category}</span>
                <span className="px-2 py-1 rounded text-xs bg-dark-700 text-dark-300">
                  {mockLearningPath.difficulty === 'beginner' ? '入门级' : 
                   mockLearningPath.difficulty === 'intermediate' ? '进阶级' : '高级'}
                </span>
              </div>
              <h2 className="text-xl font-bold mb-2">{mockLearningPath.title}</h2>
              <p className="text-dark-400">{mockLearningPath.description}</p>
              <div className="flex items-center gap-2 mt-3 text-dark-500 text-sm">
                <Clock className="w-4 h-4" />
                预计时长: {mockLearningPath.estimatedTime}
              </div>
            </div>
          </div>

          {/* Progress */}
          <div className="mb-8">
            <div className="flex items-center justify-between text-sm mb-2">
              <span className="text-dark-400">学习进度</span>
              <span className="text-primary-400">1 / {mockLearningPath.stages.length} 阶段</span>
            </div>
            <div className="h-2 bg-dark-800 rounded-full overflow-hidden">
              <div 
                className="h-full bg-gradient-to-r from-primary-500 to-primary-400 rounded-full transition-all"
                style={{ width: `${(1 / mockLearningPath.stages.length) * 100}%` }}
              />
            </div>
          </div>

          {/* Stages */}
          <div className="space-y-4">
            {mockLearningPath.stages.map((stage, index) => (
              <div 
                key={stage.order}
                className={`p-5 rounded-xl border transition-all ${
                  stage.completed 
                    ? 'bg-green-500/5 border-green-500/20' 
                    : index === 1 
                    ? 'bg-primary-500/5 border-primary-500/20'
                    : 'bg-dark-800/50 border-dark-700/50'
                }`}
              >
                <div className="flex items-start gap-4">
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 ${
                    stage.completed 
                      ? 'bg-green-500 text-white' 
                      : index === 1
                      ? 'bg-primary-500 text-white'
                      : 'bg-dark-700 text-dark-400'
                  }`}>
                    {stage.completed ? (
                      <CheckCircle className="w-5 h-5" />
                    ) : (
                      <span className="font-semibold">{stage.order}</span>
                    )}
                  </div>
                  
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h3 className="font-semibold">{stage.title}</h3>
                      <span className="text-dark-500 text-sm">{stage.duration}</span>
                    </div>
                    <p className="text-dark-400 text-sm mb-4">{stage.description}</p>
                    
                    {/* Resources */}
                    <div className="space-y-2">
                      {stage.resources.map((resource, rIndex) => {
                        const ResourceIcon = resourceIcons[resource.type as keyof typeof resourceIcons] || FileText
                        return (
                          <a
                            key={rIndex}
                            href={resource.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-3 p-3 rounded-lg bg-dark-900/50 hover:bg-dark-900 transition-colors group"
                          >
                            <ResourceIcon className="w-5 h-5 text-dark-500" />
                            <div className="flex-1 min-w-0">
                              <p className="font-medium text-sm group-hover:text-primary-400 transition-colors">
                                {resource.title}
                              </p>
                              <p className="text-dark-500 text-xs">{resource.provider}</p>
                            </div>
                            {resource.isFree && (
                              <span className="px-2 py-0.5 rounded text-xs bg-green-500/10 text-green-400">
                                免费
                              </span>
                            )}
                            <ExternalLink className="w-4 h-4 text-dark-500 group-hover:text-primary-400" />
                          </a>
                        )
                      })}
                    </div>
                  </div>

                  <ChevronRight className="w-5 h-5 text-dark-600 flex-shrink-0" />
                </div>
              </div>
            ))}
          </div>
        </motion.div>
      )}

      {!hasPath && !showGenerator && (
        <div className="text-center py-16">
          <GraduationCap className="w-16 h-16 text-dark-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">还没有学习路线</h3>
          <p className="text-dark-400 mb-6">让 AI 根据您的研究方向生成个性化学习路径</p>
          <button 
            onClick={() => setShowGenerator(true)}
            className="btn-primary inline-flex items-center gap-2"
          >
            <Sparkles className="w-5 h-5" />
            生成学习路线
          </button>
        </div>
      )}
    </div>
  )
}

