import { cookies } from 'next/headers'

export const getAuthTokenServer = async () => {
    const cookieStore = await cookies()
    return cookieStore.get("auth-token")?.value
}

export const setAuthTokenServer = async (token: string) => {
    const cookieStore = await cookies()
    cookieStore.set("auth-token", token)
}

export const removeAuthTokenServer = async () => {
    const cookieStore = await cookies()
    cookieStore.delete("auth-token")
}
