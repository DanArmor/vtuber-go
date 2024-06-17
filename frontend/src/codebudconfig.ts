import { CodeBud, Instruction, InstructionGroup } from "@appklaar/codebud";
import { RootStore } from "./store/store";
import { NetworkInterceptorXHR } from "@appklaar/codebud/Network/NetworkInterceptorXHR";
import { RootState } from "./store/rootReducer";
import { apiKey } from "./secrets";

export const CodeBudSetup = (store: RootStore) => {
    console.log("Init of CodeBud");
    const instructions: (Instruction | InstructionGroup)[] = [
    ];
    CodeBud.init(apiKey, instructions, {
        Interceptor: NetworkInterceptorXHR,
    })
    console.log('CodeBud init:', CodeBud.isInit);
    console.log('CodeBud state:', CodeBud.state);
    function select(state: RootState) {
        return state;
    }
    const unsubscribeCodeBudFromStoreChanges = store.subscribe(
        CodeBud.createReduxStoreChangeHandler(store, select, 500),
    );

    CodeBud.enableLocalStorageMonitor(localStorage);
}