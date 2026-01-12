import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from "recharts";

const data = [
  { name: "Mon", visitors: 4000, bandwidth: 2400 },
  { name: "Tue", visitors: 3000, bandwidth: 1398 },
  { name: "Wed", visitors: 2000, bandwidth: 9800 },
  { name: "Thu", visitors: 2780, bandwidth: 3908 },
  { name: "Fri", visitors: 1890, bandwidth: 4800 },
  { name: "Sat", visitors: 2390, bandwidth: 3800 },
  { name: "Sun", visitors: 3490, bandwidth: 4300 },
];

export function AnalyticsCharts() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div className="bg-card rounded-xl p-6 border border-border shadow-xl">
        <h3 className="text-foreground font-semibold mb-6">Visitor Traffic</h3>
        <div className="h-[300px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="colorVisitors" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ccff00" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#ccff00" stopOpacity={0}/>
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
        <h3 className="text-foreground font-semibold mb-6">Bandwidth Consumption (GB)</h3>
        <div className="h-[300px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="colorBandwidth" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#a78bfa" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#a78bfa" stopOpacity={0}/>
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
