export default function CurrentInputFileSection({ t, form, setForm }) {
    const toggles = [
        { key: 'flash', label: t('settings.currentInputFileFlash') },
        { key: 'pro', label: t('settings.currentInputFilePro') },
        { key: 'vision', label: t('settings.currentInputFileVision') },
    ]

    return (
        <div className="bg-card border border-border rounded-xl p-5 space-y-4">
            <div className="space-y-1">
                <h3 className="font-semibold">{t('settings.currentInputFileTitle')}</h3>
                <p className="text-sm text-muted-foreground">{t('settings.currentInputFileDesc')}</p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {toggles.map((toggle) => (
                    <label key={toggle.key} className="flex items-start gap-3 rounded-lg border border-border bg-background/60 p-4">
                        <input
                            type="checkbox"
                            checked={Boolean(form.current_input_file?.[toggle.key])}
                            onChange={(e) => setForm((prev) => ({
                                ...prev,
                                current_input_file: {
                                    ...prev.current_input_file,
                                    [toggle.key]: e.target.checked,
                                },
                            }))}
                            className="mt-1 h-4 w-4 rounded border-border"
                        />
                        <div className="space-y-1">
                            <span className="text-sm font-medium block">{toggle.label}</span>
                            <span className="text-xs text-muted-foreground block">{t('settings.currentInputFileModelHelp')}</span>
                        </div>
                    </label>
                ))}
            </div>
        </div>
    )
}
