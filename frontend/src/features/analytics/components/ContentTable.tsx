export function ContentTable() {
  const topContent = [
    { path: "/home", visits: "8,432", unique: "6,120", avgTime: "2m 15s", bounce: "24.5%" },
    { path: "/products/analytics-pro", visits: "3,210", unique: "2,980", avgTime: "4m 32s", bounce: "18.2%" },
    { path: "/blog/optimizing-latency", visits: "2,105", unique: "1,890", avgTime: "8m 12s", bounce: "45.1%" },
    { path: "/pricing", visits: "1,540", unique: "1,200", avgTime: "1m 05s", bounce: "32.0%" },
    { path: "/docs/api-v2", visits: "980", unique: "560", avgTime: "12m 40s", bounce: "12.5%" },
  ];

  return (
    <div className="bg-card rounded-xl overflow-hidden border border-border shadow-xl">
      <div className="px-6 py-4 border-b border-white/5 flex justify-between items-center bg-white/5">
        <h3 className="text-foreground font-bold text-lg">Top Content</h3>
        <button className="text-xs font-bold text-primary hover:text-foreground hover:underline transition-colors">
          View All Reports
        </button>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm text-muted-foreground">
          <thead className="text-xs uppercase bg-black/40 text-foreground border-b border-border">
            <tr>
              <th className="px-6 py-3 font-semibold tracking-wider text-primary" scope="col">Page Path</th>
              <th className="px-6 py-3 font-semibold tracking-wider" scope="col">Visits</th>
              <th className="px-6 py-3 font-semibold tracking-wider" scope="col">Unique</th>
              <th className="px-6 py-3 font-semibold tracking-wider" scope="col">Avg. Time</th>
              <th className="px-6 py-3 font-semibold tracking-wider" scope="col">Bounce Rate</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border bg-white/2">
            {topContent.map((row) => (
              <tr key={row.path} className="hover:bg-primary/5 transition-colors group">
                <td className="px-6 py-4 font-medium text-foreground group-hover:text-primary transition-colors">{row.path}</td>
                <td className="px-6 py-4">{row.visits}</td>
                <td className="px-6 py-4">{row.unique}</td>
                <td className="px-6 py-4">{row.avgTime}</td>
                <td className={`px-6 py-4 font-bold ${row.bounce.replace('%', '') > '40' ? 'text-red-400' : 'text-primary'}`}>{row.bounce}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
