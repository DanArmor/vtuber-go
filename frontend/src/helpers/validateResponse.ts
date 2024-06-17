import { ApiResponse } from "apisauce";
import { DefaultVtuberGoResult, VtuberGoError } from "../types/api";

export function isResponseOk<T extends DefaultVtuberGoResult<any> & VtuberGoError>(response: ApiResponse<T>) { return response.status === 200 && response.data?.result }