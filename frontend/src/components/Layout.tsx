import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { 
  LayoutDashboard, 
  TrendingUp, 
  FileText, 
  GraduationCap, 
  BookOpen,
  Settings,
  LogOut,
  Menu,
  X
} from 'lucide-react'
import { useState } from 'react'
import { useAuthStore } from '@/stores/authStore'

const navItems = [
  { path: '/app', icon: LayoutDashboard, label: '概览', end: true },
  { path: '/app/trending', icon: TrendingUp, label: '热点' },
  { path: '/app/papers', icon: FileText, label: '论文' },
  { path: '/app/learning', icon: GraduationCap, label: '学习路线' },
  { path: '/app/reviews', icon: BookOpen, label: '综述' },
  { path: '/app/settings', icon: Settings, label: '设置' },
]

export default function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen flex">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div 
          className="fixed inset-0 bg-dark-950/80 backdrop-blur-sm z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <motion.aside
        initial={{ x: -280 }}
        animate={{ x: sidebarOpen ? 0 : -280 }}
        className="fixed lg:static inset-y-0 left-0 z-50 w-64 bg-dark-900/80 backdrop-blur-xl border-r border-dark-700/50 flex flex-col lg:translate-x-0"
      >
        {/* Logo */}
        <div className="h-16 flex items-center justify-between px-6 border-b border-dark-700/50">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center">
              <BookOpen className="w-5 h-5 text-white" />
            </div>
            <span className="font-semibold text-lg gradient-text">PaperBeginner</span>
          </div>
          <button 
            onClick={() => setSidebarOpen(false)}
            className="lg:hidden text-dark-400 hover:text-dark-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-1">
          {navItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              end={item.end}
              onClick={() => setSidebarOpen(false)}
              className={({ isActive }) => `
                flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200
                ${isActive 
                  ? 'bg-primary-500/10 text-primary-400 border border-primary-500/20' 
                  : 'text-dark-300 hover:text-dark-100 hover:bg-dark-800/50'
                }
              `}
            >
              <item.icon className="w-5 h-5" />
              <span className="font-medium">{item.label}</span>
            </NavLink>
          ))}
        </nav>

        {/* User section */}
        <div className="p-4 border-t border-dark-700/50">
          <div className="flex items-center gap-3 px-4 py-3 rounded-xl bg-dark-800/50">
            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center text-white font-semibold">
              {user?.name?.charAt(0).toUpperCase() || 'U'}
            </div>
            <div className="flex-1 min-w-0">
              <p className="font-medium text-dark-100 truncate">{user?.name}</p>
              <p className="text-sm text-dark-400 truncate">{user?.email}</p>
            </div>
            <button 
              onClick={handleLogout}
              className="p-2 text-dark-400 hover:text-red-400 transition-colors"
              title="Logout"
            >
              <LogOut className="w-5 h-5" />
            </button>
          </div>
        </div>
      </motion.aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Top bar */}
        <header className="h-16 flex items-center gap-4 px-6 border-b border-dark-700/50 bg-dark-900/50 backdrop-blur-xl">
          <button 
            onClick={() => setSidebarOpen(true)}
            className="lg:hidden text-dark-400 hover:text-dark-100"
          >
            <Menu className="w-6 h-6" />
          </button>
          <div className="flex-1" />
          {/* Add search or other header items here */}
        </header>

        {/* Page content */}
        <main className="flex-1 p-6 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

