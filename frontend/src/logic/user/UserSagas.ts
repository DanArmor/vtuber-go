import { Action } from "@reduxjs/toolkit";
import { UserActions } from "./UserSlice";
import { ApiResponse } from "apisauce";
import { call, put, takeLatest } from "redux-saga/effects";
import { api, updateAuthorizationHeader } from "../../api";
import { isResponseOk } from "../../helpers/validateResponse";
import { createErrorFromResponse } from "../../helpers/formatError";
import { CodeBud } from "@appklaar/codebud";
import { AuthResponse } from "../../types/api";
import { setTokenToLocalStorage } from "../../helpers/localStorage";

function* authRequest(action: Action) {
  if (UserActions.auth.request.match(action)) {
    try {
      const authResponse: ApiResponse<AuthResponse> = yield call(
        api.authGet,
        action.payload
      );
      if (isResponseOk(authResponse)) {
        const token: string = authResponse.data!.result.token;
        yield put(UserActions.setToken(token));
        CodeBud.captureEvent("User got new token(auth)", { token: token });
        yield call(updateAuthorizationHeader, token)
        yield call(setTokenToLocalStorage, token);
        yield put(UserActions.auth.success());
      } else {
        throw createErrorFromResponse(authResponse);
      }
    } catch (error) {
      yield put(UserActions.auth.failure(error));
    }
  }
}

export function* UserSaga() {
  yield* [
    takeLatest(UserActions.auth.request.type, authRequest),
  ];
}