'use server'

import { apiFetch } from "@/lib/api-fetch";
import { LoginRequestDto } from "../dtos/request/login.request.dto";
import { RegisterRequestDto } from "../dtos/request/register.request.dto";
import { AuthResponseDto } from "../dtos/response/auth.response.dto";
import { setAuthTokenServer, removeAuthTokenServer } from "@/lib/cookies/cookies-actions";
import { redirect } from "next/navigation";

export async function login(prevState: any, loginRequestDto: LoginRequestDto) {
    const res = await apiFetch<AuthResponseDto>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(loginRequestDto),
    });
    console.log(res.data);

    if (res.success && res.data?.token) {

        await setAuthTokenServer(res.data.token);
    }

    return res;
}

export async function logout() {
    await removeAuthTokenServer();
    redirect("/");
}

export async function register(prevState: any, registerRequestDto: RegisterRequestDto) {
    const res = await apiFetch<AuthResponseDto>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(registerRequestDto),
    });

    if (res.success && res.data?.token) {
        await setAuthTokenServer(res.data.token);
    }

    return res;
}
