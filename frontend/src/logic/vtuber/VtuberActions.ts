import { restActionCreatorHelper } from './../../helpers/restActionCreatorHelper';

const vtuberRestAction = restActionCreatorHelper(`vtuber`);

export const VtuberSearchRestActions = vtuberRestAction("vtuberSearch");

export const VtuberOrgsRestActions = vtuberRestAction("vtuberOrgs");

export const VtuberSelectRestActions = vtuberRestAction("vtuberSelect");
