import {
    ActionCreatorWithPayload,
    ActionCreatorWithPreparedPayload,
  } from '@reduxjs/toolkit';
  
  type RestActions = 'request' | 'success' | 'failure' | 'needUpdate';
  
  export type RestActionTypes = { [K in RestActions]: string };
  export type BaseFieldsType = { [key: string]: string };
  
  export type RestStateType<D = any> = {
    error?: any;
    data?: D;
    fetching: boolean;
    needUpdate: boolean;
  };
  
  export type NodeRestStateType<T extends string, R> = {
    [K in T]: RestStateType;
  } &
    R;
  
  export type CreateErrorPreparedAction = (
    data: any,
  ) => {
    payload: any;
    error: any;
  };
  
  export type RestActionsType = {
    request: ActionCreatorWithPayload<any, RestActionTypes['request']>;
    success: ActionCreatorWithPayload<any, RestActionTypes['success']>;
    needUpdate: ActionCreatorWithPayload<any, RestActionTypes['needUpdate']>;
    failure: ActionCreatorWithPreparedPayload<
      [any],
      any,
      RestActionTypes['failure'],
      any
    >;
  };