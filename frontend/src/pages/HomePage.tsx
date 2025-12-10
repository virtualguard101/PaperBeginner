import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { 
  TrendingUp, 
  FileText, 
  GraduationCap, 
  BookOpen, 
  Sparkles,
  ArrowRight,
  Github,
  Brain,
  Target
} from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'

const features = [
  {
    icon: TrendingUp,
    title: '行业前沿热点嗅探',
    description: '实时追踪 GitHub 热门项目和 CCF A 类会议论文，AI 智能分类研究领域',
    color: 'from-blue-500 to-cyan-500',
  },
  {
    icon: GraduationCap,
    title: '个性化学习路线',
    description: '根据您的研究方向，AI 生成定制化学习路径，整合顶尖院校公开课资源',
    color: 'from-purple-500 to-pink-500',
  },
  {
    icon: FileText,
    title: '论文智能分析',
    description: '上传论文，AI 自动提取摘要、方法、贡献点，快速掌握论文核心内容',
    color: 'from-amber-500 to-orange-500',
  },
  {
    icon: BookOpen,
    title: '综述撰写与评分',
    description: 'AI 辅助撰写高质量论文综述，并提供专业评分和改进建议',
    color: 'from-emerald-500 to-teal-500',
  },
]

const ccfCategories = [
  '计算机体系结构',
  '计算机网络',
  '网络与信息安全',
  '软件工程',
  '数据库/数据挖掘',
  '人工智能',
  '计算机图形学',
  '人机交互',
]

export default function HomePage() {
  const { isAuthenticated } = useAuthStore()

  return (
    <div className="min-h-screen">
      {/* Navigation */}
      <nav className="fixed top-0 left-0 right-0 z-50 bg-dark-950/80 backdrop-blur-xl border-b border-dark-700/50">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center">
              <BookOpen className="w-5 h-5 text-white" />
            </div>
            <span className="font-semibold text-lg gradient-text">PaperBeginner</span>
          </div>
          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <Link to="/app" className="btn-primary">
                进入应用
              </Link>
            ) : (
              <>
                <Link to="/login" className="nav-link">登录</Link>
                <Link to="/register" className="btn-primary">免费注册</Link>
              </>
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-6">
        <div className="max-w-7xl mx-auto">
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="text-center max-w-4xl mx-auto"
          >
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary-500/10 border border-primary-500/20 text-primary-400 text-sm font-medium mb-6">
              <Sparkles className="w-4 h-4" />
              <span>AI 驱动的学术研究助手</span>
            </div>
            
            <h1 className="text-5xl md:text-6xl font-bold mb-6 leading-tight">
              开启您的
              <span className="gradient-text"> 计算机学术研究 </span>
              之旅
            </h1>
            
            <p className="text-xl text-dark-300 mb-10 max-w-2xl mx-auto leading-relaxed">
              PaperBeginner 是您的 AI 学术导师，帮助学术新人快速掌握研究领域前沿动态，
              构建专业知识体系，提升论文阅读与写作能力。
            </p>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link to={isAuthenticated ? '/app' : '/register'} className="btn-primary flex items-center gap-2">
                开始使用
                <ArrowRight className="w-5 h-5" />
              </Link>
              <a 
                href="https://github.com/virtualguard/paperbeginner" 
                target="_blank" 
                rel="noopener noreferrer"
                className="btn-secondary flex items-center gap-2"
              >
                <Github className="w-5 h-5" />
                GitHub
              </a>
            </div>
          </motion.div>

          {/* Animated gradient orbs */}
          <div className="relative mt-20">
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="w-96 h-96 bg-primary-500/20 rounded-full blur-3xl animate-pulse-slow" />
              <div className="w-72 h-72 bg-accent-500/20 rounded-full blur-3xl animate-pulse-slow animation-delay-500 -ml-32" />
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 px-6">
        <div className="max-w-7xl mx-auto">
          <motion.div 
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl md:text-4xl font-bold mb-4">核心功能</h2>
            <p className="text-dark-400 text-lg">全方位助力您的学术研究之路</p>
          </motion.div>

          <div className="grid md:grid-cols-2 gap-6">
            {features.map((feature, index) => (
              <motion.div
                key={feature.title}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="glass-card p-8 hover:border-primary-500/30 transition-all duration-300 group"
              >
                <div className={`w-14 h-14 rounded-2xl bg-gradient-to-br ${feature.color} flex items-center justify-center mb-6 group-hover:scale-110 transition-transform`}>
                  <feature.icon className="w-7 h-7 text-white" />
                </div>
                <h3 className="text-xl font-semibold mb-3">{feature.title}</h3>
                <p className="text-dark-400 leading-relaxed">{feature.description}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CCF Categories Section */}
      <section className="py-20 px-6 bg-dark-900/30">
        <div className="max-w-7xl mx-auto">
          <motion.div 
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-center mb-12"
          >
            <h2 className="text-3xl md:text-4xl font-bold mb-4">CCF 推荐领域全覆盖</h2>
            <p className="text-dark-400 text-lg">支持中国计算机学会推荐的所有研究领域</p>
          </motion.div>

          <div className="flex flex-wrap justify-center gap-3">
            {ccfCategories.map((category, index) => (
              <motion.span
                key={category}
                initial={{ opacity: 0, scale: 0.9 }}
                whileInView={{ opacity: 1, scale: 1 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.05 }}
                className="category-badge text-base px-5 py-2"
              >
                {category}
              </motion.span>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-20 px-6">
        <div className="max-w-7xl mx-auto">
          <motion.div 
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl md:text-4xl font-bold mb-4">使用流程</h2>
            <p className="text-dark-400 text-lg">三步开启您的学术研究之旅</p>
          </motion.div>

          <div className="grid md:grid-cols-3 gap-8">
            {[
              { icon: Target, title: '选择研究领域', desc: '浏览热点趋势，确定您感兴趣的研究方向' },
              { icon: Brain, title: 'AI 智能分析', desc: '上传论文，获取智能摘要和深度分析' },
              { icon: GraduationCap, title: '系统化学习', desc: '获取定制学习路线，逐步构建专业知识体系' },
            ].map((step, index) => (
              <motion.div
                key={step.title}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.15 }}
                className="text-center"
              >
                <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center mx-auto mb-6 shadow-lg shadow-primary-500/20">
                  <step.icon className="w-8 h-8 text-white" />
                </div>
                <div className="text-4xl font-bold text-primary-500 mb-2">0{index + 1}</div>
                <h3 className="text-xl font-semibold mb-2">{step.title}</h3>
                <p className="text-dark-400">{step.desc}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-6">
        <div className="max-w-4xl mx-auto">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="glass-card p-12 text-center relative overflow-hidden"
          >
            <div className="absolute inset-0 bg-gradient-to-br from-primary-500/10 to-accent-500/10" />
            <div className="relative">
              <h2 className="text-3xl md:text-4xl font-bold mb-4">准备好开始了吗？</h2>
              <p className="text-dark-300 text-lg mb-8 max-w-2xl mx-auto">
                立即注册，免费体验 AI 驱动的学术研究辅助工具
              </p>
              <Link to="/register" className="btn-primary inline-flex items-center gap-2">
                免费开始
                <ArrowRight className="w-5 h-5" />
              </Link>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 px-6 border-t border-dark-700/50">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center">
              <BookOpen className="w-5 h-5 text-white" />
            </div>
            <span className="font-semibold gradient-text">PaperBeginner</span>
          </div>
          <p className="text-dark-500 text-sm">
            © 2024 PaperBeginner. 让学术研究更简单。
          </p>
        </div>
      </footer>
    </div>
  )
}

