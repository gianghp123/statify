'use server'

import { apiFetch } from "@/lib/api-fetch";
import { LoginRequestDto } from "../dtos/request/login.request.dto";
import { RegisterRequestDto } from "../dtos/request/register.request.dto";
import { AuthResponseDto } from "../dtos/response/auth.response.dto";

export function login(loginRequestDto: LoginRequestDto) {
    return apiFetch<AuthResponseDto>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(loginRequestDto),
    })
}

// export function logout() {
//     return apiFetch('/auth/logout', {
//         method: 'POST',
//     })
// }


export function register(registerRequestDto: RegisterRequestDto) {
    return apiFetch<AuthResponseDto>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(registerRequestDto),
    })
}