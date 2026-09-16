import { useCallback, useMemo, useState } from 'react'
import { motion } from 'framer-motion'
import { useDropzone } from 'react-dropzone'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  FileText,
  Upload,
  Search,
  Eye,
  Trash2,
  Brain,
  Clock,
  CheckCircle,
  AlertCircle,
  Loader2,
  X,
} from 'lucide-react'
import toast from 'react-hot-toast'
import { paperApi } from '@/services/api'
import type { Paper, PaperAnalysis } from '@/types'

const statusConfig = {
  pending: { icon: Clock, label: '等待处理', color: 'text-yellow-400 bg-yellow-400/10', animate: false },
  processing: { icon: Loader2, label: '分析中', color: 'text-blue-400 bg-blue-400/10', animate: true },
  completed: { icon: CheckCircle, label: '已完成', color: 'text-green-400 bg-green-400/10', animate: false },
  failed: { icon: AlertCircle, label: '处理失败', color: 'text-red-400 bg-red-400/10', animate: false },
}

function errMsg(error: unknown) {
  const err = error as { response?: { data?: { error?: { message?: string } } } }
  return err.response?.data?.error?.message || '操作失败'
}

export default function PapersPage() {
  const queryClient = useQueryClient()
  const [searchQuery, setSearchQuery] = useState('')
  const [selected, setSelected] = useState<Paper | null>(null)

  const { data: papers = [], isLoading } = useQuery({
    queryKey: ['papers'],
    queryFn: async () => {
      const res = await paperApi.list({ page: 1, per_page: 50 })
      return (res.data.data || []) as Paper[]
    },
  })

  const uploadMutation = useMutation({
    mutationFn: async (file: File) => {
      const form = new FormData()
      form.append('file', file)
      form.append('title', file.name.replace(/\.pdf$/i, ''))
      return paperApi.upload(form)
    },
    onSuccess: () => {
      toast.success('上传成功')
      queryClient.invalidateQueries({ queryKey: ['papers'] })
    },
    onError: (e) => toast.error(errMsg(e)),
  })

  const analyzeMutation = useMutation({
    mutationFn: (id: string) => paperApi.analyze(id, ['summary', 'methodology', 'contributions']),
    onSuccess: () => {
      toast.success('分析完成')
      queryClient.invalidateQueries({ queryKey: ['papers'] })
    },
    onError: (e) => toast.error(errMsg(e)),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => paperApi.delete(id),
    onSuccess: () => {
      toast.success('已删除')
      queryClient.invalidateQueries({ queryKey: ['papers'] })
    },
    onError: (e) => toast.error(errMsg(e)),
  })

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles[0]) {
      uploadMutation.mutate(acceptedFiles[0])
    }
  }, [uploadMutation])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { 'application/pdf': ['.pdf'] },
    maxFiles: 1,
  })

  const filteredPapers = useMemo(() => {
    const q = searchQuery.toLowerCase()
    return papers.filter((paper) =>
      paper.title.toLowerCase().includes(q) ||
      (paper.authors || []).some((a) => a.toLowerCase().includes(q))
    )
  }, [papers, searchQuery])

  const openDetail = async (paper: Paper) => {
    try {
      const res = await paperApi.get(paper.id)
      setSelected(res.data.data as Paper)
    } catch (e) {
      toast.error(errMsg(e))
    }
  }

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-3">
            <FileText className="w-7 h-7 text-primary-400" />
            我的论文
          </h1>
          <p className="text-dark-400 mt-1">上传论文进行 AI 智能分析</p>
        </div>
      </motion.div>

      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
        <div
          {...getRootProps()}
          className={`glass-card p-8 border-2 border-dashed cursor-pointer transition-all ${
            isDragActive ? 'border-primary-500 bg-primary-500/5' : 'border-dark-600 hover:border-primary-500/50'
          }`}
        >
          <input {...getInputProps()} />
          <div className="text-center">
            {uploadMutation.isPending ? (
              <Loader2 className="w-12 h-12 text-primary-400 mx-auto mb-4 animate-spin" />
            ) : (
              <Upload className="w-12 h-12 text-dark-500 mx-auto mb-4" />
            )}
            <p className="text-lg font-medium mb-2">
              {isDragActive ? '释放以上传文件' : uploadMutation.isPending ? '正在上传...' : '拖放 PDF 文件到这里'}
            </p>
            <p className="text-dark-400 text-sm">或者点击选择文件 • 支持 PDF 格式</p>
          </div>
        </div>
      </motion.div>

      <div className="relative">
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-400" />
        <input
          type="text"
          placeholder="搜索论文标题或作者..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="input-field pl-12"
        />
      </div>

      {isLoading && <p className="text-dark-400">加载中...</p>}

      <div className="space-y-4">
        {filteredPapers.map((paper) => {
          const status = statusConfig[paper.status] || statusConfig.pending
          const StatusIcon = status.icon
          return (
            <div key={paper.id} className="glass-card p-6 hover:border-primary-500/30 transition-all">
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 mb-2">
                    <span className={`px-2 py-1 rounded text-xs font-medium flex items-center gap-1 ${status.color}`}>
                      <StatusIcon className={`w-3 h-3 ${status.animate ? 'animate-spin' : ''}`} />
                      {status.label}
                    </span>
                    {paper.ccf_category && <span className="category-badge">{paper.ccf_category.name}</span>}
                  </div>
                  <h3 className="text-lg font-semibold mb-2 line-clamp-1">{paper.title}</h3>
                  <p className="text-dark-400 text-sm mb-3">{(paper.authors || []).join(', ') || '作者未知'}</p>
                  <p className="text-dark-500 text-sm line-clamp-2">{paper.abstract}</p>
                  <p className="text-dark-500 text-sm mt-4">上传于 {paper.created_at?.slice(0, 10)}</p>
                </div>
                <div className="flex items-center gap-2">
                  <button onClick={() => openDetail(paper)} className="p-2 rounded-lg bg-primary-500/10 text-primary-400 hover:bg-primary-500/20" title="查看">
                    <Eye className="w-5 h-5" />
                  </button>
                  {paper.status !== 'processing' && (
                    <button
                      onClick={() => analyzeMutation.mutate(paper.id)}
                      disabled={analyzeMutation.isPending}
                      className="p-2 rounded-lg bg-dark-800 text-dark-400 hover:text-primary-400"
                      title="分析"
                    >
                      {analyzeMutation.isPending ? <Loader2 className="w-5 h-5 animate-spin" /> : <Brain className="w-5 h-5" />}
                    </button>
                  )}
                  <button onClick={() => deleteMutation.mutate(paper.id)} className="p-2 rounded-lg bg-dark-800 text-dark-400 hover:text-red-400" title="删除">
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {!isLoading && filteredPapers.length === 0 && (
        <div className="text-center py-12">
          <FileText className="w-12 h-12 text-dark-600 mx-auto mb-4" />
          <p className="text-dark-400">{searchQuery ? '没有找到匹配的论文' : '还没有上传论文，开始上传您的第一篇论文吧！'}</p>
        </div>
      )}

      {selected && (
        <div className="fixed inset-0 z-50 bg-dark-950/80 backdrop-blur-sm flex items-center justify-center p-6" onClick={() => setSelected(null)}>
          <div className="glass-card max-w-3xl w-full max-h-[80vh] overflow-auto p-6" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-start justify-between mb-4">
              <h2 className="text-xl font-bold pr-4">{selected.title}</h2>
              <button onClick={() => setSelected(null)}><X className="w-5 h-5 text-dark-400" /></button>
            </div>
            <p className="text-dark-400 text-sm mb-4">{selected.abstract}</p>
            {(selected.analyses || []).length === 0 && (
              <p className="text-dark-500 text-sm">尚未分析。点击大脑图标生成摘要、方法与贡献。</p>
            )}
            {(selected.analyses || []).map((a: PaperAnalysis) => (
              <div key={a.id} className="mb-4 p-4 rounded-xl bg-dark-800/50">
                <h3 className="font-medium mb-2 text-primary-400">{a.analysis_type}</h3>
                <p className="text-sm text-dark-200 whitespace-pre-wrap">{a.content}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
