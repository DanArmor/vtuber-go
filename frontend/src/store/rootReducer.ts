import { vtuberReducer } from "../logic/vtuber/VtuberSlice";
import { userReducer } from "../logic/user/UserSlice";
import { combineReducers } from '@reduxjs/toolkit'

export const rootReducer = combineReducers({
    user: userReducer,
    vtuber: vtuberReducer
});

export type RootState = ReturnType<typeof rootReducer>;