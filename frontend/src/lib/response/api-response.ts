export interface BaseResponse<T> {
    success: boolean;
    code: number;
    message?: string;
    data?: T;
}

export interface BasePaginatedResponse<T> extends BaseResponse<T> {
    pagination?: Pagination;
}

export interface Pagination {
    totalCount: number;
    page: number;
    limit: number;
}