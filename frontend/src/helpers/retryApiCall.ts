import { ApiResponse } from "apisauce";
import { isResponseOk } from "./validateResponse";
import { call, delay, put, select } from "redux-saga/effects";
import { ActionCreatorWithPayload } from "@reduxjs/toolkit";
import { UserActions } from "../logic/user/UserSlice";
import { userInitDataSelector } from "../logic/user/UserSelectors";

export type RetryApiCallParams = {
    maxTries?: number;
    delayMs?: number;
    apiRequest: (...args: any[]) => Promise<ApiResponse<any>>;
    args?: any[];
    apiResponseValidator?: (resp: any) => boolean,
}

export function* retryApiCall(
    {
        maxTries = 2,
        delayMs = 1e3,
        apiRequest,
        args = [],
        apiResponseValidator = isResponseOk
    }: RetryApiCallParams
) {
    for (let i = 0; i < maxTries; i++) {
        const apiResponse: ApiResponse<any> = yield call(apiRequest, ...args);
        if (apiResponseValidator(apiResponse)) {
            return apiResponse;
        } else {
            if (i < (maxTries - 1)) {
                if (apiResponse.status === 401) {
                    const initData: string = yield select(userInitDataSelector);
                    if (initData && initData !== "") {
                        yield put(UserActions.auth.request(initData));
                    }
                }
                yield delay(delayMs);
            } else {
                return apiResponse;
            }
        }
    }
}