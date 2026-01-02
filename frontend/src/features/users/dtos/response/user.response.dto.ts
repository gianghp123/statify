import { UserRole } from "@/lib/enums/user-role.enum";

export interface UserDto {
  id: number;
  username: string;
  email: string;
  role: UserRole;
  createdAt: string;
  updatedAt: string;
}
