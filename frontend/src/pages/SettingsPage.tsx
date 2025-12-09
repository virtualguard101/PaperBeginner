import { useState } from 'react'
import { motion } from 'framer-motion'
import { 
  Settings, 
  User,
  Key,
  Bell,
  Palette,
  Save,
  Eye,
  EyeOff,
  Check
} from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import toast from 'react-hot-toast'

const llmProviders = [
  { id: 'openai', name: 'OpenAI', models: ['gpt-4o', 'gpt-4-turbo', 'gpt-3.5-turbo'] },
  { id: 'anthropic', name: 'Anthropic Claude', models: ['claude-3-5-sonnet', 'claude-3-opus'] },
  { id: 'deepseek', name: 'DeepSeek', models: ['deepseek-chat', 'deepseek-coder'] },
  { id: 'ollama', name: 'Ollama (本地)', models: ['llama3.1', 'mistral', 'qwen2'] },
]

export default function SettingsPage() {
  const { user } = useAuthStore()
  const [activeTab, setActiveTab] = useState('profile')
  
  // Profile settings
  const [name, setName] = useState(user?.name || '')
  const [email] = useState(user?.email || '')
  
  // API Keys
  const [showKeys, setShowKeys] = useState<Record<string, boolean>>({})
  const [apiKeys, setApiKeys] = useState({
    openai: '',
    anthropic: '',
    deepseek: '',
  })
  
  // Preferences
  const [preferredProvider, setPreferredProvider] = useState('openai')
  const [notifyTrending, setNotifyTrending] = useState(true)
  const [language, setLanguage] = useState('zh')

  const handleSave = () => {
    toast.success('设置已保存')
  }

  const tabs = [
    { id: 'profile', label: '个人资料', icon: User },
    { id: 'api', label: 'API 密钥', icon: Key },
    { id: 'preferences', label: '偏好设置', icon: Palette },
    { id: 'notifications', label: '通知设置', icon: Bell },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-2xl font-bold flex items-center gap-3">
          <Settings className="w-7 h-7 text-primary-400" />
          设置
        </h1>
        <p className="text-dark-400 mt-1">管理您的账户和应用偏好</p>
      </motion.div>

      <div className="flex flex-col lg:flex-row gap-6">
        {/* Tabs */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="lg:w-64 flex-shrink-0"
        >
          <div className="glass-card p-2">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                  activeTab === tab.id
                    ? 'bg-primary-500/10 text-primary-400'
                    : 'text-dark-300 hover:bg-dark-800/50'
                }`}
              >
                <tab.icon className="w-5 h-5" />
                {tab.label}
              </button>
            ))}
          </div>
        </motion.div>

        {/* Content */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex-1"
        >
          <div className="glass-card p-6">
            {/* Profile Tab */}
            {activeTab === 'profile' && (
              <div className="space-y-6">
                <h2 className="text-lg font-semibold">个人资料</h2>
                
                <div className="flex items-center gap-6">
                  <div className="w-20 h-20 rounded-full bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center text-3xl font-bold text-white">
                    {name.charAt(0).toUpperCase()}
                  </div>
                  <button className="btn-secondary">更换头像</button>
                </div>

                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-dark-300 mb-2">姓名</label>
                    <input
                      type="text"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="input-field"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-dark-300 mb-2">邮箱</label>
                    <input
                      type="email"
                      value={email}
                      disabled
                      className="input-field bg-dark-800 cursor-not-allowed"
                    />
                  </div>
                </div>

                <button onClick={handleSave} className="btn-primary flex items-center gap-2">
                  <Save className="w-4 h-4" />
                  保存更改
                </button>
              </div>
            )}

            {/* API Keys Tab */}
            {activeTab === 'api' && (
              <div className="space-y-6">
                <div>
                  <h2 className="text-lg font-semibold mb-2">API 密钥</h2>
                  <p className="text-dark-400 text-sm">
                    配置您自己的 API 密钥以使用各 LLM 服务。密钥将安全存储。
                  </p>
                </div>

                <div className="space-y-4">
                  {llmProviders.filter(p => p.id !== 'ollama').map((provider) => (
                    <div key={provider.id} className="p-4 rounded-xl bg-dark-800/50">
                      <div className="flex items-center justify-between mb-3">
                        <div>
                          <h3 className="font-medium">{provider.name}</h3>
                          <p className="text-dark-500 text-sm">
                            支持模型: {provider.models.join(', ')}
                          </p>
                        </div>
                      </div>
                      <div className="relative">
                        <input
                          type={showKeys[provider.id] ? 'text' : 'password'}
                          value={apiKeys[provider.id as keyof typeof apiKeys]}
                          onChange={(e) => setApiKeys({
                            ...apiKeys,
                            [provider.id]: e.target.value
                          })}
                          placeholder={`输入 ${provider.name} API Key`}
                          className="input-field pr-12"
                        />
                        <button
                          onClick={() => setShowKeys({
                            ...showKeys,
                            [provider.id]: !showKeys[provider.id]
                          })}
                          className="absolute right-4 top-1/2 -translate-y-1/2 text-dark-400 hover:text-dark-200"
                        >
                          {showKeys[provider.id] ? (
                            <EyeOff className="w-5 h-5" />
                          ) : (
                            <Eye className="w-5 h-5" />
                          )}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>

                <button onClick={handleSave} className="btn-primary flex items-center gap-2">
                  <Save className="w-4 h-4" />
                  保存密钥
                </button>
              </div>
            )}

            {/* Preferences Tab */}
            {activeTab === 'preferences' && (
              <div className="space-y-6">
                <h2 className="text-lg font-semibold">偏好设置</h2>

                <div>
                  <label className="block text-sm font-medium text-dark-300 mb-3">
                    默认 LLM 提供商
                  </label>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    {llmProviders.map((provider) => (
                      <button
                        key={provider.id}
                        onClick={() => setPreferredProvider(provider.id)}
                        className={`p-4 rounded-xl text-center transition-all ${
                          preferredProvider === provider.id
                            ? 'bg-primary-500/20 border border-primary-500/30'
                            : 'bg-dark-800/50 border border-transparent hover:bg-dark-800'
                        }`}
                      >
                        <p className="font-medium text-sm">{provider.name}</p>
                        {preferredProvider === provider.id && (
                          <Check className="w-4 h-4 text-primary-400 mx-auto mt-2" />
                        )}
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-dark-300 mb-3">
                    界面语言
                  </label>
                  <div className="flex gap-3">
                    {[
                      { id: 'zh', name: '中文' },
                      { id: 'en', name: 'English' },
                    ].map((lang) => (
                      <button
                        key={lang.id}
                        onClick={() => setLanguage(lang.id)}
                        className={`px-6 py-3 rounded-xl transition-all ${
                          language === lang.id
                            ? 'bg-primary-500/20 border border-primary-500/30 text-primary-400'
                            : 'bg-dark-800/50 border border-transparent hover:bg-dark-800'
                        }`}
                      >
                        {lang.name}
                      </button>
                    ))}
                  </div>
                </div>

                <button onClick={handleSave} className="btn-primary flex items-center gap-2">
                  <Save className="w-4 h-4" />
                  保存设置
                </button>
              </div>
            )}

            {/* Notifications Tab */}
            {activeTab === 'notifications' && (
              <div className="space-y-6">
                <h2 className="text-lg font-semibold">通知设置</h2>

                <div className="space-y-4">
                  <div className="flex items-center justify-between p-4 rounded-xl bg-dark-800/50">
                    <div>
                      <p className="font-medium">热点趋势更新</p>
                      <p className="text-dark-400 text-sm">当您关注的领域有新热点时通知您</p>
                    </div>
                    <button
                      onClick={() => setNotifyTrending(!notifyTrending)}
                      className={`w-12 h-6 rounded-full transition-all ${
                        notifyTrending ? 'bg-primary-500' : 'bg-dark-600'
                      }`}
                    >
                      <div className={`w-5 h-5 rounded-full bg-white transition-transform ${
                        notifyTrending ? 'translate-x-6' : 'translate-x-0.5'
                      }`} />
                    </button>
                  </div>

                  <div className="flex items-center justify-between p-4 rounded-xl bg-dark-800/50">
                    <div>
                      <p className="font-medium">论文分析完成</p>
                      <p className="text-dark-400 text-sm">当您上传的论文分析完成时通知您</p>
                    </div>
                    <button className="w-12 h-6 rounded-full bg-primary-500">
                      <div className="w-5 h-5 rounded-full bg-white translate-x-6" />
                    </button>
                  </div>

                  <div className="flex items-center justify-between p-4 rounded-xl bg-dark-800/50">
                    <div>
                      <p className="font-medium">周报摘要</p>
                      <p className="text-dark-400 text-sm">每周发送一份研究领域热点摘要</p>
                    </div>
                    <button className="w-12 h-6 rounded-full bg-dark-600">
                      <div className="w-5 h-5 rounded-full bg-white translate-x-0.5" />
                    </button>
                  </div>
                </div>

                <button onClick={handleSave} className="btn-primary flex items-center gap-2">
                  <Save className="w-4 h-4" />
                  保存设置
                </button>
              </div>
            )}
          </div>
        </motion.div>
      </div>
    </div>
  )
}

