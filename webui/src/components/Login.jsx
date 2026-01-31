import { useState } from 'react'
import { Key, ArrowRight, ShieldCheck, Lock } from 'lucide-react'
import clsx from 'clsx'

export default function Login({ onLogin, onMessage }) {
    const [adminKey, setAdminKey] = useState('')
    const [loading, setLoading] = useState(false)
    const [remember, setRemember] = useState(true)

    const handleLogin = async (e) => {
        e.preventDefault()
        if (!adminKey.trim()) return

        setLoading(true)

        try {
            const res = await fetch('/admin/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ admin_key: adminKey }),
            })

            const data = await res.json()

            if (res.ok && data.success) {
                const storage = remember ? localStorage : sessionStorage
                storage.setItem('ds2api_token', data.token)
                storage.setItem('ds2api_token_expires', Date.now() + data.expires_in * 1000)

                onLogin(data.token)
                if (data.message) {
                    onMessage('warning', data.message)
                }
            } else {
                onMessage('error', data.detail || '登录失败')
            }
        } catch (e) {
            onMessage('error', '网络错误: ' + e.message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="flex flex-col items-center justify-center p-4 w-full max-w-md relative z-10">
            <div className="w-full bg-card/50 backdrop-blur-xl border border-white/10 shadow-2xl rounded-2xl p-8 space-y-8 animate-in fade-in zoom-in-95 duration-300">
                <div className="text-center space-y-2">
                    <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-primary/20 text-primary mb-4 ring-1 ring-white/10 shadow-inner">
                        <Lock className="w-8 h-8" />
                    </div>
                    <h1 className="text-2xl font-bold tracking-tight text-white">
                        欢迎回来
                    </h1>
                    <p className="text-muted-foreground text-sm">
                        请输入管理员密钥访问控制面板
                    </p>
                </div>

                <form onSubmit={handleLogin} className="space-y-6">
                    <div className="space-y-4">
                        <div className="space-y-2">
                            <label className="text-xs font-medium text-muted-foreground ml-1 uppercase tracking-wider">
                                管理员密钥
                            </label>
                            <div className="relative group">
                                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-muted-foreground group-focus-within:text-primary transition-colors">
                                    <Key className="w-4 h-4" />
                                </div>
                                <input
                                    type="password"
                                    className="block w-full pl-10 pr-3 py-2.5 bg-secondary/50 border border-border rounded-lg text-sm placeholder:text-muted-foreground/50 focus:border-primary/50 focus:ring-2 focus:ring-primary/20 focus:bg-secondary focus:outline-none transition-all duration-200"
                                    placeholder="输入您的 DS2API_ADMIN_KEY"
                                    value={adminKey}
                                    onChange={e => setAdminKey(e.target.value)}
                                    autoFocus
                                />
                            </div>
                        </div>

                        <div className="flex items-center space-x-2">
                            <button
                                type="button"
                                role="checkbox"
                                aria-checked={remember}
                                onClick={() => setRemember(!remember)}
                                className={clsx(
                                    "w-4 h-4 rounded border flex items-center justify-center transition-all duration-200",
                                    remember ? "bg-primary border-primary text-primary-foreground" : "border-muted-foreground/50 bg-transparent"
                                )}
                            >
                                {remember && <div className="w-2 h-2 rounded-[1px] bg-current" />}
                            </button>
                            <span
                                onClick={() => setRemember(!remember)}
                                className="text-sm text-muted-foreground cursor-pointer select-none hover:text-foreground transition-colors"
                            >
                                记住登录状态
                            </span>
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full flex items-center justify-center py-2.5 px-4 rounded-lg bg-primary hover:bg-primary/90 text-primary-foreground font-medium text-sm transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary/20 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-primary/25 hover:shadow-primary/40 hover:-translate-y-0.5 active:translate-y-0"
                    >
                        {loading ? (
                            <div className="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin" />
                        ) : (
                            <>
                                <span>进入控制面板</span>
                                <ArrowRight className="w-4 h-4 ml-2" />
                            </>
                        )}
                    </button>
                </form>

                <div className="pt-6 border-t border-border/50 text-center">
                    <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground/70">
                        <ShieldCheck className="w-3 h-3" />
                        <span>安全会话 • 24小时有效</span>
                    </div>
                </div>
            </div>
        </div>
    )
}
