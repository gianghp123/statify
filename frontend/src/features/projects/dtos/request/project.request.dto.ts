export interface CreateProjectRequestDto {
  name: string;
  subdomain: string;
}

export interface UpdateProjectRequestDto {
  name?: string;
  subdomain?: string;
}
