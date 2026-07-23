import { memo } from 'react'
import { formatMoney, formatDate } from '../utils/format'

// ----- Линейный график -----
interface LineChartProps {
  data: number[][]
  labels: string[]
  colors: string[]
  max: number
}

export const LineChart = memo(({ data, labels, colors, max }: LineChartProps) => {
  const padding = { top: 30, bottom: 30, left: 60, right: 40 }
  const width = 600
  const height = 250
  const innerWidth = width - padding.left - padding.right
  const innerHeight = height - padding.top - padding.bottom

  const steps = 5
  const stepValue = Math.ceil(max / steps)
  const maxScaled = stepValue * steps

  const points = data.map((dataset) =>
    dataset.map((value, i) => ({
      x: padding.left + (i / (labels.length - 1 || 1)) * innerWidth,
      y: padding.top + innerHeight - (value / maxScaled) * innerHeight,
    }))
  )

  const pathD = points.map((p) =>
    p.map((pt, i) => (i === 0 ? `M ${pt.x} ${pt.y}` : `L ${pt.x} ${pt.y}`)).join(' ')
  )

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      preserveAspectRatio="xMidYMid meet"
      style={{ width: '100%', height: '100%' }}
    >
      {/* Сетка */}
      {Array.from({ length: steps + 1 }, (_, i) => {
        const y = padding.top + innerHeight - (i / steps) * innerHeight
        const val = (i / steps) * maxScaled
        return (
          <g key={i}>
            <line x1={padding.left} y1={y} x2={width - padding.right} y2={y} stroke="#333" strokeWidth="1" strokeDasharray="4" />
            <text x={padding.left - 8} y={y + 4} textAnchor="end" fontSize="10" fill="#666">
              {formatMoney(val)}
            </text>
          </g>
        )
      })}
      {/* Подписи X */}
      {labels.map((label, i) => {
        const x = padding.left + (i / (labels.length - 1 || 1)) * innerWidth
        return (
          <text key={i} x={x} y={height - 4} textAnchor="middle" fontSize="10" fill="#666">
            {label}
          </text>
        )
      })}
      {/* Линии */}
      {pathD.map((d, idx) => (
        <path key={idx} d={d} fill="none" stroke={colors[idx]} strokeWidth="2.5" strokeLinejoin="round" strokeLinecap="round" />
      ))}
      {/* Точки */}
      {points.map((dataset, idx) =>
        dataset.map((pt, i) => (
          <circle key={`${idx}-${i}`} cx={pt.x} cy={pt.y} r="4" fill={colors[idx]} stroke="#0a0a0f" strokeWidth="1.5">
            <title>
              {labels[i]}: {formatMoney(data[idx][i])}
            </title>
          </circle>
        ))
      )}
    </svg>
  )
})

// ----- Накопленный баланс -----
interface CumulativeChartProps {
  data: { label: string; value: number }[]
}

export const CumulativeChart = memo(({ data }: CumulativeChartProps) => {
  const padding = { top: 30, bottom: 30, left: 60, right: 40 }
  const width = 600
  const height = 180
  const innerWidth = width - padding.left - padding.right
  const innerHeight = height - padding.top - padding.bottom

  const maxAbs = data.reduce((m, d) => Math.max(m, Math.abs(d.value)), 1)
  const magnitude = Math.pow(10, Math.floor(Math.log10(maxAbs)))
  let maxScaled = Math.ceil(maxAbs / magnitude) * magnitude
  if (maxScaled === 0) maxScaled = 1000

  const centerY = padding.top + innerHeight / 2

  const points = data.map((d, i) => ({
    x: padding.left + (i / (data.length - 1 || 1)) * innerWidth,
    y: centerY - (d.value / maxScaled) * (innerHeight / 2),
  }))

  const path = points.map((p, i) => (i === 0 ? `M ${p.x} ${p.y}` : `L ${p.x} ${p.y}`)).join(' ')

  const steps = 4
  const stepVal = maxScaled / steps
  const yValues: number[] = []
  for (let i = -steps; i <= steps; i++) {
    yValues.push(i * stepVal)
  }

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      preserveAspectRatio="xMidYMid meet"
      style={{ width: '100%', height: '100%' }}
    >
      {/* Горизонтальные линии */}
      {yValues.map((val) => {
        const y = centerY - (val / maxScaled) * (innerHeight / 2)
        if (Math.abs(val) < 0.001) return null
        return (
          <g key={val}>
            <line x1={padding.left} y1={y} x2={width - padding.right} y2={y} stroke="#333" strokeWidth="1" strokeDasharray="4" />
            <text x={padding.left - 8} y={y + 4} textAnchor="end" fontSize="10" fill="#666">
              {formatMoney(val)}
            </text>
          </g>
        )
      })}
      {/* Нулевая линия */}
      <line x1={padding.left} y1={centerY} x2={width - padding.right} y2={centerY} stroke="#555" strokeWidth="1.5" strokeDasharray="6" />
      <text x={padding.left - 8} y={centerY + 4} textAnchor="end" fontSize="10" fill="#888">
        {formatMoney(0)}
      </text>

      {/* Область под линией */}
      {points.length > 1 && (
        <path
          d={`${path} L ${points[points.length - 1].x} ${centerY} L ${points[0].x} ${centerY} Z`}
          fill="url(#cumulativeGradient)"
          opacity="0.3"
        />
      )}
      <defs>
        <linearGradient id="cumulativeGradient" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="#4ade80" stopOpacity="0.4" />
          <stop offset="100%" stopColor="#4ade80" stopOpacity="0" />
        </linearGradient>
      </defs>

      {/* Линия */}
      <path d={path} fill="none" stroke="#4ade80" strokeWidth="2.5" strokeLinejoin="round" strokeLinecap="round" />
      {/* Точки */}
      {points.map((p, i) => (
        <circle key={i} cx={p.x} cy={p.y} r="4" fill="#4ade80" stroke="#0a0a0f" strokeWidth="1.5">
          <title>
            {formatDate(data[i].label)}: {formatMoney(data[i].value)}
          </title>
        </circle>
      ))}
      {/* Подписи дат */}
      {data.map((d, i) => {
        const x = padding.left + (i / (data.length - 1 || 1)) * innerWidth
        return (
          <text key={i} x={x} y={height - 4} textAnchor="middle" fontSize="9" fill="#666">
            {formatDate(d.label)}
          </text>
        )
      })}
    </svg>
  )
})