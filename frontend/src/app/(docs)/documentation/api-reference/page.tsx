import { BreadScrum } from "@/components/bread-scrum";
import { Copy, FileArchive } from "lucide-react";

export default function ApiReferencePage() {
  return (
    <div className="flex flex-col xl:flex-row gap-10">
      <div className="flex-1 min-w-0">
        <BreadScrum items={[
          { name: "Docs", href: "/documentation" },
          { name: "API Reference", isCurrent: true }
        ]} />

        <article className="prose prose-slate dark:prose-invert max-w-none">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
            Create Deployment
          </h1>
          <p className="leading-7 not-first:mt-6 text-xl text-muted-foreground mb-8">
            Create a new deployment for a specific project. This endpoint accepts a multipart form-data upload containing a .zip file of your static assets.
          </p>

          <section className="mb-10 not-prose">
            <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground mb-4 border-b pb-2">
              Endpoint
            </h2>
            <div className="flex items-center gap-2 font-mono text-sm bg-muted/50 p-3 rounded-lg border border-border">
              <span className="text-muted-foreground">https://api.statify.online/v1</span>
              <span className="text-primary font-bold">/projects/{`{project_id}`}/deployments</span>
            </div>
            <p className="text-xs text-muted-foreground mt-2 italic">* Replace `{`{project_id}`}` with your numeric project identifier.</p>
          </section>

          <section className="mb-10 not-prose">
            <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground mb-4 border-b pb-2">
              Request Format
            </h2>
            <p className="text-sm text-muted-foreground mb-4">
              The request must be sent as `multipart/form-data`.
            </p>
            <div className="space-y-6">
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-sm font-semibold">file</span>
                  <span className="text-xs px-1.5 py-0.5 bg-muted rounded text-muted-foreground uppercase">file (.zip)</span>
                  <span className="text-xs text-orange-500 font-medium">Required</span>
                </div>
                <p className="text-sm text-muted-foreground">
                  A binary .zip file containing all the static assets for your website.
                </p>
              </div>
            </div>
          </section>

          <section className="mb-10 not-prose">
            <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground mb-4 border-b pb-2">
              Responses
            </h2>
            <div className="bg-card rounded-xl border overflow-hidden shadow-sm">
              <table className="w-full text-left text-sm">
                <thead className="bg-muted/50 border-b">
                  <tr>
                    <th className="px-6 py-3 font-semibold text-muted-foreground">Code</th>
                    <th className="px-6 py-3 font-semibold text-muted-foreground">Description</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  <tr>
                    <td className="px-6 py-4 font-mono font-medium text-emerald-500">201 Created</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      Deployment created and processing started.
                    </td>
                  </tr>
                  <tr>
                    <td className="px-6 py-4 font-mono font-medium text-amber-500">400 Bad Request</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      No file provided or invalid project ID.
                    </td>
                  </tr>
                  <tr>
                    <td className="px-6 py-4 font-mono font-medium text-red-500">401 Unauthorized</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      Missing or invalid authentication token.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>
        </article>
      </div>

      <div className="lg:col-span-5 space-y-6 xl:w-[450px]">
        <div className="sticky top-16 pt-8">
          <div className="bg-slate-950 dark:bg-black rounded-2xl overflow-hidden shadow-2xl border border-slate-800">
            <div className="flex items-center justify-between px-4 py-3 bg-slate-900 border-b border-slate-800">
              <div className="flex items-center gap-2">
                <div className="flex gap-1.5">
                  <div className="w-2.5 h-2.5 rounded-full bg-slate-700"></div>
                  <div className="w-2.5 h-2.5 rounded-full bg-slate-700"></div>
                  <div className="w-2.5 h-2.5 rounded-full bg-slate-700"></div>
                </div>
                <span className="text-xs font-medium text-slate-400 ml-2">curl --request POST</span>
              </div>
              <button className="flex items-center gap-1 text-[10px] font-bold text-slate-400 hover:text-white transition-colors">
                <Copy className="w-3 h-3" />
                COPY
              </button>
            </div>
            <div className="p-6 overflow-x-auto">
              <pre className="font-mono text-[13px] leading-relaxed text-slate-300">
                <code>
                  <span className="text-primary">curl</span> -X POST https://api.statify.online/v1/projects/1/deployments \{"\n"}
                  {"  "}-H <span className="text-emerald-400">"Authorization: Bearer <span className="text-slate-500">YOUR_TOKEN"</span></span> \{"\n"}
                  {"  "}-F <span className="text-emerald-400">"file=@./site.zip"</span>
                </code>
              </pre>
            </div>
          </div>

          <div className="mt-6 bg-slate-950 dark:bg-black rounded-2xl overflow-hidden shadow-2xl border border-slate-800">
            <div className="px-4 py-3 bg-slate-900 border-b border-slate-800 flex justify-between items-center">
              <span className="text-xs font-bold text-emerald-400 uppercase tracking-widest">Example Response</span>
              <span className="text-[10px] text-slate-500 font-mono">201 CREATED</span>
            </div>
            <div className="p-6">
              <pre className="font-mono text-[13px] leading-relaxed text-slate-300">
                <code>
                  {"{"}\n
                  {"  "}<span className="text-primary">"code"</span>: 201,\n
                  {"  "}<span className="text-primary">"message"</span>: <span className="text-emerald-400">"Deployment created"</span>,\n
                  {"  "}<span className="text-primary">"data"</span>: {"{"}\n
                  {"    "}<span className="text-primary">"id"</span>: 42,\n
                  {"    "}<span className="text-primary">"status"</span>: <span className="text-emerald-400">"READY"</span>,\n
                  {"    "}<span className="text-primary">"is_SPA"</span>: <span className="text-amber-400">false</span>\n
                  {"  "}{"}"}\n
                  {"}"}
                </code>
              </pre>
            </div>
          </div>

          <div className="mt-8 bg-primary/5 border border-primary/20 rounded-2xl p-5">
            <div className="flex items-center gap-3 mb-2 text-primary">
              <FileArchive className="w-5 h-5" />
              <h4 className="font-bold text-sm">Deployment Tips</h4>
            </div>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Ensure your `.zip` file contains an `index.html` at the root. Statify will automatically detect it and serve it as the entry point.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
