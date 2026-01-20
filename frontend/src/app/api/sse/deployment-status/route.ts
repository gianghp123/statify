import { getAuthTokenServer } from "@/lib/cookies/cookies-actions"

export const runtime = "nodejs"

export async function GET() {
  const token = await getAuthTokenServer()

  const backendRes = await fetch(
    `${process.env.API_URL}/deployments/stream-status`,
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