import { RootState } from './../../store/rootReducer';

export const vtuberSelector = (state: RootState) => state.vtuber;
export const vtubersListSelector = (state: RootState) => state.vtuber.vtubers;