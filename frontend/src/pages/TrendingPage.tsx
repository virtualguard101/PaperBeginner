import { useState } from 'react'
import { motion } from 'framer-motion'
import { 
  TrendingUp, 
  Github, 
  BookOpen, 
  Filter,
  ExternalLink,
  Star,
  GitFork,
  Sparkles,
  RefreshCw
} from 'lucide-react'

const categories = [
  { id: 0, name: '全部领域' },
  { id: 1, name: '计算机体系结构' },
  { id: 2, name: '计算机网络' },
  { id: 3, name: '网络与信息安全' },
  { id: 4, name: '软件工程' },
  { id: 5, name: '数据库/数据挖掘' },
  { id: 6, name: '计算机科学理论' },
  { id: 7, name: '计算机图形学' },
  { id: 8, name: '人工智能' },
  { id: 9, name: '人机交互' },
  { id: 10, name: '交叉/新兴' },
]

const sources = [
  { id: 'all', name: '全部来源', icon: Filter },
  { id: 'github', name: 'GitHub', icon: Github },
  { id: 'ccf', name: 'CCF 会议', icon: BookOpen },
]

// Mock data
const mockTrendingItems = [
  {
    id: '1',
    source: 'github',
    title: 'openai/whisper',
    description: 'Robust Speech Recognition via Large-Scale Weak Supervision',
    url: 'https://github.com/openai/whisper',
    stars: 52000,
    forks: 5800,
    language: 'Python',
    topics: ['speech-recognition', 'deep-learning', 'transformer'],
    category: '人工智能',
    trendScore: 98.5,
  },
  {
    id: '2',
    source: 'github',
    title: 'langchain-ai/langchain',
    description: 'Building applications with LLMs through composability',
    url: 'https://github.com/langchain-ai/langchain',
    stars: 78000,
    forks: 12000,
    language: 'Python',
    topics: ['llm', 'agents', 'rag'],
    category: '人工智能',
    trendScore: 97.2,
  },
  {
    id: '3',
    source: 'ccf',
    title: 'Attention Is All You Need',
    description: 'The Transformer architecture paper that revolutionized NLP and beyond',
    url: 'https://arxiv.org/abs/1706.03762',
    category: '人工智能',
    trendScore: 99.9,
  },
  {
    id: '4',
    source: 'github',
    title: 'AUTOMATIC1111/stable-diffusion-webui',
    description: 'Stable Diffusion web UI',
    url: 'https://github.com/AUTOMATIC1111/stable-diffusion-webui',
    stars: 125000,
    forks: 24000,
    language: 'Python',
    topics: ['stable-diffusion', 'ai-art', 'image-generation'],
    category: '计算机图形学',
    trendScore: 96.8,
  },
]

export default function TrendingPage() {
  const [selectedCategory, setSelectedCategory] = useState(0)
  const [selectedSource, setSelectedSource] = useState('all')
  const [loading] = useState(false)

  const filteredItems = mockTrendingItems.filter(item => {
    if (selectedSource !== 'all' && item.source !== selectedSource) return false
    return true
  })

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
            <TrendingUp className="w-7 h-7 text-primary-400" />
            行业前沿热点
          </h1>
          <p className="text-dark-400 mt-1">追踪 GitHub 热门项目和 CCF 顶会论文</p>
        </div>
        <button className="btn-secondary flex items-center gap-2 self-start">
          <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          刷新数据
        </button>
      </motion.div>

      {/* Filters */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="glass-card p-4"
      >
        {/* Source Filter */}
        <div className="flex flex-wrap gap-2 mb-4">
          {sources.map((source) => (
            <button
              key={source.id}
              onClick={() => setSelectedSource(source.id)}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-all ${
                selectedSource === source.id
                  ? 'bg-primary-500/20 text-primary-400 border border-primary-500/30'
                  : 'bg-dark-800/50 text-dark-300 hover:bg-dark-800 border border-transparent'
              }`}
            >
              <source.icon className="w-4 h-4" />
              {source.name}
            </button>
          ))}
        </div>

        {/* Category Filter */}
        <div className="flex flex-wrap gap-2">
          {categories.map((category) => (
            <button
              key={category.id}
              onClick={() => setSelectedCategory(category.id)}
              className={`px-3 py-1.5 rounded-lg text-sm transition-all ${
                selectedCategory === category.id
                  ? 'bg-primary-500 text-white'
                  : 'bg-dark-800/50 text-dark-300 hover:bg-dark-800'
              }`}
            >
              {category.name}
            </button>
          ))}
        </div>
      </motion.div>

      {/* Trending Items */}
      <div className="grid gap-4">
        {filteredItems.map((item, index) => (
          <motion.div
            key={item.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 + index * 0.05 }}
            className="glass-card p-6 hover:border-primary-500/30 transition-all"
          >
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-3 mb-2">
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    item.source === 'github' 
                      ? 'bg-dark-700 text-dark-200' 
                      : 'bg-primary-500/20 text-primary-400'
                  }`}>
                    {item.source === 'github' ? (
                      <span className="flex items-center gap-1">
                        <Github className="w-3 h-3" /> GitHub
                      </span>
                    ) : (
                      <span className="flex items-center gap-1">
                        <BookOpen className="w-3 h-3" /> CCF
                      </span>
                    )}
                  </span>
                  <span className="category-badge">{item.category}</span>
                  <span className="flex items-center gap-1 text-accent-400 text-sm">
                    <Sparkles className="w-4 h-4" />
                    {item.trendScore}
                  </span>
                </div>

                <a 
                  href={item.url} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-lg font-semibold hover:text-primary-400 transition-colors flex items-center gap-2"
                >
                  {item.title}
                  <ExternalLink className="w-4 h-4" />
                </a>
                
                <p className="text-dark-400 mt-2 line-clamp-2">{item.description}</p>

                {item.source === 'github' && (
                  <div className="flex items-center gap-4 mt-4">
                    <span className="flex items-center gap-1 text-dark-300 text-sm">
                      <Star className="w-4 h-4 text-yellow-500" />
                      {(item.stars! / 1000).toFixed(1)}k
                    </span>
                    <span className="flex items-center gap-1 text-dark-300 text-sm">
                      <GitFork className="w-4 h-4" />
                      {(item.forks! / 1000).toFixed(1)}k
                    </span>
                    {item.language && (
                      <span className="flex items-center gap-1 text-dark-300 text-sm">
                        <span className="w-3 h-3 rounded-full bg-blue-500" />
                        {item.language}
                      </span>
                    )}
                  </div>
                )}

                {item.topics && item.topics.length > 0 && (
                  <div className="flex flex-wrap gap-2 mt-3">
                    {item.topics.map((topic) => (
                      <span 
                        key={topic}
                        className="px-2 py-1 rounded text-xs bg-dark-800 text-dark-300"
                      >
                        {topic}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </motion.div>
        ))}
      </div>

      {filteredItems.length === 0 && (
        <div className="text-center py-12">
          <TrendingUp className="w-12 h-12 text-dark-600 mx-auto mb-4" />
          <p className="text-dark-400">暂无符合条件的热点数据</p>
        </div>
      )}
    </div>
  )
}

