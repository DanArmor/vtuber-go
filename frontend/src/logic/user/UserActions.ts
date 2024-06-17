import { restActionCreatorHelper } from './../../helpers/restActionCreatorHelper';

const userRestAction = restActionCreatorHelper(`user`);

export const AuthRestActions = userRestAction("auth");
