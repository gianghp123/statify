import { getAuthTokenServer } from "@/lib/cookies/cookies-actions";

export async function GET(req: Request, ctx: RouteContext<'/api/sse/project/[projectId]/stream-analytics'>) {
  const { projectId } = await ctx.params
  const { searchParams } = new URL(req.url);

  const startTime = searchParams.get("startTime")
  const endTime = searchParams.get("endTime")

  const query = new URLSearchParams()
  if (startTime) query.set("startTime", startTime)
  if (endTime) query.set("endTime", endTime)

  const token = await getAuthTokenServer()

  const backendRes = await fetch(
    `${process.env.API_URL}/projects/${projectId}/stream-analytics${query.toString() ? `?${query.toString()}` : ""}`,
    {
      headers: {
        Authorization: `Bearer ${token}`
      },
      cache: "no-store",
    }
  )

  if (!backendRes.body) {
    return new Response("No stream", { status: 500 })
  }

  return new Response(backendRes.body, {
    status: backendRes.status,
    headers: {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      "Connection": "keep-alive",
    },
  })

}