import { cookies } from 'next/headers'

export const getAuthTokenServer = async () => {
    const cookieStore = await cookies()
    return cookieStore.get("token")?.value
}

export const setAuthTokenServer = async (token: string) => {
    const cookieStore = await cookies()
    cookieStore.set("token", token)
}

export const removeAuthTokenServer = async () => {
    const cookieStore = await cookies()
    cookieStore.delete("token")
}
