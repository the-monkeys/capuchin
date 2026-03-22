interface EmptyStateProps {
  message?: string
}

export function EmptyState({ message = "No tasks yet" }: EmptyStateProps) {
  return (
    <div className="text-center py-14 px-6 text-[#b0b7c3]">
      <div className="text-[2.6rem] mb-2.5 opacity-50">🐒</div>
      <div className="text-[0.88rem] font-medium">{message}</div>
    </div>
  )
}