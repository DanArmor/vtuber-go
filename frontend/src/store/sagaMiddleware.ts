import createSagaMiddleware, {
  SagaMiddleware
} from 'redux-saga';

let sagaMiddleware: SagaMiddleware = createSagaMiddleware();

export function getSagaMiddleware(): SagaMiddleware {
  return sagaMiddleware;
}