import { ActionReducerMapBuilder, createAction } from '@reduxjs/toolkit';
import _ from 'lodash';
import {
  CreateErrorPreparedAction,
  RestActionsType,
  RestActionTypes,
  RestStateType,
} from './restHelper.d';

function createRestActions<
  SuccessPayload = void,
  RequestPayload = void,
  NeedUpdatePayload = void,
  T extends RestActionTypes = RestActionTypes
>(actions: T) {
  return {
    request: createAction<RequestPayload, T['request']>(actions.request),
    success: createAction<SuccessPayload, T['success']>(actions.success),
    needUpdate: createAction<NeedUpdatePayload, T['needUpdate']>(
      actions.needUpdate,
    ),
    failure: createAction<CreateErrorPreparedAction, T['failure']>(
      actions.failure,
      (data: any) => {
        return {
          payload: undefined,
          error: data,
        };
      },
    ),
  };
}

function getLens(path?: string, payloadId?: string) {
  function payloadLens(state: any, payload: any): RestStateType {
    if (path || (payloadId && _.hasIn(payload, `payload.${payloadId}`))) {
      const node = path ? _.get(state, path) : state;
      if (payloadId) {
        const currentNodeId = _.get(payload, `payload.${payloadId}`);
        if (!node[currentNodeId]) {
          node[currentNodeId] = getDefaultRestState();
        }
        return node[currentNodeId];
      }
      return node;
    }
    return state;
  }

  return function updater<A, S>(
    updateFunc: (state: RestStateType, action: A) => void,
  ) {
    return (state: S, action: A) => {
      updateFunc(payloadLens(state, action), action);
    };
  };
}

function createNodeRestReducers(
  builder: ActionReducerMapBuilder<any>,
  restActions: RestActionsType,
  path?: string,
  payloadId?: string,
): ActionReducerMapBuilder<any> {
  const lens = getLens(path, payloadId);

  return builder
    .addCase(
      restActions.request,
      lens(state => {
        state.fetching = true;
        state.error = undefined;
        state.needUpdate = false;
      }),
    )
    .addCase(
      restActions.success,
      lens((state, action) => {
        state.fetching = false;
        state.data = action.payload;
      }),
    )
    .addCase(
      restActions.failure,
      lens((state, action) => {
        state.fetching = false;
        state.error = action.error;
      }),
    )
    .addCase(
      restActions.needUpdate,
      lens(state => {
        state.needUpdate = true;
      }),
    );
}

function addRestReducers(
  actions: RestActionsType,
  path: string,
  payloadId?: string,
) {
  return (builder: ActionReducerMapBuilder<any>) =>
    createNodeRestReducers(builder, actions, path, payloadId);
}

function getDefaultRestState<D>(defaultData: any = undefined): RestStateType<D> {
  return {
    data: defaultData,
    fetching: false,
    error: undefined,
    needUpdate: true,
  };
}

function needUpdateSelector(state: any, path: string, id?: string): boolean {
  const nodePath = id ? `${path}.${id}` : path;
  const node = _.get(state, nodePath) as RestStateType;

  if (node) {
    return node.needUpdate || (node.data === undefined && !node.fetching);
  }
  return true;
}

function addAllRestReducers<T extends { [key: string]: any }>(restActionsTable: T) {
  return (Object.keys(restActionsTable) as Array<string>).map((path) => addRestReducers(restActionsTable[path], path));
}

export {
  createRestActions,
  addRestReducers,
  createNodeRestReducers,
  getDefaultRestState,
  needUpdateSelector,
  addAllRestReducers
};