import { motion } from 'framer-motion'
import { 
  FileText, 
  TrendingUp, 
  GraduationCap, 
  BookOpen,
  ArrowUpRight,
  Sparkles,
  Clock,
  Star
} from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

const stats = [
  { label: '我的论文', value: '12', icon: FileText, color: 'from-blue-500 to-cyan-500', link: '/app/papers' },
  { label: '热点追踪', value: '156', icon: TrendingUp, color: 'from-purple-500 to-pink-500', link: '/app/trending' },
  { label: '学习路线', value: '3', icon: GraduationCap, color: 'from-amber-500 to-orange-500', link: '/app/learning' },
  { label: '论文综述', value: '5', icon: BookOpen, color: 'from-emerald-500 to-teal-500', link: '/app/reviews' },
]

const recentActivities = [
  { type: 'paper', title: 'Attention Is All You Need', time: '2 小时前', status: '分析完成' },
  { type: 'learning', title: '深度学习入门路线', time: '1 天前', status: '进行中' },
  { type: 'trending', title: '人工智能领域热点报告', time: '2 天前', status: '已生成' },
  { type: 'review', title: 'Transformer 架构综述', time: '3 天前', status: '已评分' },
]

const quickActions = [
  { label: '上传论文', icon: FileText, link: '/app/papers', desc: '上传 PDF 进行智能分析' },
  { label: '查看热点', icon: TrendingUp, link: '/app/trending', desc: '探索研究领域最新趋势' },
  { label: '生成学习路线', icon: GraduationCap, link: '/app/learning', desc: '获取个性化学习建议' },
  { label: '撰写综述', icon: BookOpen, link: '/app/reviews', desc: 'AI 辅助生成论文综述' },
]

export default function DashboardPage() {
  const { user } = useAuthStore()

  return (
    <div className="space-y-8">
      {/* Welcome Header */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="glass-card p-8"
      >
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold mb-2">
              欢迎回来，{user?.name || '研究者'} 👋
            </h1>
            <p className="text-dark-400">
              今天想要探索哪个研究领域？让 AI 助手帮助您快速入门。
            </p>
          </div>
          <div className="hidden md:flex items-center gap-2 px-4 py-2 rounded-xl bg-accent-500/10 border border-accent-500/20">
            <Sparkles className="w-5 h-5 text-accent-400" />
            <span className="text-accent-400 font-medium">AI 助手在线</span>
          </div>
        </div>
      </motion.div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => (
          <motion.div
            key={stat.label}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <Link
              to={stat.link}
              className="stat-card group hover:border-primary-500/30 transition-all duration-300"
            >
              <div className="flex items-center justify-between mb-4">
                <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${stat.color} flex items-center justify-center`}>
                  <stat.icon className="w-6 h-6 text-white" />
                </div>
                <ArrowUpRight className="w-5 h-5 text-dark-500 group-hover:text-primary-400 transition-colors" />
              </div>
              <p className="text-3xl font-bold mb-1">{stat.value}</p>
              <p className="text-dark-400 text-sm">{stat.label}</p>
            </Link>
          </motion.div>
        ))}
      </div>

      {/* Quick Actions */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
      >
        <h2 className="text-lg font-semibold mb-4">快速操作</h2>
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
          {quickActions.map((action) => (
            <Link
              key={action.label}
              to={action.link}
              className="glass-card p-5 hover:border-primary-500/30 transition-all duration-300 group"
            >
              <div className="w-10 h-10 rounded-lg bg-primary-500/10 flex items-center justify-center mb-3 group-hover:bg-primary-500/20 transition-colors">
                <action.icon className="w-5 h-5 text-primary-400" />
              </div>
              <h3 className="font-medium mb-1">{action.label}</h3>
              <p className="text-dark-400 text-sm">{action.desc}</p>
            </Link>
          ))}
        </div>
      </motion.div>

      {/* Recent Activities */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5 }}
        className="glass-card p-6"
      >
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold">最近活动</h2>
          <span className="text-dark-400 text-sm">最近 7 天</span>
        </div>
        <div className="space-y-4">
          {recentActivities.map((activity, index) => (
            <div
              key={index}
              className="flex items-center gap-4 p-4 rounded-xl bg-dark-800/50 hover:bg-dark-800 transition-colors"
            >
              <div className="w-10 h-10 rounded-lg bg-primary-500/10 flex items-center justify-center">
                {activity.type === 'paper' && <FileText className="w-5 h-5 text-primary-400" />}
                {activity.type === 'learning' && <GraduationCap className="w-5 h-5 text-purple-400" />}
                {activity.type === 'trending' && <TrendingUp className="w-5 h-5 text-amber-400" />}
                {activity.type === 'review' && <BookOpen className="w-5 h-5 text-emerald-400" />}
              </div>
              <div className="flex-1 min-w-0">
                <p className="font-medium truncate">{activity.title}</p>
                <div className="flex items-center gap-2 text-sm text-dark-400">
                  <Clock className="w-4 h-4" />
                  <span>{activity.time}</span>
                </div>
              </div>
              <span className="px-3 py-1 rounded-full text-xs font-medium bg-primary-500/10 text-primary-400">
                {activity.status}
              </span>
            </div>
          ))}
        </div>
      </motion.div>

      {/* Tip Card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.6 }}
        className="glass-card p-6 border-accent-500/20 bg-gradient-to-br from-accent-500/5 to-transparent"
      >
        <div className="flex items-start gap-4">
          <div className="w-10 h-10 rounded-lg bg-accent-500/20 flex items-center justify-center flex-shrink-0">
            <Star className="w-5 h-5 text-accent-400" />
          </div>
          <div>
            <h3 className="font-semibold mb-1">研究小贴士</h3>
            <p className="text-dark-400 text-sm">
              建议从 CCF A 类会议的最新论文入手，了解领域前沿。您可以在"热点追踪"中查看各领域的热门研究方向，
              或者上传感兴趣的论文让 AI 帮您快速理解论文核心内容。
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  )
}

