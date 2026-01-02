import { UserDto } from "@/features/users/dtos/response/user.response.dto";

export interface AuthResponseDto {
  user: UserDto;
  token: string;
}
