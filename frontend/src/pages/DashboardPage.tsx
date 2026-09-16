import { motion } from 'framer-motion'
import {
  FileText,
  TrendingUp,
  GraduationCap,
  BookOpen,
  ArrowUpRight,
  Sparkles,
  Clock,
  Star,
} from 'lucide-react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/authStore'
import { learningApi, paperApi, reviewApi, trendingApi } from '@/services/api'

const quickActions = [
  { label: '上传论文', icon: FileText, link: '/app/papers', desc: '上传 PDF 进行智能分析' },
  { label: '查看热点', icon: TrendingUp, link: '/app/trending', desc: '探索研究领域最新趋势' },
  { label: '生成学习路线', icon: GraduationCap, link: '/app/learning', desc: '获取个性化学习建议' },
  { label: '撰写综述', icon: BookOpen, link: '/app/reviews', desc: 'AI 辅助生成论文综述' },
]

export default function DashboardPage() {
  const { user } = useAuthStore()

  const { data: paperCount = 0 } = useQuery({
    queryKey: ['papers-count'],
    queryFn: async () => {
      const res = await paperApi.list({ page: 1, per_page: 1 })
      return res.data.meta?.total ?? (res.data.data?.length || 0)
    },
  })
  const { data: trendCount = 0 } = useQuery({
    queryKey: ['trending-count'],
    queryFn: async () => {
      const res = await trendingApi.getItems({ limit: 50 })
      return (res.data.data || []).length
    },
  })
  const { data: pathCount = 0 } = useQuery({
    queryKey: ['learning-count'],
    queryFn: async () => {
      const res = await learningApi.getPaths()
      return (res.data.data || []).length
    },
  })
  const { data: reviewCount = 0 } = useQuery({
    queryKey: ['reviews-count'],
    queryFn: async () => {
      const res = await reviewApi.list({ page: 1, per_page: 1 })
      return res.data.meta?.total ?? (res.data.data?.length || 0)
    },
  })

  const stats = [
    { label: '我的论文', value: String(paperCount), icon: FileText, color: 'from-blue-500 to-cyan-500', link: '/app/papers' },
    { label: '热点追踪', value: String(trendCount), icon: TrendingUp, color: 'from-purple-500 to-pink-500', link: '/app/trending' },
    { label: '学习路线', value: String(pathCount), icon: GraduationCap, color: 'from-amber-500 to-orange-500', link: '/app/learning' },
    { label: '论文综述', value: String(reviewCount), icon: BookOpen, color: 'from-emerald-500 to-teal-500', link: '/app/reviews' },
  ]

  return (
    <div className="space-y-8">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="glass-card p-8">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold mb-2">欢迎回来，{user?.name || '研究者'}</h1>
            <p className="text-dark-400">今天想要探索哪个研究领域？让 AI 助手帮助您快速入门。</p>
          </div>
          <div className="hidden md:flex items-center gap-2 px-4 py-2 rounded-xl bg-accent-500/10 border border-accent-500/20">
            <Sparkles className="w-5 h-5 text-accent-400" />
            <span className="text-accent-400 font-medium">演示模式已就绪</span>
          </div>
        </div>
      </motion.div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Link key={stat.label} to={stat.link} className="stat-card group hover:border-primary-500/30 transition-all duration-300">
            <div className="flex items-center justify-between mb-4">
              <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${stat.color} flex items-center justify-center`}>
                <stat.icon className="w-6 h-6 text-white" />
              </div>
              <ArrowUpRight className="w-5 h-5 text-dark-500 group-hover:text-primary-400 transition-colors" />
            </div>
            <p className="text-3xl font-bold mb-1">{stat.value}</p>
            <p className="text-dark-400 text-sm">{stat.label}</p>
          </Link>
        ))}
      </div>

      <div>
        <h2 className="text-lg font-semibold mb-4">快速操作</h2>
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
          {quickActions.map((action) => (
            <Link key={action.label} to={action.link} className="glass-card p-5 hover:border-primary-500/30 transition-all duration-300 group">
              <div className="w-10 h-10 rounded-lg bg-primary-500/10 flex items-center justify-center mb-3 group-hover:bg-primary-500/20">
                <action.icon className="w-5 h-5 text-primary-400" />
              </div>
              <h3 className="font-medium mb-1">{action.label}</h3>
              <p className="text-dark-400 text-sm">{action.desc}</p>
            </Link>
          ))}
        </div>
      </div>

      <div className="glass-card p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">演示建议</h2>
          <span className="text-dark-400 text-sm flex items-center gap-1"><Clock className="w-4 h-4" />约 5 分钟</span>
        </div>
        <ol className="list-decimal list-inside space-y-2 text-dark-300 text-sm">
          <li>在热点页查看 GitHub / CCF 条目并阅读周报摘要</li>
          <li>上传一份 PDF，点击分析查看摘要、方法、贡献</li>
          <li>生成一条学习路径</li>
          <li>用已分析论文生成综述并评分</li>
        </ol>
      </div>

      <div className="glass-card p-6 border-accent-500/20 bg-gradient-to-br from-accent-500/5 to-transparent">
        <div className="flex items-start gap-4">
          <div className="w-10 h-10 rounded-lg bg-accent-500/20 flex items-center justify-center flex-shrink-0">
            <Star className="w-5 h-5 text-accent-400" />
          </div>
          <div>
            <h3 className="font-semibold mb-1">研究小贴士</h3>
            <p className="text-dark-400 text-sm">
              未配置 LLM_API_KEY 时，分析/综述会返回离线模板，保证演示流程可走完。
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
