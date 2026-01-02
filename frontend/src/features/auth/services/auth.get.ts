import { apiFetch } from "@/lib/api-fetch";
import { UserDto } from "@/features/users/dtos/response/user.response.dto";

export function getCurrentUser() {
    return apiFetch<UserDto>("/auth/current-user");
}