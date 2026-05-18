export default function PromptSection({ t, form, setForm }) {
    const updatePrompt = (key, value) => {
        setForm((prev) => ({
            ...prev,
            prompts: { ...prev.prompts, [key]: value },
        }))
    }

    const dynamicPlaceholder = t('settings.promptDynamicPlaceholder')

    const fields = [
        { key: 'output_integrity_guard', label: t('settings.outputIntegrityGuard'), desc: t('settings.outputIntegrityGuardDesc'), placeholder: '', rows: 4 },
        { key: 'tool_call_instructions', label: t('settings.toolCallInstructions'), desc: t('settings.toolCallInstructionsDesc'), placeholder: dynamicPlaceholder, rows: 8 },
        { key: 'tool_descriptions_prefix', label: t('settings.toolDescriptionsPrefix'), desc: t('settings.toolDescriptionsPrefixDesc'), placeholder: '', rows: 2 },
        { key: 'read_tool_cache_guard', label: t('settings.readToolCacheGuard'), desc: t('settings.readToolCacheGuardDesc'), placeholder: '', rows: 4 },
        { key: 'tool_call_examples', label: t('settings.toolCallExamples'), desc: t('settings.toolCallExamplesDesc'), placeholder: dynamicPlaceholder, rows: 6 },
    ]

    return (
        <div className="bg-card border border-border rounded-xl p-5 space-y-4">
            <h3 className="font-semibold">{t('settings.promptSectionTitle')}</h3>
            <p className="text-sm text-muted-foreground">{t('settings.promptSectionDesc')}</p>
            <div className="space-y-4">
                {fields.map(({ key, label, desc, placeholder, rows }) => (
                    <label key={key} className="text-sm space-y-1">
                        <span className="text-muted-foreground">{label}</span>
                        <p className="text-xs text-muted-foreground/70">{desc}</p>
                        <textarea
                            value={form.prompts?.[key] || ''}
                            onChange={(e) => updatePrompt(key, e.target.value)}
                            placeholder={placeholder || undefined}
                            rows={rows}
                            className="w-full bg-background border border-border rounded-lg px-3 py-2 font-mono text-xs resize-y placeholder:text-muted-foreground/40"
                        />
                    </label>
                ))}
            </div>
        </div>
    )
}
