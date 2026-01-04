"use client";

import { useActionState, useEffect } from "react";
import { toast } from "sonner";
import { BaseResponse } from "@/lib/response/api-response";

export type Action<TInput, TOutput> = (
  prevState: any,
  data: TInput
) => Promise<BaseResponse<TOutput>>;

export function useAction<TInput, TOutput>(
  action: Action<TInput, TOutput>,
  options?: {
    onSuccess?: (data: TOutput) => void;
    onError?: (error: string) => void;
  }
) {
  const [state, dispatch, isPending] = useActionState(action, {
    success: false,
    code: 0,
    data: undefined,
  } as BaseResponse<TOutput>);

  useEffect(() => {
    if (state.success && state.data) {
      options?.onSuccess?.(state.data);
    } else if (!state.success && state.message && state.code !== 0) {
      options?.onError?.(state.message);
      toast.error(state.message);
    }
  }, [state, options]);

  return {
    execute: dispatch,
    result: state,
    isPending,
  };
}
