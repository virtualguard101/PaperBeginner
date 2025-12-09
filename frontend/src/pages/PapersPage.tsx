import { useState, useCallback } from 'react'
import { motion } from 'framer-motion'
import { useDropzone } from 'react-dropzone'
import { 
  FileText, 
  Upload, 
  Search, 
  MoreVertical,
  Eye,
  Trash2,
  Brain,
  Clock,
  CheckCircle,
  AlertCircle,
  Loader2
} from 'lucide-react'

// Mock data
const mockPapers = [
  {
    id: '1',
    title: 'Attention Is All You Need',
    authors: ['Ashish Vaswani', 'Noam Shazeer', 'Niki Parmar'],
    abstract: 'The dominant sequence transduction models are based on complex recurrent or convolutional neural networks...',
    status: 'completed',
    category: '人工智能',
    createdAt: '2024-01-15',
  },
  {
    id: '2',
    title: 'BERT: Pre-training of Deep Bidirectional Transformers',
    authors: ['Jacob Devlin', 'Ming-Wei Chang', 'Kenton Lee'],
    abstract: 'We introduce a new language representation model called BERT...',
    status: 'processing',
    category: '人工智能',
    createdAt: '2024-01-14',
  },
  {
    id: '3',
    title: 'Deep Residual Learning for Image Recognition',
    authors: ['Kaiming He', 'Xiangyu Zhang', 'Shaoqing Ren'],
    abstract: 'Deeper neural networks are more difficult to train...',
    status: 'pending',
    category: '计算机图形学',
    createdAt: '2024-01-13',
  },
]

const statusConfig = {
  pending: { icon: Clock, label: '等待处理', color: 'text-yellow-400 bg-yellow-400/10' },
  processing: { icon: Loader2, label: '分析中', color: 'text-blue-400 bg-blue-400/10', animate: true },
  completed: { icon: CheckCircle, label: '已完成', color: 'text-green-400 bg-green-400/10' },
  failed: { icon: AlertCircle, label: '处理失败', color: 'text-red-400 bg-red-400/10' },
}

export default function PapersPage() {
  const [papers] = useState(mockPapers)
  const [searchQuery, setSearchQuery] = useState('')
  const [uploading, setUploading] = useState(false)

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      setUploading(true)
      // Simulate upload
      setTimeout(() => {
        setUploading(false)
      }, 2000)
    }
  }, [])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'application/pdf': ['.pdf']
    },
    maxFiles: 1,
  })

  const filteredPapers = papers.filter(paper =>
    paper.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    paper.authors.some(a => a.toLowerCase().includes(searchQuery.toLowerCase()))
  )

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
            <FileText className="w-7 h-7 text-primary-400" />
            我的论文
          </h1>
          <p className="text-dark-400 mt-1">上传论文进行 AI 智能分析</p>
        </div>
      </motion.div>

      {/* Upload Area */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        {...getRootProps()}
        className={`glass-card p-8 border-2 border-dashed cursor-pointer transition-all ${
          isDragActive 
            ? 'border-primary-500 bg-primary-500/5' 
            : 'border-dark-600 hover:border-primary-500/50'
        }`}
      >
        <input {...getInputProps()} />
        <div className="text-center">
          {uploading ? (
            <Loader2 className="w-12 h-12 text-primary-400 mx-auto mb-4 animate-spin" />
          ) : (
            <Upload className="w-12 h-12 text-dark-500 mx-auto mb-4" />
          )}
          <p className="text-lg font-medium mb-2">
            {isDragActive ? '释放以上传文件' : uploading ? '正在上传...' : '拖放 PDF 文件到这里'}
          </p>
          <p className="text-dark-400 text-sm">
            或者点击选择文件 • 支持 PDF 格式
          </p>
        </div>
      </motion.div>

      {/* Search */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="relative"
      >
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-400" />
        <input
          type="text"
          placeholder="搜索论文标题或作者..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="input-field pl-12"
        />
      </motion.div>

      {/* Papers List */}
      <div className="space-y-4">
        {filteredPapers.map((paper, index) => {
          const status = statusConfig[paper.status as keyof typeof statusConfig]
          const StatusIcon = status.icon

          return (
            <motion.div
              key={paper.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 + index * 0.05 }}
              className="glass-card p-6 hover:border-primary-500/30 transition-all"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 mb-2">
                    <span className={`px-2 py-1 rounded text-xs font-medium flex items-center gap-1 ${status.color}`}>
                      <StatusIcon className={`w-3 h-3 ${status.animate ? 'animate-spin' : ''}`} />
                      {status.label}
                    </span>
                    <span className="category-badge">{paper.category}</span>
                  </div>

                  <h3 className="text-lg font-semibold mb-2 line-clamp-1">{paper.title}</h3>
                  
                  <p className="text-dark-400 text-sm mb-3">
                    {paper.authors.join(', ')}
                  </p>
                  
                  <p className="text-dark-500 text-sm line-clamp-2">{paper.abstract}</p>

                  <div className="flex items-center gap-4 mt-4">
                    <span className="text-dark-500 text-sm">
                      上传于 {paper.createdAt}
                    </span>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {paper.status === 'completed' && (
                    <button className="p-2 rounded-lg bg-primary-500/10 text-primary-400 hover:bg-primary-500/20 transition-colors">
                      <Eye className="w-5 h-5" />
                    </button>
                  )}
                  {paper.status !== 'processing' && (
                    <button className="p-2 rounded-lg bg-dark-800 text-dark-400 hover:bg-dark-700 hover:text-primary-400 transition-colors">
                      <Brain className="w-5 h-5" />
                    </button>
                  )}
                  <button className="p-2 rounded-lg bg-dark-800 text-dark-400 hover:bg-dark-700 hover:text-red-400 transition-colors">
                    <Trash2 className="w-5 h-5" />
                  </button>
                  <button className="p-2 rounded-lg bg-dark-800 text-dark-400 hover:bg-dark-700 transition-colors">
                    <MoreVertical className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </motion.div>
          )
        })}
      </div>

      {filteredPapers.length === 0 && (
        <div className="text-center py-12">
          <FileText className="w-12 h-12 text-dark-600 mx-auto mb-4" />
          <p className="text-dark-400">
            {searchQuery ? '没有找到匹配的论文' : '还没有上传论文，开始上传您的第一篇论文吧！'}
          </p>
        </div>
      )}
    </div>
  )
}

