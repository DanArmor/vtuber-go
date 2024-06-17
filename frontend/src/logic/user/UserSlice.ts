import { createAction, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { addAllRestReducers, createRestActions, getDefaultRestState } from '../../store/restHelper';
import fp from 'lodash/fp';
import { AuthRestActions } from './UserActions';

type AuthPayload = string;

const authRestActions = createRestActions<void, AuthPayload>(AuthRestActions);

const UserRestActions = {
    auth: authRestActions,
};

type UserState = {
    token: string | undefined,
    initData: string | undefined,
}

const initialUserState: UserState = {
    token: undefined,
    initData: undefined
}

const initialRestState = {
    auth: getDefaultRestState(),
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

