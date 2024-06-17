import { configureStore } from "@reduxjs/toolkit";
import { getSagaMiddleware } from "./sagaMiddleware";
import { rootReducer } from "./rootReducer";
import { CONFIG } from "../config";
import { CodeBud } from "@appklaar/codebud";

function getOurMiddleware() {
    const middleware = [];
    middleware.push(getSagaMiddleware());
    if (CONFIG.USE_CODEBUD) {
        middleware.push(CodeBud.createReduxActionMonitorMiddleware());
    }
    return middleware;
}

function getOurEnhancers() {
    const enhancers: any[] = [];

    return enhancers;
}

export const store = configureStore({
    reducer: rootReducer,
    middleware: (gDM) => gDM().concat(getOurMiddleware()),
    enhancers: (gDE) => gDE().concat(getOurEnhancers())
})

export type RootStore = typeof store;