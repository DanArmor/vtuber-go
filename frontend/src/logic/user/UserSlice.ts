import { createAction, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { addAllRestReducers, createRestActions, getDefaultRestState } from '../../store/restHelper';
import fp from 'lodash/fp';
import { AuthRestActions, TimezoneChangeRestActions, TimezoneGetRestActions } from './UserActions';

type AuthPayload = string;

type TimezoneChangePayload = number;

const authRestActions = createRestActions<void, AuthPayload>(AuthRestActions);
const timezoneGetRestActions = createRestActions<void, void>(TimezoneGetRestActions);
const timezoneChangeRestActions = createRestActions<void, TimezoneChangePayload>(TimezoneChangeRestActions);

const UserRestActions = {
    auth: authRestActions,
    timezone_get: timezoneGetRestActions,
    timezone_change: timezoneChangeRestActions
};

type UserState = {
    token: string | undefined,
    initData: string | undefined,
    timezone_shift: number | undefined
}

const initialUserState: UserState = {
    token: undefined,
    initData: undefined,
    timezone_shift: 0
}

const initialRestState = {
    auth: getDefaultRestState(),
    timezone_get: getDefaultRestState(),
    timezone_change: getDefaultRestState(),
}

const userSlice = createSlice({
    name: "user",
    initialState: { ...initialUserState, ...initialRestState },
    reducers: {
        setToken(state, action: PayloadAction<string | undefined>) {
            state.token = action.payload;
        },
        setInitData(state, action: PayloadAction<string | undefined>) {
            state.initData = action.payload;
        },
        setTimezoneShift(state, action: PayloadAction<number | undefined>) {
            state.timezone_shift = action.payload;
        }
    },
    extraReducers: (builder) => fp.flow(addAllRestReducers<typeof UserRestActions>(UserRestActions))(builder)
})

const userReducer = userSlice.reducer;
const UserActions = {
    ...UserRestActions,
    ...userSlice.actions,
}

export { userReducer, UserActions };

