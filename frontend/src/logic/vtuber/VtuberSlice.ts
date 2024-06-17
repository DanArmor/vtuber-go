import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { addAllRestReducers, createRestActions, getDefaultRestState } from '../../store/restHelper';
import fp from 'lodash/fp';
import { Vtuber, VtuberOrg } from '../../types/types';
import { VtuberOrgsRestActions, VtuberSearchRestActions, VtuberSelectRestActions } from './VtuberActions';
import { PageMetaInfo, SelectionType } from '../../types/api';

type VtuberSearchPayload = {
    name?: string,
    orgs?: number[],
    waves?: number[],
    selected?: SelectionType,
    offset: number,
    page_size: number
};
type VtuberSearchResponse = {
    vtubers: Vtuber[],
    page_meta: PageMetaInfo
};


type VtuberOrgsPayload = void;
type VtuberOrgsResponse = VtuberOrg[];

type VtuberSelectPayload = {
    vtuber_id: number
}
type VtuberSelectResponse = {
    vtuber_id: number
    selected: boolean
}

const vtuberSearchRestActions = createRestActions<VtuberSearchResponse, VtuberSearchPayload>(VtuberSearchRestActions);
const vtuberOrgsRestActions = createRestActions<VtuberOrgsResponse, VtuberOrgsPayload>(VtuberOrgsRestActions);
const vtuberSelectRestActions = createRestActions<VtuberSelectResponse, VtuberSelectPayload>(VtuberSelectRestActions);

const VtuberRestActions = {
    vtuberSearch: vtuberSearchRestActions,
    vtuberOrgs: vtuberOrgsRestActions,
    vtuberSelect: vtuberSelectRestActions
};

type VtuberState = {
    vtubers: Vtuber[],
    page_meta: PageMetaInfo | undefined,
    orgs: VtuberOrg[]
}

const initialVtuberState: VtuberState = {
    vtubers: [],
    orgs: [],
    page_meta: undefined
}

const initialRestState = {
    vtuberSearch: getDefaultRestState(),
    vtuberOrgs: getDefaultRestState(),
    vtuberSelect: getDefaultRestState()
}

const vtuberSlice = createSlice({
    name: "vtuber",
    initialState: { ...initialVtuberState, ...initialRestState },
    reducers: {
        setVtubers(state, action: PayloadAction<Vtuber[]>) {
            state.vtubers = action.payload.map((val) => {
                return {
                    ...val,
                    isSelected: val.isSelected ?? !!val.edges.users
                }
            });
        },
        setOrgs(state, action: PayloadAction<VtuberOrg[]>) {
            state.orgs = action.payload;
        },
        setPageMeta(state, action: PayloadAction<PageMetaInfo | undefined>) {
            state.page_meta = action.payload;
        }
    },
    extraReducers: (builder) => fp.flow(addAllRestReducers<typeof VtuberRestActions>(VtuberRestActions))(builder)
})

const vtuberReducer = vtuberSlice.reducer;
const VtuberActions = { ...VtuberRestActions, ...vtuberSlice.actions }

export { vtuberReducer, VtuberActions };

