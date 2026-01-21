import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { TimeSeriesPointDTO } from "@/features/analytics/dtos/response/analytic-metrics.dto";


interface AnalyticsChartsProps {
  data: TimeSeriesPointDTO[];
}

export function AnalyticsCharts({ data }: AnalyticsChartsProps) {
  const chartData = data.map(point => ({
    name: new Date(point.timestamp).toLocaleDateString(),
    visitors: point.requests,
    bandwidth: point.bandwidth.toFixed(5),
  }));

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div className="bg-card rounded-xl p-6 border border-border shadow-xl">
        <h3 className="text-foreground font-semibold mb-6">Visitor Traffic</h3>
        <div className="h-[300px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorVisitors" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ccff00" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#ccff00" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#2a1b3d" vertical={false} />
              <XAxis dataKey="name" stroke="#aba3bf" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis stroke="#aba3bf" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(value) => `${value}`} />
              <Tooltip
                contentStyle={{ backgroundColor: '#180f26', border: '1px solid rgba(255,255,255,0.1)', borderRadius: '8px' }}
                itemStyle={{ color: '#ccff00' }}
              />
              <Area type="monotone" dataKey="visitors" stroke="#ccff00" fillOpacity={1} fill="url(#colorVisitors)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="bg-card rounded-xl p-6 border border-border shadow-xl">
        <h3 className="text-foreground font-semibold mb-6">Bandwidth Consumption (MB)</h3>
        <div className="h-[300px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorBandwidth" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#a78bfa" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#a78bfa" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#2a1b3d" vertical={false} />
              <XAxis dataKey="name" stroke="#aba3bf" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis stroke="#aba3bf" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(value) => `${value}`} />
              <Tooltip
                contentStyle={{ backgroundColor: '#180f26', border: '1px solid rgba(255,255,255,0.1)', borderRadius: '8px' }}
                itemStyle={{ color: '#a78bfa' }}
              />
              <Area type="monotone" dataKey="bandwidth" stroke="#a78bfa" fillOpacity={1} fill="url(#colorBandwidth)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}
