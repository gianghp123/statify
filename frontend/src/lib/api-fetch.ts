import "server-only";
import { getAuthTokenServer } from "./cookies/cookies-actions";
import { BasePaginatedResponse, BaseResponse } from "./response/api-response";
import { snakeToCamel } from "./case";

type ApiFetchOptions = {
  baseUrl?: string;
  withCredentials?: boolean;
  query?: Record<string, any>;
} & RequestInit;

export async function apiFetch<T = any>(
  url: string,
  options?: ApiFetchOptions
): Promise<BaseResponse<T> | BasePaginatedResponse<T>> {

  try {
    const {
      withCredentials = false,
      baseUrl = process.env.API_URL,
      query,
      ...fetchOptions
    } = options || {};

    if (!baseUrl) {
      throw new Error(
        "Server API_URL is not configured. Please set API_URL environment variable."
      );
    }

    const headers: Record<string, any> = {
      ...fetchOptions?.headers,
      apikey: process.env.API_KEY || "",
    };

    if (withCredentials) {
      const accessToken = await getAuthTokenServer();
      if (accessToken) headers["Authorization"] = `Bearer ${accessToken}`;
      else
        console.log('Unauthorized')
        return {
          success: false,
          message: "Unauthorized",
          code: 401,
          data: undefined,
        };
    }
    let queryString = "";
    if (query && Object.keys(query).length > 0) {
      const searchParams = new URLSearchParams();
      Object.entries(query).forEach(([key, value]) => {
        if (
          value !== undefined &&
          value !== null &&
          value !== "" &&
          (Array.isArray(value) ? value.length > 0 : true)
        ) {
          if (Array.isArray(value)) {
            value.forEach((v) => searchParams.append(key, String(v)));
          } else {
            searchParams.append(key, String(value));
          }
        }
      });

      queryString = `?${searchParams.toString()}`;
    }

    const fullUrl = `${baseUrl}${url}${queryString}`;

    console.log(fullUrl, headers)

    const response = await fetch(fullUrl, { ...fetchOptions, headers });

    console.log(response)
    if (!response.ok) {
      let message = "Unknown error";
      try {
        const errorData = await response.json();
        message = errorData.message || message;
      } catch (_) {}
      return {
        success: false,
        message,
        code: response.status,
        data: undefined,
      };
    }

    let rawData: any;
    const contentType = response.headers.get("content-type");

    if (contentType && contentType.includes("application/json")) {
      const text = await response.text();
      rawData = text ? JSON.parse(text) : {};
    } else {
      rawData = {};
    }

    const data = snakeToCamel<any>(rawData);

    if (data.pagination || data.meta) {
      return {
        success: true,
        data: data.data ? (data.data as T) : ([] as T),
        pagination: data.pagination || data.meta?.pagination,
      } as BasePaginatedResponse<T>;
    }
    return {
      success: true,
      code: response.status,
      data: data.data !== undefined ? (data.data as T) : (data as T),
    } as BaseResponse<T>;
  } catch (error: any) {
    console.error(error);
    return {
      success: false,
      message: error.message || "Unknown error",
      code: error.code || 500,
      data: undefined,
    };
  }
}
