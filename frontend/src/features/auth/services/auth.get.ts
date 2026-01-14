import { apiFetch } from "@/lib/api-fetch";
import { UserDto } from "@/features/users/dtos/response/user.response.dto";
import { BaseResponse } from "@/lib/response/api-response";

export function getCurrentUser() {
    return apiFetch<BaseResponse<UserDto>>("/auth/me", {
        withCredentials: true,
    });
}