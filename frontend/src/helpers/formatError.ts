import { ApiResponse } from "apisauce";
import { DefaultVtuberGoResult, VtuberGoError } from "../types/api";

export type VtuberGoErrorExtended = {
    httpCode?: number,
    code?: number,
    msg?: string
    description?: string
}

export function createErrorFromResponse<T extends VtuberGoError>(response: ApiResponse<T>, description?: string): VtuberGoErrorExtended {
    if (response.data) {
        return {
            httpCode: response.status,
            code: response.data?.error.code,
            msg: response.data?.error.msg,
            description
        }
    } else {
        return {
            code: 0,
            httpCode: response.status,
            msg: response.originalError?.stack,
            description: response.originalError?.message
        }
    }
}