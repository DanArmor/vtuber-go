import { decodeToken } from 'react-jwt';
import { RootState } from './../../store/rootReducer';

export const userSelector = (state: RootState) => state.user;
export const userInitDataSelector = (state: RootState) => state.user.initData;