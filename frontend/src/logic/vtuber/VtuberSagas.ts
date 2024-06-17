import { Action } from "@reduxjs/toolkit";
import { VtuberActions } from "./VtuberSlice";
import { ApiResponse } from "apisauce";
import { call, put, select, takeLatest } from "redux-saga/effects";
import { api } from "../../api";
import { isResponseOk } from "../../helpers/validateResponse";
import { createErrorFromResponse } from "../../helpers/formatError";
import { VtuberOrgsResponse, VtuberSearchResponse, VtuberSelectResponse } from "../../types/api";
import { vtubersListSelector } from "./VtuberSelectors";
import { Vtuber } from "../../types/types";
import { retryApiCall } from "../../helpers/retryApiCall";

function* vtuberSearchRequest(action: Action) {
    if (VtuberActions.vtuberSearch.request.match(action)) {
        try {
            const vtuberSearchResponse: ApiResponse<VtuberSearchResponse> = yield call(
                retryApiCall,
                { apiRequest: api.vtuberSearchPost, args: [action.payload] }
            );
            if (isResponseOk(vtuberSearchResponse)) {
                const vtubers = vtuberSearchResponse.data!.result.vtubers;
                const oldVtubers: Vtuber[] = yield select(vtubersListSelector);
                yield put(VtuberActions.setPageMeta(vtuberSearchResponse.data!.result.page_meta));
                yield put(VtuberActions.setVtubers([...oldVtubers, ...vtubers]));
                yield put(VtuberActions.vtuberSearch.success(vtuberSearchResponse.data!.result));
            } else {
                throw createErrorFromResponse(vtuberSearchResponse);
            }
        } catch (error) {
            yield put(VtuberActions.vtuberSearch.failure(error));
        }
    }
}

function* vtuberOrgsRequest(action: Action) {
    if (VtuberActions.vtuberOrgs.request.match(action)) {
        try {
            const vtuberOrgsResponse: ApiResponse<VtuberOrgsResponse> = yield call(
                retryApiCall, { apiRequest: api.vtuberOrgsGet }
            );
            if (isResponseOk(vtuberOrgsResponse)) {
                const orgs = vtuberOrgsResponse.data!.result.orgs;
                yield put(VtuberActions.vtuberOrgs.success(orgs));
                yield put(VtuberActions.setOrgs(orgs));
            } else {
                throw createErrorFromResponse(vtuberOrgsResponse);
            }
        } catch (error) {
            yield put(VtuberActions.vtuberOrgs.failure(error));
        }
    }
}

function* vtuberSelectRequest(action: Action) {
    if (VtuberActions.vtuberSelect.request.match(action)) {
        try {
            const vtuberSelectResponse: ApiResponse<VtuberSelectResponse> = yield call(
                retryApiCall,
                { apiRequest: api.vtuberSelectPost, args: [action.payload] }
            );
            if (isResponseOk(vtuberSelectResponse)) {
                const { selected, vtuber_id } = vtuberSelectResponse.data!.result;
                const oldVtubers: Vtuber[] = yield select(vtubersListSelector);
                yield put(VtuberActions.setVtubers(oldVtubers.map((val) => {
                    if (val.id === vtuber_id) {
                        return { ...val, isSelected: selected };
                    } else {
                        return val;
                    }
                })));
                yield put(VtuberActions.vtuberSelect.success({ selected, vtuber_id }));
            } else {
                throw createErrorFromResponse(vtuberSelectResponse);
            }
        } catch (error) {
            yield put(VtuberActions.vtuberSelect.failure(error));
        }
    }
}

export function* VtuberSaga() {
    yield* [
        takeLatest(VtuberActions.vtuberSearch.request.type, vtuberSearchRequest),
        takeLatest(VtuberActions.vtuberOrgs.request.type, vtuberOrgsRequest),
        takeLatest(VtuberActions.vtuberSelect.request.type, vtuberSelectRequest),
    ];
}