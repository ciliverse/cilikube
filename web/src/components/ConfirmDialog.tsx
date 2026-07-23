import { useEffect, useState, type ReactNode } from 'react'
import { Button, Modal } from '@/components/ui'

type Props = {
  open: boolean
  title: string
  description?: ReactNode
  /** When set, user must type this exact string to enable confirm. */
  confirmText?: string
  confirmLabel?: string
  cancelLabel?: string
  danger?: boolean
  busy?: boolean
  onConfirm: () => void | Promise<void>
  onClose: () => void
}

export function ConfirmDialog({
  open,
  title,
  description,
  confirmText,
  confirmLabel = 'Confirm',
  cancelLabel = 'Cancel',
  danger = true,
  busy = false,
  onConfirm,
  onClose,
}: Props) {
  const [typed, setTyped] = useState('')

  useEffect(() => {
    if (open) setTyped('')
  }, [open, confirmText])

  const needMatch = Boolean(confirmText)
  const matched = !needMatch || typed === confirmText

  return (
    <Modal open={open} title={title} onClose={onClose}>
      <div className="space-y-4 px-5 py-4">
        {description ? <div className="text-sm text-text-dim">{description}</div> : null}
        {needMatch ? (
          <label className="block space-y-2">
            <span className="hud-label">
              Type <span className="text-cyan">{confirmText}</span> to confirm
            </span>
            <input
              className="hud-field"
              value={typed}
              autoFocus
              spellCheck={false}
              onChange={(e) => setTyped(e.target.value)}
              placeholder={confirmText}
            />
          </label>
        ) : null}
        <div className="flex justify-end gap-2 pt-2">
          <Button variant="ghost" type="button" disabled={busy} onClick={onClose}>
            {cancelLabel}
          </Button>
          <Button
            variant={danger ? 'danger' : 'primary'}
            type="button"
            disabled={!matched || busy}
            onClick={() => void onConfirm()}
          >
            {busy ? 'Working…' : confirmLabel}
          </Button>
        </div>
      </div>
    </Modal>
  )
}
