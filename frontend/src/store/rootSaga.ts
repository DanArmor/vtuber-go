import { call, spawn, put, all, race, take } from 'redux-saga/effects'
import { UserSaga } from '../logic/user/UserSagas';
import { UserActions } from '../logic/user/UserSlice';
import { updateAuthorizationHeader } from '../api';
import { VtuberSaga } from '../logic/vtuber/VtuberSagas';
import { useInitData } from '@vkruglikov/react-telegram-web-app';
import { useDispatch } from 'react-redux';
import { getTokenFromLocalStorage } from '../helpers/localStorage';

function* AppSaga() {
    const token: string | null = yield call(getTokenFromLocalStorage);
    if (token) {
        yield put(UserActions.setToken(token));
        yield call(updateAuthorizationHeader, token);
    }
}

function* rootSaga() {
    const sagas = [AppSaga, UserSaga, VtuberSaga];

    yield* sagas.map(saga => {
        return spawn(function* () {
            while (true) {
                try {
                    yield call(saga);
                    break;
                } catch (e) {
                    console.log("rootSaga error:", e);
                }
            }
        });
    });
}

export { rootSaga };