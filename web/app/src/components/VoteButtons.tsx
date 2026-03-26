type VoteColor = 'green' | 'yellow' | 'red'

interface Props {
  selected: VoteColor | null
  onSelect: (color: VoteColor) => void
}

export function VoteButtons({ selected, onSelect }: Props) {
  const colors: VoteColor[] = ['green', 'yellow', 'red']

  return (
    <div className="vote-buttons">
      {colors.map((color) => (
        <button
          key={color}
          type="button"
          className={`vote-circle vote-circle-${color} ${selected === color ? 'selected' : ''}`}
          onClick={() => onSelect(color)}
          aria-label={`Vote ${color}`}
        />
      ))}
    </div>
  )
}
