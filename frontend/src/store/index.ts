import { CONFIG } from '../config';
import {rootSaga} from './rootSaga';
import {getSagaMiddleware} from './sagaMiddleware';
import { store } from './store';

Promise.resolve().then(() => {
  getSagaMiddleware().run(rootSaga);
  if (CONFIG.USE_CODEBUD) {
    const codebud = require('./../codebudconfig');
    codebud.CodeBudSetup(store);
  }
});