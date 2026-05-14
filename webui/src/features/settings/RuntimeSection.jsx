export default function RuntimeSection({ t, form, setForm }) {
    const scheduleMode = form.runtime.account_schedule_mode || 'fast_round_robin'

    return (
        <div className="bg-card border border-border rounded-xl p-5 space-y-4">
            <h3 className="font-semibold">{t('settings.runtimeTitle')}</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.accountMaxInflight')}</span>
                    <input
                        type="number"
                        min={1}
                        value={form.runtime.account_max_inflight}
                        onChange={(e) => setForm((prev) => ({
                            ...prev,
                            runtime: { ...prev.runtime, account_max_inflight: Number(e.target.value || 1) },
                        }))}
                        className="w-full bg-background border border-border rounded-lg px-3 py-2"
                    />
                </label>
                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.accountMaxQueue')}</span>
                    <input
                        type="number"
                        min={1}
                        value={form.runtime.account_max_queue}
                        onChange={(e) => setForm((prev) => ({
                            ...prev,
                            runtime: { ...prev.runtime, account_max_queue: Number(e.target.value || 1) },
                        }))}
                        className="w-full bg-background border border-border rounded-lg px-3 py-2"
                    />
                </label>
                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.globalMaxInflight')}</span>
                    <input
                        type="number"
                        min={1}
                        value={form.runtime.global_max_inflight}
                        onChange={(e) => setForm((prev) => ({
                            ...prev,
                            runtime: { ...prev.runtime, global_max_inflight: Number(e.target.value || 1) },
                        }))}
                        className="w-full bg-background border border-border rounded-lg px-3 py-2"
                    />
                </label>
                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.tokenRefreshIntervalHours')}</span>
                    <input
                        type="number"
                        min={1}
                        max={720}
                        step={1}
                        value={form.runtime.token_refresh_interval_hours}
                        onChange={(e) => setForm((prev) => ({
                            ...prev,
                            runtime: { ...prev.runtime, token_refresh_interval_hours: Number(e.target.value || 1) },
                        }))}
                        className="w-full bg-background border border-border rounded-lg px-3 py-2"
                    />
                </label>
                <label className="text-sm space-y-2">
                    <span className="text-muted-foreground">{t('settings.accountScheduleMode')}</span>
                    <select
                        value={scheduleMode}
                        onChange={(e) => setForm((prev) => ({
                            ...prev,
                            runtime: {
                                ...prev.runtime,
                                account_schedule_mode: e.target.value,
                            },
                        }))}
                        className="w-full bg-background border border-border rounded-lg px-3 py-2"
                    >
                        <option value="fast_round_robin">{t('settings.accountScheduleModeFast')}</option>
                        <option value="sticky_round_robin">{t('settings.accountScheduleModeSticky')}</option>
                    </select>
                </label>
                {scheduleMode === 'sticky_round_robin' && (
                    <label className="text-sm space-y-2">
                        <span className="text-muted-foreground">{t('settings.accountStickyReuseCount')}</span>
                        <input
                            type="number"
                            min={1}
                            value={form.runtime.account_sticky_reuse_count}
                            onChange={(e) => setForm((prev) => ({
                                ...prev,
                                runtime: { ...prev.runtime, account_sticky_reuse_count: Number(e.target.value || 1) },
                            }))}
                            className="w-full bg-background border border-border rounded-lg px-3 py-2"
                        />
                    </label>
                )}
            </div>
        </div>
    )
}
